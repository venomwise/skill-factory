Clarification: added a `--dry-run` CLI option to `scripts/export.py` so the script reports what it would export without creating directories or writing files.

Assumptions/defaults:
- Existing defaults remain unchanged: `--input-dir input`, `--output-dir out`.
- Dry run output uses `would export <path>`.
- Normal export behavior remains `exported <path>` and still writes copied `.txt` files.

Implementation summary:
- Updated `export_files(input_dir, output_dir, dry_run=False)`.
- Skips `output_dir.mkdir(...)` during dry run.
- Skips `target.write_text(...)` during dry run.
- Still returns the target paths that would be exported.
- Added argparse flag:

```bash
--dry-run
```

Validation:
- Ran dry run:

```bash
python3 scripts/export.py --dry-run --output-dir dry-out
```

Output:

```text
would export dry-out/a.txt
would export dry-out/b.txt
```

Confirmed `dry-out` was not created.

- Ran normal export:

```bash
python3 scripts/export.py --output-dir out
```

Output:

```text
exported out/a.txt
exported out/b.txt
```

Confirmed files were written with expected contents.

Changed file:
- `scripts/export.py`
