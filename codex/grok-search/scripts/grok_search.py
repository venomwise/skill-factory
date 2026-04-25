import argparse
import json
import os
import re
import sys
import time
import urllib.error
import urllib.request
from typing import Any


PLACEHOLDER_KEYS = {
    "YOUR_API_KEY",
    "YOUR_GROK_API_KEY",
    "API_KEY",
    "CHANGE_ME",
    "REPLACE_ME",
    "<YOUR_GROK_API_KEY>",
}
FAILOVER_STATUS_CODES = {401, 403, 429}
FAILOVER_TEXT_PATTERNS = [
    "rate limit",
    "rate_limit",
    "too many requests",
    "no available tokens",
    "token unavailable",
    "invalid api key",
    "api key invalid",
    "unauthorized",
    "forbidden",
    "quota",
    "credits",
    "billing",
    "exhaust",
    "usage limit",
]
MODE_SYSTEM_PROMPTS = {
    "general": (
        "You are a real-time web research assistant. Use live web search/browsing when answering. "
        "Return ONLY a single JSON object with keys: content (string), sources (array of objects with url/title/snippet when possible). "
        "Keep content concise and evidence-backed."
    ),
    "news": (
        "You are a breaking-news research assistant. Prioritize the freshest reliable web information, recent developments, dates, and what changed. "
        "Return ONLY a single JSON object with keys: content (string), sources (array of objects with url/title/snippet when possible)."
    ),
    "social": (
        "You are a social/discourse research assistant. Focus on what people are saying now across live web sources, especially social and community discussion when available. "
        "Return ONLY a single JSON object with keys: content (string), sources (array of objects with url/title/snippet when possible)."
    ),
    "research": (
        "You are a multi-source research assistant. Use live web search/browsing to synthesize the most relevant viewpoints with evidence. "
        "Return ONLY a single JSON object with keys: content (string), sources (array of objects with url/title/snippet when possible)."
    ),
    "docs-compare": (
        "You are a research assistant comparing official documentation with community interpretation and recent discussion. "
        "Use live web search/browsing. In the content field, produce four short labeled sections exactly in this order: "
        "Official docs:, Community interpretation:, Agreement/conflict:, Bottom line:. "
        "Treat official documentation as the source of factual claims. Treat community discussion as interpretation, speculation, or operational experience unless it is directly supported by official docs. "
        "When the two disagree, say so explicitly. If official docs are missing or ambiguous, say that clearly instead of pretending certainty. "
        "Return ONLY a single JSON object with keys: content (string), sources (array of objects with url,title,snippet when possible)."
    ),
}
DEFAULT_COOLDOWN = {
    "enabled": True,
    "state_file": "runtime/cooldowns.json",
    "default_minutes": 15,
    "rate_limit_minutes": 20,
    "quota_minutes": 60,
    "auth_minutes": 360,
}


def _compact_json(data: Any) -> str:
    return json.dumps(data, ensure_ascii=False, separators=(",", ":"), sort_keys=False)


def _skill_root() -> str:
    return os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir))


def _default_config_paths() -> list[str]:
    root = _skill_root()
    home = os.path.expanduser("~")
    return [
        os.path.join(home, ".codex", "config", "grok-search.json"),
        os.path.join(root, "config.json"),
        os.path.join(root, "config.local.json"),
    ]


def _load_json_file(path: str) -> dict[str, Any]:
    try:
        with open(path, "r", encoding="utf-8-sig") as f:
            value = json.load(f)
    except FileNotFoundError:
        return {}
    if not isinstance(value, dict):
        raise ValueError(f"config at {path} must be a JSON object")
    return value


def _resolve_config(explicit_path: str) -> tuple[dict[str, Any], list[str], str]:
    if explicit_path:
        return _load_json_file(explicit_path), [explicit_path], explicit_path

    merged: dict[str, Any] = {}
    used: list[str] = []
    for path in _default_config_paths():
        if os.path.exists(path):
            merged.update(_load_json_file(path))
            used.append(path)
    primary = used[-1] if used else _default_config_paths()[-1]
    return merged, used, primary


