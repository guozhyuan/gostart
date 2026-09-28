# 隐藏接口倒推 Checklist（Network 面板看不到请求时用）

> 适用场景：网页数据明明加载出来了，但开发者工具 Network 里找不到对应的 XHR/fetch。
> 典型原因：**SSR 直出**、构建时内联、懒加载分块未触发、请求被 Service Worker/缓存接管。

---

## 0. 先判断数据到底怎么来的（30 秒定性）

```bash
curl -s -A "Mozilla/5.0 ... Chrome/124.0 Safari/537.36" "<URL>" -o page.html
grep -c 'href="/[0-9]\{3,\}"' page.html    # 命中多 => SSR 直出
```

| 现象 | 结论 | 下一步 |
|---|---|---|
| HTML 里已有数据 | SSR 直出 | 走 §1 §2（状态对象里通常直接写着接口） |
| HTML 是空壳/只有骨架 | CSR 渲染 | 走 §3（读 JS bundle 找 fetch 调用） |
| 滚动才出现 | 懒加载分块 | 走 §3，注意 chunk 是动态注入的 |

**关键认知**：Network 面板只能"观察别人怎么调"；读源码 + 实测才能"自己推导出怎么调"。

---

## 1. 找全局状态对象（最高性价比的一步）

SSR 框架会把首屏数据序列化到 `window` 上，**接口路径往往也一起被序列化进去**。

```bash
grep -o 'window\.[A-Za-z_$][A-Za-z0-9_$]*' page.html | sort -u
```

常见变量名（本工具已内置这些模式）：

| 变量 | 框架 |
|---|---|
| `window.$DATA` | 斗鱼 Shark |
| `window.__NUXT__` | Nuxt (Vue) |
| `<script id="__NEXT_DATA__">` | Next.js |
| `window.__INITIAL_STATE__` | Vue/Vuex SSR |
| `window._INIT_DATA` / `window.PAGE_DATA` | 各类自研 |
| `window.__APOLLO_STATE__` | Apollo GraphQL |

```bash
# 提取并格式化
python - <<'EOF'
import re, json, io
t = io.open('page.html', encoding='utf-8').read()
m = re.search(r'window\.\$DATA\s*=\s*(\{.*?\});\s*var', t, re.S)
d = json.loads(m.group(1))
print(list(d.keys()))
for k in ('pagePath','tabTagPath','apiUrl','listUrl'):
    if k in d: print(k, '=', d[k])
EOF
```

**典型收获**（斗鱼实例）：
```json
"pagePath":   "/gapi/rknc/directory/mixListV1/2_4133/",
"tabTagPath": "/gapi/rkc/directory/c_tag/2_4133/list",
"cateInfo":   {"cate2Name":"三角洲行动","cid2":4133}, "currentPage": 1
```
→ 接口命名规律一眼可见：`mixListV1/{分类}_{id}/{页码}`。

---

## 2. 判定哪个候选才是"列表接口"

状态里常有多个路径，按以下优先级挑：

| 特征 | 权重 | 说明 |
|---|---|---|
| 含 `mixList` / `roomList` / `feed` / `list` | ★★★ | 语义上就是列表 |
| 路径**以 `/` 结尾** | ★★★ | 页码通常直接拼在最后一段 |
| 含 `/gapi/` `/api/` | ★★ | 数据接口命名空间 |
| 含 `数字_数字`（如 `2_4133`） | ★ | 分类 id 参数 |
| 含 `c_tag` / `tag` / `user` / `rank` / `login` | 排除 | 是标签/榜单/用户接口，不是列表 |

---

## 3. 下载 JS bundle，反查"路径怎么用"

拿到路径还不够，必须知道**页码拼哪、参数叫啥、响应取哪个字段**。

```bash
# 列出该页所有 script（含动态 loader）
grep -o '<script[^>]*src="[^"]*"' page.html

# loader 会把该页所有 chunk 列出来，全部拉下来后按关键字定位
grep -o 'https://[^"]*\.js' page.html | sort -u
```

在 bundle 里搜这几类模式：

