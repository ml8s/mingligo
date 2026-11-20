package bazi

// RelationPosition is the engine-owned strength hint for a paired relation:
// adjacent pairs are tight combinations, others are candidates.
type RelationPosition uint8

const (
	positionSeparated RelationPosition = iota
	positionAdjacent
	positionRemote
)

func (p RelationPosition) String() string {
	switch p {
	case positionAdjacent:
		return "adjacent"
	case positionRemote:
		return "remote"
	default:
		return "separated"
	}
}

func EncodeRelationPosition(name string) (RelationPosition, bool) {
	switch name {
	case "adjacent":
		return positionAdjacent, true
	case "remote":
		return positionRemote, true
	case "separated":
		return positionSeparated, true
	default:
		return 0, false
	}
}
