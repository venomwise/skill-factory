Implemented `--dry-run` in `scripts/export.py`.

Changes:
- Added `--dry-run` CLI flag.
- In dry-run mode:
  - Does not create the output directory.
  - Does not write files.
  - Prints `would export <path>` for each file.
- Normal export behavior remains unchanged.

Verified with:
- `python3 scripts/export.py --dry-run`
- `python3 scripts/export.py`