def _normalize_api_key(api_key: str | None) -> str:
    key = (api_key or "").strip()
    if not key or key.upper() in PLACEHOLDER_KEYS:
        return ""
    return key


def _normalize_base_url_value(base_url: str | None) -> str:
    value = (base_url or "").strip()
    if not value:
        return ""
    placeholder = {
        "https://your-grok-endpoint.example",
        "YOUR_BASE_URL",
        "BASE_URL",
        "CHANGE_ME",
        "REPLACE_ME",
    }
    if value.upper() in placeholder:
        return ""
    return value


def _normalize_base_url(base_url: str) -> str:
    base_url = base_url.strip().rstrip("/")
    if base_url.endswith("/v1"):
        return base_url[: -len("/v1")]
    return base_url


def _coerce_json_object(text: str) -> dict[str, Any] | None:
    text = text.strip()
    if not text:
        return None
    if text.startswith("{") and text.endswith("}"):
        try:
            value = json.loads(text)
            return value if isinstance(value, dict) else None
        except json.JSONDecodeError:
            return None
    return None


def _extract_urls(text: str) -> list[str]:
    urls = re.findall(r"https?://[^\s)\]}>\"']+", text)
    seen: set[str] = set()
    out: list[str] = []
    for url in urls:
        url = url.rstrip(".,;:!?\'\"")
        if url and url not in seen:
            seen.add(url)
            out.append(url)
    return out


def _load_json_env(var_name: str) -> dict[str, Any]:
    raw = os.environ.get(var_name, "").strip()
    if not raw:
        return {}
    value = json.loads(raw)
    if not isinstance(value, dict):
        raise ValueError(f"{var_name} must be a JSON object")
    return value


def _parse_json_object(raw: str, *, label: str) -> dict[str, Any]:
    raw = raw.strip()
    if not raw:
        return {}
    value = json.loads(raw)
    if not isinstance(value, dict):
        raise ValueError(f"{label} must be a JSON object")
    return value


def _shorten(text: str | None, limit: int = 240) -> str:
    if not text:
        return ""
    text = text.strip()
    if len(text) <= limit:
        return text
    return text[: limit - 3] + "..."


def _format_ts(ts: float | int) -> str:
    return time.strftime("%Y-%m-%d %H:%M:%S UTC", time.gmtime(float(ts)))


def _normalize_cooldown_config(config: dict[str, Any]) -> dict[str, Any]:
    raw = config.get("cooldown")
    data = raw if isinstance(raw, dict) else {}
    merged = dict(DEFAULT_COOLDOWN)
    merged.update(data)
    return {
        "enabled": bool(merged.get("enabled", True)),
        "state_file": str(merged.get("state_file") or merged.get("stateFile") or DEFAULT_COOLDOWN["state_file"]),
        "default_seconds": int(float(merged.get("default_minutes", merged.get("defaultMinutes", DEFAULT_COOLDOWN["default_minutes"]))) * 60),
        "rate_limit_seconds": int(float(merged.get("rate_limit_minutes", merged.get("rateLimitMinutes", DEFAULT_COOLDOWN["rate_limit_minutes"]))) * 60),
        "quota_seconds": int(float(merged.get("quota_minutes", merged.get("quotaMinutes", DEFAULT_COOLDOWN["quota_minutes"]))) * 60),
        "auth_seconds": int(float(merged.get("auth_minutes", merged.get("authMinutes", DEFAULT_COOLDOWN["auth_minutes"]))) * 60),
    }


def _resolve_state_path(path_value: str) -> str:
    if os.path.isabs(path_value):
        return path_value
    return os.path.join(_skill_root(), path_value)


def _load_cooldown_state(path: str) -> dict[str, Any]:
    try:
        with open(path, "r", encoding="utf-8") as f:
            value = json.load(f)
    except FileNotFoundError:
        return {"profiles": {}}
    if not isinstance(value, dict):
        return {"profiles": {}}
    profiles = value.get("profiles")
    if not isinstance(profiles, dict):
        value["profiles"] = {}
    return value


