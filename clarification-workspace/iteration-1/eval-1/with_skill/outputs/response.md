Clarification: Added a `--dry-run` option to `scripts/export.py` that shows which files would be exported without creating directories or writing files.

Assumptions/defaults:
- Dry run should preserve existing defaults: `--input-dir input`, `--output-dir out`.
- Dry run output reports target paths only, matching the existing export reporting style.
- In dry-run mode, the output directory is not created.

Implementation summary:
- Updated `export_files()` to accept `dry_run: bool = False`.
- Skips `mkdir()` and `write_text()` when `dry_run` is enabled.
- Added CLI flag:

```bash
--dry-run
```

- Dry-run output now prints:

```text
would export <path>
```

instead of:

```text
exported <path>
```

Validation:
- Ran dry-run export and confirmed no output directory was created:

```bash
python3 scripts/export.py --output-dir /tmp/export-dry-run-test --dry-run
test ! -e /tmp/export-dry-run-test
```

- Ran normal export and confirmed files were written:

```bash
python3 scripts/export.py --output-dir /tmp/export-run-test
test -f /tmp/export-run-test/a.txt
test -f /tmp/export-run-test/b.txt
```
