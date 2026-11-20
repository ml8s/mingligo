package qimen

import "github.com/ml8s/mingligo/pkg/ganzhi"

// ganEntry holds the named interaction data (without gan fields set at runtime).
type ganEntry struct {
	PatternName string
	Auspicious  bool
}

// computeGanInteractions returns the 十干克应 for each gong.
func computeGanInteractions(board Board) []GanInteraction {
	result := []GanInteraction{}
	for i := 0; i < 9; i++ {
		p := board.GongWei[i]
		for _, heavenGan := range heavenGans(p) {
			key := [2]ganzhi.Gan{p.DiPanGan, heavenGan}
			if entry, ok := ganInteractionTable[key]; ok {
				result = append(result, GanInteraction{
					Gong:       GongIndex(i + 1),
					DiPanGan:   p.DiPanGan,
					TianPanGan: heavenGan,
					Name:       ganzhi.DecodeGan(heavenGan) + "+" + ganzhi.DecodeGan(p.DiPanGan),
					Auspicious: entry.Auspicious,
				})
			}
		}
	}
	return result
}
