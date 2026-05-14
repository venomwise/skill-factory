Implemented `--dry-run` in `scripts/export.py`.

What changed:
- Added CLI flag: `--dry-run`
- Dry run now:
  - Lists files that would be exported
  - Does not create the output directory
  - Does not write any files
- Normal export behavior is unchanged

Example:

```bash
python3 scripts/export.py --dry-run
```

Output:

```text
would export out/a.txt
would export out/b.txt
```

I also verified normal export still writes files correctly.
