#!/usr/bin/env python3
"""
对比 db-explorer 当前版本与基线版本（git HEAD）的正负样例表现，并输出评分。

默认运行 SQLite 测试。
如果设置了以下环境变量，还会运行真实 MySQL / PostgreSQL 测试：
- DBX_MYSQL_URL
- DBX_POSTGRES_URL
"""

from __future__ import annotations

import json
import os
import sqlite3
import stat
import subprocess
import sys
import tempfile
from dataclasses import dataclass
from pathlib import Path
from typing import Callable
from urllib.parse import urlparse


ROOT = Path(__file__).resolve().parents[3]
SCRIPT_REL = "claude-code/db-explorer/scripts/db_query.py"
CURRENT_SCRIPT = ROOT / SCRIPT_REL
GRADES_PATH = ROOT / "claude-code/db-explorer/evals/grades.json"
MYSQL_TABLES = ("dbx_users", "dbx_orders")
POSTGRES_TABLES = ("dbx_users", "dbx_orders")


@dataclass
class RunResult:
    returncode: int
    stdout: str
    stderr: str
    cmd: list[str]

    @property
    def combined(self) -> str:
        return (self.stdout + "\n" + self.stderr).strip()


@dataclass
class CaseResult:
    id: int
    name: str
    category: str
    result: str
    justification: str


def extract_baseline_script(target: Path) -> None:
    proc = subprocess.run(
        ["git", "show", f"HEAD:{SCRIPT_REL}"],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=True,
    )
    target.write_text(proc.stdout, encoding="utf-8")
    target.chmod(target.stat().st_mode | stat.S_IXUSR)


def create_sqlite_test_db(db_path: Path) -> None:
    conn = sqlite3.connect(db_path)
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


def parse_db_url(url: str) -> dict[str, str | int | None]:
    parsed = urlparse(url)
    return {
        "scheme": parsed.scheme,
        "host": parsed.hostname or "127.0.0.1",
        "port": parsed.port,
        "user": parsed.username,
        "password": parsed.password,
        "database": parsed.path.lstrip("/"),
    }


def setup_postgres(url: str) -> None:
    import psycopg2

    params = parse_db_url(url)
    conn = psycopg2.connect(url)
    conn.autocommit = True
    try:
        with conn.cursor() as cur:
            for table in reversed(POSTGRES_TABLES):
                cur.execute(f'DROP TABLE IF EXISTS "{table}" CASCADE')

            cur.execute(
                """
                CREATE TABLE "dbx_users" (
                    id INTEGER PRIMARY KEY,
                    name TEXT NOT NULL,
                    email TEXT NOT NULL UNIQUE,
                    age INTEGER,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                )
                """
            )
            cur.execute(
                """
                CREATE TABLE "dbx_orders" (
                    id INTEGER PRIMARY KEY,
                    user_id INTEGER NOT NULL REFERENCES "dbx_users"(id),
                    amount NUMERIC(10, 2) NOT NULL,
                    status TEXT NOT NULL,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                )
                """
            )
            cur.execute('CREATE INDEX "idx_dbx_orders_user" ON "dbx_orders"(user_id)')
            cur.execute(
                """
                INSERT INTO "dbx_users" (id, name, email, age) VALUES
                    (1, 'Alice', 'alice@example.com', 30),
                    (2, 'Bob', 'bob@example.com', 28)
                """
            )
            cur.execute(
                """
                INSERT INTO "dbx_orders" (id, user_id, amount, status) VALUES
                    (1, 1, 100.00, 'paid'),
                    (2, 1, 45.49, 'paid'),
                    (3, 2, 120.00, 'paid')
                """
            )
    finally:
        conn.close()


