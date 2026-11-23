package mcpcontract

import "github.com/ml8s/mingligo/pkg/bazi"

// bazi.bond wire contract — bazi_chart now uses the domain-owned bazi.ChartView.

// ── bazi.bond input ──

type BaziBondInput struct {
	A BaziBondSideInput `json:"a" jsonschema:"甲方：已归一化的出生主体。"`
	B BaziBondSideInput `json:"b" jsonschema:"乙方：已归一化的出生主体。"`
}

type BaziBondSideInput struct {
	Subject NatalSubject `json:"subject" jsonschema:"已归一化的出生主体。"`
}

// ── bazi.bond output ──

// ── bazi.liupan input ──

type BaziLiupanInput struct {
	Scope        BaziLiupanScopeName `json:"scope" jsonschema:"流盘层级。"`
	Subject      NatalSubject        `json:"subject" jsonschema:"已归一化的出生主体。"`
	StandardTime StandardTime        `json:"standard_time,omitempty" jsonschema:"流盘选择起点；liu_nian/liu_yue/liu_ri/liu_shi 必填，local standard time，RFC3339 fixed offset。"`
	Count        *SeriesCount        `json:"count,omitempty" jsonschema:"连续查询数量。"`
	Longitude    *Longitude          `json:"longitude,omitempty" jsonschema:"地点经度；liu_ri/liu_shi 必填（-180..180），真太阳时校正用于日柱/时柱。"`
}

// ── bazi.liupan union output ──

// BaziLiupanOutput carries exactly one selected flow layer. The explicit
// output schema keeps the current single-root wire face while allowing one
// typed handler to serve the five scope forms.
type BaziLiupanOutput struct {
	LiuNian []bazi.LiuNianOutputView `json:"liu_nian,omitempty" jsonschema:"流年数组，长度 = count。"`
	LiuYue  []bazi.LiuYueOutputView  `json:"liu_yue,omitempty" jsonschema:"流月数组，长度 = count。"`
	LiuRi   []bazi.LiuRiOutputView   `json:"liu_ri,omitempty" jsonschema:"流日数组，长度 = count。"`
	LiuShi  []bazi.LiuShiOutputView  `json:"liu_shi,omitempty" jsonschema:"流时数组，长度 = count。"`
	XiaoYun []bazi.XiaoYunItemView   `json:"xiao_yun,omitempty" jsonschema:"各学派小运列表；恒为两个学派，count 控制各学派柱数。"`
}
