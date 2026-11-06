"use client";

import { useRef, useState, useSyncExternalStore } from "react";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { formatRelativeTime } from "@/lib/format";
import { dictionaries } from "@/lib/i18n";
import {
  MAX_PLAYLIST_NAME_LEN,
  isValidPlaylistName,
  readUnavailableVideos,
  subscribeUnavailable,
} from "@/lib/playlist";
import type { Playlist, PlaylistItem } from "@/types/playlist";
import { apiRequest, useAuth } from "../auth/AuthProvider";
import { ConfirmDialog } from "../auth/Modal";
import { PlayIcon } from "../icons";
import { PlaylistCover } from "./PlaylistCover";

function pickRandom<T>(list: T[]): T {
  return list[Math.floor(Math.random() * list.length)];
}

function moveItem<T>(list: T[], from: number, to: number): T[] {
  const next = [...list];
  const [moved] = next.splice(from, 1);
  next.splice(to, 0, moved);
  return next;
}

export function PlaylistDetailClient({
  initialPlaylist,
  country,
}: Readonly<{ initialPlaylist: Playlist; country: string }>) {
  const { lang, handleUnauthorized } = useAuth();
  const router = useRouter();
  const t = dictionaries[lang].playlist;

  const [name, setName] = useState(initialPlaylist.name);
  const [updatedAt, setUpdatedAt] = useState(initialPlaylist.updatedAt);
  const [items, setItems] = useState<PlaylistItem[]>(initialPlaylist.items ?? []);
  const [renaming, setRenaming] = useState(false);
  const [draftName, setDraftName] = useState(initialPlaylist.name);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [dragIndex, setDragIndex] = useState<number | null>(null);
  const [overIndex, setOverIndex] = useState<number | null>(null);
  const handleRefs = useRef(new Map<string, HTMLButtonElement>());

  const unavailableRaw = useSyncExternalStore(
    subscribeUnavailable,
    readUnavailableVideos,
    () => "",
  );
  const unavailable = new Set(unavailableRaw.split(",").filter(Boolean));
  const isUnavailable = (item: PlaylistItem) => unavailable.has(item.videoId) || !item.title;

  const id = initialPlaylist.id;
  const playable = items.filter((i) => !isUnavailable(i));
  const lower = country.toLowerCase();

  function fail(status: number) {
    if (status === 401) handleUnauthorized();
    setError(t.actionFailed);
  }

  async function onRename(e: React.FormEvent) {
    e.preventDefault();
    if (!isValidPlaylistName(draftName)) {
      setError(t.errorName);
      return;
    }
    setBusy(true);
    setError(null);
    const res = await apiRequest<Playlist>(`/api/me/playlists/${id}`, {
      method: "PATCH",
      body: { name: draftName.trim() },
    });
    setBusy(false);
    if (!res.ok || !res.data) return fail(res.status);
    setName(res.data.name);
    setUpdatedAt(res.data.updatedAt);
    setRenaming(false);
  }

  async function onDelete() {
    setBusy(true);
    setError(null);
    const res = await apiRequest(`/api/me/playlists/${id}`, { method: "DELETE" });
    setBusy(false);
    if (!res.ok) {
      setConfirmDelete(false);
      return fail(res.status);
    }
    router.push("/me/playlists");
    router.refresh();
  }

  async function onRemove(videoId: string) {
    const previous = items;
    setItems((list) => list.filter((i) => i.videoId !== videoId));
    setError(null);
    const res = await apiRequest(`/api/me/playlists/${id}/items/${videoId}`, { method: "DELETE" });
    if (!res.ok && res.status !== 404) {
      setItems(previous);
      fail(res.status);
    }
  }

  async function reorder(from: number, to: number) {
    if (from === to || to < 0 || to >= items.length) return;
    const previous = items;
    const next = moveItem(items, from, to);
    const movedId = items[from].videoId;
    setItems(next);
    setError(null);
    requestAnimationFrame(() => handleRefs.current.get(movedId)?.focus());

    const res = await apiRequest<Playlist>(`/api/me/playlists/${id}/order`, {
      method: "PUT",
      body: { videoIds: next.map((i) => i.videoId) },
    });
    if (!res.ok) {
      setItems(previous);
      return fail(res.status);
    }
    if (res.data?.items) setItems(res.data.items);
    if (res.data?.updatedAt) setUpdatedAt(res.data.updatedAt);
  }

  const watchHref = (videoId: string, shuffle = false) =>
    `/${lower}/watch/${videoId}?playlist=${id}${shuffle ? "&shuffle=1" : ""}`;

  function onShuffle() {
    if (playable.length === 0) return;
    const start = pickRandom(playable);
    router.push(watchHref(start.videoId, true));
  }

  return (
    <div className="px-4 pb-10 pt-5 md:px-8">
      <Link href="/me/playlists" className="text-sm font-semibold">
        {t.backToList}
      </Link>

      <div className="mt-4 flex flex-wrap items-end gap-5">
        <PlaylistCover items={items} className="w-[132px] flex-shrink-0 md:w-[160px]" />
        <div className="min-w-[200px] flex-1">
          <span className="badge">{t.privateBadge}</span>
          {renaming ? (
            <form onSubmit={onRename} className="my-2 flex max-w-md flex-wrap gap-2">
              <input
                className="text-input"
                value={draftName}
                autoFocus
                maxLength={MAX_PLAYLIST_NAME_LEN}
                aria-label={t.nameLabel}
                onChange={(e) => setDraftName(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Escape") setRenaming(false);
                }}
              />
              <button type="submit" className="btn btn--accent btn--sm" disabled={busy}>
                {t.save}
              </button>
              <button type="button" className="btn btn--sm" onClick={() => setRenaming(false)}>
                {t.cancel}
              </button>
            </form>
          ) : (
            <h1 className="font-display my-1 break-words text-2xl font-bold md:text-[28px]">
              {name}
            </h1>
          )}
          <p className="mb-3 text-[12.5px] text-[var(--text-3)]">
            {t.songCount(items.length)} · {t.updated(formatRelativeTime(updatedAt, lang))}
          </p>
          <div className="flex flex-wrap gap-2">
            {playable.length > 0 ? (
              <Link href={watchHref(playable[0].videoId)} className="btn btn--accent">
                {t.playAll}
              </Link>
            ) : (
              <button type="button" className="btn btn--accent" disabled>
                {t.playAll}
              </button>
            )}
            <button type="button" className="btn" onClick={onShuffle} disabled={!playable.length}>
              {t.shuffle}
            </button>
            <button
              type="button"
              className="btn"
              onClick={() => {
                setDraftName(name);
                setError(null);
                setRenaming(true);
              }}
            >
              {t.rename}
            </button>
            <button
              type="button"
              className="btn btn--danger"
              onClick={() => {
                setError(null);
                setConfirmDelete(true);
              }}
            >
              {t.delete}
            </button>
          </div>
        </div>
      </div>

      {error ? (
        <p role="alert" className="mt-4 text-sm text-[var(--error-text)]">
          {error}
        </p>
      ) : null}

      {items.length === 0 ? (
        <div className="mt-8 text-center">
          <p className="state-title">{t.emptyTitle}</p>
          <p className="mt-2 text-sm text-[var(--text-2)]">{t.emptyBody}</p>
        </div>
      ) : (
        <ol className="mt-5 flex flex-col gap-0.5">
          {items.map((item, index) => {
            const off = isUnavailable(item);
            return (
              <li
                key={item.videoId}
                className={`row w-full gap-2 md:gap-3 ${off ? "row--unavailable" : ""} ${
                  overIndex === index && dragIndex !== null && dragIndex !== index
                    ? "row--drag-over"
                    : ""
                }`}
                onDragOver={(e) => {
                  if (dragIndex === null) return;
                  e.preventDefault();
                  setOverIndex(index);
                }}
                onDrop={(e) => {
                  e.preventDefault();
                  if (dragIndex !== null) void reorder(dragIndex, index);
                  setDragIndex(null);
                  setOverIndex(null);
                }}
              >
                <button
                  ref={(el) => {
                    if (el) handleRefs.current.set(item.videoId, el);
                    else handleRefs.current.delete(item.videoId);
                  }}
                  type="button"
                  draggable
                  className="drag-handle flex h-8 w-6 flex-shrink-0 items-center justify-center border-0 bg-transparent"
                  aria-label={`${t.dragHandle}: ${item.title || t.unavailableTitle}`}
                  title={t.dragHandle}
                  onDragStart={(e) => {
                    e.dataTransfer.effectAllowed = "move";
                    e.dataTransfer.setData("text/plain", item.videoId);
                    setDragIndex(index);
                  }}
                  onDragEnd={() => {
                    setDragIndex(null);
                    setOverIndex(null);
                  }}
                  onKeyDown={(e) => {
                    if (e.key === "ArrowUp") {
                      e.preventDefault();
                      void reorder(index, index - 1);
                    } else if (e.key === "ArrowDown") {
                      e.preventDefault();
                      void reorder(index, index + 1);
                    }
                  }}
                >
                  ⋮⋮
                </button>

                <Link
                  href={off ? "#" : watchHref(item.videoId)}
                  aria-disabled={off || undefined}
                  tabIndex={off ? -1 : undefined}
                  onClick={off ? (e) => e.preventDefault() : undefined}
                  className="row-link"
                >
                  <span className="thumb relative h-[47px] w-[84px] overflow-hidden md:h-[63px] md:w-[112px]">
                    {item.thumbnailUrl && !off ? (
                      <Image
                        src={item.thumbnailUrl}
                        alt=""
                        fill
                        sizes="(min-width: 768px) 112px, 84px"
                        className="object-cover"
                      />
                    ) : null}
                    {!off ? (
                      <span className="thumb-play absolute inset-0 flex items-center justify-center bg-black/0 transition-colors">
                        <PlayIcon size={16} />
                      </span>
                    ) : null}
                  </span>
                  <div className="min-w-0 flex-1">
                    <div className="title">{off ? t.unavailableTitle : item.title}</div>
                    <div className="channel">{off ? t.unavailableBody : item.channelTitle}</div>
                  </div>
                  <div className="hidden w-[120px] flex-shrink-0 text-right text-[12px] text-[var(--text-3)] md:block">
                    {off ? t.skipOnPlay : formatRelativeTime(item.addedAt, lang)}
                  </div>
                </Link>

                <div className="flex flex-shrink-0 items-center gap-1">
                  <button
                    type="button"
                    className="icon-btn hidden sm:inline-flex"
                    aria-label={t.moveUp}
                    title={t.moveUp}
                    disabled={index === 0}
                    onClick={() => void reorder(index, index - 1)}
                  >
                    ↑
                  </button>
                  <button
                    type="button"
                    className="icon-btn hidden sm:inline-flex"
                    aria-label={t.moveDown}
                    title={t.moveDown}
                    disabled={index === items.length - 1}
                    onClick={() => void reorder(index, index + 1)}
                  >
                    ↓
                  </button>
                  <button
                    type="button"
                    className="icon-btn"
                    aria-label={`${t.removeItem}: ${item.title}`}
                    title={t.removeItem}
                    onClick={() => void onRemove(item.videoId)}
                  >
                    ✕
                  </button>
                </div>
              </li>
            );
          })}
        </ol>
      )}

      {confirmDelete ? (
        <ConfirmDialog
          title={t.deleteTitle}
          body={t.deleteBody(name)}
          confirmLabel={t.delete}
          cancelLabel={t.cancel}
          busy={busy}
          onConfirm={() => void onDelete()}
          onCancel={() => setConfirmDelete(false)}
        />
      ) : null}
    </div>
  );
}
