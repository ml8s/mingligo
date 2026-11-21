package huangli

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ml8s/mingligo/pkg/ganzhi"

	"gopkg.in/yaml.v3"
)

var huangDaoStars [12]huangDaoStar
var qingLongStart map[ganzhi.Zhi]ganzhi.Zhi

//go:embed data/jianchu.yaml
var jianchuYAML []byte

//go:embed data/huangdao.json
var huangdaoJSON []byte

var jianChuCfg jianchuConfig

var (
	loadOnce sync.Once
	loadErr  error
)

// Load parses the embedded huangli tables. Idempotent; call from main/TestMain.
func Load() error {
	loadOnce.Do(func() {
		if err := ganzhi.Load(); err != nil {
			loadErr = err
			return
		}
		for _, load := range []struct {
			name string
			fn   func() error
		}{
			{"jianchu", loadJianchu},
			{"huangdao", loadHuangdao},
		} {
			if err := load.fn(); err != nil {
				loadErr = fmt.Errorf("huangli: load %s: %w", load.name, err)
				return
			}
		}
	})
	return loadErr
}

func loadJianchu() error {
	// Parse into anonymous struct — yaml.v3 skips map values with named types.
	var raw struct {
		Jianchu struct {
			Sequence   []string            `yaml:"sequence"`
			Suitable   map[string][]string `yaml:"suitable"`
			Forbidden  map[string][]string `yaml:"forbidden"`
			EventRules map[string]struct {
				Label     string   `yaml:"label"`
				Suitable  []string `yaml:"suitable"`
				Forbidden []string `yaml:"forbidden"`
			} `yaml:"event_rules"`
		} `yaml:"jianchu"`
	}
	if err := yaml.Unmarshal(jianchuYAML, &raw); err != nil {
		return err
	}

	cfg := raw.Jianchu
	c := &jianChuCfg
	seq, err := encodeJianChuList(cfg.Sequence)
	if err != nil {
		return err
	}
	c.Sequence = seq
	c.Suitable = make(map[JianChu][]string, len(cfg.Suitable))
	for k, v := range cfg.Suitable {
		jc, err := EncodeJianChu(k)
		if err != nil {
			return err
		}
		c.Suitable[jc] = v
	}
	c.Forbidden = make(map[JianChu][]string, len(cfg.Forbidden))
	for k, v := range cfg.Forbidden {
		jc, err := EncodeJianChu(k)
		if err != nil {
			return err
		}
		c.Forbidden[jc] = v
	}

	c.EventRules = make(map[string]eventRule, len(cfg.EventRules)*2)
	for k, v := range cfg.EventRules {
		suitable, err := encodeJianChuList(v.Suitable)
		if err != nil {
			return err
		}
		forbidden, err := encodeJianChuList(v.Forbidden)
		if err != nil {
			return err
		}
		r := eventRule{Label: v.Label, Suitable: suitable, Forbidden: forbidden}
		c.EventRules[k] = r
		c.EventRules[v.Label] = r
	}

	return nil
}

func encodeJianChuList(values []string) ([]JianChu, error) {
	out := make([]JianChu, 0, len(values))
	for _, v := range values {
		jc, err := EncodeJianChu(v)
		if err != nil {
			return nil, err
		}
		out = append(out, jc)
	}
	return out, nil
}

func loadHuangdao() error {
	var data struct {
		Stars []struct {
			Name string `json:"name"`
			Path string `json:"path"`
		} `json:"stars"`
		QingLongStart map[string]string `json:"qing_long_start"`
	}
	if err := json.Unmarshal(huangdaoJSON, &data); err != nil {
		return err
	}

	for i, s := range data.Stars {
		huangDaoStars[i] = huangDaoStar{
			Index:    i,
			Name:     s.Name,
			Path:     s.Path,
			Sequence: i,
		}
	}

	qingLongStart = make(map[ganzhi.Zhi]ganzhi.Zhi, len(data.QingLongStart))
	for zhiStr, startStr := range data.QingLongStart {
		zhi, err := ganzhi.EncodeZhi(zhiStr)
		if err != nil {
			return err
		}
		start, err := ganzhi.EncodeZhi(startStr)
		if err != nil {
			return err
		}
		qingLongStart[zhi] = start
	}

	return nil
}
