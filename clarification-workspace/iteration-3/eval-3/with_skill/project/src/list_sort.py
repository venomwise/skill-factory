PRIORITY_ORDER = {
    'urgent': 0,
    'high': 1,
    'normal': 2,
    'low': 3,
}


def sort_items(items):
    """Return items sorted by business priority for display.

    Items with the same priority keep the previous newest-first ordering.
    Unknown priorities are placed after known priorities.
    """
    return sorted(
        items,
        key=lambda item: (
            PRIORITY_ORDER.get(item.get('priority'), len(PRIORITY_ORDER)),
            -item['created_at'],
        ),
    )
