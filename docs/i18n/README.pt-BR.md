<h1 align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="../../assets/logo-dark.svg">
    <img alt="MingliGo logo" src="../../assets/logo.svg" width="72" align="absmiddle">
  </picture>
  MingliGo
</h1>

> **O README em inglês é a documentação canônica.**
> Esta página é um resumo em português do Brasil, não uma tradução completa.

MingliGo é um núcleo determinista de fatos em Go puro para cálculos de metafísica chinesa. Ele expõe pacotes de domínio estáveis para uso direto. O servidor `examples/mingligo-mcp` é um adapter MCP de exemplo. O projeto Liki usa esse exemplo, mas ele não faz parte do produto da biblioteca.

## Limites de design

| Limite | Responsabilidade |
| --- | --- |
| `pkg/*` | Fatos deterministas do domínio, tabelas e validação |
| `examples/mingligo-mcp` | Typed tool adapter do SDK oficial e composition root |
| `examples/mingligo-mcp/contracts/mcp-tools.json` | Catálogo de ferramentas MCP do adapter de exemplo |

O núcleo de fatos nunca importa código de MCP, transporte, produto ou implantação. Os DTOs de MCP e o mapeamento de transporte permanecem no adapter de referência.

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
- [Decisão do adapter MCP](../decisions/0003-official-mcp-adapter.md)
- [Licença MIT](../../LICENSE)
