#!/usr/bin/env python3
"""Export data to JSON file."""
import json
import sys
from pathlib import Path


def export_data(output_file):
    """Export user data to JSON file.
    
    Args:
        output_file: Path to output JSON file
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
    output_path.parent.mkdir(parents=True, exist_ok=True)
    
    with open(output_path, 'w') as f:
        json.dump(data, f, indent=2)
    
    print(f"✓ Exported {data['total_count']} users to {output_file}")
    return output_path


def main():
    if len(sys.argv) < 2:
        print("Usage: export.py <output_file>")
        print("Example: export.py data/users.json")
        sys.exit(1)
    
    output_file = sys.argv[1]
    export_data(output_file)


if __name__ == "__main__":
    main()
