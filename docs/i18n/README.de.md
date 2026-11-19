<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="../../assets/logo-dark.svg">
    <img alt="MingliGo logo" src="../../assets/logo.svg" width="72">
  </picture>
</p>

# MingliGo

> **Das englische README ist die kanonische Dokumentation.**
> Diese Seite ist eine lokalisierte Übersicht und keine vollständige Übersetzung.

MingliGo ist ein deterministischer Datenkern in pure Go für Berechnungen der chinesischen Metaphysik. Er stellt stabile Domänenpakete für die direkte Nutzung bereit. Der Server `examples/mingligo-mcp` ist ein MCP-Beispieladapter zur Validierung dieser Oberfläche und gehört nicht zum Datenkern.

## Designgrenzen

| Grenze | Verantwortung |
| --- | --- |
| `pkg/*` | Deterministische Domänenfakten, Tabellen und Validierung |
| `examples/mingligo-mcp` | Typed Tool Adapter des offiziellen SDK und Composition Root |
| `contracts/mcp-tools.json` | Öffentlicher MCP-Toolkatalog |

Der Datenkern importiert keinen MCP-, Transport-, Produkt- oder Deployment-Code. MCP-DTOs und Wire-Mapping bleiben im Referenzadapter.

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
- [MCP-Toolkatalog](../../contracts/mcp-tools.json)
- [Entscheidung zum MCP-Adapter](../decisions/0003-official-mcp-adapter.md)
- [MIT-Lizenz](../../LICENSE)
