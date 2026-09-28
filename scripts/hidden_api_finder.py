# -*- coding: utf-8 -*-
"""
hidden_api_finder.py —— 从 SSR/SPA 页面倒推隐藏数据接口的通用小工具

背景
   很多网站的列表首屏是"服务端直出(SSR)"或"构建时内联"的, 打开浏览器
   开发者工具看不到任何 XHR 请求, 但数据确实在页面上。本工具把
   "curl 页面 -> 找全局状态 -> 找接口候选 -> 反查 JS 调用点 -> 实测分页"
   这套人工流程自动化。

用法
   python hidden_api_finder.py <URL> [选项]

   python hidden_api_finder.py https://www.douyu.com/g_DF --json-key rl --probe
   python hidden_api_finder.py https://www.douyu.com/g_DF -A 20 --csv

选项
   -o, --out DIR      输出目录(默认: 以域名为名的目录)
   -k, --json-key K   房间/数据列表在 JSON 里的字段名(默认: rl,list,data,items,rows,rooms)
   -p, --probe        对 /gapi /api 候选接口发一次 GET 实测
   -A, --all-rows N   翻页抓取直到空页(最多 N 页)
   --csv              额外导出表格 CSV
   --top N            接口候选最多输出 N 条(默认 60)
   --timeout S        单请求超时秒数(默认 25)

产出(在输出目录内)
   index.html          页面源码
   apis.txt            接口候选清单(按类型分组)
   state.json          从页面提取到的全局状态对象(若有)
   rooms.csv           列表数据(仅当 -A/--all-rows 或 --csv)
   report.txt          完整分析报告

依赖: requests
   python -m pip install requests
"""
from __future__ import print_function

import argparse
import csv
import io
import json
import os
import re
import sys
import time

try:
    import requests
except ImportError:
    sys.exit("missing dependency -> python -m pip install requests")

UA = ("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
      "(KHTML, like Gecko) Chrome/124.0 Safari/537.36")

DEFAULT_LIST_KEYS = ["rl", "list", "data", "items", "rows", "rooms", "room_list", "videoList"]

# ---------------------------------------------------------------- 基础工具

def log(msg=""):
    """ASCII 安全输出, 避免 Windows 控制台 GBK 报错"""
    try:
        print(msg)
    except UnicodeEncodeError:
        print(msg.encode("utf-8", "replace").decode("ascii", "replace"))


def fetch(url, referer=None, timeout=25):
    headers = {"User-Agent": UA, "Accept": "text/html,application/json,*/*"}
    if referer:
        headers["Referer"] = referer
    r = requests.get(url, headers=headers, timeout=timeout)
    return r


def domain_of(url):
    m = re.match(r"https?://([^/]+)", url)
    return m.group(1).replace(":", "_") if m else "output"


# ------------------------------------------------------ 1. 全局状态提取

STATE_PATTERNS = [
    # (名字, 正则, 是否 JSON)
    ("$DATA",        r"window\.\$DATA\s*=\s*(\{.*?\});\s*(?:var|</script>)", True),
    ("__NUXT__",     r"window\.__NUXT__\s*=\s*(\{.*?\})\s*;?\s*</script>", True),
    ("__NEXT_DATA__", r'<script id="__NEXT_DATA__"[^>]*>(\{.*?\})</script>', True),
    ("__INITIAL_STATE__", r"window\.__INITIAL_STATE__\s*=\s*(\{.*?\});\s*(?:var|</script>)", True),
    ("_INIT_DATA",   r"window\._INIT_DATA\s*=\s*(\{.*?\});\s*(?:var|</script>)", True),
    ("PAGE_DATA",    r"window\.PAGE_DATA\s*=\s*(\{.*?\});\s*(?:var|</script>)", True),
    ("__APOLLO_STATE__", r"window\.__APOLLO_STATE__\s*=\s*(\{.*?\});\s*</script>", True),
]


def extract_states(html):
    """返回 [(名字, dict)]"""
    out = []
    for name, pat, is_json in STATE_PATTERNS:
        for m in re.finditer(pat, html, re.S):
            raw = m.group(1)
            if is_json:
                try:
                    out.append((name, json.loads(raw)))
                except Exception:
                    # 截断兜底: 退化为字符串, 仍可正则取路径
                    out.append((name + "(unparsed)", raw))
    return out


