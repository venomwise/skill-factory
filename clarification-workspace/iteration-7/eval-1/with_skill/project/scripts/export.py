#!/usr/bin/env python3
"""Export data to JSON file."""
import argparse
import json
import sys
from pathlib import Path


def export_data(output_file, dry_run=False):
    """Export user data to JSON file.
    
    Args:
        output_file: Path to output JSON file
        dry_run: If True, show what would be exported without writing
    """
    # Simulate fetching data from database
    data = {
        "users": [
            {"id": 1, "name": "Alice", "email": "alice@example.com"},
            {"id": 2, "name": "Bob", "email": "bob@example.com"},
            {"id": 3, "name": "Charlie", "email": "charlie@example.com"}
        ],
        "exported_at": "2024-01-15T10:30:00Z",
        "total_count": 3
    }
    
    # Write to file
    output_path = Path(output_file)
    
    if dry_run:
        print(f"[DRY RUN] Would export {data['total_count']} users to {output_file}")
        print(f"[DRY RUN] Parent directory: {output_path.parent.absolute()}")
        print(f"[DRY RUN] File size: ~{len(json.dumps(data, indent=2))} bytes")
        print("\n[DRY RUN] Preview of data:")
        print(json.dumps(data, indent=2)[:500] + "..." if len(json.dumps(data, indent=2)) > 500 else json.dumps(data, indent=2))
    else:
        output_path.parent.mkdir(parents=True, exist_ok=True)
        
        with open(output_path, 'w') as f:
            json.dump(data, f, indent=2)
        
        print(f"✓ Exported {data['total_count']} users to {output_file}")
    
    return output_path


def main():
    parser = argparse.ArgumentParser(
        description="Export user data to JSON file",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="Example: export.py data/users.json\n"
               "         export.py data/users.json --dry-run"
    )
    parser.add_argument(
        "output_file",
        help="Path to output JSON file"
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Show what would be exported without writing files"
    )
    
    args = parser.parse_args()
    export_data(args.output_file, dry_run=args.dry_run)


if __name__ == "__main__":
    main()
