#!/usr/bin/env python3
import argparse
from pathlib import Path


def collect_files(input_dir: Path):
    return sorted(p for p in input_dir.glob('*.txt') if p.is_file())


def export_files(input_dir: Path, output_dir: Path, dry_run: bool = False):
    if not dry_run:
        output_dir.mkdir(parents=True, exist_ok=True)
    exported = []
    for source in collect_files(input_dir):
        target = output_dir / source.name
        if not dry_run:
            target.write_text(source.read_text(encoding='utf-8'), encoding='utf-8')
        exported.append(target)
    return exported


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--input-dir', default='input')
    parser.add_argument('--output-dir', default='out')
    parser.add_argument(
        '--dry-run',
        action='store_true',
        help='show what would be exported without writing files',
    )
    args = parser.parse_args()

    exported = export_files(Path(args.input_dir), Path(args.output_dir), dry_run=args.dry_run)
    for path in exported:
        if args.dry_run:
            print(f'would export {path}')
        else:
            print(f'exported {path}')


if __name__ == '__main__':
    main()
