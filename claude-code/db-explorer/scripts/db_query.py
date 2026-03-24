#!/usr/bin/env python3
"""
db_query.py - 轻量级数据库查询工具
支持 PostgreSQL、MySQL、SQLite 的表结构和数据查询。
所有操作只读，拒绝任何写入语句。
"""

import argparse
import json
import os
import re
import sqlite3
import sys
import importlib
from contextlib import contextmanager


# ---------------------------------------------------------------------------
# 危险语句检测
# ---------------------------------------------------------------------------

DANGEROUS_KEYWORDS = re.compile(
    r"\b(INSERT|UPDATE|DELETE|DROP|ALTER|TRUNCATE|CREATE|GRANT|REVOKE|EXEC|MERGE|CALL|VACUUM|REINDEX|ATTACH|DETACH)\b",
    re.IGNORECASE,
)

ALLOWED_PREFIXES = ("SELECT", "WITH", "SHOW", "DESCRIBE", "DESC", "EXPLAIN", "PRAGMA")
SAFE_PRAGMA_NAMES = {"table_info", "index_list", "index_info", "foreign_key_list"}


def strip_sql_comments(sql: str) -> str:
    """去掉 SQL 注释，便于做简单只读校验"""
    sql = re.sub(r"/\*.*?\*/", " ", sql, flags=re.DOTALL)
    sql = re.sub(r"--.*?$", " ", sql, flags=re.MULTILINE)
    return sql.strip()


def has_multiple_statements(sql: str) -> bool:
    """检测是否包含多条语句（允许末尾单个分号）"""
    in_single = False
    in_double = False
    in_backtick = False

    for i, ch in enumerate(sql):
        prev = sql[i - 1] if i > 0 else ""

        if ch == "'" and not in_double and not in_backtick and prev != "\\":
            in_single = not in_single
        elif ch == '"' and not in_single and not in_backtick and prev != "\\":
            in_double = not in_double
        elif ch == "`" and not in_single and not in_double and prev != "\\":
            in_backtick = not in_backtick
        elif ch == ";" and not in_single and not in_double and not in_backtick:
            if sql[i + 1 :].strip():
                return True

    return False


def validate_readonly(sql: str) -> None:
    """拒绝任何非只读语句"""
    stripped = strip_sql_comments(sql).rstrip(";").strip()
    if not stripped:
        print("ERROR: SQL 语句不能为空", file=sys.stderr)
        sys.exit(1)

    if has_multiple_statements(stripped):
        print("ERROR: 只允许执行单条只读 SQL 语句", file=sys.stderr)
        sys.exit(1)

    if DANGEROUS_KEYWORDS.search(stripped):
        print(f"ERROR: 拒绝执行危险语句: {stripped[:80]}...", file=sys.stderr)
        sys.exit(1)

    first_token = re.match(r"^[A-Za-z]+", stripped)
    if not first_token or first_token.group(0).upper() not in ALLOWED_PREFIXES:
        print(
            "ERROR: 仅允许只读查询（SELECT/WITH/SHOW/DESCRIBE/EXPLAIN/PRAGMA）",
            file=sys.stderr,
        )
        sys.exit(1)

    if first_token.group(0).upper() == "PRAGMA":
        pragma_match = re.match(r"^PRAGMA\s+([A-Za-z_][A-Za-z0-9_]*)", stripped, re.IGNORECASE)
        pragma_name = pragma_match.group(1).lower() if pragma_match else None
        if pragma_name not in SAFE_PRAGMA_NAMES:
            print(
                "ERROR: 仅允许只读元数据 PRAGMA（table_info/index_list/index_info/foreign_key_list）",
                file=sys.stderr,
            )
            sys.exit(1)


def validate_table_name(name: str) -> str:
    """简单校验表名，防止注入"""
    if not re.match(r"^[a-zA-Z_][a-zA-Z0-9_]*$", name):
        print(f"ERROR: 非法表名: {name}", file=sys.stderr)
        sys.exit(1)
    return name


def quote_identifier(name: str, db_type: str) -> str:
    """按数据库类型安全包裹标识符"""
    if db_type == "mysql":
        return f"`{name}`"
    return f'"{name}"'


def positive_int(value: str) -> int:
    try:
        parsed = int(value)
    except ValueError:
        raise argparse.ArgumentTypeError(f"非法整数: {value}")
    if parsed <= 0:
        raise argparse.ArgumentTypeError("必须是大于 0 的整数")
    return parsed


# ---------------------------------------------------------------------------
# 数据库连接
# ---------------------------------------------------------------------------

def parse_url(url: str) -> dict:
    """从 URL 解析连接参数"""
    from urllib.parse import urlparse
    parsed = urlparse(url)
    return {
        "host": parsed.hostname or "localhost",
        "port": parsed.port,
        "user": parsed.username,
        "password": parsed.password,
        "database": parsed.path.lstrip("/"),
    }