def _save_cooldown_state(path: str, state: dict[str, Any]) -> None:
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w", encoding="utf-8") as f:
        json.dump(state, f, ensure_ascii=False, indent=2)


def _prune_cooldown_state(state: dict[str, Any], now: float) -> bool:
    profiles = state.setdefault("profiles", {})
    changed = False
    for key in list(profiles.keys()):
        item = profiles.get(key)
        if not isinstance(item, dict):
            profiles.pop(key, None)
            changed = True
            continue
        until = float(item.get("until", 0) or 0)
        if until <= now:
            profiles.pop(key, None)
            changed = True
    return changed


def _get_active_cooldown(state: dict[str, Any], profile_id: str, now: float) -> dict[str, Any] | None:
    profiles = state.get("profiles") or {}
    item = profiles.get(profile_id)
    if not isinstance(item, dict):
        return None
    until = float(item.get("until", 0) or 0)
    if until <= now:
        return None
    return item


def _clear_profile_cooldown(state: dict[str, Any], profile_id: str) -> bool:
    profiles = state.setdefault("profiles", {})
    if profile_id in profiles:
        profiles.pop(profile_id, None)
        return True
    return False


def _cooldown_seconds_for_failure(status_code: int | None, detail: str, cooldown: dict[str, Any]) -> int:
    if not cooldown.get("enabled"):
        return 0
    lower = (detail or "").lower()
    if status_code in {401, 403} or any(token in lower for token in ["invalid api key", "authentication_error", "unauthorized", "forbidden"]):
        return cooldown["auth_seconds"]
    if any(token in lower for token in ["quota", "credits", "billing", "usage limit", "exhaust", "insufficient"]):
        return cooldown["quota_seconds"]
    if status_code == 429 or any(token in lower for token in ["rate limit", "rate_limit", "too many requests", "no available tokens", "token unavailable"]):
        return cooldown["rate_limit_seconds"]
    return cooldown["default_seconds"]


def _set_profile_cooldown(
    state: dict[str, Any],
    profile_id: str,
    *,
    seconds: int,
    reason: str,
    status: int | None,
    now: float,
) -> dict[str, Any]:
    until = now + max(0, seconds)
    entry = {
        "until": until,
        "untilText": _format_ts(until),
        "seconds": seconds,
        "reason": reason,
        "status": status,
        "setAt": now,
        "setAtText": _format_ts(now),
    }
    state.setdefault("profiles", {})[profile_id] = entry
    return entry


def _extract_profiles(
    config: dict[str, Any],
    *,
    cli_api_key: str,
    env_api_key: str,
    env_api_keys: list[str],
    forced_profile: str,
) -> list[dict[str, Any]]:
    def apply_filter(items: list[dict[str, Any]]) -> list[dict[str, Any]]:
        if not forced_profile:
            return items
        return [p for p in items if p.get("id") == forced_profile]

    if cli_api_key:
        return [{"id": "cli", "api_key": cli_api_key, "source": "--api-key", "base_url": "", "model": ""}]

    if env_api_keys:
        return apply_filter([
            {"id": f"env-{i}", "api_key": key, "source": "GROK_API_KEYS", "base_url": "", "model": ""}
            for i, key in enumerate(env_api_keys, start=1)
            if key
        ])

    if env_api_key:
        return apply_filter([
            {"id": "env", "api_key": env_api_key, "source": "GROK_API_KEY", "base_url": "", "model": ""}
        ])

    profiles: list[dict[str, Any]] = []
    raw_profiles = config.get("profiles")
    if isinstance(raw_profiles, list) and raw_profiles:
        for i, item in enumerate(raw_profiles, start=1):
            if isinstance(item, str):
                api_key = _normalize_api_key(item)
                if api_key:
                    profiles.append({
                        "id": f"profile-{i}",
                        "api_key": api_key,
                        "source": "config.profiles",
                        "base_url": "",
                        "model": "",
                    })
                continue
            if not isinstance(item, dict):
                continue
            if item.get("enabled", True) is False:
                continue
            api_key = _normalize_api_key(str(item.get("api_key") or item.get("key") or ""))
            if not api_key:
                continue
            profiles.append({
                "id": str(item.get("id") or f"profile-{i}"),
                "api_key": api_key,
                "source": "config.profiles",
                "base_url": _normalize_base_url_value(str(item.get("base_url") or "")),
                "model": str(item.get("model") or "").strip(),
            })
        return apply_filter(profiles)

    raw_api_keys = config.get("api_keys")
    if isinstance(raw_api_keys, list) and raw_api_keys:
        profiles = [
            {"id": f"key-{i}", "api_key": _normalize_api_key(str(key)), "source": "config.api_keys", "base_url": "", "model": ""}
            for i, key in enumerate(raw_api_keys, start=1)
        ]
        profiles = [p for p in profiles if p.get("api_key")]
        return apply_filter(profiles)

    single_key = _normalize_api_key(str(config.get("api_key") or ""))
    if single_key:
        return apply_filter([
            {
                "id": str(config.get("profile_id") or "main"),
                "api_key": single_key,
                "source": "config.api_key",
                "base_url": "",
                "model": "",
            }
        ])

    return []


