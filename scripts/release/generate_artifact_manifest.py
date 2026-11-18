#!/usr/bin/env python3
"""Emit the minimal approved MingliGo artifact binding."""

from __future__ import annotations

import argparse
import os
import pathlib
import re
import subprocess


DIGEST_RE = re.compile(r"^[^@]+@sha256:[0-9a-f]{64}$")


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
    if not DIGEST_RE.fullmatch(args.image):
        raise SystemExit("image must be digest pinned")
    if "hk.liki" in args.identity or args.identity == "hk.liki/engine":
        raise SystemExit("default identity must be neutral")
    document = {
        "kind": "mingligo-artifact",
        "source_commit": source_commit(repo),
        "image": args.image,
        "default_server_identity": args.identity,
    }
    print(json.dumps(document, ensure_ascii=False, indent=2, sort_keys=True))
    return 0


if __name__ == "__main__":
    import json
    raise SystemExit(main())
