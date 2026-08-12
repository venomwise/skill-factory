#!/usr/bin/env python3
"""Hacker News digest tool — fetch stories/comments/users from the official Firebase API.

Usage:
  hn.py categories
  hn.py top [--count N] [--min-score N] [--summary] [--vibe] [--one-liner]
  hn.py <category> [--count N] [--min-score N] [--summary] [--vibe]
  hn.py item <id>
  hn.py story <id> [--comments]
  hn.py user <username>
  hn.py search <query> [--count N]

Categories: top, best, new, ask (Ask HN), show (Show HN), job

Flags for digest output:
  --summary    include auto-generated 1-line content summary for each item
  --vibe       analyze comment sentiment (positive / mixed / debate, top concerns)
  --one-liner  append a single recommended-takeaway sentence (implied by --summary too)

Output: articles are printed as a terse ranked list.
"""
import argparse
import concurrent.futures
import html
import json
import re
import sys
import time
import urllib.parse
import urllib.request

BASE = "https://hacker-news.firebaseio.com/v0"
WORKERS = 8
CATEGORY_KEY = {
    "top": "topstories", "best": "beststories", "new": "newstories",
    "ask": "askstories", "show": "showstories", "job": "jobstories",
}

def get_many(getter, ids):
    """Concurrently fetch a list of ids through a getter function."""
    results = []
    with concurrent.futures.ThreadPoolExecutor(max_workers=WORKERS) as ex:
        for r in ex.map(getter, ids):
            results.append(r)
    return results

def get(path):
    url = BASE + path
    for attempt in range(3):
        try:
            req = urllib.request.Request(url, headers={"User-Agent": "minis-hn-skill"})
            with urllib.request.urlopen(req, timeout=20) as r:
                return json.loads(r.read())
        except Exception as e:
            if attempt == 2:
                raise SystemExit(f"API error after 3 retries: {e}")
            time.sleep(1.5 * (attempt + 1))


def strip_html(s):
    if not s:
        return ""
    s = re.sub(r"<[^>]+>", " ", s)
    s = html.unescape(s)
    return re.sub(r"\s+", " ", s).strip()


def cat_ids(category):
    return get(f"/{CATEGORY_KEY[category]}.json") or []


def fetch_items(ids):
    return get_many(lambda i: get(f"/item/{i}.json") or {}, ids)


# ----------------------------------------------------------------------------
# Summary: try Algolia search API for intro text, else DIY from item fields
# ----------------------------------------------------------------------------
def al_summary(item):
    """Return a short neutral intro for a linked story. Algolia stores story_text."""
    oid = item.get("id")
    if not oid:
        return ""
    try:
        url = f"https://hn.algolia.com/api/v1/items/{oid}"
        with urllib.request.urlopen(url, timeout=20) as r:
            d = json.loads(r.read())
        t = d.get("text") or (d.get("children") or [{}])[0].get("text", "")
        t = strip_html(t)
        if len(t) > 8:
            return t
    except Exception:
        pass
    return ""


def make_summary(item):
    """Build a 1-line content summary for a story item."""
    txt = item.get("text") or ""
    txt = strip_html(txt)

    if item.get("type") == "ask":
        if txt:
            return txt[:160] + ("…" if len(txt) > 160 else "")
        return "Ask HN discussion — no body text."
    if item.get("type") == "show":
        return make_show_summary(item, txt)
    if item.get("type") == "job":
        return (item.get("title") or "")[:100]

    # Linked story: use Algolia text as a community intro, tag it as snippet
    t = al_summary(item)
    if t:
        t = re.sub(r"https?://\S+", "", t).strip()
        if t:
            return "📌 摘要: " + t[:130] + ("…" if len(t) > 130 else "")
    if txt:
        return "📝 " + txt[:140] + ("…" if len(txt) > 140 else "")
    return "🔗 外部文章（点链接看原文）"


def make_show_summary(item, txt):
    base = (item.get("title") or "")[:120]
    if txt:
        return f"{base} — {txt[:130]}"
    return base


# ----------------------------------------------------------------------------
# Comment sentiment / vibe analysis (lightweight, no external deps)
# ----------------------------------------------------------------------------
POS = {"great","love","amazing","excellent","awesome","impressive","useful","helpful",
       "brilliant","fantastic","cool","nice","good","best","like ","thank ","works well","solid",
       "recommend","clean","elegant","powerful","fast","win","big deal","groundbreaking","genius"}