def _should_failover(status_code: int | None, detail: str) -> bool:
    if status_code in FAILOVER_STATUS_CODES:
        return True
    lower = (detail or "").lower()
    return any(token in lower for token in FAILOVER_TEXT_PATTERNS)


def _request_chat_completions(
    *,
    base_url: str,
    api_key: str,
    model: str,
    query: str,
    timeout_seconds: float,
    extra_headers: dict[str, Any],
    extra_body: dict[str, Any],
    system_prompt: str,
) -> dict[str, Any]:
    url = f"{_normalize_base_url(base_url)}/v1/chat/completions"

    body: dict[str, Any] = {
        "model": model,
        "messages": [
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": query},
        ],
        "temperature": 0.2,
        "stream": False,
    }
    body.update(extra_body)

    headers: dict[str, str] = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {api_key}",
        "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
    }
    for key, value in extra_headers.items():
        headers[str(key)] = str(value)

    req = urllib.request.Request(
        url=url,
        data=_compact_json(body).encode("utf-8"),
        headers=headers,
        method="POST",
    )
    with urllib.request.urlopen(req, timeout=timeout_seconds) as resp:
        raw = resp.read().decode("utf-8", errors="replace")
        
        # Handle SSE (Server-Sent Events) streaming response
        if raw.startswith('data: '):
            content_parts = []
            for line in raw.strip().split('\n'):
                if line.startswith('data: '):
                    data_str = line[6:]  # Remove 'data: ' prefix
                    if data_str == '[DONE]':
                        break
                    try:
                        chunk = json.loads(data_str)
                        if 'choices' in chunk and chunk['choices']:
                            delta = chunk['choices'][0].get('delta', )
                            if 'content' in delta:
                                content_parts.append(delta['content'])
                    except json.JSONDecodeError:
                        pass
            
            # Return in OpenAI format
            return {
                "choices": [
                    {
                        "message": {
                            "role": "assistant",
                            "content": ''.join(content_parts)
                        }
                    }
                ]
            }
        
        # Handle regular JSON response
        return json.loads(raw)


