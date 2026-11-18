package ziwei

import "github.com/ml8s/mingligo/pkg/ganzhi"

// determineJuShu returns the bureau number from ming gong gan and zhi.
func determineJuShu(mingGan Gan, mingZhi Zhi) juShu {
	nayinName := ganzhi.NayinLabel(mingGan, mingZhi)
	wx := ganzhi.NayinWuxing(nayinName)
	if wx == 0 {
		return 0
	}
	return juShuFromWuxing(wx)
}
