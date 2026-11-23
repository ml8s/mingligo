#!/usr/bin/env python3
"""Compare every deterministic MingliGo method across API and MCP transports.

This is a release-time evidence tool, not runtime dependency: start both
example servers, then run it against their local endpoints. time_now is
compared structurally because its value is intentionally wall-clock dependent.
"""

import argparse
import json
import os
import sys
import urllib.request
from pathlib import Path


API_DIR = Path(__file__).resolve().parents[1]
FIXTURE_DIR = API_DIR / "testdata" / "wire_inputs"


def api_call(base_url: str, method: str, arguments: dict) -> object:
    try:
        return _api_call(base_url, method, arguments)
    except urllib.error.HTTPError as error:
        detail = error.read().decode("utf-8", "replace")
        raise RuntimeError(f"{method} API returned {error.code}: {detail}") from error


def _api_call(base_url: str, method: str, arguments: dict) -> object:
    api_url = base_url.rstrip("/")
    if not api_url.endswith("/mingligo/api"):
        api_url += "/mingligo/api"
    request = urllib.request.Request(
        f"{api_url}/{method}",
        data=json.dumps(arguments, separators=(",", ":")).encode("utf-8"),
        headers={"Content-Type": "application/json", "Accept": "application/json"},
        method="POST",
    )
    with urllib.request.urlopen(request, timeout=15) as response:
        return json.load(response)


def mcp_call(base_url: str, method: str, arguments: dict, request_id: int) -> object:
    params = {
        "name": method,
        "arguments": arguments,
        "_meta": {
            "io.modelcontextprotocol/protocolVersion": "2026-07-28",
            "io.modelcontextprotocol/clientCapabilities": {},
        },
    }
    body = json.dumps({
        "jsonrpc": "2.0", "id": request_id, "method": "tools/call", "params": params,
    }).encode("utf-8")
    request = urllib.request.Request(
        f"{base_url.rstrip('/')}/mcp",
        data=body,
        headers={
            "Accept": "application/json, text/event-stream",
            "Content-Type": "application/json",
            "Mcp-Method": "tools/call",
            "Mcp-Name": method,
            "MCP-Protocol-Version": "2026-07-28",
        },
        method="POST",
    )
    with urllib.request.urlopen(request, timeout=15) as response:
        raw = response.read().decode("utf-8")
    document = json.loads(raw) if raw.lstrip().startswith("{") else None
    if document is None:
        for line in raw.splitlines():
            if line.startswith("data: "):
                document = json.loads(line[6:])
                break
    if document is None:
        raise RuntimeError(f"unparseable MCP response for {method}: {raw[:200]!r}")
    if document.get("error") or document.get("result", {}).get("isError"):
        raise RuntimeError(f"{method} failed: {document}")
    payload = document.get("result", {}).get("structuredContent")
    if payload is None:
        raise RuntimeError(f"{method} returned no structuredContent")
    return payload


def fixture_method(name: str) -> str:
    domain, _, rest = name.partition("_")
    if domain == "qimen":
        return "qimen_chart"
    tokens = rest.split("_")
    for token in reversed(tokens):
        if token in {"bond", "chart"}:
            return f"{domain}_{token}"
        if token in {"liunian", "liuyue", "liuri", "liushi", "xiaoyun"}:
            return f"{domain}_liupan"
    raise ValueError(f"cannot map fixture {name}")


