"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import YouTube, { type YouTubePlayer, type YouTubeEvent } from "react-youtube";
import { dictionaries, type Lang } from "@/lib/i18n";
import { PLAY_STORE_URL } from "@/lib/installPrompt";
import { MUSIC_CATEGORY_ID, TABS, type TabKey } from "@/lib/tabs";
import type { CountryCode } from "@/lib/countries";
import { countryLabel } from "@/lib/countries";
import type { Song } from "@/types/song";
import type { DiscoverItem } from "@/types/discover";
import { StateMessage } from "../StateMessage";
import { AutoPlayPromptDialog } from "./AutoPlayPromptDialog";
import { DiscoverCarousel } from "./DiscoverCarousel";
import { ViewCountNoticeButton } from "./ViewCountNoticeButton";
import { WatchQueueRow } from "./WatchQueueRow";

const AUTO_PLAY_DELAY_SECONDS = 5;

export function WatchPageClient({
  country,
  tab,
  lang,
  videoId: initialVideoId,
  title: initialTitle,
  categoryId: initialCategoryId,
  initialSongs,
  discoverItems,
}: Readonly<{
  country: CountryCode;
  tab: TabKey;
  lang: Lang;
  videoId: string;
  title: string;
  categoryId: string;
  initialSongs: Song[];
  discoverItems: DiscoverItem[];
}>) {
  const router = useRouter();
  const t = dictionaries[lang].watch;

  const [videoId, setVideoId] = useState(initialVideoId);
  const [title, setTitle] = useState(initialTitle);
  const [categoryId, setCategoryId] = useState(initialCategoryId);
  const [relatedTab, setRelatedTab] = useState<TabKey>(tab);
  const [songs, setSongs] = useState<Song[]>(initialSongs);
  const [loadingRelated, setLoadingRelated] = useState(false);
  const [relatedError, setRelatedError] = useState(false);

  const [autoPlayChoice, setAutoPlayChoice] = useState<boolean | null>(null);
  const [showAutoPlayPrompt, setShowAutoPlayPrompt] = useState(false);
  const [autoAdvanceSecondsLeft, setAutoAdvanceSecondsLeft] = useState<number | null>(null);
  const [copiedFeedback, setCopiedFeedback] = useState(false);

  const playerRef = useRef<YouTubePlayer | null>(null);
  const handledEndRef = useRef(false);
  const autoPlayGenerationRef = useRef(0);
  const nowPlayingRef = useRef<HTMLButtonElement | null>(null);

  const showMv = categoryId === "" || categoryId === MUSIC_CATEGORY_ID;

  const updateUrl = useCallback(
    (nextVideoId: string, nextTab: TabKey) => {
      window.history.replaceState(null, "", `/${country}/watch/${nextVideoId}?tab=${nextTab}`);
    },
    [country],
  );

  const loadRelated = useCallback(
    async (nextTab: TabKey) => {
      setLoadingRelated(true);
      setRelatedError(false);
      try {
        const res = await fetch(`/api/songs?country=${country}&tab=${nextTab}`);
        const body = await res.json();
        if (!body.success) throw new Error(body.error ?? "failed");
        setSongs(body.data ?? []);
      } catch {
        setRelatedError(true);
      } finally {
        setLoadingRelated(false);
      }
    },
    [country],
  );

  useEffect(() => {
    if (!nowPlayingRef.current) return;
    nowPlayingRef.current.scrollIntoView({ behavior: "smooth", block: "center" });
  }, [videoId, songs]);

  function selectRelatedTab(nextTab: TabKey) {
    if (nextTab === relatedTab) return;
    const wasAutoPlaying = autoPlayChoice === true;
    setRelatedTab(nextTab);
    autoPlayGenerationRef.current += 1;
    setAutoAdvanceSecondsLeft(null);
    if (wasAutoPlaying) setAutoPlayChoice(null);
    void loadRelated(nextTab);
  }

  const advanceToNext = useCallback(async () => {
    if (songs.length === 0) return;
    const currentIndex = songs.findIndex((s) => s.id === videoId);
    const next = currentIndex === -1 ? songs[0] : songs[currentIndex + 1];
    if (!next) return;

    const generation = autoPlayGenerationRef.current;
    for (let secondsLeft = AUTO_PLAY_DELAY_SECONDS; secondsLeft > 0; secondsLeft--) {
      if (generation !== autoPlayGenerationRef.current) return;
      setAutoAdvanceSecondsLeft(secondsLeft);
      await new Promise((resolve) => setTimeout(resolve, 1000));
    }
    if (generation !== autoPlayGenerationRef.current) return;

    setVideoId(next.id);
    setTitle(next.title);
    setCategoryId(next.categoryId);
    setAutoAdvanceSecondsLeft(null);
    handledEndRef.current = false;
    playerRef.current?.loadVideoById(next.id);
    updateUrl(next.id, relatedTab);
  }, [songs, videoId, relatedTab, updateUrl]);

  function selectSong(song: Song) {
    if (song.id === videoId) return;
    autoPlayGenerationRef.current += 1;
    setAutoAdvanceSecondsLeft(null);
    setVideoId(song.id);
    setTitle(song.title);
    setCategoryId(song.categoryId);
    handledEndRef.current = false;
    playerRef.current?.loadVideoById(song.id);
    updateUrl(song.id, relatedTab);
  }

  function onPlayerReady(event: YouTubeEvent) {
    playerRef.current = event.target;
  }

  async function onPlayerStateChange(event: YouTubeEvent<number>) {
    if (event.data !== YouTube.PlayerState.ENDED || handledEndRef.current) return;
    handledEndRef.current = true;

    if (autoPlayChoice === null) {
      setShowAutoPlayPrompt(true);
      return;
    }
    if (autoPlayChoice) void advanceToNext();
  }

  function onAutoPlayDecision(accepted: boolean) {
    setShowAutoPlayPrompt(false);
    setAutoPlayChoice(accepted);
    if (accepted) void advanceToNext();
  }

  function handleBack() {
    if (window.history.length > 1) {
      router.back();
    } else {
      router.push(`/${country}?tab=${relatedTab}`);
    }
  }

  async function handleShare() {
    const url = `${window.location.origin}/${country}/watch/${videoId}?tab=${relatedTab}`;
    const text = t.shareMessage(title, url, PLAY_STORE_URL);

    if (navigator.share) {
      try {
        // Only pass `text` — it already embeds the title, song link, and Play
        // Store link. Passing `title`/`url` as separate fields too makes most
        // share targets append them a second time on top of `text`.
        await navigator.share({ text });
      } catch {
        // user cancelled — no-op
      }
      return;
    }

    await navigator.clipboard.writeText(text);
    setCopiedFeedback(true);
    setTimeout(() => setCopiedFeedback(false), 2000);
  }

  const visibleTabs = showMv ? TABS : TABS.filter((key) => key !== "mv");

  return (
    <div className="flex min-h-full flex-col bg-[var(--bg)] md:grid md:h-screen md:grid-cols-[75%_25%] md:grid-rows-[70%_30%]">
      {/* Column 1, row 1 (70%): header + player stacked — fits together with
          the queue sidebar in one screen, no outer page scroll on desktop */}
      <div className="flex min-h-0 flex-col md:col-start-1 md:row-start-1">
        {/* Header: back / title / share — same background as the tabs/queue
            column so it reads as one continuous surface instead of a bar
            that blends into (or clashes with) the video below it */}
        <div className="flex items-center gap-2 bg-[var(--bg)] px-3 py-2.5">
          <button
            type="button"
            onClick={handleBack}
            aria-label={t.back}
            className="flex h-8 w-8 flex-shrink-0 items-center justify-center text-[var(--text)]"
          >
            ←
          </button>
          <div className="min-w-0 flex-1 truncate text-[14.5px] font-semibold text-[var(--text)]">
            {title}
          </div>
          <button
            type="button"
            onClick={handleShare}
            aria-label={t.shareTooltip}
            title={t.shareTooltip}
            className="flex h-8 w-8 flex-shrink-0 items-center justify-center text-[var(--text-2)]"
          >
            ⇪
          </button>
        </div>

        {/* Player + auto-advance banner — deliberately inverted from the
            current theme (like the now-playing overlay) so the player area
            stays visually distinct: dark in light theme, light in dark theme */}
        <div className="min-h-0 flex-1 bg-[var(--player-bg)]">
          {copiedFeedback && (
            <div className="px-3 pb-1 text-xs text-[var(--accent)]">{t.shareCopied}</div>
          )}

          <div className="aspect-video w-full">
            <YouTube
              videoId={videoId}
              opts={{ width: "100%", height: "100%", playerVars: { autoplay: 1 } }}
              className="h-full w-full"
              iframeClassName="h-full w-full"
              onReady={onPlayerReady}
              onStateChange={onPlayerStateChange}
            />
          </div>

          {autoAdvanceSecondsLeft !== null && (
            <div className="bg-[var(--player-bg)] px-4 py-2.5">
              <div className="inline-flex items-center gap-1.5 rounded-[10px] bg-[var(--player-banner-bg)] px-3 py-2">
                <span className="text-[var(--accent)]">▶</span>
                <span className="text-[12.5px] text-[var(--player-fg)]">
                  {t.autoAdvanceCountdown(autoAdvanceSecondsLeft)}
                </span>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Tabs + queue sidebar — desktop: column 2, row 1 (same 70% height as header+player) */}
      <div className="flex min-h-0 flex-1 flex-col bg-[var(--bg)] md:col-start-2 md:row-start-1 md:h-full">
        <div className="flex gap-6 px-4 pt-3.5">
          {visibleTabs.map((key) => (
            <button
              key={key}
              type="button"
              onClick={() => selectRelatedTab(key)}
              className={key === relatedTab ? "tab tab--active" : "tab"}
            >
              {dictionaries[lang].tabs[key]}
            </button>
          ))}
        </div>
        <div className="mt-3 border-t border-[var(--border)]" />

        <div className="flex-1 overflow-y-auto px-2 py-2">
          {relatedError ? (
            <StateMessage
              variant="error"
              title={dictionaries[lang].error.title}
              description={dictionaries[lang].error.description}
              retryLabel={dictionaries[lang].error.retry}
              onRetry={() => void loadRelated(relatedTab)}
            />
          ) : loadingRelated ? (
            <div className="py-10 text-center text-sm text-[var(--text-3)]">…</div>
          ) : songs.length === 0 ? (
            <StateMessage
              variant="empty"
              title={dictionaries[lang].songList.emptyTitle}
              description={dictionaries[lang].songList.emptyDescription(
                countryLabel(country, lang),
              )}
            />
          ) : (
            songs.map((song, index) => {
              const isPlaying = song.id === videoId;
              return (
                <WatchQueueRow
                  key={`${song.id}-${song.countryCode}`}
                  ref={isPlaying ? nowPlayingRef : undefined}
                  song={song}
                  rank={index + 1}
                  lang={lang}
                  isPlaying={isPlaying}
                  onSelect={() => selectSong(song)}
                />
              );
            })
          )}
        </div>
      </div>

      {/* Discover carousel — desktop-only, row 2 (30%), spans both columns */}
      <div className="hidden min-h-0 md:col-start-1 md:col-span-2 md:row-start-2 md:block md:overflow-hidden">
        <DiscoverCarousel items={discoverItems} lang={lang} />
      </div>

      <ViewCountNoticeButton lang={lang} />

      {showAutoPlayPrompt && <AutoPlayPromptDialog lang={lang} onDecide={onAutoPlayDecision} />}
    </div>
  );
}
