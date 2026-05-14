Perfect! I've added the `--dry-run` option to `scripts/export.py`. Here's what changed:

**Key changes:**
- Replaced simple `sys.argv` parsing with `argparse` for better argument handling
- Added `--dry-run` flag that shows what would be exported without writing files
- Modified `export_data()` to accept a `dry_run` parameter
- When dry-run is enabled, it prints a preview of the data and the file path that would be created

**Usage:**
```bash
# Normal export (writes file)
python scripts/export.py data/users.json

# Dry run (preview only, no files written)
python scripts/export.py data/users.json --dry-run

# Help
python scripts/export.py --help
```

The dry-run mode will show you the directory that would be created, the number of records, and a full JSON preview of the data without actually writing anything to disk.