def setup_mysql(url: str) -> None:
    import mysql.connector

    params = parse_db_url(url)
    conn = mysql.connector.connect(
        host=params["host"],
        port=params["port"] or 3306,
        user=params["user"],
        password=params["password"],
        database=params["database"],
    )
    conn.autocommit = True
    try:
        cur = conn.cursor()
        for table in reversed(MYSQL_TABLES):
            cur.execute(f"DROP TABLE IF EXISTS `{table}`")

        cur.execute(
            """
            CREATE TABLE `dbx_users` (
                id INT PRIMARY KEY,
                name VARCHAR(255) NOT NULL,
                email VARCHAR(255) NOT NULL UNIQUE,
                age INT,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
            """
        )
        cur.execute(
            """
            CREATE TABLE `dbx_orders` (
                id INT PRIMARY KEY,
                user_id INT NOT NULL,
                amount DECIMAL(10, 2) NOT NULL,
                status VARCHAR(64) NOT NULL,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                CONSTRAINT fk_dbx_orders_user
                    FOREIGN KEY (user_id) REFERENCES `dbx_users`(id)
            )
            """
        )
        cur.execute("CREATE INDEX `idx_dbx_orders_user` ON `dbx_orders`(user_id)")
        cur.execute(
            """
            INSERT INTO `dbx_users` (id, name, email, age) VALUES
                (1, 'Alice', 'alice@example.com', 30),
                (2, 'Bob', 'bob@example.com', 28)
            """
        )
        cur.execute(
            """
            INSERT INTO `dbx_orders` (id, user_id, amount, status) VALUES
                (1, 1, 100.00, 'paid'),
                (2, 1, 45.49, 'paid'),
                (3, 2, 120.00, 'paid')
            """
        )
    finally:
        conn.close()


def setup_live_databases() -> list[str]:
    enabled: list[str] = []
    postgres_url = os.environ.get("DBX_POSTGRES_URL")
    mysql_url = os.environ.get("DBX_MYSQL_URL")

    if postgres_url:
        setup_postgres(postgres_url)
        enabled.append("postgres")
    if mysql_url:
        setup_mysql(mysql_url)
        enabled.append("mysql")
    return enabled


def run(script: Path, *args: str, extra_env: dict[str, str] | None = None) -> RunResult:
    env = os.environ.copy()
    if extra_env:
        env.update(extra_env)

    cmd = [sys.executable, str(script), *args]
    proc = subprocess.run(
        cmd,
        cwd=ROOT,
        text=True,
        capture_output=True,
        env=env,
    )
    return RunResult(
        returncode=proc.returncode,
        stdout=proc.stdout.strip(),
        stderr=proc.stderr.strip(),
        cmd=cmd,
    )


def contains_all(text: str, parts: list[str]) -> bool:
    return all(part in text for part in parts)


def score_results(results: list[CaseResult]) -> float:
    if not results:
        return 0.0
    passed = sum(1 for item in results if item.result == "PASS")
    return round(passed / len(results), 2)


def make_case(
    case_id: int,
    name: str,
    category: str,
    ok: bool,
    success_reason: str,
    failure_detail: str,
) -> CaseResult:
    return CaseResult(
        case_id,
        name,
        category,
        "PASS" if ok else "FAIL",
        success_reason if ok else failure_detail,
    )


def case_sqlite_tables_and_schema(script: Path, sqlite_url: str) -> CaseResult:
    tables = run(script, "--db-type", "sqlite", "--url", sqlite_url, "tables")
    schema = run(script, "--db-type", "sqlite", "--url", sqlite_url, "schema", "users")
    ok = (
        tables.returncode == 0
        and contains_all(tables.stdout, ["users", "orders"])
        and schema.returncode == 0
        and contains_all(schema.stdout, ["id", "name", "email", "age", "created_at"])
    )
    return make_case(
        1,
        "sqlite tables + schema",
        "positive",
        ok,
        "SQLite 的 tables 和 schema 都成功，输出包含 users/orders 与 users 的核心字段。",
        f"tables={tables.combined}\n\nschema={schema.combined}",
    )


def case_sqlite_data_and_aggregate(script: Path, sqlite_url: str) -> CaseResult:
    data = run(
        script,
        "--db-type",
        "sqlite",
        "--url",
        sqlite_url,
        "--format",
        "markdown",
        "data",
        "orders",
        "--limit",
        "3",
    )
    query = run(
        script,
        "--db-type",
        "sqlite",
        "--url",
        sqlite_url,
        "--format",
        "markdown",
        "query",
        (
            "SELECT u.name, ROUND(SUM(o.amount), 2) AS total_amount "
            "FROM users u JOIN orders o ON u.id = o.user_id "
            "GROUP BY u.name ORDER BY u.name"
        ),
    )
    ok = (
        data.returncode == 0
        and "(3 行)" in data.stdout
        and query.returncode == 0
        and contains_all(query.stdout, ["Alice", "145.49", "Bob", "120.0"])
    )
    return make_case(
        2,
        "sqlite sample data + aggregate query",
        "positive",
        ok,
        "SQLite 的采样与聚合查询都成功，金额结果符合预期。",
        f"data={data.combined}\n\nquery={query.combined}",
    )


