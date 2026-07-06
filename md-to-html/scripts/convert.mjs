#!/usr/bin/env node
/**
 * Markdown → HTML 增强转换脚本
 * 在基础转换之上增加：内容分析、智能摘要、目录导航、关键信息高亮
 * 用法: node scripts/convert.mjs <markdown-file-path>
 * 输出: 同目录下同名 .html 文件
 */

import { existsSync, readFileSync, writeFileSync } from "fs";
import { Marked } from "marked";
import { basename, dirname, extname, resolve } from "path";

// ─── 参数解析 ─────────────────────────────────────────────────────────────────
const args = process.argv.slice(2);
if (args.length === 0) {
  console.error("❌ 请提供 Markdown 文件路径");
  console.error("   用法: node convert.mjs <markdown-file>");
  process.exit(1);
}

const inputPath = resolve(args[0]);
if (!existsSync(inputPath)) {
  console.error(`❌ 文件不存在: ${inputPath}`);
  process.exit(1);
}
if (!/\.(md|markdown)$/i.test(inputPath)) {
  console.error(`⚠️  文件不是 Markdown 格式: ${inputPath}`);
  process.exit(1);
}

const markdown = readFileSync(inputPath, "utf-8");

// ─── 内容分析引擎 ───────────────────────────────────────────────────────────
function analyzeDocument(md) {
  // 统一换行符为 \n，避免 Windows \r\n 干扰正则匹配
  md = md.replace(/\r\n/g, "\n");

  // 1. YAML frontmatter
  const meta = {};
  const fmMatch = md.match(/^---\n([\s\S]*?)\n---/);
  if (fmMatch) {
    fmMatch[1].split("\n").forEach(line => {
      const m = line.match(/^(\w[\w_-]*):\s*"?([^"]*)"?\s*$/);
      if (m) meta[m[1]] = m[2];
    });
  }

  // 2. 标题层级（用于 TOC）
  const headings = [];
  const lines = md.split("\n");
  let inCodeBlock = false;
  lines.forEach(line => {
    if (/^```/.test(line.trim())) { inCodeBlock = !inCodeBlock; return; }
    if (inCodeBlock) return;
    const m = line.match(/^(#{1,6})\s+(.+)$/);
    if (m) {
      const level = m[1].length;
      const text = m[2].replace(/[*_`\[\]]/g, "").replace(/~~([^~]+)~~/g, "$1").trim();
      const id = text
        .toLowerCase()
        .replace(/[^\w\u4e00-\u9fff]+/g, "-")
        .replace(/^-+|-+$/g, "");
      headings.push({ level, text, id });
    }
  });

  // 3. 统计数据
  const featureMatches = md.match(/\bF\d{2}\b/g) || [];
  const brMatches = md.match(/\bBR-\d{2}\b/g) || [];
  const acMatches = md.match(/\bAC-\d{2}\b/g) || [];
  const tableBlocks = md.split(/\n\n+/);
  const tableCount = tableBlocks.filter(block => {
    const lines = block.trim().split('\n');
    return lines.length >= 2 && /^\|/.test(lines[0]) && /^\|[-\s|:]+\|$/.test(lines[1].trim());
  }).length;
  const codeBlockCount = (md.match(/^```/gm) || []).length / 2;
  const wordCount = md.replace(/```[\s\S]*?```/g, "").replace(/\|.*\|/g, "")
    .replace(/[#*_`>\-\[\]()]/g, " ").split(/\s+/).filter(Boolean).length;

  // 4. 章节摘要（取 h2 下第一段非空文本）
  const sections = [];
  let currentH2 = null;
  let captureNext = false;
  lines.forEach(line => {
    if (/^##\s+/.test(line) && !/^###/.test(line)) {
      currentH2 = line.replace(/^##\s+/, "").trim();
      captureNext = true;
      return;
    }
    if (captureNext && line.trim() && !/^[#|>`\-]/.test(line.trim()) && !/^```/.test(line.trim())) {
      const summary = line.trim().replace(/[*_`]/g, "").slice(0, 120);
      sections.push({ title: currentH2, summary });
      captureNext = false;
    }
    if (line.trim() === "") captureNext = false;
  });

  // 5. 关键标签提取（BR/AC 去重统计）
  const uniqueBR = [...new Set(brMatches)];
  const uniqueAC = [...new Set(acMatches)];
  const uniqueFeatures = [...new Set(featureMatches)];

  return {
    meta, headings, tableCount, codeBlockCount: Math.floor(codeBlockCount),
    wordCount, sections, uniqueBR, uniqueAC, uniqueFeatures,
    brCount: uniqueBR.length, acCount: uniqueAC.length, featureCount: uniqueFeatures.length
  };
}

const analysis = analyzeDocument(markdown);

// ─── 文档质量预检引擎 ──────────────────────────────────────────────────────
function lintDocument(md) {
  const lines = md.replace(/\r\n/g, "\n").split("\n");
  const issues = [];
  let inCode = false;

  // 1. 未闭合代码块
  const fenceCount = (md.match(/^```/gm) || []).length;
  if (fenceCount % 2 !== 0) {
    issues.push({ severity: "error", type: "代码块", message: `代码块标记 \`\`\` 出现 ${fenceCount} 次（奇数），疑似未闭合` });
  }

  // 2. 标题层级跳跃
  let prevLevel = 0;
  lines.forEach((line, i) => {
    if (/^```/.test(line.trim())) { inCode = !inCode; return; }
    if (inCode) return;
    const m = line.match(/^(#{1,6})\s+(.+)$/);
    if (m) {
      const level = m[1].length;
      if (prevLevel > 0 && level > prevLevel + 1) {
        issues.push({ severity: "warning", type: "标题层级", line: i + 1, message: `第 ${i + 1} 行: h${prevLevel} 直接跳至 h${level}，建议逐级递进` });
      }
      prevLevel = level;
    }
  });

  // 3. 空链接 / 断链
  inCode = false;
  lines.forEach((line, i) => {
    if (/^```/.test(line.trim())) { inCode = !inCode; return; }
    if (inCode) return;
    // [text]() 或 [text](#)
    const emptyLinks = line.match(/\[([^\]]*)\]\(\s*#?\s*\)/g);
    if (emptyLinks) {
      emptyLinks.forEach(link => {
        issues.push({ severity: "warning", type: "空链接", line: i + 1, message: `第 ${i + 1} 行: 空链接 ${link}` });
      });
    }
  });

  // 4. 表格列数不一致
  inCode = false;
  let tableStart = -1;
  let headerCols = 0;
  lines.forEach((line, i) => {
    if (/^```/.test(line.trim())) { inCode = !inCode; return; }
    if (inCode) return;
    if (/^\|/.test(line.trim())) {
      const cols = line.split("|").filter(c => c.trim() !== "").length;
      if (tableStart === -1) {
        tableStart = i;
        headerCols = cols;
      } else if (i === tableStart + 1) {
        // separator line, skip
      } else if (cols !== headerCols && cols > 0) {
        issues.push({ severity: "warning", type: "表格", line: i + 1, message: `第 ${i + 1} 行: 表格列数 ${cols} 与表头 ${headerCols} 列不一致` });
      }
    } else {
      if (line.trim() === "" || !/^\|/.test(line.trim())) {
        tableStart = -1;
        headerCols = 0;
      }
    }
  });

  // 5. 重复编号（BR-xx / AC-xx）
  const brMap = {}, acMap = {};
  inCode = false;
  lines.forEach((line, i) => {
    if (/^```/.test(line.trim())) { inCode = !inCode; return; }
    if (inCode) return;
    (line.match(/\bBR-\d{2}\b/g) || []).forEach(id => {
      if (!brMap[id]) brMap[id] = [];
      brMap[id].push(i + 1);
    });
    (line.match(/\bAC-\d{2}\b/g) || []).forEach(id => {
      if (!acMap[id]) acMap[id] = [];
      acMap[id].push(i + 1);
    });
  });
  Object.entries(brMap).forEach(([id, lns]) => {
    if (lns.length > 1) {
      issues.push({ severity: "info", type: "重复编号", message: `${id} 出现 ${lns.length} 次（行 ${lns.join(", ")}）` });
    }
  });
  Object.entries(acMap).forEach(([id, lns]) => {
    if (lns.length > 1) {
      issues.push({ severity: "info", type: "重复编号", message: `${id} 出现 ${lns.length} 次（行 ${lns.join(", ")}）` });
    }
  });

  // 6. YAML frontmatter 缺失
  if (!md.startsWith("---")) {
    issues.push({ severity: "info", type: "Frontmatter", message: "文件缺少 YAML frontmatter（--- 开头）" });
  }

  // 7. 中英文混排缺少空格
  inCode = false;
  let spacingCount = 0;
  lines.forEach((line, i) => {
    if (/^```/.test(line.trim())) { inCode = !inCode; return; }
    if (inCode) return;
    if (line.match(/[\u4e00-\u9fff][A-Za-z0-9]|[A-Za-z0-9][\u4e00-\u9fff]/)) {
      // 排除代码内联 `xxx` 里的匹配
      const cleaned = line.replace(/`[^`]+`/g, "");
      if (cleaned.match(/[\u4e00-\u9fff][A-Za-z0-9]|[A-Za-z0-9][\u4e00-\u9fff]/)) {
        spacingCount++;
        if (spacingCount <= 5) {
          issues.push({ severity: "info", type: "排版", line: i + 1, message: `第 ${i + 1} 行: 中英文/数字之间建议添加空格` });
        }
      }
    }
  });
  if (spacingCount > 5) {
    issues.push({ severity: "info", type: "排版", message: `另有 ${spacingCount - 5} 处中英文混排缺少空格（仅展示前 5 处）` });
  }

  // 8. 连续空行过多（>3 行）
  let blankRun = 0;
  lines.forEach((line, i) => {
    if (line.trim() === "") {
      blankRun++;
      if (blankRun === 4) {
        issues.push({ severity: "info", type: "排版", line: i + 1, message: `第 ${i - 2}–${i + 1} 行: 连续 ${blankRun}+ 空行，建议精简` });
      }
    } else {
      blankRun = 0;
    }
  });

  // 汇总
  const errors = issues.filter(i => i.severity === "error").length;
  const warnings = issues.filter(i => i.severity === "warning").length;
  const infos = issues.filter(i => i.severity === "info").length;
  return { issues, errors, warnings, infos };
}

const lintResult = lintDocument(markdown);

// ─── 配置 Marked ──────────────────────────────────────────────────────────────
const marked = new Marked({ gfm: true, breaks: false, pedantic: false });

function escapeAttr(value) {
  return String(value ?? "")
    .replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;").replace(/'/g, "&#39;");
}

function embedLocalImage(href) {
  if (!href || /^(?:[a-z][a-z\d+.-]*:|\/\/|#)/i.test(href)) return href;

  const filePart = href.split(/[?#]/)[0];
  let decodedFilePart = filePart;
  try {
    decodedFilePart = decodeURIComponent(filePart);
  } catch {
    // Keep the original path when it is not URL-encoded.
  }

  const imagePath = resolve(dirname(inputPath), decodedFilePart);
  if (!existsSync(imagePath)) return href;

  const mimeByExt = {
    ".png": "image/png",
    ".jpg": "image/jpeg",
    ".jpeg": "image/jpeg",
    ".gif": "image/gif",
    ".webp": "image/webp",
    ".svg": "image/svg+xml",
  };
  const mime = mimeByExt[extname(imagePath).toLowerCase()];
  if (!mime) return href;

  return `data:${mime};base64,${readFileSync(imagePath).toString("base64")}`;
}

const renderer = {
  code({ text, lang: infostring }) {
    const code = text || "";
    const lang = (infostring || "").trim().split(/\s+/)[0];
    const escaped = code
      .replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;").replace(/'/g, "&#39;");
    const langLabel = lang ? `<div class="code-lang">${lang}</div>` : "";
    return `<div class="code-block">${langLabel}<pre><code class="language-${lang || "text"}">${escaped}</code></pre></div>`;
  },
  table({ header, rows }) {
    const headerHtml = header.map(c => `<th>${this.parser.parseInline(c.tokens)}</th>`).join("");
    const bodyHtml = rows.map(row =>
      "<tr>" + row.map(c => `<td>${this.parser.parseInline(c.tokens)}</td>`).join("") + "</tr>"
    ).join("\n");
    return `<div class="table-wrapper"><table>\n<thead>\n<tr>${headerHtml}</tr>\n</thead>\n<tbody>\n${bodyHtml}\n</tbody>\n</table></div>`;
  },
  heading({ tokens, text, depth }) {
    const html = (tokens && this.parser) ? this.parser.parseInline(tokens) : (text || '');
    const id = html.replace(/<[^>]*>/g, "")
      .toLowerCase().replace(/[^\w\u4e00-\u9fff]+/g, "-").replace(/^-+|-+$/g, "");
    return `<h${depth} id="${id}">${html}</h${depth}>`;
  },
  image({ href, title, text }) {
    const src = embedLocalImage(href);
    const titleAttr = title ? ` title="${escapeAttr(title)}"` : "";
    return `<img src="${escapeAttr(src)}" alt="${escapeAttr(text)}"${titleAttr}>`;
  },
};

marked.use({ renderer });

// ─── 高亮关键模式（在 HTML 转换后处理）────────────────────────────────────────
function highlightPatterns(html) {
  // BR-xx / AC-xx / F0x 标签高亮
  html = html.replace(/\b(BR-\d{2})\b/g, '<span class="tag tag-br">$1</span>');
  html = html.replace(/\b(AC-\d{2})\b/g, '<span class="tag tag-ac">$1</span>');
  html = html.replace(/\b(F\d{2})\b(?!\s*<\/)/g, '<span class="tag tag-feature">$1</span>');
  // Given/When/Then 关键字高亮
  html = html.replace(/^(Given|When|Then)\b/gm, '<strong class="keyword">$1</strong>');
  return html;
}

// 剥离 YAML frontmatter，避免被当作正文渲染（需先统一换行符）
const mdNormalized = markdown.replace(/\r\n/g, "\n");
const mdForConversion = mdNormalized.replace(/^---\n[\s\S]*?\n---\n?/, "");
let htmlBody = marked.parse(mdForConversion);
htmlBody = highlightPatterns(htmlBody);

// ─── 提取页面标题 ─────────────────────────────────────────────────────────────
const titleMatch = markdown.match(/^#\s+(.+)$/m);
const pageTitle = titleMatch ? titleMatch[1].replace(/[*_`]/g, "") : "Document";

// ─── 生成 TOC 侧边栏 HTML ────────────────────────────────────────────────────
function generateTOC(headings) {
  if (headings.length === 0) return "";
  // 只展示 h1-h3，h4+ 太多会臃肿
  const items = headings.filter(h => h.level <= 3);
  let html = '<nav class="toc" id="toc">';
  html += '<div class="toc-title">目录导航</div>';
  html += '<ul class="toc-list">';
  items.forEach(h => {
    const indent = `toc-level-${h.level}`;
    html += `<li class="${indent}"><a href="#${h.id}">${h.text}</a></li>`;
  });
  html += '</ul></nav>';
  return html;
}

const tocHtml = generateTOC(analysis.headings);

// ─── 语法高亮脚本（String.raw 保留正则反斜杠，避免模板字面量转义）────────────
const syntaxHighlightScript = String.raw`
// ─── 代码语法高亮（内联实现，无外部依赖）─────────────────────────────
(function(){
  var E={'&':'&amp;','<':'&lt;','>':'&gt;'};
  function esc(s){return s.replace(/[&<>]/g,function(c){return E[c]});}
  function hl(txt,rules){
    var hits=[];
    rules.forEach(function(p){
      p[0].lastIndex=0;
      var m;while((m=p[0].exec(txt))!==null)
        hits.push({s:m.index,e:m.index+m[0].length,c:p[1],t:m[0]});
    });
    hits.sort(function(a,b){return a.s-b.s;});
    var out='',pos=0;
    for(var i=0;i<hits.length;i++){
      var h=hits[i];
      if(h.s<pos)continue;
      if(h.s>pos)out+=esc(txt.slice(pos,h.s));
      out+='<span class="t'+h.c+'">'+esc(h.t)+'</span>';
      pos=h.e;
    }
    return out+esc(txt.slice(pos));
  }
  var R={
    json:[
      [/"(?:[^"\\]|\\.)*"(?=\s*:)/g,'k'],
      [/"(?:[^"\\]|\\.)*"/g,'s'],
      [/\b(true|false|null)\b/g,'b'],
      [/-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?/g,'n']
    ],
    python:[
      [/#[^\n]*/g,'c'],
      [/"""[\s\S]*?"""|'''[\s\S]*?'''|"(?:[^"\\]|\\.)*"|'(?:[^'\\]|\\.)*'/g,'s'],
      [/\b(import|from|def|class|if|elif|else|for|while|return|in|not|and|or|True|False|None|try|except|with|as|pass|lambda|self)\b/g,'b'],
      [/\b\d+(?:\.\d+)?\b/g,'n']
    ],
    java:[
      [/\/\/[^\n]*|\/\*[\s\S]*?\*\//g,'c'],
      [/"(?:[^"\\]|\\.)*"/g,'s'],
      [/@\w+/g,'a'],
      [/\b(public|private|protected|static|final|class|interface|new|return|void|if|else|for|while|try|catch|throw|throws|import|package|extends|implements|String|int|long|boolean|double|float|Integer|Long|List|Map|null|true|false)\b/g,'b'],
      [/\b\d+[LlFfDd]?\b/g,'n']
    ],
    http:[
      [/^(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\b/gm,'b'],
      [/^[\w-]+(?=:)/gm,'k'],
      [/HTTP\/[\d.]+/g,'n']
    ]
  };
  document.querySelectorAll('pre>code[class*="language-"]').forEach(function(b){
    var lang=(b.className.match(/language-(\w+)/)||[])[1];
    if(R[lang])b.innerHTML=hl(b.textContent,R[lang]);
    else if(lang==='text')b.innerHTML=hl(b.textContent,R.http);
  });
})();
`;

// ─── HTML 输出 ────────────────────────────────────────────────────────────────
const html = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<title>${pageTitle}</title>
<style>
:root {
  --bg:#ffffff; --fg:#1a1a2e; --muted:#6b7280; --border:#e5e7eb;
  --accent:#2563eb; --accent-bg:#eff6ff; --code-bg:#f3f4f6;
  --pre-bg:#1e293b; --pre-fg:#e2e8f0; --blockquote-border:#3b82f6;
  --blockquote-bg:#f0f9ff; --table-header-bg:#f8fafc; --table-stripe:#fafafa;
  --shadow:0 1px 3px rgba(0,0,0,.08); --radius:8px;
  --sidebar-w:270px;
}
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
html{font-size:16px;-webkit-font-smoothing:antialiased;scroll-behavior:smooth}
body{
  font-family:-apple-system,BlinkMacSystemFont,"Segoe UI","Noto Sans SC","PingFang SC","Microsoft YaHei",sans-serif;
  color:var(--fg);background:#f4f6f8;line-height:1.75;
}

/* ─── 布局 ────────────────────────────────────────────────────────────── */
.layout{display:flex;min-height:100vh}
.sidebar{
  width:var(--sidebar-w);flex-shrink:0;position:fixed;top:0;left:0;bottom:0;
  background:#f8fafc;border-right:1px solid var(--border);overflow-y:auto;z-index:10;
  padding:24px 0;
}
.sidebar::-webkit-scrollbar{width:4px}
.sidebar::-webkit-scrollbar-thumb{background:#cbd5e1;border-radius:4px}
.main{margin-left:var(--sidebar-w);flex:1;min-width:0;padding:40px 48px 80px}
.container{max-width:880px;margin:0 auto;background:var(--bg);border-radius:12px;
  box-shadow:0 4px 24px rgba(0,0,0,.07);padding:56px 64px}
@media(max-width:1080px){
  .sidebar{display:none} .main{margin-left:0;padding:24px 16px 60px}
  .container{width:100%;padding:32px 24px}
}

/* ─── TOC 侧边栏 ──────────────────────────────────────────────────────── */
.toc-title{
  font-size:.75rem;font-weight:700;text-transform:uppercase;letter-spacing:.08em;
  color:var(--muted);padding:0 20px 12px;border-bottom:1px solid var(--border);margin-bottom:8px;
}
.toc-list{list-style:none;padding:0}
.toc-list li{margin:0}
.toc-list a{
  display:block;padding:5px 20px;font-size:.82rem;color:#374151;
  text-decoration:none;border-left:3px solid transparent;transition:all .15s;
  line-height:1.5;word-break:break-word;
}
.toc-list a:hover{color:var(--accent);background:var(--accent-bg);border-left-color:var(--accent)}
.toc-level-1>a{font-weight:700;font-size:.88rem;padding-left:20px}
.toc-level-2>a{font-weight:600;padding-left:28px}
.toc-level-3>a{padding-left:42px;color:var(--muted);font-size:.79rem}

/* ─── 标题 ─────────────────────────────────────────────────────────────── */
h1,h2,h3,h4,h5,h6{color:var(--fg);font-weight:700;line-height:1.35;margin-top:2em;margin-bottom:.6em}
h1{font-size:2rem;margin-top:0;padding-bottom:.5em;border-bottom:2px solid var(--border)}
h2{font-size:1.55rem;padding-bottom:.35em;border-bottom:1px solid var(--border);scroll-margin-top:20px}
h3{font-size:1.25rem;scroll-margin-top:20px}
h4{font-size:1.1rem;scroll-margin-top:20px}
h5{font-size:1rem} h6{font-size:.9rem;color:var(--muted)}

/* ─── 段落 & 行内 ──────────────────────────────────────────────────────── */
p{margin:1em 0}
a{color:var(--accent);text-decoration:none;border-bottom:1px solid transparent;transition:border-color .15s}
a:hover{border-bottom-color:var(--accent)}
strong{font-weight:700} em{font-style:italic}
code{
  font-family:"JetBrains Mono","Fira Code",Consolas,monospace;font-size:.88em;
  background:var(--code-bg);color:#b91c1c;padding:2px 6px;border-radius:4px;
}

/* ─── 代码块 ─────────────────────────────────────────────────────────── */
.code-block{position:relative;margin:1.4em 0;border-radius:var(--radius);overflow:hidden;box-shadow:var(--shadow)}
.code-lang{
  position:absolute;top:0;right:0;background:rgba(255,255,255,.1);color:#94a3b8;
  font-size:.72rem;padding:3px 10px;border-radius:0 0 0 6px;font-family:sans-serif;
  text-transform:uppercase;letter-spacing:.05em;
}
pre{background:var(--pre-bg);color:var(--pre-fg);padding:20px 24px;overflow-x:auto;font-size:.88rem;line-height:1.65;margin:0}
pre code{background:none;color:inherit;padding:0;border-radius:0;font-size:inherit}

/* ─── 列表 ─────────────────────────────────────────────────────────────── */
ul,ol{padding-left:1.6em;margin:.8em 0} li{margin:.3em 0}
li>ul,li>ol{margin:.2em 0}
ul li input[type="checkbox"]{margin-right:.4em;accent-color:var(--accent)}

/* ─── 引用块 ─────────────────────────────────────────────────────────── */
blockquote{
  border-left:4px solid var(--blockquote-border);background:var(--blockquote-bg);
  padding:14px 20px;margin:1.4em 0;border-radius:0 var(--radius) var(--radius) 0;color:#374151;
}
blockquote p{margin:.4em 0} blockquote blockquote{margin:.8em 0}

/* ─── 表格 ─────────────────────────────────────────────────────────────── */
.table-wrapper{overflow-x:auto;margin:1.4em 0;border-radius:var(--radius);box-shadow:var(--shadow)}
table{width:100%;border-collapse:collapse;font-size:.92rem}
thead{background:var(--table-header-bg)}
th,td{padding:10px 16px;text-align:left;border:1px solid var(--border)}
th{font-weight:600;color:var(--fg)}
tbody tr:nth-child(even){background:var(--table-stripe)}
tbody tr:hover{background:var(--accent-bg)}

/* ─── 分隔线 & 图片 ──────────────────────────────────────────────────── */
hr{border:none;border-top:1.5px solid var(--border);margin:2.5em 0}
img{max-width:100%;height:auto;border-radius:var(--radius);box-shadow:var(--shadow);margin:1em 0}

/* ─── Details 折叠面板 ───────────────────────────────────────────────── */
.content details{
  border-bottom:1px solid var(--border);margin:0;
}
.content details summary{
  display:flex;align-items:center;justify-content:space-between;gap:20px;
  min-height:74px;padding:18px 4px 18px 0;cursor:pointer;list-style:none;
  color:#111827;font-size:1rem;font-weight:700;line-height:1.5;
}
.content details summary::-webkit-details-marker{display:none}
.content details summary::after{
  content:"";width:9px;height:9px;flex:0 0 auto;border-right:2px solid #111827;
  border-bottom:2px solid #111827;transform:rotate(45deg);transition:transform .18s ease;
}
.content details[open] summary::after{transform:rotate(225deg)}
.content details summary:hover{color:#0f172a}
.content details p,
.content details ul{
  margin:0;padding:0 34px 20px 0;color:#4b5563;
}
.content details ul{padding-left:1.4em}

/* ─── 关键标签高亮 ──────────────────────────────────────────────────── */
.tag{
  display:inline-flex;align-items:center;padding:1px 7px;border-radius:4px;
  font-size:.82em;font-weight:700;font-family:monospace;letter-spacing:.02em;
  vertical-align:middle;line-height:1.7;
}
.tag-br{background:#fee2e2;color:#991b1b;border:1px solid #fca5a5}
.tag-ac{background:#dcfce7;color:#166534;border:1px solid #86efac}
.tag-feature{background:#dbeafe;color:#1e40af;border:1px solid #93c5fd}
.keyword{color:#7c3aed;font-weight:700}

/* ─── 返回顶部 ─────────────────────────────────────────────────────────── */
.back-to-top{
  position:fixed;bottom:28px;right:28px;width:40px;height:40px;border-radius:50%;
  background:var(--accent);color:#fff;border:none;cursor:pointer;font-size:1.1rem;
  box-shadow:0 4px 12px rgba(37,99,235,.35);display:none;align-items:center;
  justify-content:center;transition:opacity .2s;z-index:20;
}
.back-to-top.visible{display:flex}

/* ─── 打印 ─────────────────────────────────────────────────────────────── */
@media print{
  body{background:#fff} .sidebar,.back-to-top{display:none!important}
  .main{margin-left:0} .container{box-shadow:none;max-width:100%;padding:20px}
  pre{white-space:pre-wrap;word-break:break-all}
  a{color:var(--fg);border-bottom:1px dashed var(--muted)}
}

/* ─── 语法高亮 token ─────────────────────────────────────────────────── */
.tk{color:#61afef}.ts{color:#98c379}.tn{color:#d19a66}
.tb{color:#c678dd}.tc{color:#5c6370;font-style:italic}.ta{color:#e06c75}
</style>
</head>
<body>
<div class="layout">
  ${tocHtml ? `<aside class="sidebar">${tocHtml}</aside>` : ""}
  <main class="main">
    <div class="container">
      <article class="content">
        ${htmlBody}
      </article>
    </div>
  </main>
</div>
<button class="back-to-top" id="backToTop" title="返回顶部">↑</button>
<script>
// 返回顶部按钮
const btn=document.getElementById('backToTop');
window.addEventListener('scroll',()=>{
  btn.classList.toggle('visible',window.scrollY>400);
},{ passive:true });
btn.addEventListener('click',()=>window.scrollTo({top:0,behavior:'smooth'}));

// TOC 当前高亮
const tocLinks=document.querySelectorAll('.toc-list a');
const sections=[];
tocLinks.forEach(a=>{
  const id=a.getAttribute('href').slice(1);
  const el=document.getElementById(id);
  if(el) sections.push({el,a});
});
let ticking=false;
window.addEventListener('scroll',()=>{
  if(!ticking){
    window.requestAnimationFrame(()=>{
      let current=null;
      sections.forEach(s=>{
        if(s.el.getBoundingClientRect().top<=80) current=s;
      });
      tocLinks.forEach(a=>a.classList.remove('active'));
      if(current) current.a.classList.add('active');
      ticking=false;
    });
    ticking=true;
  }
},{passive:true});

${syntaxHighlightScript}
</script>
</body>
</html>`;

// ─── 写入文件 ─────────────────────────────────────────────────────────────────
const dir = dirname(inputPath);
const name = basename(inputPath, extname(inputPath));
const outputPath = resolve(dir, `${name}.html`);
writeFileSync(outputPath, html, "utf-8");
console.log(`✅ 转换成功: ${outputPath}`);
console.log(`   📊 功能模块: ${analysis.featureCount} | 业务规则: ${analysis.brCount} | 验收标准: ${analysis.acCount} | 表格: ${analysis.tableCount}`);

// ─── 控制台输出预检报告 ──────────────────────────────────────────────────────
if (lintResult.issues.length === 0) {
  console.log(`   ✅ 文档质量检验: 全部通过`);
} else {
  console.log(`   🔍 文档质量检验: ${lintResult.errors} 错误 / ${lintResult.warnings} 警告 / ${lintResult.infos} 提示`);
  lintResult.issues.forEach(issue => {
    const icon = issue.severity === "error" ? "✗" : issue.severity === "warning" ? "⚠" : "ℹ";
    console.log(`      ${icon} [${issue.type}] ${issue.message}`);
  });
}
