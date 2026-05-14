PRIORITY_ORDER = {
    'urgent': 0,
    'high': 1,
    'normal': 2,
    'low': 3,
}


def sort_items(items):
    """Return items sorted for display by business priority.

    Items with the same priority keep the previous newest-first ordering.
    Unknown priorities are placed after known priorities.
    """
    newest_first = sorted(items, key=lambda item: item['created_at'], reverse=True)
    return sorted(
        newest_first,
        key=lambda item: PRIORITY_ORDER.get(item.get('priority'), len(PRIORITY_ORDER)),
    )