def case_sqlite_schema_json(script: Path, sqlite_url: str) -> CaseResult:
    result = run(
        script,
        "--db-type",
        "sqlite",
        "--url",
        sqlite_url,
        "--format",
        "json",
        "schema",
        "users",
    )
    ok = False
    detail = result.combined
    if result.returncode == 0:
        try:
            parsed = json.loads(result.stdout)
            ok = any(item.get("column_name") == "email" for item in parsed)
        except json.JSONDecodeError as exc:
            detail = f"JSON 解析失败: {exc}\nstdout={result.stdout}\nstderr={result.stderr}"
    return make_case(
        3,
        "sqlite schema json is parseable",
        "positive",
        ok,
        "SQLite 的 schema JSON 可解析，且索引信息未污染 stdout。",
        detail,
    )


def case_sqlite_url_env(script: Path, sqlite_url: str) -> CaseResult:
    result = run(
        script,
        "--db-type",
        "sqlite",
        "--url-env",
        "TEST_DB_URL",
        "tables",
        extra_env={"TEST_DB_URL": sqlite_url},
    )
    ok = result.returncode == 0 and contains_all(result.stdout, ["users", "orders"])
    return make_case(
        4,
        "sqlite read url from named env var",
        "positive",
        ok,
        "SQLite 支持 --url-env，从指定环境变量读取连接并成功列出表。",
        result.combined,
    )


def case_sqlite_reject_vacuum(script: Path, sqlite_url: str) -> CaseResult:
    result = run(
        script,
        "--db-type",
        "sqlite",
        "--url",
        sqlite_url,
        "query",
        "VACUUM",
    )
    ok = result.returncode != 0 and ("拒绝" in result.combined or "只读" in result.combined)
    return make_case(
        5,
        "sqlite reject VACUUM",
        "negative",
        ok,
        "SQLite 场景下正确拒绝 VACUUM 这类非只读语句。",
        result.combined,
    )


def case_sqlite_reject_write_pragma(script: Path, sqlite_url: str) -> CaseResult:
    result = run(
        script,
        "--db-type",
        "sqlite",
        "--url",
        sqlite_url,
        "query",
        "PRAGMA journal_mode=WAL",
    )
    ok = result.returncode != 0 and "元数据 PRAGMA" in result.combined
    return make_case(
        6,
        "sqlite reject write-ish PRAGMA",
        "negative",
        ok,
        "SQLite 场景下正确拒绝会改数据库状态的 PRAGMA。",
        result.combined,
    )


def case_sqlite_reject_multi_statement(script: Path, sqlite_url: str) -> CaseResult:
    result = run(
        script,
        "--db-type",
        "sqlite",
        "--url",
        sqlite_url,
        "query",
        "SELECT 1; DROP TABLE users",
    )
    ok = result.returncode != 0 and "单条只读 SQL" in result.combined
    return make_case(
        7,
        "sqlite reject multi-statement SQL",
        "negative",
        ok,
        "SQLite 场景下正确拒绝多语句 SQL。",
        result.combined,
    )


def case_postgres_tables_and_schema(script: Path, postgres_url: str) -> CaseResult:
    tables = run(script, "--db-type", "postgres", "--url", postgres_url, "tables")
    schema = run(script, "--db-type", "postgres", "--url", postgres_url, "schema", "dbx_users")
    ok = (
        tables.returncode == 0
        and contains_all(tables.stdout, ["dbx_users", "dbx_orders"])
        and schema.returncode == 0
        and contains_all(schema.stdout, ["id", "name", "email", "age", "created_at"])
    )
    return make_case(
        8,
        "postgres tables + schema",
        "positive",
        ok,
        "PostgreSQL 的 tables 和 schema 都成功，输出包含 dbx_users/dbx_orders 与核心字段。",
        f"tables={tables.combined}\n\nschema={schema.combined}",
    )


