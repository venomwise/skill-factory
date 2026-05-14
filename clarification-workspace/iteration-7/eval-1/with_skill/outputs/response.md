Done! I've added a `--dry-run` option to `scripts/export.py`. Here's what changed:

**New features:**
- `--dry-run` flag that previews the export without writing files
- Switched to `argparse` for better argument handling and help text
- Dry-run mode shows:
  - What file would be created
  - Parent directory path
  - Estimated file size
  - Preview of the data (first 500 bytes)

**Usage:**
```bash
# Normal export (writes file)
python scripts/export.py data/users.json

# Dry-run (preview only)
python scripts/export.py data/users.json --dry-run

# Help
python scripts/export.py --help
```

The dry-run mode validates everything except the actual file write, so you can safely test paths and see what data would be exported.
