#!/usr/bin/env python3
"""Enumerate Maven modules for the springcloud-init skill.

This script only scans pom.xml files and emits structured JSON. It does not
classify modules as microservices or shared modules, and it never writes files.
"""

from __future__ import annotations

import argparse
import json
import os
import sys
import xml.etree.ElementTree as ET
from pathlib import Path
from typing import Any

EXCLUDED_DIRS = {"target", ".git", "node_modules", "build", ".idea"}


def strip_namespace(tag: str) -> str:
    return tag.rsplit("}", 1)[-1] if "}" in tag else tag


def children_by_name(element: ET.Element, name: str) -> list[ET.Element]:
    return [child for child in list(element) if strip_namespace(child.tag) == name]


def child_text(element: ET.Element, name: str) -> str | None:
    for child in list(element):
        if strip_namespace(child.tag) == name:
            text = child.text.strip() if child.text else ""
            return text or None
    return None


def parse_parent(root: ET.Element) -> dict[str, str | None] | None:
    parents = children_by_name(root, "parent")
    if not parents:
        return None
    parent = parents[0]
    return {
        "groupId": child_text(parent, "groupId"),
        "artifactId": child_text(parent, "artifactId"),
        "version": child_text(parent, "version"),
        "relativePath": child_text(parent, "relativePath"),
    }


def parse_modules(root: ET.Element) -> list[str]:
    modules_nodes = children_by_name(root, "modules")
    if not modules_nodes:
        return []
    modules: list[str] = []
    for module in children_by_name(modules_nodes[0], "module"):
        if module.text and module.text.strip():
            modules.append(module.text.strip())
    return modules


def parse_dependencies(root: ET.Element) -> list[dict[str, str | None]]:
    deps: list[dict[str, str | None]] = []
    for deps_node in children_by_name(root, "dependencies"):
        for dep in children_by_name(deps_node, "dependency"):
            artifact_id = child_text(dep, "artifactId")
            group_id = child_text(dep, "groupId")
            if not artifact_id and not group_id:
                continue
            deps.append(
                {
                    "groupId": group_id,
                    "artifactId": artifact_id,
                    "version": child_text(dep, "version"),
                    "scope": child_text(dep, "scope"),
                    "type": child_text(dep, "type"),
                    "optional": child_text(dep, "optional"),
                }
            )
    return deps


def parse_pom(path: Path, project_root: Path) -> dict[str, Any]:
    try:
        tree = ET.parse(path)
    except ET.ParseError as exc:
        return {
            "path": str(path.relative_to(project_root)),
            "parseError": str(exc),
        }

    root = tree.getroot()
    parent = parse_parent(root)
    group_id = child_text(root, "groupId") or (parent or {}).get("groupId")
    artifact_id = child_text(root, "artifactId")
    version = child_text(root, "version") or (parent or {}).get("version")
    packaging = child_text(root, "packaging") or "jar"

    return {
        "path": str(path.relative_to(project_root)),
        "moduleDir": str(path.parent.relative_to(project_root)) or ".",
        "groupId": group_id,
        "artifactId": artifact_id,
        "version": version,
        "packaging": packaging,
        "parent": parent,
        "modules": parse_modules(root),
        "dependencies": parse_dependencies(root),
    }


def iter_poms(project_root: Path) -> list[Path]:
    poms: list[Path] = []
    for current_root, dirs, files in os.walk(project_root):
        dirs[:] = [d for d in dirs if d not in EXCLUDED_DIRS]
        if "pom.xml" in files:
            poms.append(Path(current_root) / "pom.xml")
    return sorted(poms)


def add_internal_dependencies(modules: list[dict[str, Any]]) -> None:
    by_artifact: dict[str, list[dict[str, Any]]] = {}
    by_ga: dict[tuple[str | None, str | None], list[dict[str, Any]]] = {}

    for module in modules:
        artifact_id = module.get("artifactId")
        if artifact_id:
            by_artifact.setdefault(artifact_id, []).append(module)
        by_ga.setdefault((module.get("groupId"), artifact_id), []).append(module)

    for module in modules:
        internal: list[dict[str, str | None]] = []
        for dep in module.get("dependencies", []):
            dep_artifact = dep.get("artifactId")
            dep_group = dep.get("groupId")
            matches = by_ga.get((dep_group, dep_artifact)) or by_artifact.get(dep_artifact or "") or []
            for match in matches:
                if match is module:
                    continue
                internal.append(
                    {
                        "groupId": match.get("groupId"),
                        "artifactId": match.get("artifactId"),
                        "path": match.get("moduleDir"),
                    }
                )
        module["internalDependencies"] = internal


def main() -> int:
    parser = argparse.ArgumentParser(description="Enumerate Maven modules as JSON")
    parser.add_argument("root", nargs="?", default=".", help="Project root directory")
    parser.add_argument("--pretty", action="store_true", help="Pretty-print JSON")
    args = parser.parse_args()

    project_root = Path(args.root).resolve()
    if not project_root.exists():
        print(f"Root does not exist: {project_root}", file=sys.stderr)
        return 2
    if not project_root.is_dir():
        print(f"Root is not a directory: {project_root}", file=sys.stderr)
        return 2

    modules = [parse_pom(path, project_root) for path in iter_poms(project_root)]
    add_internal_dependencies(modules)

    output = {
        "root": str(project_root),
        "excludedDirs": sorted(EXCLUDED_DIRS),
        "moduleCount": len(modules),
        "modules": modules,
    }
    indent = 2 if args.pretty else None
    print(json.dumps(output, ensure_ascii=False, indent=indent))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
