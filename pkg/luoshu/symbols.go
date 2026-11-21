package luoshu

// Trigram is one of the eight bagua trigrams (后天八卦). Its value IS the
// stable display text. Shared by bazhai / xuankong / qimen via luoshu.
type Trigram string

const (
	TrigramQian  Trigram = "乾"
	TrigramDui   Trigram = "兑"
	TrigramLi    Trigram = "离"
	TrigramZhen  Trigram = "震"
	TrigramXun   Trigram = "巽"
	TrigramKan   Trigram = "坎"
	TrigramGen   Trigram = "艮"
	TrigramKun   Trigram = "坤"
)

var trigramSet = map[Trigram]bool{
	TrigramQian: true, TrigramDui: true, TrigramLi: true, TrigramZhen: true,
	TrigramXun: true, TrigramKan: true, TrigramGen: true, TrigramKun: true,
}

// EncodeTrigram converts a display name to its Trigram symbol (boundary helper).
func EncodeTrigram(s string) (Trigram, bool) {
	t := Trigram(s)
	return t, trigramSet[t]
}
