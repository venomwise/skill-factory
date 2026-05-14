Clarification: Added a `--dry-run` CLI option to `scripts/export.py` so it reports planned exports without creating the output directory or writing files.

Assumptions/defaults:
- Existing defaults remain unchanged: `--input-dir input`, `--output-dir out`.
- Dry run should still discover the same `.txt` files and show their target paths.
- Dry run should not create `--output-dir` if it does not already exist.

Implementation summary:
- Updated `export_files(input_dir, output_dir, dry_run=False)`.
- Skips `output_dir.mkdir(...)` and `target.write_text(...)` when `dry_run=True`.
- Added argparse flag:

```bash
--dry-run
```

- Output changes:
  - Normal mode: `exported out/a.txt`
  - Dry-run mode: `would export out/a.txt`

Validation:
- Ran dry-run:

```bash
python3 scripts/export.py --dry-run --output-dir dry_out
```

Output:

```text
would export dry_out/a.txt
would export dry_out/b.txt
```

Confirmed `dry_out` was not created.

- Ran normal export:

```bash
python3 scripts/export.py --output-dir out
```

Confirmed files were written to `out/a.txt` and `out/b.txt`.