```python
# a) 直接调用点
r'\.(?:get|post|request|fetch)\([^,)]{0,80}?["\'](/[A-Za-z0-9_\-/\.]+)["\']'
# b) 常量路径
r'["\'](/(?:gapi|api|japi)[A-Za-z0-9_\-/\.]*)["\']'
# c) 页码拼接（★ 最关键）
r'=\s*""\s*\+\s*(\w+)\s*\+\s*(\w+)\s*;'
```

斗鱼实例找到的唯一真源（`listAll~*.js`）：
```js
var l = "" + pagePath + page;          // url = pagePath + page
listCellService.getListStream({ url: l, page: page, pagePath: pagePath, data: readList });
// 解析: y.pgcnt=总页数, y.rl=房间数组
```

---

## 4. 实测验证（不猜，只测）

对每个不确定点做 A/B，用**返回条数 / 首条 id 是否变化**判定：

| 要验证的问题 | 怎么测 | 斗鱼实测结果 |
|---|---|---|
| 页码在路径还是 query？ | `/{page}` vs `?page=` | 路径有效，`?page=`/`?limit=`/`?sort=` 全被忽略 |
| 域名前缀哪个对？ | `/gapi/rkc/` vs `/gapi/rknc/` | 两者等价（200，内容一致） |
| 分页上限？ | 递增翻到空数组 | 1–14 页×40、15 页 32、16 页空 → 上限约 600 |
| 参数化的维度？ | 换分类 id | `2_1`→英雄联盟、`2_270`→热门游戏，规律成立 |
| 响应有总页数吗？ | 看 JSON 字段 | 只有 `rl` + `userRecommendRec`，无总数，只能靠空页停 |

必备请求头：`User-Agent`（真浏览器）、`Referer`（页面自身 URL）、`Accept`。

---

## 5. 闭环：写代码跑通

只"看起来像"不算数，**能用代码稳定抓下来才算定位成功**。

```python
import requests
url = "https://www.douyu.com/gapi/rkc/directory/mixListV1/2_4133/1"
j = requests.get(url, headers={"User-Agent": "...Chrome/124.0...",
                               "Referer": "https://www.douyu.com/g_DF"}, timeout=20).json()
print(j["code"], len(j["data"]["rl"]))
```

---

## 6. 一键化：`hidden_api_finder.py`

上述 §0–§5 已封装成工具，无需手工重复：

```bash
python hidden_api_finder.py <URL> [选项]
```

| 选项 | 作用 |
|---|---|
| `-o DIR` | 输出目录（默认用域名） |
| `-k rl,list,...` | 指定数据数组字段名 |
| `-p` | 对候选接口逐个发 GET 实测（状态码/条数） |
| `-A N` | 翻页抓取直到空页（最多 N 页） |
| `--csv` | 额外导出 CSV |
| `--top N` | 候选输出上限 |

产出：`index.html`、`apis.txt`（候选清单）、`state.json`（提取出的全局状态）、
`report.txt`（完整报告）、`rooms.json` / `rooms.csv`（列表数据）。

斗鱼实测：
```bash
python hidden_api_finder.py https://www.douyu.com/g_DF -p -A 16 --csv
# => 自动判定 SSR 直出 / 自动提取 $DATA / 自动选中
#    /gapi/rknc/directory/mixListV1/2_4133/  / 翻页抓到 593 条并导出
```

---

## 7. 常见坑

| 坑 | 表现 | 对策 |
|---|---|---|
| 首屏 SSR，看不到请求 | Network 只有 document | 直接读 HTML 里的状态对象 |
| chunk 懒加载未触发 | 滚动后才有请求 | 手动滚到底 / 看 loader 列出的 chunk 列表 |
| 保留日志没开 | 跳转后记录被清 | DevTools 勾 Disable cache + Preserve log |
| `pgcnt` 含义误判 | 以为是总页数 | 确认它是列表容器的字段还是标签项的字段（斗鱼踩过：SSR 里 `pgcnt` 是标签条数，顶层 `pageCount` 才是真值且为 0） |
| 请求头缺失 | 403 / 空数组 | 带 `Referer` + 真 UA |
| 静态资源混进候选 | 一堆 `.js/.png` | 按扩展名过滤（工具已内置） |
| 频率过快 | 429 / 封 IP | 加 `sleep`，控制并发 |
