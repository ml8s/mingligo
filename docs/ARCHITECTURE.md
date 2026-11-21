# MingliGo 系统架构：命理符号空间与算子

本文是唯一的设计文档，讲清楚内核的构造方式：**一个封闭的命理符号空间，所有运算是在其上运行的纯函数算子，文本只通过双向映射进出。**

## 1. 核心模型

把命理内核建模成一个代数结构：

```
外部文本 ──encode──▶ ┌─────────────────────────────┐ ──decode──▶ 外部文本
                     │  符号空间                      │
                     │  算子: 符号ⁿ → 符号ᵐ (纯函数)   │
                     └─────────────────────────────┘
```

- **符号** = 空间里的点（命理概念的封闭集合）
- **算子** = 空间上的函数：输入符号、输出符号，结果仍在空间内（封闭性）
- **encode / decode** = 空间与文本世界之间唯一的接口

三条公理：符号唯一定义、算子纯函数、文本只在边界。

## 2. 符号空间

符号按领域归属定义，底层只定义一次、上层复用。**每个域包都有 `symbols.go`**，集中定义该域所有封闭命理概念（类型 + 常量 + `Encode`/`Decode`）：

| 域包 | 负责的符号 |
|---|---|
| `ganzhi` | 干、支、五行、阴阳、十神、旺衰、性别、刑类型 |
| `tianwen` | 节气、历法日期、时间边界 |
| `luoshu` | 九宫、八卦 `Trigram`（跨派风水/奇门共用） |
| `fengshui` | 三元龙 `YuanLong`（24山体系） |
| `bazi` | 格局、从格、日支取象 `BranchImageKind`、宫位状态 `PalaceState`、神煞名 `ShenShaName`、关系种类 |
| `qimen` | 宫位 `GongIndex`、门 `DoorIndex`、星 `StarIndex`、神 `SpiritIndex`、金函十二神 `JinhanSpirit`、五行关系 `ElementRelation` |
| `liuyao` | 六亲、六神、爻标、格局子类 |
| `ziwei` | 星曜、亮度 `Brightness`/分组 `BrightnessGroup`、主星态 `MajorStarState`、四化、宫位 |
| `bazhai` | 命卦组 `LifeGroup`、卦组 `TrigramGroup`、游星 `YouXing`、吉凶 `Rating`、匹配 `Match` |
| `huangli` | 建除 `JianChu`（十二建除神） |
| `xuankong` | 元 `Yuan`、运名 `YunName`、伏吟/反吟层 `RecitationLayer` |

符号在归属域定义、跨包复用（如 `GongIndex`/`Trigram`/`ganzhi.YinYang`），不重复造。**所有命运算子都在符号空间内计算**——全局扫描确认零中文字符串参与运算；文本只在输出边界 decode 出来（string 底层符号直接输出文本，bool 符号经 MarshalJSON）。

## 3. 双向映射（唯一边界）

以**符号系统为中心**定义方向：

```
encode: 文本 → 符号   （人类语言编码进系统；输入边界）
decode: 符号 → 文本   （系统解码出人类语言；输出边界）
```

- `DecodeXxx` 是查表的权威实现；`MarshalJSON`、`String()` 薄壳都调它，映射逻辑只有一份。
- 数据文件（JSON）加载时就把中文 key/value 经 `EncodeXxx` 转成符号，**不把显示文本当内部 key**。
- 规则表用符号做代理键：`map[StrengthLevel]rule`、`map[ganzhi.Wuxing]...`，改名/国际化不影响内核。

## 4. 算子规范

**所有命理计算算子都是纯函数**——数学函数语义：

- 输入输出严格是符号（或值对象、`int`/`bool` 等数量标志），不出现表示命理概念的裸 `string`。
- 无副作用：不写指针参数指向的数据、不写全局可变状态；确定性。
- 装配阶段用 `with*(x T) T` 值入值出，不用 `apply*(x *T)`。
- 命中集用 `merge*` 纯累积（如 38 个神煞 helper 返回 `[4][]...` 再由 `mergeShenSha` 合并）。
- 数据加载后的规则表只读；节气等 memo 缓存属实现优化豁免。

## 5. 时间与数据

时间是天文事实，不是产品约定。`tianwen` 提供标准时/真太阳时的**投影函数与边界解析**、节气精确时刻与农历投影，以及 `JieQi`/`LunarBasis` 符号（四个时间尺度本身是 `time.Time` + 命名约定，不是 Go 符号类型）；八字等域只消费这些事实。所有需要出生时刻的入口用 `Compute*At`（标准时），时间解析先走 `tianwen.ResolveStandard`。

