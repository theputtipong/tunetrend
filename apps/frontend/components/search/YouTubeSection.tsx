"use client";

import { useCallback, useEffect, useState } from "react";
import { dictionaries } from "@/lib/i18n";
import { formatLocalTime, formatRelativeTime, formatViewCount } from "@/lib/format";
import type { SearchItem, YouTubeSearchError } from "@/types/search";
import { useAuth } from "../auth/AuthProvider";
import { VideoTypeBadge } from "../VideoTypeBadge";
import { SearchItemRow } from "./SearchItemRow";

type State =
  | { status: "idle" }
  | { status: "loading" }
  | { status: "done"; items: SearchItem[]; remaining: number }
  | { status: "error"; error: YouTubeSearchError; retryAt?: string; remaining?: number };

const KNOWN_ERRORS = new Set<YouTubeSearchError>([
  "quota_exhausted",
  "user_limit",
  "search_unavailable",
  "youtube_error",
]);

async function fetchYouTube(
  q: string,
  country: string,
  category: string,
): Promise<State | "unauthorized"> {
  const params = new URLSearchParams({ q, country, category });
  let res: Response;
  let body: Record<string, unknown> | null = null;
  try {
    res = await fetch(`/api/search/youtube?${params.toString()}`, { cache: "no-store" });
    body = await res.json().catch(() => null);
  } catch {
    return { status: "error", error: "youtube_error" };
  }

  if (res.ok && body?.success) {
    const data = body.data as { items?: SearchItem[]; remaining?: number } | undefined;
    return { status: "done", items: data?.items ?? [], remaining: data?.remaining ?? 0 };
  }
  if (res.status === 401) return "unauthorized";
  const code = body?.error as YouTubeSearchError | undefined;
  return {
    status: "error",
    error: code && KNOWN_ERRORS.has(code) ? code : "youtube_error",
    retryAt: typeof body?.retryAt === "string" ? body.retryAt : undefined,
    remaining: typeof body?.remaining === "number" ? body.remaining : undefined,
  };
}

export function YouTubeSection({
  q,
  country,
  category,
  suggestYouTube,
  excludeIds,
}: Readonly<{
  q: string;
  country: string;
  category: string;
  suggestYouTube: boolean;
  excludeIds: string[];
}>) {
  const { user, lang, loginEnabled, openLogin, handleUnauthorized } = useAuth();
  const t = dictionaries[lang].search;
  const autoStart = Boolean(user) && suggestYouTube;
  const [state, setState] = useState<State>(() =>
    autoStart ? { status: "loading" } : { status: "idle" },
  );

  const apply = useCallback(
    (next: State | "unauthorized") => {
      if (next === "unauthorized") {
        handleUnauthorized();
        setState({ status: "idle" });
      } else {
        setState(next);
      }
    },
    [handleUnauthorized],
  );

  useEffect(() => {
    if (!autoStart) return;
    let cancelled = false;
    fetchYouTube(q, country, category).then((next) => {
      if (!cancelled) apply(next);
    });
    return () => {
      cancelled = true;
    };
  }, [autoStart, q, country, category, apply]);

  function start() {
    setState({ status: "loading" });
    void fetchYouTube(q, country, category).then(apply);
  }

  if (!user) {
    if (!loginEnabled) return null;
    return (
      <div className="yt-cta">
        <div>
          <b>{t.loginTitle}</b>
          {t.loginBody}
        </div>
        <button type="button" className="btn" onClick={() => openLogin("youtube")}>
          <span className="g-dot" aria-hidden="true" />
          {t.loginButton}
        </button>
      </div>
    );
  }

  if (state.status === "idle") {
    return (
      <div className="yt-cta">
        <div>
          <b>{t.ctaTitle}</b>
          {t.ctaBody}
        </div>
        <button type="button" className="btn btn--accent" onClick={start}>
          {t.ctaButton}
        </button>
      </div>
    );
  }

  if (state.status === "error") {
    if (state.error === "quota_exhausted") {
      return (
        <div className="yt-cta yt-cta--warn" role="status">
          <div>
            <b>{t.quotaTitle}</b>
            {t.quotaBody(state.retryAt ? formatLocalTime(state.retryAt, lang) : "—")}
          </div>
        </div>
      );
    }
    if (state.error === "user_limit") {
      return (
        <div className="yt-cta yt-cta--muted" role="status">
          <div>
            <b>{t.userLimitTitle}</b>
            {t.userLimitBody}
          </div>
        </div>
      );
    }
    if (state.error === "search_unavailable") {
      return (
        <div className="yt-cta yt-cta--muted" role="status">
          <div>
            <b>{t.unavailableTitle}</b>
            {t.unavailableBody}
          </div>
        </div>
      );
    }
    return (
      <div className="yt-cta yt-cta--warn" role="alert">
        <div>
          <b>{t.youtubeErrorTitle}</b>
          {t.youtubeErrorBody}
        </div>
        <button type="button" className="btn" onClick={start}>
          {t.retry}
        </button>
      </div>
    );
  }

  const exclude = new Set(excludeIds);
  const items = state.status === "done" ? state.items.filter((i) => !exclude.has(i.id)) : [];

  return (
    <section aria-busy={state.status === "loading"}>
      <div className="sec-title">
        <b>{t.youtubeSection}</b>
        <span>{state.status === "loading" ? t.loading : t.remaining(state.remaining)}</span>
      </div>
      {state.status === "loading" ? (
        <div className="flex flex-col gap-0.5">
          {[0, 1, 2].map((i) => (
            <div key={i} className="row animate-pulse">
              <span className="thumb h-[47px] w-[84px] md:h-[63px] md:w-[112px]" />
              <div className="min-w-0 flex-1">
                <div className="h-4 w-2/5 rounded bg-[var(--bg-elev-2)]" />
                <div className="mt-2 h-3 w-1/4 rounded bg-[var(--bg-elev-2)]" />
              </div>
            </div>
          ))}
        </div>
      ) : items.length === 0 ? (
        <p className="px-[14px] py-4 text-sm text-[var(--text-3)]">{t.youtubeEmpty}</p>
      ) : (
        items.map((item) => (
          <SearchItemRow
            key={item.id}
            item={item}
            href={`/${country.toLowerCase()}/watch/${item.id}`}
            hit={{ country, q }}
            badges={item.videoType ? <VideoTypeBadge videoType={item.videoType} /> : undefined}
            statPrimary={formatViewCount(item.viewCount, lang)}
            statSecondary={formatRelativeTime(item.publishedAt, lang)}
          />
        ))
      )}
    </section>
  );
}
