#!/usr/bin/env node

import { readFile, readdir, stat } from "node:fs/promises";
import { extname, resolve } from "node:path";

const args = process.argv.slice(2);
let max = 120;
const inputs = [];

for (let index = 0; index < args.length; index += 1) {
  if (args[index] === "--max") {
    max = Number(args[index + 1]);
    index += 1;
    continue;
  }
  inputs.push(args[index]);
}

if (!Number.isInteger(max) || max < 1 || inputs.length === 0) {
  console.error("Usage: check-markdown-lines.mjs [--max 120] <file-or-directory> [...]");
  process.exit(2);
}

async function collect(input) {
  const path = resolve(input);
  const info = await stat(path);
  if (info.isFile()) {
    return extname(path).toLowerCase() === ".md" ? [path] : [];
  }

  const files = [];
  for (const entry of await readdir(path, { withFileTypes: true })) {
    if (entry.name === ".git" || entry.name === "node_modules") continue;
    files.push(...(await collect(resolve(path, entry.name))));
  }
  return files;
}

function isStandaloneException(line) {
  const value = line.trim();
  return /^https?:\/\/\S+$/.test(value) || /^[0-9a-f]{40,}$/i.test(value);
}

const files = (await Promise.all(inputs.map(collect))).flat().sort();
let violations = 0;

for (const file of files) {
  const lines = (await readFile(file, "utf8")).split(/\r?\n/);
  let fenced = false;

  for (let index = 0; index < lines.length; index += 1) {
    const line = lines[index];
    if (/^\s*(```|~~~)/.test(line)) {
      fenced = !fenced;
      continue;
    }
    if (fenced || isStandaloneException(line)) continue;

    const width = [...line].length;
    if (width > max) {
      console.error(`${file}:${index + 1}: line width ${width} exceeds ${max}`);
      violations += 1;
    }
  }
}

if (violations > 0) {
  console.error(`Found ${violations} Markdown line-width violation(s).`);
  process.exit(1);
}

console.log(`Checked ${files.length} Markdown file(s); no lines exceed ${max} characters.`);
