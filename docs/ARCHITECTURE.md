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

符号按领域归属定义，底层只定义一次、上层复用：

| 域包 | 负责的符号 |
|---|---|
| `ganzhi` | 干、支、五行、阴阳、十神、旺衰、性别、刑类型 |
| `tianwen` | 节气、历法日期、时间边界 |
| `bazi` `qimen` `liuyao` `ziwei` | 各域专属符号（格局、宫位、门星神、六亲、杂曜等） |

需要符号化的域包以 `symbols.go` 集中定义该域所有符号：**类型 + 常量 + 名字表 + `Encode`/`Decode`**（目前 ganzhi/tianwen/bazi/qimen/liuyao/ziwei 六包有；其它包无封闭命理符号集）。同名概念跨包共用（如 `GongIndex`/`StarIndex`），不重复造。

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

## 6. 输出与适配器

- MCP 是**官方 Go SDK 的 typed 示例适配器**，不是内核：内核是领域事实层，传输层只是它的一个消费者。换 REST 等其他适配器不需要改领域层。
- domain 输出**结构化 object**，字段该 string 就 string、该 int 就 int，**不产出 JSON**。
- MCP 等适配器只做两件事：透传 domain 结果、从 wire DTO 生成 schema。wire DTO 只承载 schema 标签，不是第二套数据模型。
- **适配器不约束 domain**：若 domain 字段符号化与适配器的 schema 推断冲突，在适配器侧用独立 wire struct / 注册 schema 适配，绝不回退 domain 的符号化。

## 7. 事实层与分析层的边界

MingliGo 只做**确定性排盘计算**：领域输出是符号空间上的**命理领域事实**。上层内容服务（counsel）负责**拟合分析**——多因素加权、规则冲突消解、概率性倾向，其规则表用经典依据做拟合包装。二者判据：

| 判据 | 归属 | 例子 |
|---|---|---|
| **直推**（经典明确规定的固定映射，无权衡拟合） | 命理领域事实，留在内核 | 神煞取象「天乙贵人→遇难有救」、奇门克应「庚+丙→白入荧贼即来」、玄空加会「一二同宫→土克水」、建除宜忌 |
| **拟合**（多因素加权、冲突消解、概率倾向） | 分析层，归 counsel | 用神选取、格局成败、综合吉凶判断 |
| **冗余 description**（结构化字段已表达的构成白话） | 删除 | 调候「某月某干用某调候」、紫微格局「太阳居午宫」 |

**直推断语是符号空间的例外：它们无法（也无必要）符号化。** 其 key（神煞名、干组合、星曜组合）已符号化；文本是查表的 **value**，不参与任何运算，只在输出边界作为原生文本 decode 出来。给它编符号 ID 再 decode 回文本没有收益。这与第 3 节一致：断语是 decode 方向的终点，不是空间内的点。

经典原文自带的适用边界声明（如「客主吉凶须合局与用神另断」「断语两说并存」「须核元运与峦头」）是对固定映射的忠实转录，仍属直推事实，不是拟合。

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