对外只需标准时 + 经度：`Compute*At(standard, longitude, ...)` 内部完成真太阳时投影（年/月按标准时定节气，日/时按真太阳时定日界与时辰）。接收已投影 `gt, st` 的双时间原语不导出，仅是内部实现细节——调用方不需要、也不应自行提供真太阳时。例外：`tianwen.ComputeBazi(gt, st)` 是投影函数本身（四柱原子算子），保留导出。

## 6. 输出与适配器

- MCP 是**官方 Go SDK 的 typed 示例适配器**，不是内核：内核是领域事实层，传输层只是它的一个消费者。换 REST 等其他适配器不需要改领域层。
- domain 输出**结构化 object**，字段该 string 就 string、该 int 就 int，**不产出 JSON**。
- 每个领域提供 `XxxView`（Go struct，按命理子域聚合 + `json` 字段名 + `jsonschema` 描述惰性元数据）和 `Chart.View()` 方法：把计算形态（Chart）的碎片按领域语义归位到输出契约。View 是唯一的领域输出模型。
- MCP 等适配器**只透传 View 对象**，输出 schema 由消费者从 View 类型反射生成（`jsonschema.For[XxxView]()`）——不再有手写 wire DTO 第二套模型。符号类型通过 `TypeSchemas` 注册为 string enum，符号的 `MarshalJSON` 负责边界 decode。
- **View 与 Chart 解耦**：Chart 保持最适计算的形态（可平行数组/内部符号），View 定义稳定的输出契约（可读字段、领域聚合）。填充分次序自由，但同一实体的碎片必须归位到同一对象。
- **适配器不约束 domain**：schema 推断所需的符号枚举在适配器侧注册，绝不回退 domain 的符号化。
- 注：bond（合婚）、liupan（流盘序列）、liuyao_qigua（起卦）工具均已 View 化——合婚用 `bazi.BondView`/`ziwei.BondView`，流盘用 `bazi.LiuNianOutputView` 等系列 View，起卦用 `liuyao.QiguaView`（domain 系列/事件对象以 JSON roundtrip 映射到带 `jsonschema` tag 的输出契约，schema 从 View 反射；enum 约束在适配器侧 patch）。time_now/city_info 是纯 MCP 适配层的信息工具（当前时间、地点候选），无领域对象映射，其输出 DTO 属于适配器产物，不属于领域 View。

## 7. 事实层与分析层的边界

MingliGo 只做**确定性排盘计算**：领域输出是符号空间上的**命理领域事实**。上层内容服务（counsel）负责**拟合分析**——多因素加权、规则冲突消解、概率性倾向，其规则表用经典依据做拟合包装。二者判据：

| 判据 | 归属 | 例子 |
|---|---|---|
| **直推**（经典明确规定的固定映射，无权衡拟合） | 命理领域事实，留在内核 | 神煞取象「天乙贵人→遇难有救」、奇门克应「庚+丙→白入荧贼即来」、玄空加会「一二同宫→土克水」、建除宜忌 |
| **拟合**（多因素加权、冲突消解、概率倾向） | 分析层，归 counsel | 用神选取、格局成败、综合吉凶判断 |
| **冗余 description**（结构化字段已表达的构成白话） | 删除 | 调候「某月某干用某调候」、紫微格局「太阳居午宫」 |

**直推断语是符号空间的例外：它们无法（也无必要）符号化。** 其 key（神煞名、干组合、星曜组合）已符号化；文本是查表的 **value**，不参与任何运算，只在输出边界作为原生文本 decode 出来。给它编符号 ID 再 decode 回文本没有收益。这与第 3 节一致：断语是 decode 方向的终点，不是空间内的点。

经典原文自带的适用边界声明（如「客主吉凶须合局与用神另断」「断语两说并存」「须核元运与峦头」）是对固定映射的忠实转录，仍属直推事实，不是拟合。

**命理逻辑落表**：客观命理知识（固定取象、属性、取用规则）以数据表承载，**由 MingliGo 权威持有**，不放分析层。例如奇门宫位人物/高低/内外、门星吉凶宜事落 `pkg/qimen/data/gong_facts.json`，以现有符号（GongIndex/DoorIndex/StarIndex）为 key、loader 经 `Encode*` 转符号、输出边界 decode 为文本；方位复用 `luoshu` 共享权威表。分析层（counsel）只读 chart 输出，不再自带镜像表。