def _emit(data: dict[str, Any], *, plain: bool = False, urls_only: bool = False) -> int:
    if urls_only:
        for item in data.get("sources", []):
            url = item.get("url")
            if url:
                print(url)
        return 0

    if plain:
        if data.get("profileId"):
            source = data.get("profileSource") or ""
            suffix = f" [{source}]" if source else ""
            print(f"Profile: {data['profileId']}{suffix}")
        attempts = data.get("attempts") or []
        if attempts:
            print("Attempts:")
            for attempt in attempts:
                parts = [f"- {attempt.get('profileId', 'unknown')}"]
                if attempt.get("cooldown"):
                    parts.append("cooldown")
                    if attempt.get("remainingSeconds") is not None:
                        parts.append(f"remaining={attempt['remainingSeconds']}s")
                    if attempt.get("untilText"):
                        parts.append(f"until={attempt['untilText']}")
                else:
                    parts.append("ok" if attempt.get("ok") else "fail")
                if attempt.get("status") is not None:
                    parts.append(f"status={attempt['status']}")
                if attempt.get("failover"):
                    parts.append("failover")
                if attempt.get("detail"):
                    parts.append(_shorten(str(attempt.get("detail")), 160))
                print(" | ".join(parts))
            print()

        if data.get("error"):
            print(f"ERROR: {data['error']}")
            if data.get("detail"):
                print(data["detail"])
            return 0

        content = data.get("content") or ""
        if content:
            print(content)
            print()
        sources = data.get("sources") or []
        if sources:
            print("Sources:")
            for idx, item in enumerate(sources, start=1):
                print(f"[{idx}] {item.get('title') or ''}")
                print(f"URL: {item.get('url') or ''}")
                snippet = item.get("snippet") or ""
                if snippet:
                    print(f"Snippet: {_shorten(snippet, 220)}")
                print()
        return 0

    print(json.dumps(data, ensure_ascii=False, indent=2))
    return 0


