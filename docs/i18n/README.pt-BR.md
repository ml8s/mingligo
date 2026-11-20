<h1 align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="../../assets/logo-dark.svg">
    <img alt="MingliGo logo" src="../../assets/logo.svg" width="72" align="absmiddle">
  </picture>
  MingliGo
</h1>

> **O README em inglês é a documentação canônica.**
> Esta página é um resumo em português do Brasil, não uma tradução completa.

MingliGo é uma biblioteca Go de código aberto e uso geral para cálculos de metafísica chinesa. Ela calcula cartas, tabuleiros e disposições reproduzíveis a partir de entradas explícitas de tempo, localidade, lançamento e escola. Os pacotes de domínio podem ser usados diretamente; o adapter MCP é apenas um exemplo e não faz parte da API.

## Limites de design

| Limite | Responsabilidade |
| --- | --- |
| `pkg/*` | Cálculos de domínio, tabelas e validação |
| `examples/mingligo-mcp` | Typed tool adapter do SDK oficial e composition root |
| `examples/mingligo-mcp/contracts/mcp-tools.json` | Catálogo de ferramentas MCP do adapter de exemplo |

O núcleo da biblioteca nunca importa código de MCP, transporte, produto ou implantação. Os DTOs de MCP e o mapeamento de transporte permanecem no adapter de referência.

## Pacotes de domínio

| Área | Pacotes |
| --- | --- |
| Cartas natais | `bazi`, `ziwei` |
| Adivinhação | `liuyao`, `qimen` |
| Feng shui | `bazhai`, `xuankong`, `fengshui` |
| Calendário e localização | `huangli`, `tianwen`, `city`, `ganzhi` |
| Base compartilhada | `tables`, `luoshu` |

## Uso essencial

```go
if err := tables.Load(); err != nil {
    panic(err)
}

chart := bazi.ComputeChartAt(birth, 116.407526, ganzhi.Male)
```

Chame `tables.Load()` primeiro para criar uma barreira fail-closed; depois use as entradas padrão de cada pacote para calcular fatos deterministas.

## Documentação principal

- [README em inglês](../../README.md)
- [Catálogo de ferramentas MCP](../../examples/mingligo-mcp/contracts/mcp-tools.json)
- [Decisão do adapter MCP](../ARCHITECTURE.md)
- [Licença MIT](../../LICENSE)
