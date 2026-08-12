# 聚合比价信息源参考

## 主源：慢慢买 Manmanbuy

- 移动端首页：`https://m.manmanbuy.com/`
- 搜索直达：`https://s.manmanbuy.com/m/search/result?keyword=<URL编码关键词>`（首选）
- 历史价格查询：`https://tool.manmanbuy.com/m/history.aspx`（单个商品历史价曲线）
- 特点：覆盖京东/天猫/淘宝/唯品会/苏宁/拼多多(百亿补贴)爆料；带「历史新低/低于双11/30天新低」标签
- 注意：直连 search.aspx 类旧参数会报 "读取错误,请重试!"，用 s.manmanbuy.com 的 result 接口

## 备选 1：什么值得买 SMZDM

- 站内搜索：`https://search.smzdm.com/?c=home&s=<关键词>` 或 `https://search.smzdm.com/?s=<关键词>`
- 特点：爆料流 + 值不值投票 + 评论；京东/天猫为主，PDD 爆料较少
- 兜底触发条件：慢慢买无结果或页面反爬时

## 备选 2：惠惠购物助手（网易）

- 搜索：`https://search.huihui.cn/search?q=<关键词>`
- 特点：直连京东/天猫比价，历史价走势完整
- 注意：商品覆盖面比慢慢买窄，适合单品历史价深挖

## 各官方平台现状（勿浪费 token 重试）

| 平台 | 未登录搜索 | 结论 |
|---|---|---|
| 京东 search.jd.com | 验证 → passport 登录页 | 不可用 |
| 淘宝 s.taobao.com | 页面开但结果区"加载中"空转 | 不可用 |
| 拼多多 mobile.yangkeduo.com | 直接跳 login.html | 不可用 |

## 经验备注

- 慢慢买搜索框提交：找包裹文本框的 form（跳过 __VIEWSTATE 隐藏域）直接 submit()
- 手机版页面是 ASP.NET WebForms，别试图手工拼 VIEWSTATE
- 抓取频率别太猛，单次查询间隔 ≥2 秒，避免触发聚合站自身风控