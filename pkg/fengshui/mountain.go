package fengshui

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/luoshu"
)

// Mountain24 represents one of the 24 directional mountains (二十四山).
// Each mountain spans 15° of the compass (360° / 24 = 15°).
// Attribute fields are domain symbols; text is produced only at the boundary.
type Mountain24 struct {
	Index    int            `json:"index"` // 0-23, starting from 子=0°
	Name     string         `json:"name"`  // 24山标识（查表 key，不参与运算）
	Angle    int            `json:"angle"` // degrees, 0-345 in 15° steps
	Element  ganzhi.Wuxing  `json:"wuxing"`
	YinYang  ganzhi.YinYang `json:"-"`
	Trigram  luoshu.Trigram `json:"-"`
	YuanLong YuanLong       `json:"-"`
}

//go:embed data/mountains24.json
var mountains24JSON []byte

var (
	mountains24Once sync.Once
	mountains24Err  error
	// Mountains24Table holds all 24 mountains in order, from 子=0° (North).
	Mountains24Table [24]Mountain24
)

// Load parses the embedded 24-mountain table (idempotent, fail-closed).
func Load() error {
	mountains24Once.Do(func() {
		var rows []struct {
			Index    int    `json:"index"`
			Name     string `json:"name"`
			Angle    int    `json:"angle"`
			Wuxing   string `json:"wuxing"`
			YinYang  string `json:"yin_yang"`
			Gua      string `json:"gua"`
			YuanLong string `json:"yuan_long"`
		}
		if err := json.Unmarshal(mountains24JSON, &rows); err != nil {
			mountains24Err = fmt.Errorf("fengshui: parse mountains24: %w", err)
			return
		}
		if len(rows) != 24 {
			mountains24Err = fmt.Errorf("fengshui: mountains24 has %d rows, want 24", len(rows))
			return
		}
		for _, row := range rows {
			if row.Index < 0 || row.Index >= 24 {
				mountains24Err = fmt.Errorf("fengshui: mountain index %d out of range", row.Index)
				return
			}
			elem, err := ganzhi.EncodeWuxing(row.Wuxing)
			if err != nil {
				mountains24Err = fmt.Errorf("fengshui: mountain %s unknown wuxing %q", row.Name, row.Wuxing)
				return
			}
			yinYang, err := ganzhi.EncodeYinYang(row.YinYang)
			if err != nil {
				mountains24Err = fmt.Errorf("fengshui: mountain %s: %w", row.Name, err)
				return
			}
			trigram, ok := luoshu.EncodeTrigram(row.Gua)
			if !ok {
				mountains24Err = fmt.Errorf("fengshui: mountain %s unknown gua %q", row.Name, row.Gua)
				return
			}
			yuanLong, ok := EncodeYuanLong(row.YuanLong)
			if !ok {
				mountains24Err = fmt.Errorf("fengshui: mountain %s unknown yuan_long %q", row.Name, row.YuanLong)
				return
			}
			Mountains24Table[row.Index] = Mountain24{
				Index:    row.Index,
				Name:     row.Name,
				Angle:    row.Angle,
				Element:  elem,
				YinYang:  yinYang,
				Trigram:  trigram,
				YuanLong: yuanLong,
			}
		}
	})
	return mountains24Err
}
