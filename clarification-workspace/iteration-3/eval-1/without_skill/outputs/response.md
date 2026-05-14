Clarification: Added `--dry-run` to `scripts/export.py` so it reports intended exports without creating directories or writing files.

Assumptions/defaults:
- Default behavior is unchanged: without `--dry-run`, files are copied from `--input-dir` to `--output-dir`.
- In dry-run mode, the script still discovers the same `*.txt` input files and prints the target paths.
- Dry-run output uses `would export ...` instead of `exported ...`.

Implementation summary:
- Updated `export_files(input_dir, output_dir, dry_run=False)`.
- Skips `output_dir.mkdir(...)` when `dry_run=True`.
- Skips `target.write_text(...)` when `dry_run=True`.
- Added CLI flag:

```bash
--dry-run
```

Validation:
- Ran:

```bash
python3 scripts/export.py --dry-run --output-dir dryrun-out
```

Output:

```text
would export dryrun-out/a.txt
would export dryrun-out/b.txt
```

Confirmed `dryrun-out` was not created.

- Also verified normal export still works:

```bash
python3 scripts/export.py --output-dir out
```

Output:

```text
exported out/a.txt
exported out/b.txt
```
