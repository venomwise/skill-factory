Done. I've added the `--dry-run` option to `scripts/export.py`. 

Changes made:
- Switched from manual `sys.argv` parsing to `argparse` for better argument handling
- Added `--dry-run` flag that shows what would be exported without writing files
- When dry-run is enabled, the script displays the output path, full data preview, and summary
- Normal execution still works as before, actually writing the file

Verified both modes work correctly:
- `--dry-run` shows the preview without creating files
- Without the flag, files are written normally