NEG = {"awful","terrible","bad","hate","useless","broken","bug","buggy","waste","fail","failure",
       "horrible","disappoint","wrong","risky","risk","dangerous","concern","worried","slow","crash",
       "overpriced","scam","garbage","meh","overhyped","dont like","not impressed","frustrat","not good",
       "hard to","disaster","ugly","unusable","overcomplicate"}
DEBATE = {"but","however","disagree","actually","false","incorrect","misleading","overblown",
          "critic","against","doubt","skeptic","overstated","why","yet","though","problem","is it",
          "stop","enough","where","how","surpris","wait","doesn't make sense","questionable",
          "could be better","i'm not sure","not so","on the other hand"}
CONCERN_KEY = {"privacy","security","cost","price","performance","speed","reliability","license",
               "battery","training data","copyright","regulation","ban","legal","gpu","memory",
               "vendor lock","proprietary","open source","closed source","ethics","censorship",
               "moderation","bias","halucinat","accuracy","quality","latency","complexity"}



def fetch_comments(item, cap=15):
    """Return top-level comment excerpts."""
    kids = (item.get("kids") or [])[:cap * 3]  # sample more to skip short ones
    if not kids:
        return []
    comments = get_many(lambda k: get(f"/item/{k}.json") or {}, kids)
    texts = []
    for c in comments:
        if (c.get("type") == "comment") and c.get("text"):
            t = strip_html(c["text"])
            if t and len(t) > 20:
                texts.append(t.lower())
        if len(texts) >= cap:
            break
    return texts[:cap]


def analyze_vibe(item, cap=15):
    """Return dict describing comment sentiment."""
    txts = fetch_comments(item, cap)
    n = len(txts)
    if n == 0:
        return {"n": 0, "mode": "no-comments",
                "pos": 0, "neg": 0, "deb": 0, "concerns": []}
    pos = neg = deb = 0
    for t in txts:
        pos_k = sum(1 for w in POS if w in t)
        neg_k = sum(1 for w in NEG if w in t)
        deb_k = sum(1 for w in DEBATE if w in t)
        if neg_k >= pos_k and neg_k > 2:
            neg += 1
        elif pos_k > 0 and neg_k > 0:
            deb += 1
        elif pos_k > 0:
            pos += 1
        else:
            deb += 1  # mostly neutral → mixed
    tot = pos + neg + deb or 1
    p = pos / tot
    mode = "positive" if p >= 0.5 else ("contentious" if neg / tot >= 0.35 else "mixed")
    p_nice = f"{pos}正/{deb}争议/{neg}负"
    concerns = sorted({w for t in txts for w in CONCERN_KEY if w in t})
    return {"n": n, "mode": mode, "pos": pos, "neg": neg, "deb": deb,
            "concerns": concerns, "pcount": p_nice}


VIBE_LABEL = {"positive": "😀 正面占主导", "mixed": "😐 褒贬不一",
              "contentious": "⚔️ 争论激烈", "no-comments": "💬 暂无评论"}


def one_liner(item, vibe):
    """Recommendation-style takeaway sentence."""
    title = item.get("title") or ""
    mode = vibe["mode"]
    if mode == "no-comments":
        return f"新帖，尚无讨论——适合当第一手资料看。"
    c = ", ".join(vibe["concerns"][:3]) if vibe["concerns"] else ""
    if mode == "positive":
        return f"社区反响正面，值得点开。"
    if mode == "contentious":
        hmm = f"核心争议在 {c}。" if c else "存在较大分歧。"
        return f"争议较大，{hmm}建议先看正反两派评论再判断。"
    return f"褒贬不一，{('关注点：' + c) if c else '讨论看法两极'}。"


def fmt_item(item, want_summary=False, want_vibe=False, ann=None):
    t = item.get("type", "?")
    title = item.get("title") or "(no title)"
    by = item.get("by", "?")
    score = item.get("score", 0)
    kids = len(item.get("kids") or [])
    oid = item.get("id")
    url = item.get("url") or f"https://news.ycombinator.com/item?id={oid}"
    lines = []
    lines.append(f"[#{oid}] {title} | score:{score} | comments:{kids} | by:{by}")
    if item.get("url"):
        lines.append(f"    {url}")
    if want_summary:
        s = (ann or {}).get("summary", make_summary(item))
        if s:
            lines.append(f"    {s}")
    if want_vibe:
        v = (ann or {}).get("vibe") or analyze_vibe(item)
        if v["n"]:
            pc = ", ".join(v["concerns"][:3]) if v["concerns"] else "无"
            lines.append(f"    💬 {VIBE_LABEL[v['mode']]} ({v.get('pcount', v['pos'])}) | 关注: {pc}")
            lines.append(f"    ✏️ 一句话: {one_liner(item, v)}")
    return "\n".join(lines)


