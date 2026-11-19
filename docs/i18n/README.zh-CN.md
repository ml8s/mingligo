<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="../../assets/logo-dark.svg">
    <img alt="MingliGo logo" src="../../assets/logo.svg" width="72">
  </picture>
</p>

# MingliGo

> **英文 README 是权威文档。**
> 本页是中文概览，不是完整翻译。

MingliGo 是中文命理计算的纯 Go 确定性事实核心。它提供可直接使用的稳定领域包；`examples/mingligo-mcp` 只是用于验证该能力的示例 MCP adapter，不属于事实核心。

## 设计边界

| 边界 | 职责 |
| --- | --- |
| `pkg/*` | 确定性领域事实、数据表和校验 |
| `examples/mingligo-mcp` | 基于 official SDK 的 typed tool adapter 和 composition root |
| `contracts/mcp-tools.json` | 公开的 MCP 工具目录 |

事实核心不会导入 MCP、transport、产品或部署代码；MCP DTO 和 wire mapping 只保留在参考 adapter 中。

## 领域包

| 领域 | 包 |
| --- | --- |
| 命盘 | `bazi`、`ziwei` |
| 占测 | `liuyao`、`qimen` |
| 风水 | `bazhai`、`xuankong`、`fengshui` |
| 历法与地点 | `huangli`、`tianwen`、`city`、`ganzhi` |
| 共享基础 | `tables`、`luoshu` |

## 快速使用

```go
if err := tables.Load(); err != nil {
    panic(err)
}

chart := bazi.ComputeChartAt(birth, 116.407526, ganzhi.Male)
```

先调用 `tables.Load()` 建立 fail-closed 启动屏障；随后使用领域包的标准入口计算确定性事实。

## 核心文档

- [英文 README](../../README.md)
- [MCP 工具目录](../../contracts/mcp-tools.json)
- [MCP adapter 决策记录](../decisions/0003-official-mcp-adapter.md)
- [MIT License](../../LICENSE)
