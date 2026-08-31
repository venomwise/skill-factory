#!/usr/bin/env node

import { readFile } from "node:fs/promises";
import { resolve } from "node:path";

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
  console.error("Usage: check-design-doc.mjs [--max 120] <design.md> [...]");
  process.exit(2);
}

const REQUIRED_SECTIONS = [
  "Summary",
  "Goals",
  "Non-Goals",
  "Context",
  "Proposed Solution",
  "Error Handling",
  "Acceptance Criteria",
  "Testing",
  "Open Questions",
];

const AC_TOKEN = /\bAC-[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*/g;
const AC_FORMAT = /^AC-[a-z0-9]+(?:-[a-z0-9]+)+$/;
const AC_DEFINITION = /^\s*-\s*\*\*(AC-[A-Za-z0-9-]+)\*\*/;
const DR_TOKEN = /\bDR-[A-Za-z0-9-]+/g;
const DR_FORMAT = /^DR-[a-z0-9]+(?:-[a-z0-9]+)*$/;
const DR_HEADING = /^#{3,4}\s+(DR-[A-Za-z0-9-]+)/;
const IF_TOKEN = /\bIF-(?:HTTP|MSG|PROTO)-[A-Za-z0-9-]+/g;
const CHANGE_TAG_HEADING = /^#{4,6}\s*\[(?:ADD|MODIFY|REMOVE)\]/;

function isStandaloneException(line) {
  const value = line.trim();
  return /^https?:\/\/\S+$/.test(value) || /^[0-9a-f]{40,}$/i.test(value);
}

function headingInfo(line) {
  const match = /^(#{1,6})\s+(.*)$/.exec(line);
  if (!match) return null;
  return { level: match[1].length, title: match[2].trim() };
}

function parseDocument(lines) {
  const active = lines.map(() => true);
  let fenced = false;
  let inComment = false;

  for (let index = 0; index < lines.length; index += 1) {
    const line = lines[index];
    if (/^\s*(```|~~~)/.test(line)) {
      fenced = !fenced;
      active[index] = false;
      continue;
    }
    if (fenced) {
      active[index] = false;
      continue;
    }
    if (inComment) {
      active[index] = false;
      if (line.includes("-->")) inComment = false;
      continue;
    }
    if (line.includes("<!--")) {
      active[index] = false;
      if (!line.includes("-->")) inComment = true;
    }
  }
  return active;
}

function findSection(lines, active, title, levels = [2]) {
  let start = -1;
  let level = 0;
  for (let index = 0; index < lines.length; index += 1) {
    if (!active[index]) continue;
    const heading = headingInfo(lines[index]);
    if (!heading) continue;
    if (start === -1) {
      if (levels.includes(heading.level) && heading.title === title) {
        start = index;
        level = heading.level;
      }
      continue;
    }
    if (heading.level <= level) return { start, end: index };
  }
  if (start === -1) return null;
  return { start, end: lines.length };
}

function checkFile(file, text) {
  const lines = text.split(/\r?\n/);
  const active = parseDocument(lines);
  const problems = [];
  const report = (lineNo, check, message) => {
    problems.push(`${file}:${lineNo}: ${check}: ${message}`);
  };

  for (let index = 0; index < lines.length; index += 1) {
    const line = lines[index];
    if (/^\s*(```|~~~)/.test(line)) continue;
    if (!active[index] || isStandaloneException(line)) continue;
    const width = [...line].length;
    if (width > max) {
      report(index + 1, "line-width", `line width ${width} exceeds ${max}`);
    }
  }

  const sectionLines = [];
  for (let index = 0; index < lines.length; index += 1) {
    if (!active[index]) continue;
    const heading = headingInfo(lines[index]);
    if (heading && heading.level === 2) sectionLines.push([index, heading.title]);
  }
  const presentSections = new Set(sectionLines.map(([, title]) => title));
  for (const section of REQUIRED_SECTIONS) {
    if (!presentSections.has(section)) {
      report(1, "section", `missing required section "## ${section}"`);
    }
  }

  const acSection = findSection(lines, active, "Acceptance Criteria");
  const acDefinitions = new Map();
  if (acSection) {
    for (let index = acSection.start; index < acSection.end; index += 1) {
      if (!active[index]) continue;
      const match = AC_DEFINITION.exec(lines[index]);
      if (!match) continue;
      const id = match[1];
      if (!AC_FORMAT.test(id)) {
        report(index + 1, "ac-format", `${id} does not match AC-<domain>-<behavior> kebab-case`);
      }
      if (acDefinitions.has(id)) {
        report(index + 1, "ac-duplicate", `${id} defined again; first definition at line ${acDefinitions.get(id)}`);
      } else {
        acDefinitions.set(id, index + 1);
      }
      let block = lines[index];
      for (let next = index + 1; next < acSection.end; next += 1) {
        const candidate = lines[next];
        if (!active[next]) break;
        if (/^\s*-\s*\*\*/.test(candidate) || /^#{1,6}\s/.test(candidate)) break;
        block += `\n${candidate}`;
      }
      const missing = [];
      if (!/\b(WHEN|IF)\b/.test(block)) missing.push("WHEN/IF");
      if (!/\bTHEN\b/.test(block)) missing.push("THEN");
      if (!/\bSHALL\b/.test(block)) missing.push("SHALL");
      if (missing.length > 0) {
        report(index + 1, "ac-structure", `${id} rule is missing ${missing.join(", ")}`);
      }
    }
  } else {
    report(1, "section", "cannot check AC definitions without an Acceptance Criteria section");
  }

  const drDefinitions = new Map();
  for (let index = 0; index < lines.length; index += 1) {
    if (!active[index]) continue;
    const match = DR_HEADING.exec(lines[index]);
    if (!match) continue;
    const id = match[1];
    if (!DR_FORMAT.test(id)) {
      report(index + 1, "dr-format", `${id} is not a lowercase kebab-case DR identifier`);
    }
    if (/\brevised\b/i.test(lines[index].slice(match.index + id.length))) {
      report(index + 1, "dr-revised", `${id} heading carries a Revised marker; update the original Decision in place`);
    }
    if (drDefinitions.has(id)) {
      report(index + 1, "dr-duplicate", `${id} defined again; first definition at line ${drDefinitions.get(id)}`);
    } else {
      drDefinitions.set(id, index + 1);
    }
  }

  const interfacesSection = findSection(lines, active, "Interfaces", [3, 4]);
  const interfaceDefinitions = new Set();
  const summaryContracts = new Set();
  const detailContracts = new Map();
  if (interfacesSection) {
    let summaryStart = -1;
    let summaryEnd = interfacesSection.end;
    for (let index = interfacesSection.start; index < interfacesSection.end; index += 1) {
      if (!active[index]) continue;
      const heading = headingInfo(lines[index]);
      if (heading && heading.title === "Change Summary" && summaryStart === -1) {
        summaryStart = index;
        continue;
      }
      if (summaryStart !== -1 && heading && heading.level <= 4) {
        summaryEnd = index;
        break;
      }
    }
    for (let index = interfacesSection.start; index < interfacesSection.end; index += 1) {
      if (!active[index]) continue;
      for (const match of lines[index].matchAll(IF_TOKEN)) {
        interfaceDefinitions.add(match[0]);
        if (summaryStart !== -1 && index > summaryStart && index < summaryEnd && /^\s*\|/.test(lines[index])) {
          summaryContracts.add(match[0]);
        }
      }
      const heading = headingInfo(lines[index]);
      if (heading && CHANGE_TAG_HEADING.test(lines[index])) {
        for (const match of lines[index].matchAll(IF_TOKEN)) {
          detailContracts.set(match[0], index + 1);
        }
      }
    }
    if (summaryStart === -1 && detailContracts.size > 0) {
      report(interfacesSection.start + 1, "contract-summary", "Interfaces has detail contracts but no Change Summary subsection");
    }
    for (const id of summaryContracts) {
      if (!detailContracts.has(id)) {
        report(interfacesSection.start + 1, "contract-summary", `${id} listed in Change Summary but has no detail contract heading`);
      }
    }
    for (const [id, lineNo] of detailContracts) {
      if (!summaryContracts.has(id)) {
        report(lineNo, "contract-summary", `${id} has a detail contract heading but is missing from Change Summary`);
      }
    }
  }

  for (let index = 0; index < lines.length; index += 1) {
    if (!active[index]) continue;
    const line = lines[index];
    for (const match of line.matchAll(AC_TOKEN)) {
      const id = match[0];
      const rest = line.slice(match.index + id.length);
      if (rest.startsWith("*") || rest.startsWith("-*")) continue;
      if (!acDefinitions.has(id)) {
        report(index + 1, "ac-undefined", `${id} is referenced but never defined in Acceptance Criteria`);
      }
    }
    for (const match of line.matchAll(DR_TOKEN)) {
      const id = match[0];
      if (!drDefinitions.has(id)) {
        report(index + 1, "dr-undefined", `${id} is referenced but has no Decision heading`);
      }
    }
    for (const match of line.matchAll(IF_TOKEN)) {
      const id = match[0];
      if (!interfaceDefinitions.has(id)) {
        report(index + 1, "contract-undefined", `${id} is referenced outside Interfaces but never defined inside it`);
      }
    }
  }

  return problems;
}

let total = 0;
for (const input of inputs) {
  const file = resolve(input);
  const text = await readFile(file, "utf8");
  const problems = checkFile(file, text);
  for (const problem of problems) console.error(problem);
  total += problems.length;
}

if (total > 0) {
  console.error(`Found ${total} design document violation(s).`);
  process.exit(1);
}

console.log(`Checked ${inputs.length} design document(s); structure, IDs, contracts and line width passed.`);
