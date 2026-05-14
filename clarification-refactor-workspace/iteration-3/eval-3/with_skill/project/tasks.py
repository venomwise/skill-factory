#!/usr/bin/env python3
"""Task management system."""
from datetime import datetime
from typing import List, Dict


TASKS = [
    {
        "id": 1,
        "name": "修复登录页面的验证bug",
        "created": "2024-01-10",
        "status": "pending"
    },
    {
        "id": 2,
        "name": "添加用户导出功能",
        "created": "2024-01-08",
        "status": "in_progress"
    },
    {
        "id": 3,
        "name": "优化数据库查询性能",
        "created": "2024-01-12",
        "status": "pending"
    },
    {
        "id": 4,
        "name": "更新API文档",
        "created": "2024-01-05",
        "status": "pending"
    },
    {
        "id": 5,
        "name": "实现支付功能",
        "created": "2024-01-15",
        "status": "pending"
    }
]


def get_sorted_tasks() -> List[Dict]:
    """Get tasks sorted by creation date (oldest first)."""
    return sorted(TASKS, key=lambda x: x['created'])


def display_tasks(tasks: List[Dict]):
    """Display tasks in a formatted list."""
    print("\n任务列表:")
    print("-" * 60)
    for i, task in enumerate(tasks, 1):
        print(f"{i}. [{task['status']}] {task['name']}")
        print(f"   创建时间: {task['created']}")
    print("-" * 60)


if __name__ == "__main__":
    tasks = get_sorted_tasks()
    display_tasks(tasks)