**神盘建模：跨派共用符号 + 门派命名 decode + 独立符号集**。转盘/飞盘/鸣法的神是**同一批传统神**（槽位 1-9），仅命名约定不同（阴阳遁异名、飞盘/鸣法固定名）；因此共用一个 `SpiritIndex`，输出**稳定槽位 key**（`shen_symbol`，如 `gouchen`，恒不随门派/阴阳遁变）与**门派显示名**（`shen`，如白虎/勾陈）两个字段，由消费端取用。金函十二神（地钥/六贼/…）是真正不同的神，用独立符号集 `JinhanSpirit`。按门派归属落在数据/命名层（`plate.json` 三张表 + `jinhan.json`），符号层不按门派拆分——否则会把"转盘白虎"和"飞盘白虎"（同一神）误建成不同符号，违反领域事实。

**盘面派生结论由引擎产出**：给定盘结果唯一的派生事实，作为 chart 字段输出，消费端只读不重算。这些字段先符号化（封闭集符号 / string 底层值即文本），算子在符号空间内计算，输出边界 decode 为文本——符合六边形架构。已落地：八字 `spouse_palace_state`/`day_branch_type`（从 guan_xi/shen_sha 派生）、紫微 `liang_du_group`/`zhu_xing_state`（从亮度/主星派生）、奇门 `gong_relations`（九宫五行关系矩阵，由 ganzhi Sheng/Ke 计算）。消费端删除对应的本地命理表与推算（constants 合化表、qimen 生克表等）。**边界**：依赖"盘中未出现的实体"（如未出现的十神的五行）的推算、以及静态知识表本身，仍留消费端——它们不是盘面结论。

## 8. 验证

- 每个工具保留 domain 输出与 wire 的等价测试，`json.Marshal` 输出字节由 golden/schema 守护。
- 改动符号或算子后跑 `make gate`（gofmt/lint/vet/test/build）。
- 主测试历法范围 1900–2100；范围外不保证历法、时区与符号表兼容。

## 9. 判定标准（审计用）

| 情形 | 判定 |
|---|---|
| 命理概念是封闭集却用裸 `string`/`map[string]` 运算 | 违规，必须符号化（即使当前只喂输出） |
| 算子写指针参数 / 全局可变状态 | 违规，改值入值出 |
| 直推的查表式断语文本（经典固定映射，作为输出 value） | 合规，属命理事实；不符号化 |
| 拟合性分析（多因素加权/概率倾向）混入领域层 | 违规，归内容层 counsel |
| 冗余 description（结构化字段已表达的白话复述） | 删除 |
| `Encode`/`Decode`、数据加载、输出展示文案 | 边界，合规 |
| wire 契约强制的 string-keyed map、复合非封闭值（如 `main_qi_not_transparent`） | 真边界，合规 |

## 附录：审计判定清单（曾争议的项）

| 项 | 判定 | 理由 |
|---|---|---|
| `NatalRelationMember`/`StemOccurrence`/`TargetedRelationEvidence` 的干支 | 已符号化 | 曾被 `EncodeGan/Zhi` 反解参与运算，是符号→string→符号往返 |
| `ElementStructure.Wuxing/Season` | 已符号化 | 用解码文本比较旺衰；改符号后直接比较 |
| `rel*` 干支关系、`cat*` 神煞吉凶 | 已符号化 | 封闭命理概念具名化 |
| `CandidateWuxing`/`TenGodSource`/`TripleRelationKind`/`GanZhiRelationKind`/`PatternSubType`/`StarLei` | 已符号化 | 封闭集具名（string 底层，值即文本），即使只喂输出 |
| `FuYiResult.WangShuai`/`WuxingCount` 的 string key | 真边界 | wire 契约强制 string-keyed map；符号 key 会输出数字破坏契约。消费端用 `fuYiWangShuaiAt` 单一边界 helper 转符号 |
| `TombSchool.EarthBranch` | 真边界 | 固定"辰"的流派标注，不参与卦理；在 `ChartView` 中 MCP 直接引用，符号化需拆整个 View |
| `TiaoHouResult.MarshalJSON` | 特例 | all-null secondary 是 domain 的输出形状定义（ADR 0004: domain owns output shape），非符号映射；wire 侧同形状对象 schema 依赖它 |
| `main_qi_not_transparent` 等复合值 | 真边界 | 非封闭集，string 合理 |
| 数据加载 `Encode*`、输入校验 | 边界 | string→符号在入口、符号→文案在出口 |
| 直推查表断语（神煞取象/奇门克应/玄空加会/建除宜忌） | 例外，留内核 | 文本是查表 value 不参与运算，key 已符号化；见第 7 节 |
| 冗余 description（调候口径/紫微格局构成） | 已删除 | 结构化字段已表达同一事实 |
| 六爻用神多现取舍（旺衰>动>不破>不空>不伤>近世应） | 直推，留内核 | 《增删卜易》明文固定优先级，确定性无权衡；候选全量输出供复核。区别于八字用神（多因素权衡，归 counsel） |
