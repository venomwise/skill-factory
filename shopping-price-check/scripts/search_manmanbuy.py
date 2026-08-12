#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""慢慢买免登录比价抓取器 (skill shopping-price-check)

用法:
    python3 search_manmanbuy.py "iphone 15"
    python3 search_manmanbuy.py "switch oled" --limit 15

输出: 编号条目列表 (商品名 | 价格 | 平台 | 时间 | 好价标签)
说明: 优先 shell 直连(SSR 服务端渲染, 稳), 内置浏览器渲染 Next.js 页反而会卡死。
"""
import re, sys, json, html, urllib.parse, urllib.request

UA = ("Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) "
      "AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1")

BASE = "https://s.manmanbuy.com/m/search/result?keyword="
PLATFORMS = ['京东自营','京东商城','天猫商城','天猫国际','天猫旗舰店','淘宝','拼多多','唯品会','苏宁','淘宝特价','抖音商城']
# 附加条件/好价标签关键字(价格后)
EXTRA = ['淘金币','88VIP','国补','补贴','券']

def fetch(keyword):
    url = BASE + urllib.parse.quote(keyword)
    req = urllib.request.Request(url, headers={'User-Agent': UA})
    return urllib.request.urlopen(req, timeout=25).read().decode('utf-8', 'ignore')

def parse(html_text):
    t = re.sub(r'<script.*?</script>', '', html_text, flags=re.S)
    t = re.sub(r'<style.*?</style>', '', t, flags=re.S)
    t = re.sub(r'<[^>]+>', '|', t)
    t = html.unescape(t)
    t = re.sub(r'[ \t\u3000]+', ' ', t)
    t = re.sub(r'\|+', '|', t)
    total_m = re.search(r'全网(\d+)条', t)
    total = total_m.group(1) if total_m else '?'
    for p in PLATFORMS:
        t = t.replace(p, '\n@@'+p)
    items = []
    for line in t.split('\n'):
        line = line.strip()
        if not line.startswith('@@'):
            continue
        parts = line[2:].strip('|').split('|')
        if not parts or not parts[0]:
            continue
        platform = parts[0]
        tm = next((x for x in parts if re.fullmatch(r'\d{2}-\d{2}( \d{2}:\d{2})?', x)), '')
        price = next((x for x in reversed(parts) if re.match(r'^\d+(\.\d+)?元', x)), '')
        pi = parts.index(price) if price else -1
        pre = parts[1:pi] if pi >= 0 else parts[1:]
        name_cands = [x for x in pre if x and len(x) >= 4 and x != tm]
        name = max(name_cands, key=len) if name_cands else ''
        tags = [x for x in pre if x and len(x) < 4 and x != tm]
        if pi >= 0:
            tags += [x for x in parts[pi+1:] if re.search(r'新低|历史|均价|双11|精选|已结束|补贴|爆料|再降价', x)]
        if not name:
            continue
        items.append({'name': name, 'price': price, 'platform': platform,
                      'time': tm, 'tags': ' '.join(t for t in tags if t)})
    return total, items

def main():
    args = [a for a in sys.argv[1:] if not a.startswith('--')]
    limit = 15
    if '--limit' in sys.argv:
        try: limit = int(sys.argv[sys.argv.index('--limit')+1])
        except Exception: pass
    keyword = ' '.join(args) if args else 'iphone 15'
    try:
        raw = fetch(keyword)
    except Exception as e:
        print(f"抓取失败: {e}"); sys.exit(1)
    total, items = parse(raw)
    print(f"关键词: {keyword} | 慢慢买全网共 {total} 条爆料好价")
    print(("-"*70))
    for i, it in enumerate(items[:limit], 1):
        extra = ' '.join(x for x in [it['price'][3:]] if False)
        tail = it['price']
        print(f"{i:2d}. {it['name'][:42]}\n    {tail}  |  {it['platform']}  |  {it['time']}  |  {it['tags']}")
    print(f"(显示前 {min(limit,len(items))} 条, 共 {len(items)} 条可解析; 标注'已结束'为过期爆料)")

if __name__ == '__main__':
    main()