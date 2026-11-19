<h1 align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="../../assets/logo-dark.svg">
    <img alt="MingliGo logo" src="../../assets/logo.svg" width="72" align="absmiddle">
  </picture>
  MingliGo
</h1>

> **英語版 README が正本です。**
> このページは日本語の概要であり、全文翻訳ではありません。

MingliGo は、中国命理計算のための pure Go 製 deterministic fact core です。安定したドメインパッケージを直接利用できます。`examples/mingligo-mcp` はサンプル MCP adapter であり、Liki プロジェクトが利用していますが、ライブラリの製品ではありません。

## 設計境界

| 境界 | 責務 |
| --- | --- |
| `pkg/*` | 決定論的なドメインファクト、テーブル、検証 |
| `examples/mingligo-mcp` | official SDK の typed tool adapter と composition root |
| `examples/mingligo-mcp/contracts/mcp-tools.json` | サンプル adapter の MCP ツールカタログ |

fact core は MCP、transport、プロダクト、デプロイのコードを import しません。MCP DTO と wire mapping はリファレンス adapter 内にだけ置きます。

## ドメインパッケージ

| 分野 | パッケージ |
| --- | --- |
| 命盤 | `bazi`、`ziwei` |
| 占断 | `liuyao`、`qimen` |
| 風水 | `bazhai`、`xuankong`、`fengshui` |
| 暦と位置 | `huangli`、`tianwen`、`city`、`ganzhi` |
| 共通基盤 | `tables`、`luoshu` |

## 利用の要点

```go
if err := tables.Load(); err != nil {
    panic(err)
}

chart := bazi.ComputeChartAt(birth, 116.407526, ganzhi.Male)
```

最初に `tables.Load()` を呼び出して fail-closed の起動障壁を作り、その後に各ドメインパッケージの標準エントリポイントで決定論的なファクトを計算します。

## 主要ドキュメント

- [英語 README](../../README.md)
- [MCP ツールカタログ](../../examples/mingligo-mcp/contracts/mcp-tools.json)
- [MCP adapter の意思決定記録](../decisions/0003-official-mcp-adapter.md)
- [MIT License](../../LICENSE)
