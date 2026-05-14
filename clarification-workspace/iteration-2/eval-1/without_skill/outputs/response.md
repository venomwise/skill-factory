Clarification: Added a `--dry-run` CLI option to `scripts/export.py` so the script reports intended exports without creating directories or writing files.

Assumptions/defaults:
- Default behavior is unchanged: without `--dry-run`, files are copied from `input/` to `out/`.
- Dry run still uses the configured `--input-dir` and `--output-dir`.
- Dry run output uses `would export <path>` to distinguish it from actual exports.
- No files or output directories are created during dry run.

Implementation summary:
- Updated `export_files()` to accept `dry_run: bool = False`.
- Skips `output_dir.mkdir(...)` and `target.write_text(...)` when dry run is enabled.
- Added argparse option:

```bash
--dry-run
```

- Updated CLI messages:
  - Normal mode: `exported out/a.txt`
  - Dry-run mode: `would export out/a.txt`

Validation:
- Ran dry run:

```bash
python3 scripts/export.py --dry-run --output-dir dryout
```

Output:

```text
would export dryout/a.txt
would export dryout/b.txt
```

Confirmed `dryout` was not created.

- Ran normal export:

```bash
python3 scripts/export.py --output-dir out
```

Output:

```text
exported out/a.txt
exported out/b.txt
```

Confirmed `out/a.txt` and `out/b.txt` were created.
