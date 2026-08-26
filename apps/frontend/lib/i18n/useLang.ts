"use client";

import { useSyncExternalStore } from "react";
import { DEFAULT_LANG, LANG_COOKIE, isValidLang, type Lang } from "./index";

const CHANGE_EVENT = "tunetrend-lang-change";

function readCookie(): Lang {
  const match = document.cookie.match(new RegExp(`(?:^|; )${LANG_COOKIE}=([^;]*)`));
  const value = match ? decodeURIComponent(match[1]) : undefined;
  return isValidLang(value) ? value : DEFAULT_LANG;
}

function subscribe(callback: () => void) {
  window.addEventListener(CHANGE_EVENT, callback);
  return () => window.removeEventListener(CHANGE_EVENT, callback);
}

export function useLang(): Lang {
  return useSyncExternalStore(subscribe, readCookie, () => DEFAULT_LANG);
}

export function setLang(lang: Lang) {
  document.cookie = `${LANG_COOKIE}=${lang}; path=/; max-age=${60 * 60 * 24 * 365}`;
  window.dispatchEvent(new Event(CHANGE_EVENT));
}
