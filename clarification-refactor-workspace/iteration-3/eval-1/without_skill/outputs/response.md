Perfect! I've added the `--dry-run` option to `scripts/export.py`. Here's what changed:

**Key changes:**
- Replaced simple `sys.argv` parsing with `argparse` for better argument handling
- Added `--dry-run` flag that shows what would be exported without writing files
- When dry-run is enabled, the script displays:
  - The number of users that would be exported
  - The directory that would be created
  - A preview of the actual data

**Usage:**
```bash
# Normal export (writes file)
python scripts/export.py data/users.json

# Dry run (preview only, no files written)
python scripts/export.py data/users.json --dry-run

# Help
python scripts/export.py --help
```

The dry-run mode will print `[DRY RUN]` prefixes to make it clear no actual changes are being made, and it shows the full JSON data that would be written.
