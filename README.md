# MingliGo

MingliGo is the pure Go deterministic fact core for Chinese metaphysics
computation.

It exposes stable domain packages for direct Go use and ships a reference MCP
server for quick validation.

## Library use

```go
import "github.com/ml8s/mingligo/bazi"
```

## Reference MCP

```bash
make run
```

The local reference server defaults to `local/mingligo-mcp`. Hosted deployments
inject their own service identity.