@contextmanager
def connect_sqlite(url: str):
    path = url.replace("sqlite:///", "").replace("sqlite://", "")
    if not os.path.exists(path):
        print(f"ERROR: SQLite 文件不存在: {path}", file=sys.stderr)
        sys.exit(1)
    conn = sqlite3.connect(path)
    conn.row_factory = sqlite3.Row
    try:
        yield conn
    finally:
        conn.close()


@contextmanager
def connect_postgres(url: str):
    try:
        psycopg2 = importlib.import_module("psycopg2")
    except ImportError:
        print("ERROR: 需要安装 psycopg2。请运行: pip install psycopg2-binary", file=sys.stderr)
        sys.exit(1)
    conn = psycopg2.connect(url)
    conn.set_session(readonly=True, autocommit=True)
    try:
        yield conn
    finally:
        conn.close()


@contextmanager
def connect_mysql(url: str):
    try:
        mysql_connector = importlib.import_module("mysql.connector")
    except ImportError:
        print("ERROR: 需要安装 mysql-connector-python。请运行: pip install mysql-connector-python", file=sys.stderr)
        sys.exit(1)
    params = parse_url(url)
    conn = mysql_connector.connect(
        host=params["host"],
        port=params["port"] or 3306,
        user=params["user"],
        password=params["password"],
        database=params["database"],
    )
    try:
        yield conn
    finally:
        conn.close()


def get_connection(db_type: str, url: str):
    """根据数据库类型返回连接上下文管理器"""
    if db_type == "sqlite":
        return connect_sqlite(url)
    elif db_type == "postgres":
        return connect_postgres(url)
    elif db_type == "mysql":
        return connect_mysql(url)
    else:
        print(f"ERROR: 不支持的数据库类型: {db_type}", file=sys.stderr)
        sys.exit(1)


# ---------------------------------------------------------------------------
# 查询执行
# ---------------------------------------------------------------------------

def execute_query(conn, db_type: str, sql: str, timeout: int = 30) -> tuple:
    """执行查询，返回 (columns, rows)"""
    cursor = conn.cursor()

    # 设置超时
    if db_type == "postgres":
        cursor.execute(f"SET statement_timeout = '{timeout * 1000}'")
    elif db_type == "mysql":
        cursor.execute(f"SET SESSION MAX_EXECUTION_TIME = {timeout * 1000}")

    cursor.execute(sql)
    if cursor.description is None:
        return [], []

    columns = [desc[0] for desc in cursor.description]
    rows = cursor.fetchall()

    # sqlite3.Row -> tuple
    if db_type == "sqlite":
        rows = [tuple(row) for row in rows]

    return columns, rows


# ---------------------------------------------------------------------------
# 命令实现
# ---------------------------------------------------------------------------

def cmd_test(conn, db_type: str) -> None:
    """测试连接"""
    try:
        if db_type == "sqlite":
            execute_query(conn, db_type, "SELECT 1")
        elif db_type == "postgres":
            execute_query(conn, db_type, "SELECT 1")
        elif db_type == "mysql":
            execute_query(conn, db_type, "SELECT 1")
        print("OK: 连接成功")
    except Exception as e:
        print(f"ERROR: 连接失败: {e}", file=sys.stderr)
        sys.exit(1)


def cmd_tables(conn, db_type: str) -> tuple:
    """列出所有表"""
    if db_type == "sqlite":
        sql = "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name"
    elif db_type == "postgres":
        sql = """
            SELECT table_name
            FROM information_schema.tables
            WHERE table_schema = 'public'
            ORDER BY table_name
        """
    elif db_type == "mysql":
        sql = "SHOW TABLES"

    columns, rows = execute_query(conn, db_type, sql)

    # 获取每张表的行数
    result_rows = []
    for row in rows:
        table_name = row[0]
        try:
            quoted_table = quote_identifier(table_name, db_type)
            _, count_rows = execute_query(conn, db_type, f"SELECT COUNT(*) FROM {quoted_table}")
            count = count_rows[0][0]
        except Exception:
            count = "?"
        result_rows.append((table_name, count))

    return ["table_name", "row_count"], result_rows


