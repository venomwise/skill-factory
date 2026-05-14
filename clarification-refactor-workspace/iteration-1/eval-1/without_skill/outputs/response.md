Implemented `--dry-run` in `scripts/export.py`.

What changed:
- Added `--dry-run` CLI flag.
- In dry-run mode:
  - Does not create the output directory.
  - Does not write files.
  - Prints `would export <path>`.
- Normal export behavior is unchanged.

Verified with:

```bash
python3 scripts/export.py --dry-run --output-dir dryrun-out
```

Output:

```text
would export dryrun-out/a.txt
would export dryrun-out/b.txt
```

And confirmed `dryrun-out` was not created.
