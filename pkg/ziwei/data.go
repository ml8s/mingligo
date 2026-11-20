package ziwei

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

//go:embed data/tables.json
var tablesJSON []byte

//go:embed data/miao_wang.json
var miaoWangJSON []byte

//go:embed data/xiaoxian_rules.json
var xiaoxianRulesJSON []byte

var (
	earthlyIdxTable = [12]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	nianStars       = map[ShenSha][12]int{
		ssHuaGai:   {4, 1, 10, 7, 4, 1, 10, 7, 4, 1, 10, 7},
		ssXianChi:  {9, 6, 3, 0, 9, 6, 3, 0, 9, 6, 3, 0},
		ssGuChen:   {2, 2, 5, 5, 5, 8, 8, 8, 11, 11, 11, 2},
		ssGuaSu:    {10, 10, 1, 1, 1, 4, 4, 4, 7, 7, 7, 10},
		ssPoSui:    {5, 1, 9, 5, 1, 9, 5, 1, 9, 5, 1, 9},
		ssFeiLian:  {8, 9, 10, 5, 6, 7, 2, 3, 4, 11, 0, 1},
		ssLongChi:  {4, 5, 6, 7, 8, 9, 10, 11, 0, 1, 2, 3},
		ssFengGe:   {10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 11},
		ssTianKu:   {6, 5, 4, 3, 2, 1, 0, 11, 10, 9, 8, 7},
		ssTianXu:   {6, 7, 8, 9, 10, 11, 0, 1, 2, 3, 4, 5},
		ssTianKong: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0},
		ssTianDe:   {9, 10, 11, 0, 1, 2, 3, 4, 5, 6, 7, 8},
		ssYueDe:    {5, 6, 7, 8, 9, 10, 11, 0, 1, 2, 3, 4},
	}
	yueStars = map[ShenSha][12]int{
		ssTianYao:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0},
		ssTianXing: {9, 10, 11, 0, 1, 2, 3, 4, 5, 6, 7, 8},
		ssYinSha:   {2, 0, 10, 8, 6, 4, 2, 0, 10, 8, 6, 4},
		ssJieShen:  {8, 8, 10, 10, 0, 0, 2, 2, 4, 4, 6, 6},
		ssTianYue:  {10, 5, 4, 2, 7, 3, 11, 7, 2, 6, 10, 2},
		ssTianWu:   {5, 8, 2, 11, 5, 8, 2, 11, 5, 8, 2, 11},
	}
	shiStars = map[ShenSha][12]int{
		ssTaiFu:   {6, 7, 8, 9, 10, 11, 0, 1, 2, 3, 4, 5},
		ssFengGao: {2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0, 1},
	}
	ganStars = map[ShenSha][10]int{
		ssTianChu:  {5, 6, 0, 5, 6, 8, 2, 6, 9, 11},
		ssTianGuan: {7, 4, 5, 2, 3, 9, 11, 9, 10, 6},
		ssTianFu:   {9, 8, 0, 11, 3, 2, 6, 5, 6, 5},
		ssJieLu:    {8, 6, 4, 2, 0, 8, 6, 4, 2, 0},
		ssKongWang: {9, 7, 5, 3, 1, 9, 7, 5, 3, 1},
	}
	siHuaTable    map[Gan][4]StarIndex
	ziweiStartPos map[JuShu]int
	luCunTable    [10]int
	tianKuiTable  [10]int
	tianYueTable  [10]int
	tianMaTable   [12]int
)

var (
	loadOnce sync.Once
	loadErr  error
)

// Load parses the embedded ziwei tables. Idempotent; call from main/TestMain.
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
			{"tables", loadTables},
			{"xiaoxian rules", loadXiaoXianRules},
		} {
			if err := load.fn(); err != nil {
				loadErr = fmt.Errorf("ziwei: load %s: %w", load.name, err)
				return
			}
		}
	})
	return loadErr
}

// xiaoxianStartByBranch 小限起宫：生年支 → 起宫安星序（xiaoxian_rules.json）。
var xiaoxianStartByBranch map[ganzhi.Zhi]int

