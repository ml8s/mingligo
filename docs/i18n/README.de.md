<h1 align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="../../assets/logo-dark.svg">
    <img alt="MingliGo logo" src="../../assets/logo.svg" width="72" align="absmiddle">
  </picture>
  MingliGo
</h1>

> **Das englische README ist die kanonische Dokumentation.**
> Diese Seite ist eine lokalisierte Übersicht und keine vollständige Übersetzung.

MingliGo ist eine allgemeine Open-Source-Go-Bibliothek für Berechnungen der chinesischen Metaphysik. Sie berechnet aus expliziten Zeit-, Orts-, Wurf- und Schulen-Eingaben reproduzierbare Karten, Tafeln und Layouts. Die Domänenpakete sind direkt nutzbar; der MCP-Adapter ist nur ein Beispiel und keine Bibliotheks-API.

## Designgrenzen

| Grenze | Verantwortung |
| --- | --- |
| `pkg/*` | Domänenberechnungen, Tabellen und Validierung |
| `examples/mingligo-mcp` | Typed Tool Adapter des offiziellen SDK und Composition Root |
| `examples/mingligo-mcp/contracts/mcp-tools.json` | MCP-Toolkatalog des Beispieladapters |

Der Bibliothekskern importiert keinen MCP-, Transport-, Produkt- oder Deployment-Code. MCP-DTOs und Wire-Mapping bleiben im Referenzadapter.

## Domänenpakete

| Bereich | Pakete |
| --- | --- |
| Natalcharts | `bazi`, `ziwei` |
| Weissagung | `liuyao`, `qimen` |
| Feng Shui | `bazhai`, `xuankong`, `fengshui` |
| Kalender und Ort | `huangli`, `tianwen`, `city`, `ganzhi` |
| Gemeinsame Basis | `tables`, `luoshu` |

## Wesentliche Nutzung

```go
if err := tables.Load(); err != nil {
    panic(err)
}

chart := bazi.ComputeChartAt(birth, 116.407526, ganzhi.Male)
```

Rufen Sie zuerst `tables.Load()` auf, um eine fail-closed-Startbarriere zu schaffen; verwenden Sie danach die Standardeinträge der Domänenpakete für deterministische Berechnungen.

## Hauptdokumentation

- [Englisches README](../../README.md)
- [MCP-Toolkatalog](../../examples/mingligo-mcp/contracts/mcp-tools.json)
- [Entscheidung zum MCP-Adapter](../ARCHITECTURE.md)
- [MIT-Lizenz](../../LICENSE)
