Clarification: Added a `--dry-run` option to `scripts/export.py` so it reports planned exports without creating directories or writing files.

Assumptions/defaults:
- Existing defaults remain unchanged: `--input-dir input`, `--output-dir out`.
- Dry-run output uses the resolved target path under `--output-dir`.
- In dry-run mode, the script still scans `*.txt` files from the input directory, but performs no filesystem writes.

Implementation summary:
- Updated `export_files(input_dir, output_dir, dry_run=False)`.
- Skips `output_dir.mkdir(...)` and `target.write_text(...)` when `dry_run=True`.
- Added CLI flag:

```bash
--dry-run
```

- Dry-run prints:

```text
would export <path>
```

- Normal mode still prints:

```text
exported <path>
```

Validation:
- Ran dry-run:

```bash
python3 scripts/export.py --dry-run --output-dir dryout
```

Confirmed it printed planned exports and did not create `dryout`.

- Ran normal export:

```bash
python3 scripts/export.py --output-dir out
```

Confirmed files were written.

- Ran syntax check:

```bash
python3 -m py_compile scripts/export.py
```

No errors.