def main() -> int:
    parser = argparse.ArgumentParser(description="Aggressive real-time research via OpenAI-compatible Grok endpoint")
    parser.add_argument("--query", required=True, help="Search query / research task")
    parser.add_argument("--config", default="", help="Path to config JSON file")
    parser.add_argument("--base-url", default="", help="Override base URL")
    parser.add_argument("--api-key", default="", help="Override API key (single key only)")
    parser.add_argument("--model", default="", help="Override model")
    parser.add_argument("--timeout-seconds", type=float, default=0.0, help="Override timeout (seconds)")
    parser.add_argument("--mode", choices=["general", "news", "social", "research", "docs-compare"], default="general")
    parser.add_argument("--plain", action="store_true", help="Human-readable output")
    parser.add_argument("--urls", action="store_true", help="Print source URLs only")
    parser.add_argument("--profile", default="", help="Force one configured profile id")
    parser.add_argument("--ignore-cooldown", action="store_true", help="Attempt requests even if a profile is currently cooling down")
    parser.add_argument(
        "--extra-body-json",
        default="",
        help="Extra JSON object merged into request body",
    )
    parser.add_argument(
        "--extra-headers-json",
        default="",
        help="Extra JSON object merged into request headers",
    )
    args = parser.parse_args()

    env_config_path = os.environ.get("GROK_CONFIG_PATH", "").strip()
    explicit_config_path = args.config.strip() or env_config_path

    try:
        config, config_paths, primary_config_path = _resolve_config(explicit_config_path)
    except Exception as e:
        sys.stderr.write(f"Invalid config: {e}\n")
        return 2

    cli_api_key = _normalize_api_key(args.api_key)
    env_api_key = _normalize_api_key(os.environ.get("GROK_API_KEY", ""))
    env_api_keys = [
        key for key in (_normalize_api_key(p) for p in os.environ.get("GROK_API_KEYS", "").split(",")) if key
    ]
    forced_profile = args.profile.strip()
    profiles = _extract_profiles(
        config,
        cli_api_key=cli_api_key,
        env_api_key=env_api_key,
        env_api_keys=env_api_keys,
        forced_profile=forced_profile,
    )

    base_url = _normalize_base_url_value(
        args.base_url.strip() or os.environ.get("GROK_BASE_URL", "").strip() or str(config.get("base_url") or "").strip()
    )
    model = args.model.strip() or os.environ.get("GROK_MODEL", "").strip() or str(config.get("model") or "").strip() or "grok-2-latest"

    timeout_seconds = args.timeout_seconds
    if not timeout_seconds:
        timeout_seconds = float(os.environ.get("GROK_TIMEOUT_SECONDS", "0") or "0")
    if not timeout_seconds:
        timeout_seconds = float(config.get("timeout_seconds") or 0) or 60.0

    cooldown_cfg = _normalize_cooldown_config(config)
    cooldown_state_path = _resolve_state_path(cooldown_cfg["state_file"])
    cooldown_state = _load_cooldown_state(cooldown_state_path)
    now = time.time()
    if _prune_cooldown_state(cooldown_state, now):
        _save_cooldown_state(cooldown_state_path, cooldown_state)

    if not base_url:
        return _emit({
            "ok": False,
            "error": "missing_base_url",
            "detail": "Set GROK_BASE_URL, write it to config, or pass --base-url",
            "config_path": primary_config_path,
            "config_paths": config_paths,
        }, plain=args.plain, urls_only=args.urls)

    if not profiles:
        detail = "Pass --api-key, set GROK_API_KEY/GROK_API_KEYS, or configure config.local.json/config.json"
        if forced_profile:
            detail += f" (requested profile: {forced_profile})"
        return _emit({
            "ok": False,
            "error": "missing_api_key",
            "detail": detail,
            "config_path": primary_config_path,
            "config_paths": config_paths,
        }, plain=args.plain, urls_only=args.urls)

    try:
        extra_body: dict[str, Any] = {}
        cfg_extra_body = config.get("extra_body")
        if isinstance(cfg_extra_body, dict):
            extra_body.update(cfg_extra_body)
        extra_body.update(_load_json_env("GROK_EXTRA_BODY_JSON"))
        extra_body.update(_parse_json_object(args.extra_body_json, label="--extra-body-json"))

        extra_headers: dict[str, Any] = {}
        cfg_extra_headers = config.get("extra_headers")
        if isinstance(cfg_extra_headers, dict):
            extra_headers.update(cfg_extra_headers)
        extra_headers.update(_load_json_env("GROK_EXTRA_HEADERS_JSON"))
        extra_headers.update(_parse_json_object(args.extra_headers_json, label="--extra-headers-json"))
    except Exception as e:
        return _emit({
            "ok": False,
            "error": "invalid_json",
            "detail": str(e),
            "config_path": primary_config_path,
            "config_paths": config_paths,
        }, plain=args.plain, urls_only=args.urls)

    started = time.time()
    attempts: list[dict[str, Any]] = []
    last_error: dict[str, Any] | None = None
    system_prompt = MODE_SYSTEM_PROMPTS.get(args.mode, MODE_SYSTEM_PROMPTS["general"])
    saved_state = False

    for idx, profile in enumerate(profiles):
        profile_id = profile["id"]
        now = time.time()
        active_cooldown = _get_active_cooldown(cooldown_state, profile_id, now)
        if active_cooldown and not args.ignore_cooldown:
            remaining = max(1, int(float(active_cooldown.get("until", now)) - now))
            attempts.append({
                "profileId": profile_id,
                "ok": False,
                "cooldown": True,
                "remainingSeconds": remaining,
                "untilText": active_cooldown.get("untilText") or _format_ts(float(active_cooldown.get("until", now))),
                "detail": active_cooldown.get("reason") or "cooldown active",
            })
            continue

        profile_base_url = _normalize_base_url_value(profile.get("base_url") or "") or base_url
        profile_model = profile.get("model") or model
        try:
            resp = _request_chat_completions(
                base_url=profile_base_url,
                api_key=profile["api_key"],
                model=profile_model,
                query=args.query,
                timeout_seconds=timeout_seconds,
                extra_headers=extra_headers,
                extra_body=extra_body,
                system_prompt=system_prompt,
            )
            attempts.append({"profileId": profile_id, "ok": True})
            if _clear_profile_cooldown(cooldown_state, profile_id):
                _save_cooldown_state(cooldown_state_path, cooldown_state)
                saved_state = True

            message = ""
            try:
                choice0 = (resp.get("choices") or [{}])[0]
                msg = choice0.get("message") or {}
                message = msg.get("content") or ""
            except Exception:
                message = ""

            parsed = _coerce_json_object(message)
            sources: list[dict[str, Any]] = []
            content = ""
            raw = ""

            if parsed is not None:
                content = str(parsed.get("content") or "")
                src = parsed.get("sources")
                if isinstance(src, list):
                    for item in src:
                        if isinstance(item, dict) and item.get("url"):
                            sources.append(
                                {
                                    "url": str(item.get("url")),
                                    "title": str(item.get("title") or ""),
                                    "snippet": str(item.get("snippet") or ""),
                                }
                            )
                if not sources:
                    for url in _extract_urls(content):
                        sources.append({"url": url, "title": "", "snippet": ""})
            else:
                raw = message
                for url in _extract_urls(message):
                    sources.append({"url": url, "title": "", "snippet": ""})

            out = {
                "ok": True,
                "mode": args.mode,
                "query": args.query,
                "profileId": profile_id,
                "profileSource": profile.get("source") or "",
                "attempts": attempts,
                "config_path": primary_config_path,
                "config_paths": config_paths,
                "base_url": profile_base_url,
                "model": resp.get("model") or profile_model,
                "content": content,
                "sources": sources,
                "raw": raw,
                "usage": resp.get("usage") or {},
                "elapsed_ms": int((time.time() - started) * 1000),
            }
            return _emit(out, plain=args.plain, urls_only=args.urls)
        except urllib.error.HTTPError as e:
            status_code = getattr(e, "code", None)
            detail = e.read().decode("utf-8", errors="replace") if hasattr(e, "read") else str(e)
            failover = _should_failover(status_code, detail)
            cooldown_seconds = _cooldown_seconds_for_failure(status_code, detail, cooldown_cfg) if failover else 0
            cooldown_entry = None
            if cooldown_seconds > 0:
                cooldown_entry = _set_profile_cooldown(
                    cooldown_state,
                    profile_id,
                    seconds=cooldown_seconds,
                    reason=detail,
                    status=status_code,
                    now=time.time(),
                )
                _save_cooldown_state(cooldown_state_path, cooldown_state)
                saved_state = True
            attempts.append({
                "profileId": profile_id,
                "ok": False,
                "status": status_code,
                "failover": failover,
                "cooldownSeconds": cooldown_seconds if cooldown_seconds > 0 else None,
                "untilText": cooldown_entry.get("untilText") if cooldown_entry else None,
                "detail": _shorten(detail),
            })
            last_error = {
                "error": f"http_{status_code or 'error'}",
                "detail": detail,
                "failover": failover,
            }
            if failover and idx < len(profiles) - 1:
                continue
            break
        except Exception as e:
            detail = str(e)
            attempts.append({
                "profileId": profile_id,
                "ok": False,
                "failover": False,
                "detail": _shorten(detail),
            })
            last_error = {
                "error": "request_failed",
                "detail": detail,
                "failover": False,
            }
            break

    all_in_cooldown = bool(attempts) and all(attempt.get("cooldown") for attempt in attempts)
    if all_in_cooldown:
        return _emit({
            "ok": False,
            "error": "all_profiles_in_cooldown",
            "detail": "All configured profiles are cooling down. Wait for cooldown expiry or retry with --ignore-cooldown.",
            "attempts": attempts,
            "cooldownStateFile": cooldown_state_path,
            "config_path": primary_config_path,
            "config_paths": config_paths,
            "base_url": base_url,
            "model": model,
            "elapsed_ms": int((time.time() - started) * 1000),
        }, plain=args.plain, urls_only=args.urls)

    if last_error:
        failover_exhausted = bool(last_error.get("failover")) and len([a for a in attempts if not a.get("cooldown")]) == len(profiles)
        error_code = "all_profiles_failed" if failover_exhausted and len(profiles) > 1 else str(last_error.get("error") or "request_failed")
        return _emit({
            "ok": False,
            "error": error_code,
            "detail": str(last_error.get("detail") or "Request failed"),
            "failoverExhausted": failover_exhausted,
            "attempts": attempts,
            "cooldownStateFile": cooldown_state_path if saved_state else None,
            "config_path": primary_config_path,
            "config_paths": config_paths,
            "base_url": base_url,
            "model": model,
            "elapsed_ms": int((time.time() - started) * 1000),
        }, plain=args.plain, urls_only=args.urls)

    return 1


if __name__ == "__main__":
    raise SystemExit(main())