package ziwei

import "encoding/json"

// This file defines the domain-owned, domain-aggregated output view of a Ziwei
// 合婚 (bond) result. jsonschema tags are inert metadata read only when a
// consumer reflects the type for its schema.

// BondPalaceFactsView is one side's palace facts with readable star names.
type BondPalaceFactsView struct {
	Palace     string   `json:"gong" jsonschema:"宫位。"`
	MajorStars []string `json:"major_xing" jsonschema:"主星列表。"`
}

// BondPalacePairView is a side-by-side palace pair with a relation status.
type BondPalacePairView struct {
	A      BondPalaceFactsView `json:"a" jsonschema:"甲方宫位事实。"`
	B      BondPalaceFactsView `json:"b" jsonschema:"乙方宫位事实。"`
	Status string              `json:"status" jsonschema:"宫位关系状态。"`
}

// BondView is the domain-aggregated output of a Ziwei bond.
type BondView struct {
	MingGong BondPalacePairView `json:"ming_gong" jsonschema:"双方本命宫。"`
	FuQiGong BondPalacePairView `json:"fu_qi_gong" jsonschema:"双方夫妻宫。"`
	ZiNvGong BondPalacePairView `json:"zi_nv_gong" jsonschema:"双方子女宫。"`
}

// View produces the domain-aggregated output from the computed bond. The Bond
// is a pure output value, so a JSON roundtrip maps its structurally identical
// fields (decode StarIndex to readable names) into the tagged View contract.
func (b Bond) View() BondView {
	raw, err := json.Marshal(b)
	if err != nil {
		return BondView{}
	}
	var view BondView
	if err := json.Unmarshal(raw, &view); err != nil {
		return BondView{}
	}
	return view
}
