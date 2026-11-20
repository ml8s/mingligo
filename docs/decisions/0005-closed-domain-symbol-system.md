# ADR 0005: 封闭领域符号系统（符号空间 + Encode/Decode 边界）

## Status

Accepted.

## Context

命理计算里大量概念天然是**封闭集合**：天干、地支、五行、旺衰、宫位、星曜、格局、节气、身强身弱、根类型、 caste 关系……早期实现里这些概念常以中文字符串或英文标签出现，且出现在三种位置：

1. 数据结构字段（`"旺"`、`"month_main"`、`"身强"`）
2. 内部运算的中间值（map key、switch 判据、函数参数与返回值）
3. 与外部交换的文本（JSON、数据文件、日志）

由此产生的问题：

- 同一个命理概念在不同包重复定义、命名不统一（私有/导出、`Type`/`Index` 后缀混用）。
- 内部运算拿 display label 当 key（等价于数据库里用显示名当主键）：改名波及全代码、无类型检查、无法国际化。
- `String()` / `MarshalJSON` / `ParseXxx` 各自在不同文件里手工实现，符号→文本的逻辑散落、方向语义混乱（到底哪个是 encode、哪个是 decode？）。
- 字符串往返：符号 → `String()` → 文本 → `Parse*` → 符号，出现在计算中间而非边界。

## Decision

把整个内核建模成一个**封闭的领域符号系统（符号空间）**，所有命理运算都在符号之上；文本只存在于系统边界。

### 1. 符号空间 = 代数结构

- 每个命理概念是一个**符号类型**（枚举或受约束的强类型），定义在**其归属域包**，全局共用：
  - 通用符号（干支、五行、阴阳、十神、旺衰、性别）只在 `ganzhi` 定义一次；时间符号（节气、日期）只在 `tianwen` 定义；八字专属符号在 `bazi`；奇门在 `qimen`；六爻在 `liuyao`；紫微在 `ziwei`。
  - 上层包一律 import 复用，不重造底层符号。
- 运算是空间上的**算子**：输入符号、输出符号，满足封闭性（结果仍在空间内）。内部函数签名一律符号进符号出，不出现表示命理概念的 string。
- **算子必须是纯函数**（数学函数语义）：无副作用（不写指针参数指向的数据、不写全局可变状态）、确定性。`add*(out *T, ...)` 改为返回命中集再由 `merge*` 纯累积；`apply*(x *T)` 改为 `with*(x T) T` 值入值出；slice 字段用新建 slice 避免底层数组共享写入。数据加载后的规则表全局只读，memoization 缓存（如节气 memo）属实现优化豁免。
- 同一概念跨包同名（如 `GongIndex`/`StarIndex`），命名唯一且导出 PascalCase；仅当有同名实体（如 `Gong` 结构）时保留必要的 `Index`/`Type` 区分后缀。

### 2. 一张映射表，两个方向

每个符号在**本包 `symbols.go`** 集中定义：type + 常量 + 名字表 + 双向映射。

```text
encode : 文本 → 符号   （进入系统；输入边界；原 ParseXxx 统一改名 EncodeXxx）
decode : 符号 → 文本   （离开系统；输出边界；查表权威实现 DecodeXxx）
```

- `DecodeXxx` 是**权威实现**（直接查映射表）；`MarshalJSON` 与必要的 `String()` 薄壳都调 `DecodeXxx`，不各自造映射。
- 方向命名以**符号系统为中心**：人类语言编码进系统 = encode，系统解码出人类语言 = decode。`Encode`/`Decode` 都是 object↔object 的边界动作，不产出 JSON。
- 名称相同不表示可跨领域复用：符号归属以域包为界。

### 3. 字符串只在边界，key 也一样

- **输入边界**（JSON 数据文件、外部 API 入参）：文本经 `EncodeXxx` 变成符号后进入系统。数据加载时就把中文 key/value 转成符号，**不把 display label 当内部 key**。
- **输出边界**（domain 结果 JSON）：符号经 `DecodeXxx` 填进输出字段；序列化由 SDK（`encoding/json`）完成，domain 不自行产出 JSON。
- **内部**：map key、switch 判据、函数参数/返回值一律符号。规则表用代理键语义建模——

```
规则表:    map[StrengthLevel]yongJiRule     // key 是符号 ID，不是 "身强"
规则值:    RelationKind / YongXiJiKind       // "克我者"/"dm_elem" 也是符号
季节:      ganzhi.WangShuai                  // 复用已有符号
加载:      JSON 中文 → Encode* → 符号
输出:      Decode*(符号) → 原中文
```

### 4. 与协议/架构的关系（承接 ADR 0004）

- domain 只输出结构化 object（string/int/符号的 JSON 投影），MCP 层从 wire DTO 生成 schema 并透传 domain 结果。
- JSON Schema 与 enum 约束属于协议层，不进入 `pkg/**`；本项目不维护 wire 契约作为产品——输出字节由 golden/schema 测试守护即可。
- **适配器不得约束 domain 的符号化**：若 domain 字段改为符号类型后与适配器的 schema 推断冲突（如符号的 `MarshalJSON` 输出文本但 jsonschema 从整数底层推断 integer），应在适配器侧用独立 wire struct / 注册 schema 适配，**绝不为此回退 domain 的符号化**。domain 是权威，适配器跟随。
- 领域"输出可读事实"与"符号空间闭环"不冲突：可读性是 decode 的产物，运算始终是符号。

### 5. 判定标准（如何审计）

一个函数是否合规，看**是否经过映射表**：

| 情形 | 经过映射表？ | 判定 |
|---|---|---|
| `EncodeXxx(text) → 符号` | 是 | 边界 encode，合规 |
| `DecodeXxx(符号) → text` | 是 | 边界 decode，合规 |
| 符号 `String()` 薄壳（调 Decode） | 是 | tostring 别名，合规 |
| 内部函数用 string 表示命理概念做运算/key | 否（或该走符号却用了文本） | 违规，需符号化 |
| 算子写指针参数 / 全局可变状态 | — | 违规，需改为值入值出的纯函数 |
| 值对象格式化（`Date.String()`、数字→文本） | 否 | 纯 tostring，非符号映射，合规 |

## Consequences

**正面**

- 命理概念唯一定义、全局复用；命名统一，新增概念有明确归属。
- 内部运算类型安全，非法组合在编译期被挡住；规则表改名/国际化不影响内核。
- 符号映射收敛到各包 `symbols.go`，审计只需看符号表 + 算子集。
- 输出 JSON 字节不变（decode 保形状），golden/schema 契约稳定。

**代价 / 约定**

- 新增命理概念必须：在正确域包定义符号 → 写 Encode/Decode → 更新数据加载边界 → 运算只用符号 → 输出用 Decode。
- `ganzhi`/`tianwen` 底座符号保留 `String()` 薄壳（调 Decode）以兼容 `fmt.Stringer` 与测试；逻辑仍只有一份。
- 规则表符号化涉及数据加载改造，改动必须逐项验证 golden。
- 禁止在 `pkg/**` 新增 `jsonschema` tag 或把 wire 契约当领域模型。

## Related

- ADR 0004: Domain-owned aggregated output, thin protocol adapters
- `docs/TIME_MODEL.md`：时间符号与边界
- `docs/domain/README.zh-CN.md`：领域包总览
