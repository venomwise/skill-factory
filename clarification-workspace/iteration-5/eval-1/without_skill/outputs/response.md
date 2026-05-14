Implemented `--dry-run` in `scripts/export.py`.

What changed:
- Added `--dry-run` CLI flag.
- In dry-run mode:
  - Does not create the output directory.
  - Does not write any files.
  - Prints what would be exported, e.g. `would export out/a.txt`.
- Normal export behavior is unchanged.

Verified with:

```bash
python3 scripts/export.py --dry-run
```

Output:

```text
would export out/a.txt
would export out/b.txt
```

And confirmed no `out` directory was created.
