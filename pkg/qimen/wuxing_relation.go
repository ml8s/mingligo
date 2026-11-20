package qimen

// WuxingRelation is a qimen-domain symbol for a five-element relation between
// two members (a star and a gong, or 日 and 时). The relation key is stable and
// has a reverse mapping to its catalog identifier.
type WuxingRelation uint8

const (
	wrSame             WuxingRelation = iota // same
	wrGongGeneratesXing                       // gong_generates_xing
	wrXingGeneratesGong                       // xing_generates_gong
	wrXingControlsGong                        // xing_controls_gong
	wrGongControlsXing                        // gong_controls_xing
	wrRiGeneratesShi                          // ri_generates_shi
	wrShiGeneratesRi                          // shi_generates_ri
	wrRiControlsShi                           // ri_controls_shi
	wrShiControlsRi                           // shi_controls_ri
)

var wuxingRelationNames = [...]string{
	"same", "gong_generates_xing", "xing_generates_gong",
	"xing_controls_gong", "gong_controls_xing",
	"ri_generates_shi", "shi_generates_ri",
	"ri_controls_shi", "shi_controls_ri",
}

func (r WuxingRelation) String() string {
	if int(r) < len(wuxingRelationNames) {
		return wuxingRelationNames[r]
	}
	return ""
}

func ParseWuxingRelation(name string) (WuxingRelation, bool) {
	for i, n := range wuxingRelationNames {
		if n == name {
			return WuxingRelation(i), true
		}
	}
	return 0, false
}
