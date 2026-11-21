package qimen

// GongRelation is the five-element relation between two palaces of the board.
// It is a deterministic chart fact computed by the engine; the relation value
// is a symbol (HourRelation) emitted as its stable code at the output boundary.
type GongRelation struct {
	From     GongIndex    `json:"from"`
	To       GongIndex    `json:"to"`
	Relation HourRelation `json:"relation"`
}

// computeGongRelations emits the full 9x9 palace-relation matrix. Relations
// depend only on the palaces' fixed five-elements (luoshu), computed via ganzhi
// Sheng/Ke — no local knowledge table is needed by consumers.
func computeGongRelations() []GongRelation {
	out := make([]GongRelation, 0, 81)
	for from := GongIndex(1); from <= 9; from++ {
		for to := GongIndex(1); to <= 9; to++ {
			out = append(out, GongRelation{
				From:     from,
				To:       to,
				Relation: elementRelationToHour(palaceWuxing(from), palaceWuxing(to)),
			})
		}
	}
	return out
}
