// MingliGo wire DTOs for the MCP contract layer.
//
// These types are ordinary source now, not generated artifacts. Edit them
// directly when the published MingliGo payload changes.
package mcpcontract

import "github.com/ml8s/mingligo/pkg/ziwei"

// ── ziwei.chart input ──

// ZiweiChartInput consumes a resolved natal subject. The chart domain derives
// the local astronomical lunar date and LAST shichen from the subject.
type ZiweiChartInput struct {
	Subject NatalSubject `json:"subject" jsonschema:"已归一化的出生主体。"`
}

// ── ziwei.liupan input ──

type ZiweiLiupanInput struct {
	Scope        ZiweiLiupanScopeName `json:"scope" jsonschema:"流盘层级。"`
	Subject      NatalSubject         `json:"subject" jsonschema:"已归一化的出生主体。"`
	StandardTime StandardTime         `json:"standard_time" jsonschema:"流盘选择起点：local standard time，RFC3339 fixed offset。"`
	Count        *SeriesCount         `json:"count,omitempty" jsonschema:"连续查询数量。"`
	Longitude    *Longitude           `json:"longitude,omitempty" jsonschema:"地点经度；liu_shi 必填（-180..180），真太阳时校正用于时辰。"`
}

// ── ziwei.liupan union output ──

// ZiweiLiupanOutput carries exactly one selected flow layer. The explicit
// output schema keeps the current single-root wire face while allowing one
// typed handler to serve the four scope forms.
type ZiweiLiupanOutput struct {
	LiuNian []ziwei.LiuNianOutputView `json:"liu_nian,omitempty" jsonschema:"流年数组，长度 = count。"`
	LiuYue  []ziwei.LiuYueOutputView  `json:"liu_yue,omitempty" jsonschema:"流月数组，长度 = count。"`
	LiuRi   []ziwei.LiuRiOutputView   `json:"liu_ri,omitempty" jsonschema:"流日数组，长度 = count。"`
	LiuShi  []ziwei.LiuShiOutputView  `json:"liu_shi,omitempty" jsonschema:"流时数组，长度 = count。"`
}

// ── ziwei.bond input ──

type ZiweiBondSide struct {
	Subject NatalSubject `json:"subject" jsonschema:"已归一化的出生主体。"`
}

type ZiweiBondInput struct {
	A ZiweiBondSide `json:"a" jsonschema:"甲方：已归一化的出生主体。"`
	B ZiweiBondSide `json:"b" jsonschema:"乙方：已归一化的出生主体。"`
}

// ── ziwei.bond output ──

// ── ziwei.chart output ──