def load_cases() -> list[tuple[str, dict]]:
    inline = {
		"time_now": {"timezone": "Asia/Shanghai"},
        "subject_candidates": {
            "calendar_kind": "solar", "recorded_date": "1990-05-20",
            "record_timezone": "Asia/Shanghai", "recorded_time": "11:49",
            "longitude": 116.4, "subject_gender": "male",
        },
        "city_info": {"city": "北京"},
        "huangli_days": {
            "standard_time": "2024-06-15T12:00:00+08:00", "count": 1, "event": "sign",
        },
        "bazhai_chart": {
            "subject": {
                "standard_time": "1984-06-15T12:00:00+08:00", "gender": "male",
            },
        },
        "bazhai_layout": {
            "ming_gua": "乾", "men_gua": "坤", "master_gua": "艮", "stove_gua": "兑",
        },
        "xuankong_chart": {
            "standard_time": "2026-07-31T12:00:00+08:00",
            "zuo_shan": 2, "xiang_shan": 14,
        },
        "xuankong_liunian": {
            "standard_time": "2026-06-15T12:00:00+08:00",
            "zuo_shan": 2, "xiang_shan": 14,
        },
        "liuyao_qigua": {
            "origin": {"fang_shi": "yaos", "yao_zhi": [7, 8, 7, 8, 7, 8]},
            "qi_gua_qing_jing": {
                "shi_jian": "2026-06-15T12:00:00+08:00", "jing_du": 116.4,
                "focus": {"kind": "qu_yong", "qu_yong": "官鬼"},
            },
        },
    }
    cases: dict[str, dict] = {**inline}
    for path in sorted(FIXTURE_DIR.glob("*.json")):
        method = fixture_method(path.stem)
        fixture = json.loads(path.read_text(encoding="utf-8"))
        cases.setdefault(method, fixture["input"])
    return sorted(cases.items())


def canonical(value: object) -> str:
    return json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":"))


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--api-url", default=os.environ.get("MINGLIGO_API_URL", "http://127.0.0.1:18082/mingligo/api"))
    parser.add_argument("--mcp-url", default=os.environ.get("MINGLIGO_MCP_URL", "http://127.0.0.1:8081"))
    args = parser.parse_args()

    api_base, mcp_base = args.api_url.rstrip("/"), args.mcp_url.rstrip("/")
    compared: set[str] = set()
    qigua_api: dict | None = None
    request_id = 1
    for method, arguments in load_cases():
        request_id += 1
        api_output = api_call(api_base, method, arguments)
        mcp_output = mcp_call(mcp_base, method, arguments, request_id)
        if method == "time_now":
            if set(api_output) != set(mcp_output) or not api_output.get("local_standard_time"):
                raise RuntimeError(f"time_now shape mismatch: {api_output} / {mcp_output}")
        elif canonical(api_output) != canonical(mcp_output):
            raise RuntimeError(
                f"{method} output mismatch\nAPI: {canonical(api_output)}\nMCP: {canonical(mcp_output)}"
            )
        compared.add(method)
        if method == "liuyao_qigua":
            qigua_api = api_output

    if qigua_api is None:
        raise RuntimeError("liuyao_qigua comparison was not run")
    chart_input = {"qi_gua": qigua_api["qi_gua"]}
    request_id += 1
    api_chart = api_call(api_base, "liuyao_chart", chart_input)
    mcp_chart = mcp_call(mcp_base, "liuyao_chart", chart_input, request_id)
    if canonical(api_chart) != canonical(mcp_chart):
        raise RuntimeError(f"liuyao_chart output mismatch\n{canonical(api_chart)}\n{canonical(mcp_chart)}")
    compared.add("liuyao_chart")

    expected = {
        "bazhai_chart", "bazhai_layout", "bazi_bond", "bazi_chart", "bazi_liupan",
        "city_info", "huangli_days", "liuyao_chart", "liuyao_qigua",
        "qimen_chart", "subject_candidates", "time_now", "xuankong_chart",
        "xuankong_liunian", "ziwei_bond", "ziwei_chart", "ziwei_liupan",
    }
    if compared != expected:
        raise RuntimeError(f"coverage mismatch: missing={sorted(expected-compared)} extra={sorted(compared-expected)}")
    print(f"✓ API↔MCP: {len(compared)} methods compared")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as error:
        print(f"API/MCP comparison failed: {error}", file=sys.stderr)
        raise SystemExit(1)
