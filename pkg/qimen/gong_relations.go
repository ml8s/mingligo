package qimen

// GongRelation is the five-element relation between two palaces of the board.
// It is a deterministic chart fact computed by the engine; the relation value
// is a symbol (ElementRelation) emitted as its stable code at the output boundary.
type GongRelation struct {
	From     GongIndex       `json:"from"`
	To       GongIndex       `json:"to"`
	Relation ElementRelation `json:"relation"`
}

// computeGongRelations emits the full palace-relation matrix over the nine
// GongIndex symbols. Relations depend only on the palaces' fixed five-elements
// (luoshu), computed via ganzhi Sheng/Ke — no local knowledge table is needed.
func computeGongRelations() []GongRelation {
	n := int(GongLi)
	out := make([]GongRelation, 0, n*n)
	for from := GongKan; from <= GongLi; from++ {
		for to := GongKan; to <= GongLi; to++ {
			out = append(out, GongRelation{
				From:     from,
				To:       to,
				Relation: elementRelation(palaceWuxing(from), palaceWuxing(to)),
			})
		}
	}
	return out
}
