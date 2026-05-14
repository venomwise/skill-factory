Clarification: Added a `--dry-run` option to `scripts/export.py` so exports can be previewed without creating directories or writing files.

Assumptions/defaults:
- Existing defaults remain unchanged:
  - `--input-dir input`
  - `--output-dir out`
- Dry run reports the target paths that would be exported.
- In dry-run mode, the output directory is not created and no files are written.

Implementation summary:
- Updated `export_files()` to accept `dry_run: bool = False`.
- Skips `output_dir.mkdir(...)` and `target.write_text(...)` when `dry_run` is enabled.
- Added CLI flag:

```bash
--dry-run
```

- Output now says:
  - `would export ...` during dry runs
  - `exported ...` during normal runs

Validation:
- Ran dry run:

```bash
python3 scripts/export.py --output-dir /tmp/export-dry-run-test --dry-run
```

Confirmed it printed the expected files and did not create the output directory.

- Ran normal export:

```bash
python3 scripts/export.py --output-dir /tmp/export-run-test
```

Confirmed `a.txt` and `b.txt` were written successfully.