def cmd_schema(conn, db_type: str, table_name: str) -> tuple:
    """查看表结构"""
    table_name = validate_table_name(table_name)
    quoted_table = quote_identifier(table_name, db_type)

    if db_type == "sqlite":
        sql = f"PRAGMA table_info({quoted_table})"
        columns, rows = execute_query(conn, db_type, sql)
        # PRAGMA 返回: cid, name, type, notnull, dflt_value, pk
        result_columns = ["column_name", "type", "nullable", "default", "primary_key"]
        result_rows = [
            (row[1], row[2], "NO" if row[3] else "YES", row[4] or "", "YES" if row[5] else "")
            for row in rows
        ]

        # 获取索引信息
        _, index_rows = execute_query(conn, db_type, f"PRAGMA index_list({quoted_table})")
        if index_rows:
            print("\n--- 索引 ---", file=sys.stderr)
            for idx_row in index_rows:
                idx_name = idx_row[1]
                unique = "UNIQUE" if idx_row[2] else ""
                _, idx_cols = execute_query(conn, db_type, f'PRAGMA index_info("{idx_name}")')
                col_names = ", ".join(r[2] for r in idx_cols)
                print(f"  {idx_name} ({col_names}) {unique}", file=sys.stderr)

        return result_columns, result_rows

    elif db_type == "postgres":
        sql = f"""
            SELECT
                c.column_name,
                c.data_type,
                c.is_nullable,
                c.column_default,
                CASE WHEN pk.column_name IS NOT NULL THEN 'YES' ELSE '' END as primary_key
            FROM information_schema.columns c
            LEFT JOIN (
                SELECT ku.column_name
                FROM information_schema.table_constraints tc
                JOIN information_schema.key_column_usage ku
                    ON tc.constraint_name = ku.constraint_name
                WHERE tc.table_name = '{table_name}'
                    AND tc.constraint_type = 'PRIMARY KEY'
            ) pk ON c.column_name = pk.column_name
            WHERE c.table_name = '{table_name}'
                AND c.table_schema = 'public'
            ORDER BY c.ordinal_position
        """
        columns, rows = execute_query(conn, db_type, sql)

        # 获取索引
        idx_sql = f"""
            SELECT indexname, indexdef
            FROM pg_indexes
            WHERE tablename = '{table_name}' AND schemaname = 'public'
        """
        _, idx_rows = execute_query(conn, db_type, idx_sql)
        if idx_rows:
            print("\n--- 索引 ---", file=sys.stderr)
            for idx_row in idx_rows:
                print(f"  {idx_row[0]}: {idx_row[1]}", file=sys.stderr)

        # 获取外键
        fk_sql = f"""
            SELECT
                kcu.column_name,
                ccu.table_name AS foreign_table,
                ccu.column_name AS foreign_column
            FROM information_schema.table_constraints tc
            JOIN information_schema.key_column_usage kcu
                ON tc.constraint_name = kcu.constraint_name
            JOIN information_schema.constraint_column_usage ccu
                ON tc.constraint_name = ccu.constraint_name
            WHERE tc.table_name = '{table_name}'
                AND tc.constraint_type = 'FOREIGN KEY'
        """
        _, fk_rows = execute_query(conn, db_type, fk_sql)
        if fk_rows:
            print("\n--- 外键 ---", file=sys.stderr)
            for fk_row in fk_rows:
                print(f"  {fk_row[0]} -> {fk_row[1]}.{fk_row[2]}", file=sys.stderr)

        return ["column_name", "type", "nullable", "default", "primary_key"], rows

    elif db_type == "mysql":
        sql = f"DESCRIBE `{table_name}`"
        columns, rows = execute_query(conn, db_type, sql)
        # DESCRIBE 返回: Field, Type, Null, Key, Default, Extra
        result_columns = ["column_name", "type", "nullable", "key", "default", "extra"]
        result_rows = [tuple(row) for row in rows]

        # 获取索引
        _, idx_rows = execute_query(conn, db_type, f"SHOW INDEX FROM `{table_name}`")
        if idx_rows:
            print("\n--- 索引 ---", file=sys.stderr)
            seen = set()
            for idx_row in idx_rows:
                idx_name = idx_row[2]
                if idx_name not in seen:
                    unique = "" if idx_row[1] else "UNIQUE"
                    print(f"  {idx_name} ({idx_row[4]}) {unique}", file=sys.stderr)
                    seen.add(idx_name)

        return result_columns, result_rows


def cmd_data(conn, db_type: str, table_name: str, limit: int = 10) -> tuple:
    """采样表数据"""
    table_name = validate_table_name(table_name)
    limit = min(limit, 1000)  # 硬上限
    quoted_table = quote_identifier(table_name, db_type)

    sql = f"SELECT * FROM {quoted_table} LIMIT {limit}"

    return execute_query(conn, db_type, sql)


def cmd_query(conn, db_type: str, sql: str) -> tuple:
    """执行自定义查询"""
    validate_readonly(sql)
    return execute_query(conn, db_type, sql)


# ---------------------------------------------------------------------------
# 输出格式化
# ---------------------------------------------------------------------------