def case_postgres_data_and_aggregate(script: Path, postgres_url: str) -> CaseResult:
    data = run(
        script,
        "--db-type",
        "postgres",
        "--url",
        postgres_url,
        "--format",
        "markdown",
        "data",
        "dbx_orders",
        "--limit",
        "3",
    )
    query = run(
        script,
        "--db-type",
        "postgres",
        "--url",
        postgres_url,
        "--format",
        "markdown",
        "query",
        (
            "SELECT u.name, ROUND(SUM(o.amount), 2) AS total_amount "
            "FROM dbx_users u JOIN dbx_orders o ON u.id = o.user_id "
            "GROUP BY u.name ORDER BY u.name"
        ),
    )
    ok = (
        data.returncode == 0
        and "(3 行)" in data.stdout
        and query.returncode == 0
        and contains_all(query.stdout, ["Alice", "145.49", "Bob", "120.00"])
    )
    return make_case(
        9,
        "postgres sample data + aggregate query",
        "positive",
        ok,
        "PostgreSQL 的采样与聚合查询都成功，金额结果符合预期。",
        f"data={data.combined}\n\nquery={query.combined}",
    )


def case_postgres_schema_json(script: Path, postgres_url: str) -> CaseResult:
    result = run(
        script,
        "--db-type",
        "postgres",
        "--url",
        postgres_url,
        "--format",
        "json",
        "schema",
        "dbx_users",
    )
    ok = False
    detail = result.combined
    if result.returncode == 0:
        try:
            parsed = json.loads(result.stdout)
            ok = any(item.get("column_name") == "email" for item in parsed)
        except json.JSONDecodeError as exc:
            detail = f"JSON 解析失败: {exc}\nstdout={result.stdout}\nstderr={result.stderr}"
    return make_case(
        10,
        "postgres schema json is parseable",
        "positive",
        ok,
        "PostgreSQL 的 schema JSON 可解析，且索引/外键信息未污染 stdout。",
        detail,
    )


def case_mysql_tables_and_schema(script: Path, mysql_url: str) -> CaseResult:
    tables = run(script, "--db-type", "mysql", "--url", mysql_url, "tables")
    schema = run(script, "--db-type", "mysql", "--url", mysql_url, "schema", "dbx_users")
    ok = (
        tables.returncode == 0
        and contains_all(tables.stdout, ["dbx_users", "dbx_orders"])
        and schema.returncode == 0
        and contains_all(schema.stdout, ["id", "name", "email", "age", "created_at"])
    )
    return make_case(
        11,
        "mysql tables + schema",
        "positive",
        ok,
        "MySQL 的 tables 和 schema 都成功，输出包含 dbx_users/dbx_orders 与核心字段。",
        f"tables={tables.combined}\n\nschema={schema.combined}",
    )


def case_mysql_data_and_aggregate(script: Path, mysql_url: str) -> CaseResult:
    data = run(
        script,
        "--db-type",
        "mysql",
        "--url",
        mysql_url,
        "--format",
        "markdown",
        "data",
        "dbx_orders",
        "--limit",
        "3",
    )
    query = run(
        script,
        "--db-type",
        "mysql",
        "--url",
        mysql_url,
        "--format",
        "markdown",
        "query",
        (
            "SELECT u.name, ROUND(SUM(o.amount), 2) AS total_amount "
            "FROM dbx_users u JOIN dbx_orders o ON u.id = o.user_id "
            "GROUP BY u.name ORDER BY u.name"
        ),
    )
    ok = (
        data.returncode == 0
        and "(3 行)" in data.stdout
        and query.returncode == 0
        and contains_all(query.stdout, ["Alice", "145.49", "Bob", "120.00"])
    )
    return make_case(
        12,
        "mysql sample data + aggregate query",
        "positive",
        ok,
        "MySQL 的采样与聚合查询都成功，金额结果符合预期。",
        f"data={data.combined}\n\nquery={query.combined}",
    )


