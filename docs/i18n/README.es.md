<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="../../assets/logo-dark.svg">
    <img alt="MingliGo logo" src="../../assets/logo.svg" width="72">
  </picture>
</p>

# MingliGo

> **El README en inglés es la documentación canónica.**
> Esta página es un resumen en español, no una traducción completa.

MingliGo es un núcleo determinista de hechos en Go puro para cálculos de metafísica china. Expone paquetes de dominio estables para uso directo. El servidor `examples/mingligo-mcp` es un adapter MCP de ejemplo para validar esa superficie; no forma parte del núcleo de hechos.

## Límites de diseño

| Límite | Responsabilidad |
| --- | --- |
| `pkg/*` | Hechos deterministas del dominio, tablas y validación |
| `examples/mingligo-mcp` | Typed tool adapter del SDK oficial y composition root |
| `contracts/mcp-tools.json` | Catálogo público de herramientas MCP |

El núcleo de hechos nunca importa código de MCP, transporte, producto ni despliegue. Los DTO de MCP y el mapeo de transporte permanecen en el adapter de referencia.

## Paquetes de dominio

| Área | Paquetes |
| --- | --- |
| Cartas natales | `bazi`, `ziwei` |
| Adivinación | `liuyao`, `qimen` |
| Feng shui | `bazhai`, `xuankong`, `fengshui` |
| Calendario y ubicación | `huangli`, `tianwen`, `city`, `ganzhi` |
| Base compartida | `tables`, `luoshu` |

## Uso esencial

```go
if err := tables.Load(); err != nil {
    panic(err)
}

chart := bazi.ComputeChartAt(birth, 116.407526, ganzhi.Male)
```

Llame primero a `tables.Load()` para crear una barrera fail-closed; luego use las entradas estándar de cada paquete para calcular hechos deterministas.

## Documentación principal

- [README en inglés](../../README.md)
- [Catálogo de herramientas MCP](../../contracts/mcp-tools.json)
- [Decisión del adapter MCP](../decisions/0003-official-mcp-adapter.md)
- [Licencia MIT](../../LICENSE)
