#!/usr/bin/env python3
"""Generate the neutral MingliGo engine capability artifact."""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import pathlib
import re
import subprocess


DIGEST_RE = re.compile(r"^[^@]+@sha256:[0-9a-f]{64}$")


def sha256_file(path: pathlib.Path) -> str:
    digest = hashlib.sha256()
    digest.update(path.read_bytes())
    return f"sha256:{digest.hexdigest()}"


def contracts_digest(root: pathlib.Path) -> str:
    digest = hashlib.sha256()
    for path in sorted(root.rglob("*")):
        if path.is_file():
            digest.update(path.relative_to(root).as_posix().encode("utf-8"))
            digest.update(b"\0")
            digest.update(path.read_bytes())
            digest.update(b"\0")
    return f"sha256:{digest.hexdigest()}"


def source_commit(repo: pathlib.Path) -> str:
    value = subprocess.check_output(
        ["git", "-C", str(repo), "rev-parse", "HEAD"], text=True
    ).strip()
    if not re.fullmatch(r"[0-9a-f]{40}", value):
        raise ValueError(f"invalid source commit: {value}")
    return value


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--repo", type=pathlib.Path, default=pathlib.Path.cwd())
    parser.add_argument("--image", required=True, help="pinned mingligo-mcp OCI image")
    parser.add_argument(
        "--identity",
        default=os.environ.get("MINGLIGO_SERVER_NAME", "local/mingligo-mcp"),
    )
    args = parser.parse_args()
    repo = args.repo.resolve()
    catalog = repo / "tools/mcp-tool-catalog.json"
    document = {
        "kind": "mingligo-capability",
        "project": "mingligo",
        "source_commit": source_commit(repo),
        "images": {"engine": args.image},
        "tools": {"engine": json.loads(catalog.read_text(encoding="utf-8"))["servers"]["engine"]},
        "digests": {
            "contracts": contracts_digest(repo / "contracts"),
            "tool_catalog": sha256_file(catalog),
        },
        "default_server_identity": {"engine": args.identity},
    }
    if not DIGEST_RE.fullmatch(args.image):
        raise SystemExit("image must be digest pinned")
    if "hk.liki" in json.dumps(document):
        raise SystemExit("MingliGo capability must not contain hk.liki identity")
    if args.identity == "hk.liki/engine":
        raise SystemExit("default identity must be neutral")
    print(json.dumps(document, ensure_ascii=False, indent=2, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
