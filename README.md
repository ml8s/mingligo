# MingliGo

MingliGo is the pure Go deterministic fact core for Chinese metaphysics
computation.

It exposes stable domain packages for direct Go use and ships a reference MCP
server for quick validation.

## Library use

```go
import "github.com/ml8s/mingligo/pkg/bazi"
```

## Reference MCP

```bash
make run
```

The local reference server defaults to `local/mingligo-mcp`. The `mingligo-mcp`
server is an example adapter, not the fact core. Use `make runimage` for the
same image-based validation path.
