package bazi

import "encoding/json"

// This file defines the domain-owned, domain-aggregated output view of a 合婚
// (bazi.bond) result. jsonschema tags are inert metadata read only when a
// consumer reflects the type for its schema.

type BondView struct {
	DayMaster    DayMasterView    `json:"ri_zhu" jsonschema:"双方日主及日主生克关系。"`
	SpousePalace SpousePalaceView `json:"pei_ou_gong" jsonschema:"双方配偶宫（日支）及其关系；此为日支专用视图，与 zhu_cross 的 4×4 全交叉分工不同（非重复）。"`
	ZhuCross     ZhuCrossView     `json:"zhu_cross" jsonschema:"双方四柱逐柱干支交叉。"`
	ShiShenCross ShiShenCrossView `json:"shi_shen_cross" jsonschema:"双方日主互视的十神柱位。"`
	WuxingCross  WuxingCrossView  `json:"wuxing_cross" jsonschema:"双方实盘五行数量与扶抑对照。"`
	NaYinCross   NaYinCrossView   `json:"na_yin_cross" jsonschema:"双方四柱纳音交叉。"`
	SpouseStar   SpouseStarView   `json:"pei_ou_xing" jsonschema:"双方配偶星候选与落柱。"`
	ShenshaCross ShenshaCrossView `json:"shen_sha_cross" jsonschema:"双方互见神煞事实。"`
}

type DayMasterView struct {
	A        SideIdentityView `json:"a" jsonschema:"甲方日主。"`
	B        SideIdentityView `json:"b" jsonschema:"乙方日主。"`
	Relation RelationView     `json:"relation" jsonschema:"日主关系。"`
}

type SideIdentityView struct {
	Gan    string `json:"gan" jsonschema:"天干。"`
	Wuxing string `json:"wuxing" jsonschema:"五行。"`
	Gender string `json:"gender" jsonschema:"性别。"`
}

type RelationView struct {
	Kind     string `json:"kind" jsonschema:"关系类型。"`
	Relation string `json:"relation" jsonschema:"关系描述。"`
}

type SpousePalaceView struct {
	A        ZhiSideView           `json:"a" jsonschema:"甲方日支。"`
	B        ZhiSideView           `json:"b" jsonschema:"乙方日支。"`
	Relation ExplainedRelationView `json:"relation" jsonschema:"日支关系。"`
}

type ZhiSideView struct {
	Zhi string `json:"zhi" jsonschema:"日支。"`
}

type ExplainedRelationView struct {
	Kind        string `json:"kind" jsonschema:"关系类型。"`
	Explanation string `json:"explanation" jsonschema:"关系解释。"`
}

type ZhuCrossView struct {
	Pairs []ZhuCrossPairsItemView `json:"pairs" jsonschema:"双方四柱两两全交叉（4×4）；按双方各柱 a→b 配对。"`
}

type ZhuCrossPairsItemView struct {
	A           ZhuSideView     `json:"a" jsonschema:"甲方柱。"`
	B           ZhuSideView     `json:"b" jsonschema:"乙方柱。"`
	GanRelation RelationView    `json:"gan_relation" jsonschema:"天干关系。"`
	ZhiRelation ZhiRelationView `json:"zhi_relation" jsonschema:"地支关系。"`
}

type ZhuSideView struct {
	Zhu string `json:"zhu" jsonschema:"柱位。"`
	Gan string `json:"gan" jsonschema:"天干。"`
	Zhi string `json:"zhi" jsonschema:"地支。"`
}

type ZhiRelationView struct {
	Kind        string `json:"kind" jsonschema:"关系类型。"`
	Explanation string `json:"explanation" jsonschema:"关系解释。"`
}

type ShiShenCrossView struct {
	AToB TenGodPillarsView `json:"a_to_b" jsonschema:"甲方日主视乙方四柱的十神。"`
	BToA TenGodPillarsView `json:"b_to_a" jsonschema:"乙方日主视甲方四柱的十神。"`
}

type TenGodPillarsView struct {
	NianGan string `json:"nian_gan" jsonschema:"年干。"`
	YueGan  string `json:"yue_gan" jsonschema:"月干。"`
	RiGan   string `json:"ri_gan" jsonschema:"日干。"`
	ShiGan  string `json:"shi_gan" jsonschema:"时干。"`
}

type WuxingCrossView struct {
	A           ElementCountsView          `json:"a" jsonschema:"甲方计数或候选。"`
	B           ElementCountsView          `json:"b" jsonschema:"乙方计数或候选。"`
	Combined    ElementCountsView          `json:"combined" jsonschema:"双方合并计数。"`
	FuYiInOther WuxingCrossFuYiInOtherView `json:"fu_yi_in_other" jsonschema:"一方扶抑候选在另一方实盘中的命中数量。"`
}

type ElementCountsView struct {
	Mu   int `json:"木" jsonschema:"木的数量。"`
	Huo  int `json:"火" jsonschema:"火的数量。"`
	Tu   int `json:"土" jsonschema:"土的数量。"`
	Jin  int `json:"金" jsonschema:"金的数量。"`
	Shui int `json:"水" jsonschema:"水的数量。"`
}

