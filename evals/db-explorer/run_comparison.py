#!/usr/bin/env python3
"""Regression checks for the Go db-explorer binary.

Runs deterministic SQLite cases against the current Go implementation and writes
`grades.json`. PostgreSQL/MySQL live behavior is covered by env-gated Go tests.
"""

from __future__ import annotations

import json
import sqlite3
import subprocess
import sys
import tempfile
from dataclasses import asdict, dataclass
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
GO_DIR = ROOT / "db-explorer-go"
GRADES_PATH = ROOT / "evals" / "db-explorer" / "grades.json"


@dataclass
class CaseResult:
    id: int
    name: str
    category: str
    result: str
    justification: str


def run(cmd: list[str], env: dict[str, str] | None = None, cwd: Path = ROOT) -> subprocess.CompletedProcess[str]:
    return subprocess.run(cmd, cwd=cwd, text=True, capture_output=True, env=env)


def build_binary(tmp: Path) -> Path:
    binary = tmp / ("db-explorer.exe" if sys.platform == "win32" else "db-explorer")
    proc = run(["go", "build", "-o", str(binary), "./cmd/db-explorer"], env=None, cwd=GO_DIR)
    if proc.returncode != 0:
        raise RuntimeError(proc.stdout + proc.stderr)
    return binary


def create_db(path: Path) -> None:
    conn = sqlite3.connect(path)
    cur = conn.cursor()
    cur.executescript(
        """
        PRAGMA foreign_keys = ON;
        CREATE TABLE users (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL,
            email TEXT NOT NULL UNIQUE,
            age INTEGER,
            created_at TEXT DEFAULT CURRENT_TIMESTAMP
        );
        CREATE TABLE orders (
            id INTEGER PRIMARY KEY,
            user_id INTEGER NOT NULL,
            amount REAL NOT NULL,
            status TEXT NOT NULL,
            created_at TEXT DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
        CREATE INDEX idx_orders_user ON orders(user_id);
        CREATE VIEW active_users AS SELECT id, name FROM users;
        INSERT INTO users (name, email, age) VALUES
            ('Alice', 'alice@example.com', 30),
            ('Bob', 'bob@example.com', 28);
        INSERT INTO orders (user_id, amount, status) VALUES
            (1, 100.00, 'paid'),
            (1, 45.49, 'paid'),
            (2, 120.00, 'paid');
        """
    )
    conn.commit()
    conn.close()


def invoke(binary: Path, db_path: Path, *args: str, extra_env: dict[str, str] | None = None) -> dict:
    cmd = [str(binary), *args, "--db", "sqlite", "--url", str(db_path)]
    proc = run(cmd, env=extra_env)
    try:
        obj = json.loads(proc.stdout)
    except json.JSONDecodeError as exc:
        raise AssertionError(f"stdout is not JSON: {exc}\nstdout={proc.stdout}\nstderr={proc.stderr}")
    return obj


def ok(obj: dict) -> bool:
    return obj.get("schema_version") == "1" and obj.get("ok") is True


def fail(obj: dict, code: str) -> bool:
    return obj.get("schema_version") == "1" and obj.get("ok") is False and obj.get("error", {}).get("code") == code


def make(case_id: int, name: str, category: str, passed: bool, good: str, bad: str) -> CaseResult:
    return CaseResult(case_id, name, category, "PASS" if passed else "FAIL", good if passed else bad)


