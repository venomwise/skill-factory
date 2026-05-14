PRIORITY_ORDER = {
    'urgent': 0,
    'high': 1,
    'normal': 2,
    'low': 3,
}


def sort_items(items):
    """Return items sorted for display by business priority.

    Items with the same priority keep the previous recency-based ordering.
    Unknown priorities are placed after the known business priorities.
    """
    return sorted(
        items,
        key=lambda item: (
            PRIORITY_ORDER.get(item.get('priority'), len(PRIORITY_ORDER)),
            -item['created_at'],
        ),
    )
