#!/bin/env python3

from json import load, dump
from pathlib import Path
from shutil import move
from sys import argv, stderr
from tempfile import NamedTemporaryFile
from typing import List

if __name__ == "__main__":
    if len(argv) != 3:
        print("too few args", file=stderr)
        print("", file=stderr)
        print(f"{argv[0]} <> <>", file=stderr)
        exit(1)

    messages_file = Path(argv[1])
    if not messages_file.exists():
        print(f"no file? {messages_file}", file=stderr)
        exit(1)

    out_file = Path(argv[2])
    if not out_file.exists():
        print(f"no file? {out_file}", file=stderr)
        exit(1)

    with messages_file.open(encoding="utf-8") as bf:
        merged = load(bf)
        merged_ids: List[str] = []

        for x in merged['messages']:
            if x['id'] not in merged_ids:
                merged_ids.append(x['id'])

        with out_file.open(encoding="utf-8") as mf:
            updated = load(mf)

            for x in updated['messages']:
                if x['id'] not in merged_ids:
                    print(f"  Adding: '{x['id']}'")
                    merged['messages'].append(x)

            tmpf = NamedTemporaryFile("w", prefix="merged", suffix=".json", encoding="utf-8", delete=False)
            with tmpf as f:
                dump(merged, f, indent=' ' * 8, ensure_ascii=False)

        newpath = move(tmpf.name, messages_file)
        print(tmpf.name, "->", newpath)