def analysis_for(item):
    """Compute summary + vibe for one item (thread-safe sub-calls; wrap into one)."""
    return {"summary": make_summary(item), "vibe": analyze_vibe(item)}


def list_category(category, count, min_score, want_summary, want_vibe):
    ids = cat_ids(category)
    items = fetch_items(ids[:count * 3])
    items = [it for it in items if it.get("score", 0) >= min_score][:count]
    print(f"# {category.upper()} stories (showing {len(items)}):\n")

    # Precompute analysis concurrently so echoing out stays fast
    if want_summary or want_vibe:
        ann = get_many(analysis_for, items)
    else:
        ann = [None] * len(items)

    for i, (it, a) in enumerate(zip(items, ann), 1):
        print(f"{i}. {fmt_item(it, want_summary, want_vibe, a)}\n")


def show_item(item_id, with_comments):
    it = get(f"/item/{item_id}.json") or {}
    print(fmt_item(it, want_summary=True, want_vibe=with_comments))
    if with_comments:
        print("\n  comments:")
        kids = it.get("kids") or []
        for k in kids[:20]:
            c = get(f"/item/{k}.json") or {}
            if c.get("text"):
                txt = strip_html(c["text"])[:180]
                print(f"  - [{c.get('by','?')}] {txt}")


def show_user(name):
    u = get(f"/user/{name}.json") or {}
    if not u:
        print(f"User '{name}' not found.")
        return
    created = time.strftime("%Y-%m-%d", time.gmtime(u.get("created", 0)))
    submitted = u.get("submitted") or []
    print(f"User: {name}")
    print(f"  karma: {u.get('karma', 0)} | created: {created} | about: {(u.get('about') or '').strip()[:120]}")
    print(f"  submissions: {len(submitted)} total; recent 5: {submitted[-5:] or 'none'}")


def search(term, count):
    q = urllib.parse.quote(term)
    url = f"https://hn.algolia.com/api/v1/search?query={q}&tags=story&hitsPerPage={count}"
    try:
        with urllib.request.urlopen(url, timeout=20) as r:
            data = json.loads(r.read())
    except Exception as e:
        sys.exit(f"Search error: {e}")
    print(f"# Search '{term}':\n")
    for i, h in enumerate(data.get("hits", []), 1):
        pts = h.get("points", 0); nc = h.get("num_comments", 0)
        title = h.get("title", "(no title)"); author = h.get("author", "?")
        oid = h.get("objectID")
        print(f"{i}. {title} | points:{pts} | comments:{nc} | by:{author}")
        print(f"    https://news.ycombinator.com/item?id={oid}\n")


def main():
    ap = argparse.ArgumentParser()
    sub = ap.add_subparsers(dest="cmd", required=True)
    sub.add_parser("categories")
    for c in CATEGORY_KEY:
        p = sub.add_parser(c)
        p.add_argument("--count", type=int, default=5)
        p.add_argument("--min-score", type=int, default=0)
        p.add_argument("--summary", action="store_true")
        p.add_argument("--vibe", action="store_true")
    p = sub.add_parser("item"); p.add_argument("id", type=int)
    p = sub.add_parser("story", help="story detail + optional comments")
    p.add_argument("id", type=int); p.add_argument("--comments", action="store_true")
    p = sub.add_parser("user"); p.add_argument("name")
    p = sub.add_parser("search"); p.add_argument("query"); p.add_argument("--count", type=int, default=8)
    args = ap.parse_args()

    if args.cmd == "categories":
        print("Available categories: " + ", ".join(sorted(CATEGORY_KEY.keys()))); return
    if args.cmd in CATEGORY_KEY:
        list_category(args.cmd, args.count, getattr(args, "min_score", 0),
                      getattr(args, "summary", False), getattr(args, "vibe", False))
    elif args.cmd == "item":
        show_item(args.id, False)
    elif args.cmd == "story":
        show_item(args.id, args.comments)
    elif args.cmd == "user":
        show_user(args.name)
    elif args.cmd == "search":
        search(args.query, args.count)


if __name__ == "__main__":
    main()
