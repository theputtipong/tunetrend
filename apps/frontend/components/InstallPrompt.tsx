"use client";

import { useEffect, useState, useSyncExternalStore } from "react";
import { dictionaries, type Lang } from "@/lib/i18n";
import {
  APP_STORE_URL,
  type BeforeInstallPromptEvent,
  computeEligibleKind,
  detectMobilePlatform,
  isMarkedInstalled,
  isStandaloneDisplayMode,
  markInstalled,
  markShownNow,
  PLAY_STORE_URL,
  type PromptKind,
  shouldReprompt,
} from "@/lib/installPrompt";

let deferredInstallEvent: BeforeInstallPromptEvent | null = null;

function subscribe(callback: () => void) {
  function onBeforeInstallPrompt(e: Event) {
    e.preventDefault();
    deferredInstallEvent = e as BeforeInstallPromptEvent;
    callback();
  }
  function onAppInstalled() {
    markInstalled();
    deferredInstallEvent = null;
    callback();
  }
  window.addEventListener("beforeinstallprompt", onBeforeInstallPrompt);
  window.addEventListener("appinstalled", onAppInstalled);
  return () => {
    window.removeEventListener("beforeinstallprompt", onBeforeInstallPrompt);
    window.removeEventListener("appinstalled", onAppInstalled);
  };
}

function getSnapshot(): PromptKind | null {
  if (isStandaloneDisplayMode() || isMarkedInstalled() || !shouldReprompt()) return null;

  const platform = detectMobilePlatform(navigator.userAgent);
  if (!platform) return null;

  return computeEligibleKind(platform, deferredInstallEvent !== null);
}

function getServerSnapshot(): PromptKind | null {
  return null;
}

export function InstallPrompt({ lang }: Readonly<{ lang: Lang }>) {
  const kind = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);
  const [dismissed, setDismissed] = useState(false);

  useEffect(() => {
    if (kind) markShownNow();
  }, [kind]);

  function dismiss() {
    setDismissed(true);
  }

  async function handleAndroidPwaInstall() {
    if (!deferredInstallEvent) return;
    await deferredInstallEvent.prompt();
    await deferredInstallEvent.userChoice;
    deferredInstallEvent = null;
    dismiss();
  }

  if (!kind || dismissed) return null;

  const t = dictionaries[lang].installPrompt;

  return (
    <div className="fixed inset-0 z-50 flex items-end justify-center bg-black/60 p-4 sm:items-center">
      <div className="w-full max-w-sm rounded-2xl bg-[var(--bg-elev)] p-6 text-center">
        <h2 className="font-display text-lg font-bold">{t.title}</h2>

        {kind === "ios-pwa" ? (
          <>
            <p className="mt-2 text-sm text-[var(--text-2)]">{t.iosAddToHomeScreenDescription}</p>
            <button
              type="button"
              onClick={dismiss}
              className="retry mt-5 flex w-full items-center justify-center"
            >
              {t.gotIt}
            </button>
          </>
        ) : (
          <>
            <p className="mt-2 text-sm text-[var(--text-2)]">{t.description}</p>

            {kind === "android-pwa" ? (
              <button
                type="button"
                onClick={handleAndroidPwaInstall}
                className="retry mt-5 flex w-full items-center justify-center"
              >
                {t.installButtonAndroidPwa}
              </button>
            ) : (
              <a
                href={kind === "android-store" ? PLAY_STORE_URL : APP_STORE_URL}
                target="_blank"
                rel="noopener noreferrer"
                onClick={dismiss}
                className="retry mt-5 flex w-full items-center justify-center"
              >
                {kind === "android-store" ? t.getItButtonAndroid : t.getItButtonIos}
              </a>
            )}

            <button
              type="button"
              onClick={dismiss}
              className="mt-3 text-sm text-[var(--text-3)] underline"
            >
              {t.continueInBrowser}
            </button>
          </>
        )}
      </div>
    </div>
  );
}
