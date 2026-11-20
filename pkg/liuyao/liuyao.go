package liuyao

import (
	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// Line is a single爻 in a hexagram.
type Line struct {
	Position int           `json:"position"` // 1-6, bottom to top
	Type     YaoType       `json:"yao_value"`
	Gan      ganzhi.Gan    `json:"gan"`
	Zhi      ganzhi.Zhi    `json:"zhi"`
	Wuxing   ganzhi.Wuxing `json:"wuxing"`
	LiuQin   LiuQin        `json:"liu_qin"`
	ShiYing  ShiYing       `json:"shi_ying"`
	LiuShou  LiuShou       `json:"liu_shou"`
	// Deterministic derived states (filled by computeChart):
	YuePo         bool         `json:"yue_po,omitempty"`          // 月破（月建冲该爻地支）
	DongSelf      bool         `json:"dong_self,omitempty"`       // 本爻发动（老阴/老阳）
	DongSheng     bool         `json:"dong_sheng,omitempty"`      // 有动爻生此爻
	DongKe        bool         `json:"dong_ke,omitempty"`         // 有动爻克此爻
	XunKong       bool         `json:"xun_kong,omitempty"`        // 该爻地支值日柱旬空
	MuKu          bool         `json:"mu_ku,omitempty"`           // 入墓（地支为墓库）
	MuKuBranch    string       `json:"mu_ku_zhi,omitempty"`       // 本五行墓库地支
	MuKuElement   string       `json:"mu_ku_wuxing,omitempty"`    // 被入库的五行
	MuKuTypes     []MuKuSource `json:"mu_ku_types,omitempty"`     // day/moving/transformed
	ChangShengYue LifeStage    `json:"chang_sheng_yue,omitempty"` // 本爻五行在月支的十二长生
}

// guaMeta holds static data for a hexagram.
type guaMeta struct {
	Name      string `json:"name"`     // 卦名
	PalaceIdx int    `json:"gong_idx"` // 0-7, which palace
	ShiPos    int    `json:"shi_pos"`  // 1-6,世爻 position
}

// Chart is the complete 六爻排盘 with all analysis layers.
type Chart struct {
	Name         string        `json:"gua_name"`
	BenGua       guaIndex      `json:"-"`
	BianGua      guaIndex      `json:"bian_gua,omitempty"` // 0 if no change
	Palace       string        `json:"gong"`
	PalaceWuxing ganzhi.Wuxing `json:"gong_wuxing"`
	Lines        [6]Line       `json:"lines"`
	BianLines    [6]Line       `json:"bian_yao,omitempty"`
	RiGan        ganzhi.Gan    `json:"ri_chen_gan"`
	RiZhi        ganzhi.Zhi    `json:"ri_chen_zhi"`
	YueZhi       ganzhi.Zhi    `json:"yue_jian_zhi"`
	YueGan       ganzhi.Gan    `json:"yue_jian_gan"`
	DongYao      []int         `json:"-"` // 动爻位置 1-6；RPC 审计以完整 casting 收据为准
	// Analysis layers set by ComputeChart.
	YongShen              YongShenResult         `json:"yong_shen"`
	GuaCi                 GuaCi                  `json:"gua_ci,omitempty"`
	WangShuai             [6]ganzhi.WangShuai    `json:"wang_shuai"`
	DayRelations          [6]DayRelation         `json:"ri_chen_relations"`
	XunKong               [2]ganzhi.Zhi          `json:"xun_kong"` // 日柱旬空地支（六爻断卦：用神旬空则事虚，出空填实方应）
	TimingCandidates      []TimingCandidate      `json:"timing_candidates,omitempty"`
	Patterns              []Pattern              `json:"ge_ju,omitempty"`              // 特殊格局
	DongYaoRelations      []DongYaoRelation      `json:"dong_yao_relations,omitempty"` // 动爻与用神的关系
	HiddenLines           []HiddenLine           `json:"hidden_lines,omitempty"`       // 本宫全部伏神
	BranchRelationFacts   []BranchRelationFact   `json:"zhi_relation_facts,omitempty"`
	Casting               *Casting               `json:"casting,omitempty"` // 完整起卦收据（audit only）
	TombSchool            TombSchool             `json:"tomb_school"`
	DayClashFacts         []DayClashFact         `json:"day_clash_facts,omitempty"`
	MovingTransformations []MovingTransformation `json:"moving_transformations,omitempty"`
	SanHeCandidates       []SanHeCandidate       `json:"san_he_candidates,omitempty"`
	ForceChain            *ForceChain            `json:"yong_shen_force_chain,omitempty"`
	YongShenCandidates    []YongShenCandidate    `json:"yong_shen_candidates,omitempty"`
	Conflicts             []Conflict             `json:"conflicts,omitempty"`
}

// TombSchool makes the earth-stem tomb convention explicit. Other schools
// choose Xu for earth; they must not be mixed into this chart silently.
type TombSchool struct {
	EarthBranch string `json:"di_zhi"`
	School      string `json:"school"`
}

// palaceNames.
var palaceNames = [8]string{"乾", "兑", "离", "震", "巽", "坎", "艮", "坤"}

// palaceWuxing maps palace index → five element (for六亲).
var palaceWuxing = [8]ganzhi.Wuxing{
	ganzhi.WxJin,  // 乾=金
	ganzhi.WxJin,  // 兑=金
	ganzhi.WxHuo,  // 离=火
	ganzhi.WxMu,   // 震=木
	ganzhi.WxMu,   // 巽=木
	ganzhi.WxShui, // 坎=水
	ganzhi.WxTu,   // 艮=土
	ganzhi.WxTu,   // 坤=土
}