def case_mysql_schema_json(script: Path, mysql_url: str) -> CaseResult:
    result = run(
        script,
        "--db-type",
        "mysql",
        "--url",
        mysql_url,
        "--format",
        "json",
        "schema",
        "dbx_users",
    )
    ok = False
    detail = result.combined
    if result.returncode == 0:
        try:
            parsed = json.loads(result.stdout)
            ok = any(item.get("column_name") == "email" for item in parsed)
        except json.JSONDecodeError as exc:
            detail = f"JSON 解析失败: {exc}\nstdout={result.stdout}\nstderr={result.stderr}"
    return make_case(
        13,
        "mysql schema json is parseable",
        "positive",
        ok,
        "MySQL 的 schema JSON 可解析，且索引信息未污染 stdout。",
        detail,
    )


def run_sqlite_suite(script: Path, sqlite_url: str) -> list[CaseResult]:
    cases: list[Callable[[Path, str], CaseResult]] = [
        case_sqlite_tables_and_schema,
        case_sqlite_data_and_aggregate,
        case_sqlite_schema_json,
        case_sqlite_url_env,
        case_sqlite_reject_vacuum,
        case_sqlite_reject_write_pragma,
        case_sqlite_reject_multi_statement,
    ]
    return [case(script, sqlite_url) for case in cases]


def run_postgres_suite(script: Path, postgres_url: str) -> list[CaseResult]:
    cases: list[Callable[[Path, str], CaseResult]] = [
        case_postgres_tables_and_schema,
        case_postgres_data_and_aggregate,
        case_postgres_schema_json,
    ]
    return [case(script, postgres_url) for case in cases]


def run_mysql_suite(script: Path, mysql_url: str) -> list[CaseResult]:
    cases: list[Callable[[Path, str], CaseResult]] = [
        case_mysql_tables_and_schema,
        case_mysql_data_and_aggregate,
        case_mysql_schema_json,
    ]
    return [case(script, mysql_url) for case in cases]


def format_results(results: list[CaseResult]) -> list[dict]:
    return [
        {
            "id": item.id,
            "name": item.name,
            "category": item.category,
            "result": item.result,
            "justification": item.justification,
        }
        for item in results
    ]


def main() -> int:
    with tempfile.TemporaryDirectory(prefix="db-explorer-eval-") as tmp_dir:
        tmp_path = Path(tmp_dir)
        sqlite_db_path = tmp_path / "test_db_explorer.db"
        baseline_script = tmp_path / "db_query_baseline.py"
        sqlite_url = str(sqlite_db_path)

        create_sqlite_test_db(sqlite_db_path)
        extract_baseline_script(baseline_script)

        enabled_live = setup_live_databases()

        current_results = run_sqlite_suite(CURRENT_SCRIPT, sqlite_url)
        baseline_results = run_sqlite_suite(baseline_script, sqlite_url)

        postgres_url = os.environ.get("DBX_POSTGRES_URL")
        if postgres_url:
            current_results.extend(run_postgres_suite(CURRENT_SCRIPT, postgres_url))
            baseline_results.extend(run_postgres_suite(baseline_script, postgres_url))

        mysql_url = os.environ.get("DBX_MYSQL_URL")
        if mysql_url:
            current_results.extend(run_mysql_suite(CURRENT_SCRIPT, mysql_url))
            baseline_results.extend(run_mysql_suite(baseline_script, mysql_url))

    current_score = score_results(current_results)
    baseline_score = score_results(baseline_results)
    grades = {
        "eval_name": "db-explorer-comparison",
        "live_databases": enabled_live,
        "with_skill": {
            "score": current_score,
            "assertions": format_results(current_results),
        },
        "without_skill": {
            "score": baseline_score,
            "assertions": format_results(baseline_results),
        },
        "delta": round(current_score - baseline_score, 2),
        "summary": (
            f"Current version scored {current_score:.2f}; baseline scored {baseline_score:.2f}. "
            f"Delta = {current_score - baseline_score:.2f}. "
            f"Live DBs: {', '.join(enabled_live) if enabled_live else 'sqlite-only'}."
        ),
    }

    GRADES_PATH.write_text(json.dumps(grades, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(grades, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