func loadXiaoXianRules() error {
	var raw struct {
		StartPalaces []struct {
			Triad []string `json:"triad"`
			Start string   `json:"start"`
		} `json:"start_palaces"`
	}
	if err := json.Unmarshal(xiaoxianRulesJSON, &raw); err != nil {
		return err
	}
	xiaoxianStartByBranch = make(map[ganzhi.Zhi]int)
	for _, p := range raw.StartPalaces {
		startZhi, err := ganzhi.EncodeZhi(p.Start)
		if err != nil {
			return err
		}
		startIdx := zhiIdxToAnXingIdx(zhiToZhiIdx(startZhi))
		for _, b := range p.Triad {
			zhi, err := ganzhi.EncodeZhi(b)
			if err != nil {
				return err
			}
			xiaoxianStartByBranch[zhi] = startIdx
		}
	}
	return nil
}

func loadTables() error {
	var data struct {
		SiHua      map[string]map[string]string `json:"si_hua"`
		ZiweiStart map[string]int               `json:"ziwei_start"`
		LuCun      map[string]int               `json:"lu_cun"`
		TianKui    map[string]int               `json:"tian_kui"`
		TianYue    map[string]int               `json:"tian_yue"`
		TianMa     map[string]int               `json:"tian_ma"`
	}
	if err := json.Unmarshal(tablesJSON, &data); err != nil {
		return err
	}

	// siHuaTable
	siHuaTable = make(map[Gan][4]StarIndex, 10)
	for ganName, h := range data.SiHua {
		gan, err := ganzhi.EncodeGan(ganName)
		if err != nil {
			return err
		}
		stars := [4]StarIndex{
			nameToStar(h["hua_lu"]), nameToStar(h["hua_quan"]),
			nameToStar(h["hua_ke"]), nameToStar(h["hua_ji"]),
		}
		siHuaTable[gan] = stars
	}

	// ziweiStartPos
	juNameToJuShu := map[string]JuShu{
		"水二局": JuWater, "木三局": JuWood, "金四局": JuMetal,
		"土五局": JuEarth, "火六局": JuFire,
	}
	ziweiStartPos = make(map[JuShu]int, 5)
	for name, pos := range data.ZiweiStart {
		js, ok := juNameToJuShu[name]
		if !ok {
			return fmt.Errorf("ziwei: unknown JuShu name %q", name)
		}
		ziweiStartPos[js] = pos
	}

	// luCunTable
	for ganName, pos := range data.LuCun {
		gan, err := ganzhi.EncodeGan(ganName)
		if err != nil {
			return err
		}
		luCunTable[int(gan)-1] = pos
	}

	// tianKuiTable / tianYueTable
	for ganName, pos := range data.TianKui {
		gan, err := ganzhi.EncodeGan(ganName)
		if err != nil {
			return err
		}
		tianKuiTable[int(gan)-1] = pos
	}
	for ganName, pos := range data.TianYue {
		gan, err := ganzhi.EncodeGan(ganName)
		if err != nil {
			return err
		}
		tianYueTable[int(gan)-1] = pos
	}

	// miaoWangTable
	var mwData struct {
		MiaoWang [][]string `json:"miao_wang"`
	}
	if err := json.Unmarshal(miaoWangJSON, &mwData); err != nil {
		return err
	}
	for i, row := range mwData.MiaoWang {
		for j, s := range row {
			miaoWangTable[i][j] = brightnessFrom(s)
		}
	}

	// tianMaTable — iztro formula if JSON empty
	if len(data.TianMa) == 0 {
		// 寅午戌→申, 申子辰→寅, 巳酉丑→亥, 亥卯未→巳 (zhiIdx)
		tianMaTable = [12]int{2, 11, 8, 5, 2, 11, 8, 5, 2, 11, 8, 5}
	} else {
		for zhiName, pos := range data.TianMa {
			zhi, err := ganzhi.EncodeZhi(zhiName)
			if err != nil {
				return err
			}
			tianMaTable[int(zhi)-1] = pos
		}
	}

	return nil
}

func nameToStar(s string) StarIndex {
	for si, name := range starNames {
		if name == s {
			return si
		}
	}
	return 0
}
