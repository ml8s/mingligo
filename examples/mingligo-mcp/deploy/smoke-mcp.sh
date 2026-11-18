#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
IMAGE_TAG="${IMAGE_TAG:-local}"
IMAGE_REGISTRY="${IMAGE_REGISTRY:-liki-local}"
IMAGE="${MINGLIGO_MCP_IMAGE:-$IMAGE_REGISTRY/mingligo-mcp:$IMAGE_TAG}"
PORT="${MINGLIGO_MCP_SMOKE_PORT:-18091}"
NAME="mingligo-smoke-mcp-$$"
cleanup() { docker rm -f "$NAME" >/dev/null 2>&1 || true; }
trap cleanup EXIT
docker run -d --name "$NAME" -p "$PORT:8081" "$IMAGE"
for _ in $(seq 1 30); do
  curl -fsS "http://127.0.0.1:$PORT/health" >/dev/null 2>&1 && break
  sleep 1
done
MINGLIGO_SERVER_NAME="${MINGLIGO_SERVER_NAME:-local/mingligo-mcp}" \
  python3 "$ROOT/scripts/test_mcp.py" "http://127.0.0.1:$PORT/mcp"
