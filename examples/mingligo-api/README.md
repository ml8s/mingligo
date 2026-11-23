# mingligo-api

独立的 HTTP 示例应用：17 个领域方法一对一映射为 `POST /{方法名}`。
- 与 `mingligo-mcp` 完全独立（代码复制，不共享包）
- API 层只做协议：JSON 解码、调用领域 handler、写回 View；零命理逻辑
- 每个 200 响应先按反射生成的 OpenAPI schema 校验；契约漂移返回 500，而不是静默外泄坏结构
- 无 URL 版本前缀；服务命名空间为 `/mingligo/api`；REST-ish（计算动作用 POST）。
- 错误统一为 `application/problem+json`；领域层只返回语义错误，HTTP 错误映射只在适配层。
- `openapi.json` 由 `go run ./examples/mingligo-api/cmd/gen_openapi` 从领域 View + 输入 DTO 反射生成，禁止手写

## 运行
```
go run ./examples/mingligo-api
# POST http://localhost:18082/mingligo/api/bazi_chart  {"subject":{...}}
# GET  http://localhost:18082/mingligo/api/openapi.json
```

## 验证
```
go test ./examples/mingligo-api/...
```

发布前可同时启动 `mingligo-api` 与 `mingligo-mcp`，执行一次全量双跑对照：
```
go run ./examples/mingligo-api/scripts/compare_api_mcp.py \
  --api-url http://127.0.0.1:18082/mingligo/api \
  --mcp-url http://127.0.0.1:8081
```
`--api-url` 需要包含 `/mingligo/api`。`time_now` 只比较结构和字段存在性；其余 16 个方法（含链式 `liuyao_chart`）
比较完整 JSON。

## 传输压缩

本服务不做压缩——默认、始终返回未压缩 JSON。压缩由部署边缘的 Caddy 统一处理（`encode gzip`）：

- 内网/同机（如 counsel、web → API）：不压缩，回环带宽不是瓶颈
- 公网（浏览器 → Caddy）：Caddy 在边缘 gzip，浏览器自动解压

这是关注点分离：领域事实服务只管排盘与透传，传输优化交给反向代理。

## 客户端 SDK

`openapi.json` 由 `cmd/gen_openapi` 从领域 View + 输入 DTO 反射生成，随代码自动更新。消费方可据此生成任意语言的客户端：

```
openapi-generator generate -i openapi.json -g python -o ./client
```

生成的代码含方法签名、参数/返回类型、错误处理；消费方直接调用即可，无需手写 HTTP 请求。
