package huangli

import "fmt"

// JianChu is the closed set of the twelve 建除 day gods. The symbol value is
// its stable display text, so it crosses the output boundary unchanged.
type JianChu string

const (
	JianChuJian  JianChu = "建"
	JianChuChu   JianChu = "除"
	JianChuMan   JianChu = "满"
	JianChuPing  JianChu = "平"
	JianChuDing  JianChu = "定"
	JianChuZhi   JianChu = "执"
	JianChuPo    JianChu = "破"
	JianChuWei   JianChu = "危"
	JianChuCheng JianChu = "成"
	JianChuShou  JianChu = "收"
	JianChuKai   JianChu = "开"
	JianChuBi    JianChu = "闭"
)

var jianChuAll = [...]JianChu{
	JianChuJian, JianChuChu, JianChuMan, JianChuPing,
	JianChuDing, JianChuZhi, JianChuPo, JianChuWei,
	JianChuCheng, JianChuShou, JianChuKai, JianChuBi,
}

// DecodeJianChu is the authoritative symbol-to-text mapping.
func DecodeJianChu(jc JianChu) string { return string(jc) }

// EncodeJianChu converts display text into the closed symbol set.
func EncodeJianChu(s string) (JianChu, error) {
	for _, jc := range jianChuAll {
		if string(jc) == s {
			return jc, nil
		}
	}
	return "", fmt.Errorf("huangli: unknown jianchu %q", s)
}

func (jc JianChu) String() string { return DecodeJianChu(jc) }