type WuxingCrossFuYiInOtherView struct {
	A FuYiInOtherSideView `json:"a" jsonschema:"甲方计数或候选。"`
	B FuYiInOtherSideView `json:"b" jsonschema:"乙方计数或候选。"`
}

type FuYiInOtherSideView struct {
	ModelID     string `json:"fu_yi_model_id" jsonschema:"扶抑候选模型 ID。"`
	Yong        string `json:"yong" jsonschema:"用神；缺省为 无候选。"`
	Xi          string `json:"xi" jsonschema:"喜神；缺省为 无候选。"`
	Ji          string `json:"ji" jsonschema:"忌神；缺省为 无候选。"`
	YongInOther int    `json:"yong_in_other" jsonschema:"对方用神五行在本方实盘中的数量。"`
	XiInOther   int    `json:"xi_in_other" jsonschema:"对方喜神五行在本方实盘中的数量。"`
	JiInOther   int    `json:"ji_in_other" jsonschema:"对方忌神五行在本方实盘中的数量。"`
}

type NaYinCrossView struct {
	Pairs        []NaYinCrossPairsItemView  `json:"pairs" jsonschema:"双方四柱两两全交叉（4×4）；按双方各柱 a→b 配对。"`
	WuxingCounts NaYinCrossWuxingCountsView `json:"wuxing_counts" jsonschema:"双方纳音五行计数。"`
}

type NaYinCrossPairsItemView struct {
	A        NaYinSideView `json:"a" jsonschema:"甲方柱。"`
	B        NaYinSideView `json:"b" jsonschema:"乙方柱。"`
	Relation string        `json:"relation" jsonschema:"纳音关系。"`
}

type NaYinSideView struct {
	Zhu   string `json:"zhu" jsonschema:"柱位。"`
	NaYin string `json:"na_yin" jsonschema:"纳音。"`
}

type NaYinCrossWuxingCountsView struct {
	A map[string]int `json:"a" jsonschema:"甲方纳音五行计数。"`
	B map[string]int `json:"b" jsonschema:"乙方纳音五行计数。"`
}

type SpouseStarView struct {
	A SpouseStarSideView `json:"a" jsonschema:"甲方配偶星。"`
	B SpouseStarSideView `json:"b" jsonschema:"乙方配偶星。"`
}

type SpouseStarSideView struct {
	Gender    string                  `json:"gender" jsonschema:"性别。"`
	Primary   SpouseStarCandidateView `json:"primary" jsonschema:"主配偶星候选。"`
	Secondary SpouseStarCandidateView `json:"secondary" jsonschema:"次配偶星候选。"`
}

type SpouseStarOccurrenceView struct {
	Pillar string `json:"zhu" jsonschema:"柱位。"`
	Branch string `json:"zhi" jsonschema:"地支。"`
	Stem   string `json:"gan" jsonschema:"天干。"`
	Source string `json:"source" jsonschema:"来源。"`
	IsVoid bool   `json:"is_void" jsonschema:"是否旬空。"`
}

type SpouseStarCandidateView struct {
	TenGod      string                     `json:"shi_shen" jsonschema:"配偶星十神。"`
	Occurrences []SpouseStarOccurrenceView `json:"occurrences" jsonschema:"出现位置。"`
}

type ShenshaCrossView struct {
	TianYi   MutualFactView `json:"tian_yi" jsonschema:"天乙贵人互见。"`
	Lu       MutualFactView `json:"lu" jsonschema:"禄神互见。"`
	TaoHua   MutualFactView `json:"tao_hua" jsonschema:"桃花互见。"`
	YiMa     MutualFactView `json:"yi_ma" jsonschema:"驿马互见。"`
	KongWang MutualFactView `json:"kong_wang" jsonschema:"空亡。"`
	KuiGang  MutualFactView `json:"kui_gang" jsonschema:"魁罡互见。"`
	RiDe     MutualFactView `json:"ri_de" jsonschema:"日德互见。"`
	RiGui    MutualFactView `json:"ri_gui" jsonschema:"日贵互见。"`
}

type MutualFactView struct {
	AToB bool `json:"a_to_b" jsonschema:"甲方神煞是否见于乙方。"`
	BToA bool `json:"b_to_a" jsonschema:"乙方神煞是否见于甲方。"`
}

// View produces the domain-aggregated output from the computed bond. The Bond
// is a pure output value (no computation), so a JSON roundtrip maps its
// structurally identical fields into the tagged View contract.
func (b Bond) View() BondView {
	raw, err := json.Marshal(b)
	if err != nil {
		return BondView{}
	}
	var view BondView
	if err := json.Unmarshal(raw, &view); err != nil {
		return BondView{}
	}
	// View boundary: map machine-code source values to 命理 display terms.
	for _, side := range []*SpouseStarSideView{&view.SpouseStar.A, &view.SpouseStar.B} {
		for _, cand := range []*SpouseStarCandidateView{&side.Primary, &side.Secondary} {
			for i := range cand.Occurrences {
				cand.Occurrences[i].Source = displayFromBazi(tenGodSourceDisplay, cand.Occurrences[i].Source)
			}
		}
	}
	return view
}