def run_cases(binary: Path, db_path: Path) -> list[CaseResult]:
    results: list[CaseResult] = []

    schemas = invoke(binary, db_path, "schemas")
    tables = invoke(binary, db_path, "tables")
    views = invoke(binary, db_path, "views")
    users = invoke(binary, db_path, "schema", "users")
    table_names = {item["name"] for item in tables.get("data", {}).get("tables", [])}
    view_names = {item["name"] for item in views.get("data", {}).get("views", [])}
    user_cols = {item["name"] for item in users.get("data", {}).get("columns", [])}
    results.append(make(
        1,
        "sqlite json schemas/tables/views/schema",
        "positive",
        all(map(ok, [schemas, tables, views, users]))
        and {"users", "orders"}.issubset(table_names)
        and "active_users" in view_names
        and {"id", "name", "email", "age", "created_at"}.issubset(user_cols),
        "JSON metadata commands returned expected SQLite objects.",
        f"schemas={schemas} tables={tables} views={views} users={users}",
    ))

    sample = invoke(binary, db_path, "data", "orders", "--limit", "2")
    aggregate_sql = "SELECT u.name, SUM(o.amount) AS total FROM users u JOIN orders o ON o.user_id = u.id GROUP BY u.name ORDER BY u.name"
    aggregate = invoke(binary, db_path, "query", aggregate_sql)
    totals = {row[0]: round(float(row[1]), 2) for row in aggregate.get("data", {}).get("rows", [])}
    results.append(make(
        2,
        "sqlite sample data + aggregate query",
        "positive",
        ok(sample) and len(sample.get("data", {}).get("rows", [])) <= 2 and ok(aggregate) and totals == {"Alice": 145.49, "Bob": 120.0},
        "Data sample and aggregate query returned expected bounded JSON results.",
        f"sample={sample} aggregate={aggregate}",
    ))

    orders = invoke(binary, db_path, "schema", "orders")
    indexes = {item["name"] for item in orders.get("data", {}).get("indexes", [])}
    fks = orders.get("data", {}).get("foreign_keys", [])
    results.append(make(
        3,
        "sqlite indexes and foreign keys",
        "positive",
        ok(orders) and "idx_orders_user" in indexes and any(fk.get("referenced_table") == "users" for fk in fks),
        "Schema JSON includes index and foreign-key metadata.",
        f"orders={orders}",
    ))

    env = dict(**__import__("os").environ, TEST_DB_URL=str(db_path))
    env_proc = run([str(binary), "tables", "--db", "sqlite", "--url-env", "TEST_DB_URL"], env=env)
    env_obj = json.loads(env_proc.stdout)
    env_names = {item["name"] for item in env_obj.get("data", {}).get("tables", [])}
    results.append(make(
        4,
        "sqlite read url from named env var",
        "positive",
        ok(env_obj) and {"users", "orders"}.issubset(env_names),
        "--url-env resolved TEST_DB_URL and listed tables.",
        f"env_obj={env_obj} stderr={env_proc.stderr}",
    ))

    vacuum = invoke(binary, db_path, "query", "VACUUM")
    results.append(make(5, "reject VACUUM", "negative", fail(vacuum, "SQL_NOT_READONLY"), "VACUUM rejected.", f"vacuum={vacuum}"))

    pragma = invoke(binary, db_path, "query", "PRAGMA journal_mode=WAL")
    results.append(make(6, "reject unsafe PRAGMA", "negative", fail(pragma, "UNSAFE_PRAGMA"), "Unsafe PRAGMA rejected.", f"pragma={pragma}"))

    multi = invoke(binary, db_path, "query", "SELECT 1; DROP TABLE users")
    results.append(make(7, "reject multi-statement SQL", "negative", fail(multi, "SQL_MULTIPLE_STATEMENTS"), "Multi-statement SQL rejected.", f"multi={multi}"))

    return results


def score(results: list[CaseResult]) -> float:
    return round(sum(1 for r in results if r.result == "PASS") / len(results), 2)


def main() -> None:
    with tempfile.TemporaryDirectory(prefix="db-explorer-eval-") as td:
        tmp = Path(td)
        binary = build_binary(tmp)
        db_path = tmp / "test.db"
        create_db(db_path)
        results = run_cases(binary, db_path)

    payload = {
        "eval_name": "db-explorer-go-json",
        "with_skill": {
            "score": score(results),
            "assertions": [asdict(r) for r in results],
        },
    }
    GRADES_PATH.write_text(json.dumps(payload, ensure_ascii=False, indent=2), encoding="utf-8")
    print(json.dumps(payload, ensure_ascii=False, indent=2))
    if payload["with_skill"]["score"] < 1.0:
        sys.exit(1)


if __name__ == "__main__":
    main()
