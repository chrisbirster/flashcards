#!/usr/bin/env python3

import sys
from pathlib import Path


def main() -> int:
    if len(sys.argv) != 4:
        print("usage: upsert_env.py <env_file> <key> <value>", file=sys.stderr)
        return 1

    env_file = Path(sys.argv[1])
    key = sys.argv[2].strip()
    value = sys.argv[3]
    if not key:
        print("key is required", file=sys.stderr)
        return 1

    env_file.parent.mkdir(parents=True, exist_ok=True)
    lines = env_file.read_text().splitlines() if env_file.exists() else []

    updated = False
    prefix = f"{key}="
    for index, line in enumerate(lines):
        if line.startswith(prefix):
            lines[index] = f"{prefix}{value}"
            updated = True
            break

    if not updated:
        lines.append(f"{prefix}{value}")

    env_file.write_text("\n".join(lines) + "\n")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
