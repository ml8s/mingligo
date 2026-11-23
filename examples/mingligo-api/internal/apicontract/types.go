package apicontract

// Wire-level scalar types carry JSON Schema constraints through TypeSchemas.
// Their runtime JSON representations remain ordinary strings and integers.
type (
	StandardTime     string
	DateString       string
	CalendarKindName string
	GenderName       string
	CityName         string
	HuangliEvent     string
	HuangliCount     int
	SeriesCount      int

	GuaNumber        int
	MountainIndex    int
	XuankongYear     int
	EightGuaName     string
	MountainName     string
	StarColorName    string
	ElementName      string
	YinYangName      string
	DirectionName    string
	LifeGroupName    string
	YearBoundaryName string
	FlyingStarNumber int
	FlyingStarName   string
	YouxingName      string
	BazhaiRatingName string
	TrigramGroupName string
	BazhaiMatchName  string
	BazhaiBirthYear  int

	LiuyaoQiguaMode        string
	LiuyaoCastMode         string
	LiuyaoOrderName        string
	LiuyaoCoinSymbol       string
	LiuyaoCoinValue        int
	LiuyaoYaoValue         int
	LiuyaoOptionalYaoValue int
	LiuyaoPosition         int
	LiuyaoYaoLabel         string
	LiuyaoOptionalPosition int
	LiuyaoHexagramName     string
	GanName                string
	ZhiName                string
	WangShuaiName          string
	LiuyaoYongShenName     string
)
