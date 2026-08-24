"use client";

import { useEffect, useSyncExternalStore } from "react";
import { dictionaries, type Lang } from "@/lib/i18n";
import {
  AndroidInstallMode,
  APP_STORE_URL,
  detectMobilePlatform,
  INSTALL_PROMPT_CONFIG,
  IosInstallMode,
  MobilePlatform,
  PLAY_STORE_URL,
} from "@/lib/installPrompt";

const DISMISSED_KEY = "tunetrend-install-prompt-dismissed";

function subscribe(callback: () => void) {
  window.addEventListener("storage", callback);
  return () => window.removeEventListener("storage", callback);
}

function getSnapshot(): MobilePlatform | null {
  if (localStorage.getItem(DISMISSED_KEY) === "true") return null;

  const platform = detectMobilePlatform(navigator.userAgent);
  if (platform === MobilePlatform.Android && INSTALL_PROMPT_CONFIG.android === AndroidInstallMode.PlayStore) {
    return MobilePlatform.Android;
  }
  if (platform === MobilePlatform.Ios && INSTALL_PROMPT_CONFIG.ios === IosInstallMode.AppStore && APP_STORE_URL) {
    return MobilePlatform.Ios;
  }
  return null;
}

function getServerSnapshot(): MobilePlatform | null {
  return null;
}

export function InstallPrompt({ lang }: Readonly<{ lang: Lang }>) {
  const activePrompt = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);

  useEffect(() => {
    const onBeforeInstallPrompt = (e: Event) => e.preventDefault();
    window.addEventListener("beforeinstallprompt", onBeforeInstallPrompt);
    return () => window.removeEventListener("beforeinstallprompt", onBeforeInstallPrompt);
  }, []);

  function dismiss() {
    localStorage.setItem(DISMISSED_KEY, "true");
    window.dispatchEvent(new Event("storage"));
  }

  if (!activePrompt) return null;

  const t = dictionaries[lang].installPrompt;
  const storeUrl = activePrompt === MobilePlatform.Android ? PLAY_STORE_URL : APP_STORE_URL;
  const buttonLabel = activePrompt === MobilePlatform.Android ? t.getItButtonAndroid : t.getItButtonIos;

  return (
    <div className="fixed inset-0 z-50 flex items-end justify-center bg-black/60 p-4 sm:items-center">
      <div className="w-full max-w-sm rounded-2xl bg-[var(--bg-elev)] p-6 text-center">
        <h2 className="font-display text-lg font-bold">{t.title}</h2>
        <p className="mt-2 text-sm text-[var(--text-2)]">{t.description}</p>
        <a
          href={storeUrl}
          target="_blank"
          rel="noopener noreferrer"
          onClick={dismiss}
          className="retry mt-5 flex w-full items-center justify-center"
        >
          {buttonLabel}
        </a>
        <button
          type="button"
          onClick={dismiss}
          className="mt-3 text-sm text-[var(--text-3)] underline"
        >
          {t.continueInBrowser}
        </button>
      </div>
    </div>
  );
}
