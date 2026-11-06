"use client";

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";
import { useRouter } from "next/navigation";
import type { Lang } from "@/lib/i18n";
import type { User } from "@/types/user";
import type { Playlist, PlaylistVideo } from "@/types/playlist";
import { LoginModal, type LoginReason } from "./LoginModal";

export const GOOGLE_CLIENT_ID = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID ?? "";

type ApiResult<T> = { ok: boolean; status: number; data?: T; error?: string };

export async function apiRequest<T>(
  url: string,
  init: { method?: string; body?: unknown } = {},
): Promise<ApiResult<T>> {
  try {
    const res = await fetch(url, {
      method: init.method ?? "GET",
      headers: init.body === undefined ? undefined : { "Content-Type": "application/json" },
      body: init.body === undefined ? undefined : JSON.stringify(init.body),
      cache: "no-store",
    });
    const body = await res.json().catch(() => null);
    return {
      ok: res.ok && Boolean(body?.success),
      status: res.status,
      data: body?.data as T | undefined,
      error: typeof body?.error === "string" ? body.error : undefined,
    };
  } catch {
    return { ok: false, status: 0, error: "network_error" };
  }
}

export type DeleteAccountResult = "deleted" | "unauthorized" | "failed";

type PlaylistState = {
  playlists: Playlist[];
  membership: Record<string, string[]>;
};

type AuthContextValue = {
  user: User | null;
  lang: Lang;
  loginEnabled: boolean;
  openLogin: (reason?: LoginReason) => void;
  logout: () => Promise<void>;
  deleteAccount: () => Promise<DeleteAccountResult>;
  handleUnauthorized: () => void;
  playlistState: PlaylistState | null;
  loadPlaylists: (force?: boolean) => Promise<void>;
  createPlaylist: (name: string) => Promise<ApiResult<Playlist>>;
  addToPlaylist: (playlistId: string, video: PlaylistVideo) => Promise<ApiResult<unknown>>;
  removeFromPlaylist: (playlistId: string, videoId: string) => Promise<ApiResult<unknown>>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside AuthProvider");
  return ctx;
}

export function AuthProvider({
  initialUser,
  lang,
  children,
}: Readonly<{ initialUser: User | null; lang: Lang; children: ReactNode }>) {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(initialUser);
  const [loginReason, setLoginReason] = useState<LoginReason | null>(null);
  const [playlistState, setPlaylistState] = useState<PlaylistState | null>(null);
  const loadingRef = useRef<Promise<void> | null>(null);

  const loginEnabled = GOOGLE_CLIENT_ID !== "";

  const reset = useCallback(() => {
    setUser(null);
    setPlaylistState(null);
    loadingRef.current = null;
  }, []);

  const handleUnauthorized = useCallback(() => {
    reset();
  }, [reset]);

  const openLogin = useCallback(
    (reason: LoginReason = "default") => {
      if (loginEnabled) setLoginReason(reason);
    },
    [loginEnabled],
  );

  const onLoggedIn = useCallback(
    (next: User) => {
      setUser(next);
      setPlaylistState(null);
      loadingRef.current = null;
      setLoginReason(null);
      router.refresh();
    },
    [router],
  );

  const logout = useCallback(async () => {
    await apiRequest("/api/auth/logout", { method: "POST" });
    try {
      window.google?.accounts.id.disableAutoSelect();
    } catch {}
    reset();
    router.refresh();
  }, [reset, router]);

  const deleteAccount = useCallback(async (): Promise<DeleteAccountResult> => {
    const res = await apiRequest("/api/me", { method: "DELETE" });
    if (res.status === 401) return "unauthorized";
    if (!res.ok) return "failed";
    reset();
    router.refresh();
    return "deleted";
  }, [reset, router]);

  const loadPlaylists = useCallback(
    async (force = false) => {
      if (!force && loadingRef.current) return loadingRef.current;
      const task = (async () => {
        const list = await apiRequest<Playlist[]>("/api/me/playlists");
        if (!list.ok) {
          if (list.status === 401) handleUnauthorized();
          loadingRef.current = null;
          throw new Error(list.error ?? "failed");
        }
        const playlists = list.data ?? [];
        const details = await Promise.all(
          playlists.map((p) => apiRequest<Playlist>(`/api/me/playlists/${p.id}`)),
        );
        const membership: Record<string, string[]> = {};
        playlists.forEach((p, i) => {
          membership[p.id] = (details[i].data?.items ?? []).map((item) => item.videoId);
        });
        setPlaylistState({ playlists, membership });
      })();
      loadingRef.current = task;
      return task;
    },
    [handleUnauthorized],
  );

  const createPlaylist = useCallback(
    async (name: string) => {
      const res = await apiRequest<Playlist>("/api/me/playlists", {
        method: "POST",
        body: { name },
      });
      if (res.ok && res.data) {
        const created = { ...res.data, itemCount: 0 };
        setPlaylistState((prev) =>
          prev
            ? {
                playlists: [created, ...prev.playlists],
                membership: { ...prev.membership, [created.id]: [] },
              }
            : prev,
        );
      } else if (res.status === 401) {
        handleUnauthorized();
      }
      return res;
    },
    [handleUnauthorized],
  );

  const updateMembership = useCallback((playlistId: string, videoId: string, add: boolean) => {
    setPlaylistState((prev) => {
      if (!prev) return prev;
      const current = prev.membership[playlistId] ?? [];
      const has = current.includes(videoId);
      if (add === has) return prev;
      return {
        playlists: prev.playlists.map((p) =>
          p.id === playlistId ? { ...p, itemCount: p.itemCount + (add ? 1 : -1) } : p,
        ),
        membership: {
          ...prev.membership,
          [playlistId]: add ? [...current, videoId] : current.filter((id) => id !== videoId),
        },
      };
    });
  }, []);

  const addToPlaylist = useCallback(
    async (playlistId: string, video: PlaylistVideo) => {
      const res = await apiRequest(`/api/me/playlists/${playlistId}/items`, {
        method: "POST",
        body: video,
      });
      if (res.ok || res.status === 409) updateMembership(playlistId, video.videoId, true);
      if (res.status === 401) handleUnauthorized();
      return res;
    },
    [handleUnauthorized, updateMembership],
  );

  const removeFromPlaylist = useCallback(
    async (playlistId: string, videoId: string) => {
      const res = await apiRequest(`/api/me/playlists/${playlistId}/items/${videoId}`, {
        method: "DELETE",
      });
      if (res.ok || res.status === 404) updateMembership(playlistId, videoId, false);
      if (res.status === 401) handleUnauthorized();
      return res;
    },
    [handleUnauthorized, updateMembership],
  );

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      lang,
      loginEnabled,
      openLogin,
      logout,
      deleteAccount,
      handleUnauthorized,
      playlistState,
      loadPlaylists,
      createPlaylist,
      addToPlaylist,
      removeFromPlaylist,
    }),
    [
      user,
      lang,
      loginEnabled,
      openLogin,
      logout,
      deleteAccount,
      handleUnauthorized,
      playlistState,
      loadPlaylists,
      createPlaylist,
      addToPlaylist,
      removeFromPlaylist,
    ],
  );

  return (
    <AuthContext.Provider value={value}>
      {children}
      {loginReason && loginEnabled ? (
        <LoginModal
          reason={loginReason}
          lang={lang}
          clientId={GOOGLE_CLIENT_ID}
          onClose={() => setLoginReason(null)}
          onLoggedIn={onLoggedIn}
        />
      ) : null}
    </AuthContext.Provider>
  );
}
