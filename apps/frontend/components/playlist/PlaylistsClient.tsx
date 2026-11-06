"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { dictionaries } from "@/lib/i18n";
import { MAX_PLAYLIST_NAME_LEN, isValidPlaylistName } from "@/lib/playlist";
import type { Playlist } from "@/types/playlist";
import { apiRequest, useAuth } from "../auth/AuthProvider";
import { PlaylistCover } from "./PlaylistCover";

export function PlaylistsClient({ initialPlaylists }: Readonly<{ initialPlaylists: Playlist[] }>) {
  const { lang, handleUnauthorized } = useAuth();
  const router = useRouter();
  const t = dictionaries[lang].playlist;
  const [playlists, setPlaylists] = useState(initialPlaylists);
  const [creating, setCreating] = useState(false);
  const [name, setName] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!isValidPlaylistName(name)) {
      setError(t.errorName);
      return;
    }
    setBusy(true);
    setError(null);
    const res = await apiRequest<Playlist>("/api/me/playlists", {
      method: "POST",
      body: { name: name.trim() },
    });
    setBusy(false);
    if (res.ok && res.data) {
      setPlaylists((prev) => [...prev, { ...res.data!, itemCount: 0 }]);
      setName("");
      setCreating(false);
      router.refresh();
      return;
    }
    if (res.status === 401) handleUnauthorized();
    setError(
      res.status === 422 && res.error?.includes("limit")
        ? t.errorLimit
        : res.status === 422
          ? t.errorName
          : t.errorGeneric,
    );
  }

  return (
    <div className="px-4 pb-10 pt-6 md:px-8">
      <h1 className="font-display text-2xl font-bold">{t.heading}</h1>

      {playlists.length === 0 && !creating ? (
        <p className="mt-2 text-sm text-[var(--text-2)]">
          <span className="font-semibold">{t.listEmptyTitle}</span> · {t.listEmptyBody}
        </p>
      ) : null}

      <div className="mt-5 grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
        {playlists.map((p) => (
          <Link key={p.id} href={`/me/playlists/${p.id}`} className="pl-card">
            <PlaylistCover />
            <b className="truncate text-[14.5px]">{p.name}</b>
            <small>{t.songCount(p.itemCount)}</small>
          </Link>
        ))}

        {creating ? (
          <form onSubmit={onCreate} className="pl-card justify-center gap-2">
            <label htmlFor="new-playlist" className="text-sm font-semibold">
              {t.nameLabel}
            </label>
            <input
              id="new-playlist"
              className="text-input"
              value={name}
              autoFocus
              maxLength={MAX_PLAYLIST_NAME_LEN}
              placeholder={t.newPlaceholder}
              onChange={(e) => setName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Escape") setCreating(false);
              }}
            />
            {error ? (
              <p role="alert" className="text-xs text-[var(--error-text)]">
                {error}
              </p>
            ) : null}
            <div className="flex gap-2">
              <button type="submit" className="btn btn--accent btn--sm" disabled={busy}>
                {t.create}
              </button>
              <button type="button" className="btn btn--sm" onClick={() => setCreating(false)}>
                {t.cancel}
              </button>
            </div>
          </form>
        ) : (
          <button
            type="button"
            className="pl-card min-h-[120px] items-center justify-center border-dashed font-semibold text-[var(--text-2)]"
            onClick={() => {
              setError(null);
              setCreating(true);
            }}
          >
            {t.newPlaylist}
          </button>
        )}
      </div>
    </div>
  );
}
