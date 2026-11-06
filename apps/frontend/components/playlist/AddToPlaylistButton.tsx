"use client";

import { useEffect, useId, useRef, useState } from "react";
import { dictionaries } from "@/lib/i18n";
import { MAX_PLAYLIST_NAME_LEN, isValidPlaylistName } from "@/lib/playlist";
import type { PlaylistVideo } from "@/types/playlist";
import { useAuth } from "../auth/AuthProvider";

function errorMessage(
  status: number,
  error: string | undefined,
  t: (typeof dictionaries)["en"]["playlist"],
): string {
  if (status === 422 && error?.includes("playlist limit")) return t.errorLimit;
  if (status === 422 && error?.includes("item limit")) return t.errorItemLimit;
  if (status === 422 && error?.startsWith("name")) return t.errorName;
  return t.errorGeneric;
}

export function AddToPlaylistButton({
  video,
  className = "",
}: Readonly<{ video: PlaylistVideo; className?: string }>) {
  const {
    user,
    lang,
    loginEnabled,
    openLogin,
    playlistState,
    loadPlaylists,
    createPlaylist,
    addToPlaylist,
    removeFromPlaylist,
  } = useAuth();
  const t = dictionaries[lang].playlist;
  const [isOpen, setIsOpen] = useState(false);
  const [loadError, setLoadError] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [pending, setPending] = useState<string | null>(null);
  const [newName, setNewName] = useState("");
  const ref = useRef<HTMLDivElement>(null);
  const popoverId = useId();

  useEffect(() => {
    if (!isOpen) return;
    function onClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setIsOpen(false);
    }
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === "Escape") setIsOpen(false);
    }
    document.addEventListener("mousedown", onClickOutside);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("mousedown", onClickOutside);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [isOpen]);

  if (!user && !loginEnabled) return null;

  const isSaved = Boolean(
    playlistState &&
    Object.values(playlistState.membership).some((ids) => ids.includes(video.videoId)),
  );

  function onButtonClick(e: React.MouseEvent) {
    e.preventDefault();
    e.stopPropagation();
    if (!user) {
      openLogin("playlist");
      return;
    }
    const next = !isOpen;
    setIsOpen(next);
    setError(null);
    if (next) {
      setLoadError(false);
      loadPlaylists().catch(() => setLoadError(true));
    }
  }

  async function toggle(playlistId: string, checked: boolean) {
    setPending(playlistId);
    setError(null);
    const res = checked
      ? await addToPlaylist(playlistId, video)
      : await removeFromPlaylist(playlistId, video.videoId);
    setPending(null);
    if (!res.ok && res.status !== 409 && !(res.status === 404 && !checked)) {
      setError(errorMessage(res.status, res.error, t));
    }
  }

  async function onCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!isValidPlaylistName(newName)) {
      setError(t.errorName);
      return;
    }
    setPending("new");
    setError(null);
    const created = await createPlaylist(newName.trim());
    if (!created.ok || !created.data) {
      setPending(null);
      setError(errorMessage(created.status, created.error, t));
      return;
    }
    const added = await addToPlaylist(created.data.id, video);
    setPending(null);
    if (!added.ok && added.status !== 409) setError(errorMessage(added.status, added.error, t));
    setNewName("");
  }

  return (
    <div ref={ref} className={`relative flex-shrink-0 ${className}`}>
      <button
        type="button"
        className={isSaved ? "icon-btn icon-btn--done" : "icon-btn"}
        aria-label={t.addButton}
        title={t.addButton}
        aria-haspopup="dialog"
        aria-expanded={isOpen}
        aria-controls={isOpen ? popoverId : undefined}
        onClick={onButtonClick}
      >
        {isSaved ? "✓" : "+"}
      </button>

      {isOpen && (
        <div
          id={popoverId}
          role="dialog"
          aria-label={t.popoverTitle}
          className="popover"
          onClick={(e) => e.stopPropagation()}
        >
          <h4>{t.popoverTitle}</h4>
          {loadError ? (
            <p className="px-1 py-2 text-sm text-[var(--error-text)]">{t.errorGeneric}</p>
          ) : !playlistState ? (
            <p className="px-1 py-2 text-sm text-[var(--text-3)]">{t.loading}</p>
          ) : playlistState.playlists.length === 0 ? (
            <p className="px-1 py-2 text-sm text-[var(--text-3)]">{t.noPlaylists}</p>
          ) : (
            <div className="max-h-[240px] overflow-y-auto">
              {playlistState.playlists.map((p) => {
                const checked = (playlistState.membership[p.id] ?? []).includes(video.videoId);
                return (
                  <label key={p.id} className="pl-opt">
                    <input
                      type="checkbox"
                      checked={checked}
                      disabled={pending !== null}
                      onChange={(e) => void toggle(p.id, e.target.checked)}
                    />
                    <span className="min-w-0 truncate">{p.name}</span>
                    <small>{p.itemCount}</small>
                  </label>
                );
              })}
            </div>
          )}

          {error ? (
            <p role="alert" className="px-1 pt-2 text-xs text-[var(--error-text)]">
              {error}
            </p>
          ) : null}

          {playlistState ? (
            <form
              onSubmit={onCreate}
              className="mt-2 flex gap-2 border-t border-[var(--border)] pt-3"
            >
              <input
                className="text-input"
                value={newName}
                maxLength={MAX_PLAYLIST_NAME_LEN}
                placeholder={t.newPlaceholder}
                aria-label={t.newPlaceholder}
                onChange={(e) => setNewName(e.target.value)}
              />
              <button
                type="submit"
                className="btn btn--accent btn--sm"
                disabled={pending !== null || !newName.trim()}
              >
                {t.create}
              </button>
            </form>
          ) : null}
        </div>
      )}
    </div>
  );
}
