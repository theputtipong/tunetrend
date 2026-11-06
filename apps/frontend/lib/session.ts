import { cache } from "react";
import { cookies } from "next/headers";
import type { User } from "@/types/user";
import type { Playlist } from "@/types/playlist";

export const SESSION_COOKIE = "tt_session";
export const SESSION_MAX_AGE_SECONDS = 60 * 60 * 24 * 30;

const API_BASE_URL = process.env.API_BASE_URL ?? "http://localhost:8080";
const BACKEND_TIMEOUT_MS = 10_000;
const CURRENT_USER_TIMEOUT_MS = 2_500;

export function sessionCookieOptions() {
  return {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax" as const,
    path: "/",
    maxAge: SESSION_MAX_AGE_SECONDS,
  };
}

export async function getSessionToken(): Promise<string | null> {
  const store = await cookies();
  const value = store.get(SESSION_COOKIE)?.value;
  return value && value.length <= 4096 ? value : null;
}

export function forwardedClientHeaders(clientIp: string | null | undefined): Headers {
  const headers = new Headers();
  const key = process.env.INTERNAL_PROXY_KEY;
  if (clientIp && key) {
    headers.set("X-Forwarded-For", clientIp);
    headers.set("X-Internal-Key", key);
  }
  return headers;
}

export async function backendFetch(
  path: string,
  init: RequestInit & { token?: string | null; timeoutMs?: number } = {},
): Promise<Response> {
  const { token, headers, timeoutMs = BACKEND_TIMEOUT_MS, ...rest } = init;
  const merged = new Headers(headers);
  if (token) merged.set("Authorization", `Bearer ${token}`);
  if (rest.body && !merged.has("Content-Type")) merged.set("Content-Type", "application/json");
  return fetch(`${API_BASE_URL}${path}`, {
    ...rest,
    headers: merged,
    cache: "no-store",
    signal: AbortSignal.timeout(timeoutMs),
  });
}

async function fetchAuthed<T>(path: string, timeoutMs?: number): Promise<T | null> {
  const token = await getSessionToken();
  if (!token) return null;
  try {
    const res = await backendFetch(path, { token, timeoutMs });
    if (!res.ok) return null;
    const body = await res.json();
    return body?.success ? (body.data as T) : null;
  } catch {
    return null;
  }
}

export const getCurrentUser = cache(() => fetchAuthed<User>("/me", CURRENT_USER_TIMEOUT_MS));

export async function getMyPlaylists(): Promise<Playlist[] | null> {
  return fetchAuthed<Playlist[]>("/me/playlists");
}

export async function getMyPlaylist(id: string): Promise<Playlist | null> {
  return fetchAuthed<Playlist>(`/me/playlists/${encodeURIComponent(id)}`);
}
