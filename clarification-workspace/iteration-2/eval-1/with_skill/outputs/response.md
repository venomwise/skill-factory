Clarification: Added a `--dry-run` CLI option to `scripts/export.py` so you can preview exports without creating the output directory or writing files.

Assumptions/defaults:
- Default input directory remains `input`.
- Default output directory remains `out`.
- Dry-run output uses the same target paths that a real export would use.
- In dry-run mode, no directories or files are created.

Implementation summary:
- Updated `export_files()` to accept `dry_run: bool = False`.
- Skips `output_dir.mkdir(...)` and `target.write_text(...)` when `dry_run` is enabled.
- Added argparse flag:

```bash
--dry-run
```

- Dry-run prints:

```text
would export out/a.txt
```

instead of:

```text
exported out/a.txt
```

Validation:
- Ran dry-run and confirmed it reports intended exports without creating `out`.
- Ran normal export to a temporary output directory and confirmed files are still written correctly.
