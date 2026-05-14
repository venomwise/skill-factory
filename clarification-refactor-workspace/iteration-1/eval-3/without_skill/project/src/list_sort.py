PRIORITY_ORDER = {
    'urgent': 0,
    'high': 1,
    'normal': 2,
    'low': 3,
}


def sort_items(items):
    """Return items sorted by business priority for display."""
    return sorted(
        items,
        key=lambda item: (PRIORITY_ORDER[item['priority']], item['created_at']),
        reverse=False,
    )