def walk_state(obj, path="", found=None, depth=0):
    """递归遍历状态对象, 收集所有 http/路径类字符串与列表容器"""
    if found is None:
        found = {"urls": [], "lists": []}
    if depth > 12:
        return found
    if isinstance(obj, dict):
        for k, v in obj.items():
            p = path + "." + str(k) if path else str(k)
            if isinstance(v, str) and re.match(r"^(https?:)?/[A-Za-z0-9_\-/\.\?=&%:]*$", v) and len(v) > 3:
                found["urls"].append((p, v))
            walk_state(v, p, found, depth + 1)
    elif isinstance(obj, list):
        found["lists"].append((path, len(obj)))
        if obj:
            walk_state(obj[0], path + "[0]", found, depth + 1)
    return found


# --------------------------------------------------------- 2. 接口候选

API_HINTS = ("/gapi", "/api", "/japi", "/ajax", "/json", "/xhr", "/dir_new",
             "/live", "/list", "/search", "/feed", "/graphql", "/rpc", "/v1", "/v2")

# 明显不是接口的静态资源
STATIC_EXT = re.compile(r"\.(js|css|png|jpe?g|gif|webp|avif|svg|ico|woff2?|ttf|mp4|m3u8|flv)(\?|$)", re.I)


def find_api_candidates(html, states=None):
    """
    返回 {'in_state': [(来源, 路径)], 'absolute': [...], 'relative': [...], 'js_urls': [...]}
    """
    res = {"in_state": [], "absolute": [], "relative": [], "js_urls": []}

    # a) 全局状态里的路径 —— 信息量最高
    for name, st in (states or []):
        if isinstance(st, dict):
            f = walk_state(st)
            for p, v in f["urls"]:
                if any(h in v for h in API_HINTS) and not STATIC_EXT.search(v):
                    res["in_state"].append((name + ":" + p, v))

    # b) HTML 内联脚本里的绝对地址
    for m in re.finditer(r"""["'](https?://[A-Za-z0-9\.\-]+/[A-Za-z0-9_\-/\.]*(?:gapi|api|japi|ajax|json)[A-Za-z0-9_\-/\.\?=&%]*)["']""", html):
        res["absolute"].append(m.group(1))

    # c) HTML 里的相对接口路径
    for m in re.finditer(r"""["'](/(?:gapi|api|japi|ajax|dir_new|rknc|rkc)[A-Za-z0-9_\-/\.]*)["']""", html):
        res["relative"].append(m.group(1))

    # d) 待下载的 JS bundle
    for m in re.finditer(r"""<script[^>]+src=["']([^"']+\.js[^"']*)["']""", html):
        res["js_urls"].append(m.group(1))

    for k in ("absolute", "relative"):
        seen, uniq = set(), []
        for v in res[k]:
            if v not in seen and not STATIC_EXT.search(v):
                seen.add(v)
                uniq.append(v)
        res[k] = uniq

    seen, uniq = set(), []
    for src, v in res["in_state"]:
        if v not in seen:
            seen.add(v)
            uniq.append((src, v))
    res["in_state"] = uniq
    res["js_urls"] = sorted(set(res["js_urls"]))
    return res


def scan_js_for_apis(js_text):
    """
    在 JS bundle 里找接口调用与路径拼接:
      1) httpClient.get(...,"URL",...) / axios.get("URL")
      2) 变量常量赋值  X="/gapi/..."
      3) pagePath + page 形式的拼接
    """
    hits = {"calls": [], "consts": [], "concat": []}

    for m in re.finditer(r"""\.(?:get|post|request|fetch|ajax)\s*\(\s*[^,)]{0,80}?["'](/[A-Za-z0-9_\-/\.]+)["']""", js_text):
        hits["calls"].append(m.group(1))
    for m in re.finditer(r"""["'](/(?:gapi|api|japi|dir_new|rknc|rkc)[A-Za-z0-9_\-/\.]*)["']""", js_text):
        hits["consts"].append(m.group(1))
    for m in re.finditer(r"""["']\s*\+\s*[A-Za-z_$][\w$]*\s*\+\s*["']""", js_text):
        hits["concat"].append(m.group(0))
    # 典型: var l = "" + pagePath + page
    for m in re.finditer(r'''=\s*""\s*\+\s*(\w+)\s*\+\s*(\w+)\s*;''', js_text):
        hits["concat"].append("%s + %s" % (m.group(1), m.group(2)))

    for k in hits:
        seen, uniq = set(), []
        for v in hits[k]:
            if v not in seen:
                seen.add(v)
                uniq.append(v)
        hits[k] = uniq
    return hits


