# 领域包总览

本文是 `pkg/*` 的中文领域导览，适合后续拆成 wiki 页面。实现细节以英文文档为准。

## 核心定位

MingliGo 不是排盘 SaaS，也不是命理建议系统。它是**通用的命理计算层**：

```text
标准时间 + 经度 + 明确规则
        ↓
天文 / 历法 / 干支 / 卦例 / 飞星事实
        ↓
typed struct / JSON
```

应用层可以继续做解释、建议、UI 和产品，但必须知道内核输出的是哪种规则下的事实。

## 时间是第一等事实

所有需要出生时间或事件时间的领域共用 `pkg/tianwen`：

```text
民用时间
  → 去夏令时后的标准时
    → 经度修正
      → 均时差
        → 真太阳时 / LAST
          → 时辰、日界、节气界、年柱/月柱
```

要记住三条边界：

1. 年柱用立春精确时刻，不按公历日期。
2. 月柱用节精确时刻。
3. 日柱和时柱用真太阳时；经度东移可能导致真太阳时跨日，西移可能退回前一个时辰。

`Compute*At` 系列的 `standardTime` 参数必须已经是标准时。若输入可能是夏令时壁钟、模糊时辰、时间范围或经度范围，先走 `tianwen.ResolveStandard` 做归一和候选解析。

农历投影不是单一规则。MingliGo 显式支持：

| basis | 含义 |
|---|---|
| `official` | 北京民用日期 + UTC+8 历表 |
| `local_label` | 本地日期标签 + UTC+8 历表 |
| `local_astronomical` | 本地标准日期 + 本地天文朔气规则 |

详见 [`../TIME_MODEL.md`](../TIME_MODEL.md)。

## 基础层

### `tianwen`

负责天文和历法：

- 节气精确时刻；
- 平太阳时到真太阳时的转换；
- 时辰判定；
- 日界与晚子时相关输入；
- 农历、闰月、朔日；
- 模糊出生时间解析。

`ResolveStandard` 很重要：当用户只给“某日某时辰”、时间范围、或地点经度范围时，它不会假装只有一个答案，而是返回候选时间和原因。

### `ganzhi`

干支是全库共享语言：

- 天干、地支、五行、阴阳；
- 六十甲子；
- 藏干、纳音、十二长生；
- 合会冲刑破害；
- 人元司令分野。

其他领域不再各自实现一张藏干表或冲合表。

### `tables`

聚合加载内嵌表，适合服务启动时 fail-fast。领域 API 本身也会懒加载。

### `luoshu`

洛书九宫和后天八卦的共享事实表，避免九宫方向、卦名、五行在各模块漂移。

### `city`

内嵌地点解析：

- 中国县级 2,844 条；
- 城市 337 条；
- 精选全球城市 88 条。

同名地点会返回候选列表，不会静默选择一个“朝阳区”。

## 八字 `bazi`

入口：

```go
chart := bazi.ComputeChartAt(standardTime, longitude, gender)
```

时间分工：

- 年柱、月柱：标准时间和节气；
- 日柱、时柱：真太阳时。

输出包括：

- 四柱；
- 藏干、十神、纳音；
- 旬空、自坐、十二长生；
- 神煞；
- 合会冲刑破害；
- 大运、流年、流月、流日、流时、小运；
- 胎元、命宫、身宫；
- 扶抑、调候、格局证据。

重点设计：**用神不下最终结论**。内核输出扶抑候选、调候候选和格局证据，解释层再选择流派。

详见 [`bazi.md`](bazi.md)。

## 紫微 `ziwei`

入口：

```go
chart, err := ziwei.ComputeChartAt(standardTime, longitude, gender)
```

默认农历 basis 是 `local_astronomical`；如果流派要求按北京历或本地标签查农历，用 `ComputeChartAtBasis` 显式指定。

时辰仍然来自真太阳时。海外经度和跨日案例可能改变农历日期，进而影响命盘。

支持大限、流年、流月、流日、流时与合盘。

详见 [`ziwei.md`](ziwei.md)。

## 六爻 `liuyao`

六爻从可审计的起卦事实开始：

- 三枚铜钱六次；
- 正=3，反=2；
- 6 老阴、7 少阳、8 少阴、9 老阳；
- 老阴老阳为动爻；
- 也可以直接传入六爻值。

内核输出：

- 本卦、变卦；
- 京房八宫；
- 纳甲干支；
- 六亲、六神；
- 世应；
- 伏神、飞神；
- 用神位置和旺衰；
- 旬空、月破、进退神、反吟伏吟等格局；
- 应期候选。

详见 [`liuyao.md`](liuyao.md)。

## 奇门 `qimen`

奇门用封闭方法目录，避免组合非法流派：

| 维度 | 值 |
|---|---|
| scope | `hour`, `quarter`, `day`, `month`, `year` |
| school | `zhuanpan`, `luoshu_feipan`, `mingfa_feipan`, `jinhan_yujing` |
| 定局 | `chaibu`, `zhirun`, `maoshan`, `none` |

推荐入口：

```go
method, err := qimen.ResolveMethod(...)
chart, err := qimen.ComputeChartWithMethodAt(standardTime, longitude, method)
```

盘面保留节气、局、元、符头、值符、值使、九宫、八神、用神落宫等方法元数据。

详见 [`qimen.md`](qimen.md)。

## 黄历 `huangli`

黄历是日期级领域：

```go
day, err := huangli.QueryDate("2026-06-28", "wedding")
```

输出包括：

- 日柱；
- 建除；
- 黄道/黑道；
- 喜神、财神、福神方位；
- 彭祖百忌；
- 二十八宿；
- 节气与入节天数；
- 人元司令分野；
- 十二时辰吉凶事实；
- 事项适配结果。

详见 [`huangli.md`](huangli.md)。

## 风水

### `bazhai`

- 按立春年与性别算命卦；
- 东四命 / 西四命；
- 大游年八方向；
- 门主灶；
- 年紫白飞星。

### `xuankong`

- 按立春年定三元运；
- 二十四山坐向；
- 下卦正向玄空盘；
- 运星、山星、向星；
- 旺山旺向、上山下水、伏吟反吟；
- 双星加会与收山出煞。

### `fengshui`

共享二十四山和紫白飞星表。

详见 [`fengshui.md`](fengshui.md)。

## 测试与参考

MingliGo 的“准”不是靠宣传词，而是靠三类测试：

1. 独立规则锚点；
2. 外部 golden fixture；
3. 边界和属性测试。

核心参考包括 `lunar-python`、`sxtwl`、`lunar-typescript` / ShouXingUtil、`atopx/qimen`、`deminzhang/qimen-go`、`kentang2017/kinqimen`、`iztro` 等。

版本、license、fixture 路径、容差和已知差异见 [`../VALIDATION.md`](../VALIDATION.md)。

## Wiki 拆分建议

如果要拆成站点点：

1. 时间与历法模型
2. 干支与基础表
3. 八字事实模型
4. 紫微时间基础与流派
5. 六爻起卦与装卦
6. 奇门方法目录
7. 黄历规则
8. 风水排盘
9. Golden 测试与参考项目
10. Go API 使用指南
