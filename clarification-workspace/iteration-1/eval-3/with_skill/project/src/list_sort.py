PRIORITY_ORDER = {
    'urgent': 0,
    'high': 1,
    'normal': 2,
    'low': 3,
}


def sort_items(items):
    """Return items sorted for display by business priority."""
    return sorted(
        items,
        key=lambda item: (
            PRIORITY_ORDER.get(item['priority'], len(PRIORITY_ORDER)),
            item['created_at'],
        ),
    )