# ----------------------------------------------------- 3. 数据抽取与翻页

def pick_list(payload, want_keys):
    """在响应里找房间数组"""
    if isinstance(payload, str):
        try:
            payload = json.loads(payload)
        except Exception:
            return None, None
    if not isinstance(payload, dict):
        return None, None
    containers = [payload]
    for v in payload.values():
        if isinstance(v, dict):
            containers.append(v)
    for c in containers:
        for k in want_keys:
            v = c.get(k)
            if isinstance(v, list) and v and isinstance(v[0], dict):
                return k, v
    return None, None


def sniff_fields(rows):
    """挑出常见语义字段, 便于导出"""
    if not rows:
        return []
    keys = list(rows[0].keys())
    pref = ["rid", "room_id", "roomId", "id", "rn", "title", "name", "nn", "nickname",
            "uid", "anchor", "ol", "hot", "online", "c2name", "category", "game",
            "cid2", "cid3", "url", "link", "type", "status"]
    out = [k for k in pref if k in keys]
    for k in keys:
        if k not in out and not isinstance(rows[0][k], (dict, list)):
            out.append(k)
    return out


def probe_endpoint(url, referer, timeout):
    try:
        r = fetch(url, referer=referer, timeout=timeout)
        ct = r.headers.get("Content-Type", "")
        body = r.text
        ok = r.status_code == 200 and ("json" in ct.lower() or body.lstrip()[:1] in "{[")
        n = 0
        if ok:
            try:
                k, rows = pick_list(body, DEFAULT_LIST_KEYS)
                n = len(rows) if rows else 0
            except Exception:
                n = 0
        return r.status_code, len(body), ct, ok, n
    except Exception as e:
        return None, 0, str(e), False, 0


def build_page_url(base, page):
    """base 形如 /gapi/rkc/directory/mixListV1/2_4133/  ->  末尾追加页码"""
    return base.rstrip("/") + "/" + str(page)


def join_origin(url, origin):
    if url.startswith("//"):
        return "https:" + url
    if url.startswith("/"):
        return origin + url
    if re.match(r"^https?://", url):
        return url
    return origin + "/" + url.lstrip("/")


# 列表类接口特征(加分) 与 非列表接口特征(排除)
LIST_HINT = re.compile(r"(mixList|roomList|room_list|getRoom|liveList|listRoom|feed|/list\b)", re.I)
NOT_LIST = re.compile(r"(/c_tag/|/tag|/user/|/login|/follow|/gift|/rank|videotag)", re.I)


def pick_base_url(candidates):
    """从 (来源, 路径) 候选中挑出最像"房间列表 + 可翻页"的一个"""
    best, best_score = None, -1
    for src, v in candidates:
        if STATIC_EXT.search(v) or NOT_LIST.search(v):
            continue
        score = 0
        if LIST_HINT.search(v):
            score += 5
        if v.endswith("/"):
            score += 3          # 末尾斜杠 => 页码拼在最后一段
        if "/gapi/" in v:
            score += 2
        if re.search(r"\d+_?\d*", v):
            score += 1
        if score > best_score:
            best, best_score = v, score
    return best


