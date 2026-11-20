<h1 align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="../../assets/logo-dark.svg">
    <img alt="MingliGo logo" src="../../assets/logo.svg" width="72" align="absmiddle">
  </picture>
  MingliGo
</h1>

> **El README en inglés es la documentación canónica.**
> Esta página es un resumen en español, no una traducción completa.

MingliGo es una biblioteca de Go de código abierto y uso general para cálculos de metafísica china. Calcula cartas, tableros y disposiciones reproducibles a partir de entradas explícitas de tiempo, ubicación, tirada y escuela. Los paquetes de dominio pueden usarse directamente; el adapter MCP es solo un ejemplo y no forma parte de la API.

## Límites de diseño

| Límite | Responsabilidad |
| --- | --- |
| `pkg/*` | Cálculos de dominio, tablas y validación |
| `examples/mingligo-mcp` | Typed tool adapter del SDK oficial y composition root |
| `examples/mingligo-mcp/contracts/mcp-tools.json` | Catálogo de herramientas MCP del adapter de ejemplo |

El núcleo de la biblioteca nunca importa código de MCP, transporte, producto ni despliegue. Los DTO de MCP y el mapeo de transporte permanecen en el adapter de referencia.

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
- [Catálogo de herramientas MCP](../../examples/mingligo-mcp/contracts/mcp-tools.json)
- [Decisión del adapter MCP](../ARCHITECTURE.md)
- [Licencia MIT](../../LICENSE)