def format_table(columns: list, rows: list) -> str:
    """对齐的文本表格"""
    if not columns:
        return "(无结果)"
    if not rows:
        return " | ".join(columns) + "\n(0 行)"

    # 计算列宽
    str_rows = [[str(v) if v is not None else "NULL" for v in row] for row in rows]
    widths = [max(len(col), max((len(r[i]) for r in str_rows), default=0)) for i, col in enumerate(columns)]

    # 表头
    header = " | ".join(col.ljust(w) for col, w in zip(columns, widths))
    separator = "-+-".join("-" * w for w in widths)

    # 数据行
    data_lines = []
    for row in str_rows:
        data_lines.append(" | ".join(val.ljust(w) for val, w in zip(row, widths)))

    result = f"{header}\n{separator}\n" + "\n".join(data_lines)
    result += f"\n({len(rows)} 行)"
    return result


def format_markdown(columns: list, rows: list) -> str:
    """Markdown 表格"""
    if not columns:
        return "(无结果)"
    if not rows:
        return "| " + " | ".join(columns) + " |\n" + "| " + " | ".join("---" for _ in columns) + " |\n\n*(0 行)*"

    str_rows = [[str(v) if v is not None else "NULL" for v in row] for row in rows]
    lines = ["| " + " | ".join(columns) + " |"]
    lines.append("| " + " | ".join("---" for _ in columns) + " |")
    for row in str_rows:
        lines.append("| " + " | ".join(row) + " |")
    lines.append(f"\n*({len(rows)} 行)*")
    return "\n".join(lines)


def format_json(columns: list, rows: list) -> str:
    """JSON 输出"""
    str_rows = [[v if v is not None else None for v in row] for row in rows]
    result = [dict(zip(columns, row)) for row in str_rows]
    return json.dumps(result, ensure_ascii=False, indent=2, default=str)


def format_csv(columns: list, rows: list) -> str:
    """CSV 输出"""
    import csv
    import io
    output = io.StringIO()
    writer = csv.writer(output)
    writer.writerow(columns)
    for row in rows:
        writer.writerow([str(v) if v is not None else "" for v in row])
    return output.getvalue()


FORMATTERS = {
    "table": format_table,
    "markdown": format_markdown,
    "json": format_json,
    "csv": format_csv,
}


# ---------------------------------------------------------------------------
# 主程序
# ---------------------------------------------------------------------------

def resolve_url(args) -> str:
    """从参数或环境变量解析连接 URL"""
    if args.url:
        return args.url
    if args.url_env:
        val = os.environ.get(args.url_env)
        if val:
            return val
        print(f"ERROR: 环境变量未设置: {args.url_env}", file=sys.stderr)
        sys.exit(1)
    # 检查常见环境变量
    for env_var in ["DATABASE_URL", "DB_URL", "POSTGRES_URL", "MYSQL_URL"]:
        val = os.environ.get(env_var)
        if val:
            return val
    print("ERROR: 请通过 --url 参数或 DATABASE_URL 环境变量提供连接信息", file=sys.stderr)
    sys.exit(1)


def main():
    parser = argparse.ArgumentParser(description="数据库查询工具")
    parser.add_argument("--db-type", choices=["postgres", "mysql", "sqlite"], required=True)
    parser.add_argument("--url", help="数据库连接 URL 或文件路径(SQLite)")
    parser.add_argument("--url-env", help="保存数据库连接信息的环境变量名")
    parser.add_argument("--format", choices=["table", "markdown", "json", "csv"], default="table")
    parser.add_argument("--timeout", type=positive_int, default=30, help="查询超时(秒)")

    subparsers = parser.add_subparsers(dest="command")

    subparsers.add_parser("test", help="测试连接")
    subparsers.add_parser("tables", help="列出所有表")

    schema_parser = subparsers.add_parser("schema", help="查看表结构")
    schema_parser.add_argument("table", help="表名")

    data_parser = subparsers.add_parser("data", help="采样表数据")
    data_parser.add_argument("table", help="表名")
    data_parser.add_argument("--limit", type=positive_int, default=10, help="返回行数(默认10)")

    query_parser = subparsers.add_parser("query", help="执行自定义 SQL")
    query_parser.add_argument("sql", help="SQL 语句")

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        sys.exit(1)

    url = resolve_url(args)
    formatter = FORMATTERS[args.format]

    with get_connection(args.db_type, url) as conn:
        if args.command == "test":
            cmd_test(conn, args.db_type)
            return

        if args.command == "tables":
            columns, rows = cmd_tables(conn, args.db_type)
        elif args.command == "schema":
            columns, rows = cmd_schema(conn, args.db_type, args.table)
        elif args.command == "data":
            columns, rows = cmd_data(conn, args.db_type, args.table, args.limit)
        elif args.command == "query":
            columns, rows = cmd_query(conn, args.db_type, args.sql)
        else:
            parser.print_help()
            sys.exit(1)

        print(formatter(columns, rows))


if __name__ == "__main__":
    main()
