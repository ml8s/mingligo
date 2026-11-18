# MingliGo

MingliGo is the pure Go deterministic fact core for Chinese metaphysics
computation.

It exposes stable domain packages for direct Go use and ships a reference MCP
server for quick validation.

## Library use

```go
import (
    "time"

    "github.com/ml8s/mingligo/pkg/bazi"
    "github.com/ml8s/mingligo/pkg/ganzhi"
    "github.com/ml8s/mingligo/pkg/tables"
)

func main() {
    // Embedded tables must be loaded once before domain calls.
    if err := tables.Load(); err != nil {
        panic(err)
    }

    when := time.Date(1984, 2, 4, 6, 0, 0, 0, time.FixedZone("CST", 8*3600))
    chart := bazi.ComputeChartAt(when, 120, ganzhi.Male)
    _ = chart
}
```

## Reference MCP

```bash
make run
```

The local reference server defaults to `local/mingligo-mcp`. The `mingligo-mcp`
server is an example adapter, not the fact core. Use `make runimage` for the
same image-based validation path.