def crawl(url, referer, page_size, max_pages, want_keys, timeout):
    """翻页抓取: 空页即停"""
    all_rows, page = [], 1
    while page <= max_pages:
        r = fetch(build_page_url(url, page), referer=referer, timeout=timeout)
        if r.status_code != 200:
            log("  page=%d HTTP %s (stop)" % (page, r.status_code))
            break
        k, rows = pick_list(r.text, want_keys)
        rows = rows or []
        log("  page=%-3d rows=%d" % (page, len(rows)))
        if not rows:
            break
        all_rows.extend(rows)
        page += 1
        time.sleep(0.3)
    return all_rows


# ------------------------------------------------------------------- main

def main():
    ap = argparse.ArgumentParser(description="从 SSR/SPA 页面倒推隐藏数据接口")
    ap.add_argument("url")
    ap.add_argument("-o", "--out")
    ap.add_argument("-k", "--json-key", default=",".join(DEFAULT_LIST_KEYS))
    ap.add_argument("-p", "--probe", action="store_true")
    ap.add_argument("-A", "--all-rows", type=int, default=0, metavar="N")
    ap.add_argument("--csv", action="store_true")
    ap.add_argument("--top", type=int, default=60)
    ap.add_argument("--timeout", type=float, default=25)
    args = ap.parse_args()

    want_keys = [s.strip() for s in args.json_key.split(",") if s.strip()]
    referer = args.url
    outdir = args.out or domain_of(args.url)
    if not os.path.isdir(outdir):
        os.makedirs(outdir)

    report = io.StringIO()

    def w(line=""):
        log(line)
        report.write(line + "\n")

    w("=" * 72)
    w("hidden_api_finder  target: %s" % args.url)
    w("=" * 72)

    # --- 抓页面
    r = fetch(args.url, timeout=args.timeout)
    html = r.text
    w("[1] fetch page -> HTTP %s, %d bytes" % (r.status_code, len(html)))
    io.open(os.path.join(outdir, "index.html"), "w", encoding="utf-8").write(html)

    # --- SSR 判定
    n_cards = len(re.findall(r'href="/(\d{3,})"', html))
    w("    HTML 内形如 href=\"/<数字>\" 的链接数: %d" % n_cards)
    if n_cards > 5:
        w("    => 判定: 首屏数据很可能是 SSR 直出 (所以 Network 里看不到请求)")
    else:
        w("    => 首屏可能是 JS 渲染, 接口信息主要在 JS bundle 里")

    # --- 全局状态
    states = extract_states(html)
    w("\n[2] 全局状态对象: %d 个" % len(states))
    parsed_states = []
    for name, st in states:
        if isinstance(st, dict):
            parsed_states.append((name, st))
            w("    - %-22s 顶层键: %s" % (name, ", ".join(list(st.keys())[:12])))
        else:
            w("    - %-22s (JSON 解析失败, 仅字符串匹配)" % name)
    if parsed_states:
        blob = {n: s for n, s in parsed_states}
        io.open(os.path.join(outdir, "state.json"), "w", encoding="utf-8").write(
            json.dumps(blob, ensure_ascii=False, indent=2)[:8 * 1024 * 1024])
        w("    -> state.json 已保存")

    # --- 状态里的接口路径(重点)
    cand = find_api_candidates(html, parsed_states)
    w("\n[3] 接口候选")
    w("    (a) 全局状态内直接给出的路径  ★ 最可信")
    for src, v in cand["in_state"][:args.top]:
        w("        %-46s = %s" % (src, v))
    if not cand["in_state"]:
        w("        (无)")
    w("    (b) HTML 内联的绝对接口地址: %d 条" % len(cand["absolute"]))
    for v in cand["absolute"][:args.top]:
        w("        %s" % v)
    w("    (c) HTML 内联的相对接口路径: %d 条" % len(cand["relative"]))
    for v in cand["relative"][:args.top]:
        w("        %s" % v)
    w("    (d) JS bundle: %d 个" % len(cand["js_urls"]))
    for v in cand["js_urls"][:12]:
        w("        %s" % v)

    # --- 若拿到列表接口, 抽字段
    origin = re.match(r"(https?://[^/]+)", args.url).group(1)
    base_url = pick_base_url(cand["in_state"])
    if base_url:
        w("\n[4] 选定列表接口(来自状态): %s" % base_url)
    elif cand["in_state"]:
        base_url = cand["in_state"][0][1]
        w("\n[4] 退化选取状态内首个路径: %s" % base_url)

    sample_rows, list_key = [], None
    if base_url:
        w("\n[5] 探测分页规律")
        for p in (1, 2, 3):
            try:
                rr = fetch(join_origin(build_page_url(base_url, p), origin),
                           referer=referer, timeout=args.timeout)
                k, rows = pick_list(rr.text, want_keys)
                rows = rows or []
                w("    page=%d -> HTTP %s, 列表字段=%s, 条数=%d" % (p, rr.status_code, k, len(rows)))
                if rows and not sample_rows:
                    sample_rows, list_key = rows, k
                if not rows:
                    break
            except Exception as e:
                w("    page=%d ERR %s" % (p, e))
                break
        if sample_rows:
            w("    每页条数固定(无 limit 参数时): 约 %d 条" % len(sample_rows))
            w("    字段: %s" % ", ".join(sniff_fields(sample_rows)))
            w("    建议页码写法: %s  (页码拼在路径末尾)" % build_page_url(base_url, "<page>"))

    # --- 实测其它候选
    if args.probe:
        w("\n[6] 实测候选接口 (-p)")
        tested, seen_u = [], set()
        pool = [v for _, v in cand["in_state"]] + cand["relative"][:10]
        for v in pool:
            full = join_origin(v, origin)
            if full in seen_u:
                continue
            seen_u.add(full)
            if len(tested) >= 12:
                break
            code, ln, ct, ok, n = probe_endpoint(full, referer, args.timeout)
            w("    %-6s %-10s len=%-7s rows=%-4s %s" % (code, "JSON" if ok else ct[:10], ln, n, full))
            tested.append({"url": full, "status": code, "rows": n})
        io.open(os.path.join(outdir, "probe.json"), "w", encoding="utf-8").write(
            json.dumps(tested, ensure_ascii=False, indent=2))

    # --- 翻页抓取
    if args.all_rows and base_url:
        w("\n[7] 翻页抓取 (max %d 页)" % args.all_rows)
        rows = crawl(join_origin(base_url, origin), referer, len(sample_rows) or 40,
                     args.all_rows, want_keys, args.timeout)
        w("    合计 %d 条" % len(rows))
        if rows:
            fields = sniff_fields(rows)
            simple = [{k: r.get(k) for k in fields} for r in rows]
            io.open(os.path.join(outdir, "rooms.json"), "w", encoding="utf-8").write(
                json.dumps(simple, ensure_ascii=False, indent=2))
            w("    -> rooms.json 已保存(%d 条)" % len(simple))
        if args.csv and rows:
            fields = sniff_fields(rows)
            with io.open(os.path.join(outdir, "rooms.csv"), "w", encoding="utf-8-sig", newline="") as f:
                cw = csv.DictWriter(f, fieldnames=fields, extrasaction="ignore")
                cw.writeheader()
                for r in rows:
                    cw.writerow({k: r.get(k) for k in fields})
            w("    -> rooms.csv 已保存")

    # --- 写 apis.txt
    with io.open(os.path.join(outdir, "apis.txt"), "w", encoding="utf-8") as f:
        f.write("# 全局状态内路径 (最可信)\n")
        for src, v in cand["in_state"]:
            f.write("%s\t%s\n" % (src, v))
        f.write("\n# HTML 内联绝对地址\n")
        for v in cand["absolute"]:
            f.write("%s\n" % v)
        f.write("\n# HTML 内联相对路径\n")
        for v in cand["relative"]:
            f.write("%s\n" % v)
        f.write("\n# JS bundle\n")
        for v in cand["js_urls"]:
            f.write("%s\n" % v)

    report_path = os.path.join(outdir, "report.txt")
    io.open(report_path, "w", encoding="utf-8").write(report.getvalue())

    w("\n[done] 输出目录: %s" % os.path.abspath(outdir))
    w("       apis.txt / state.json / report.txt" +
      (" / rooms.json" if args.all_rows else "") +
      (" / rooms.csv" if args.csv else ""))


if __name__ == "__main__":
    main()
