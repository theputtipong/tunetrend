"use client";

import { useEffect, useState } from "react";
import { dictionaries, type Lang } from "@/lib/i18n";

const COUNTDOWN_SECONDS = 5;

export function AutoPlayPromptDialog({
  lang,
  onDecide,
}: Readonly<{ lang: Lang; onDecide: (accepted: boolean) => void }>) {
  const t = dictionaries[lang].watch;
  const [secondsLeft, setSecondsLeft] = useState(COUNTDOWN_SECONDS);

  useEffect(() => {
    if (secondsLeft <= 1) {
      onDecide(true);
      return;
    }
    const timer = setTimeout(() => setSecondsLeft((s) => s - 1), 1000);
    return () => clearTimeout(timer);
  }, [secondsLeft, onDecide]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
      <div className="w-full max-w-sm rounded-2xl bg-[var(--bg-elev)] p-6">
        <h2 className="font-display text-lg font-bold">{t.autoPlayPromptTitle}</h2>
        <p className="mt-2 text-sm text-[var(--text-2)]">{t.autoPlayPromptDescription}</p>
        <p className="mt-2 text-xs text-[var(--text-3)]">
          {t.autoPlayPromptCountdown(secondsLeft)}
        </p>
        <div className="mt-5 flex justify-end gap-3">
          <button
            type="button"
            onClick={() => onDecide(false)}
            className="text-sm text-[var(--text-2)]"
          >
            {t.autoPlayPromptDecline}
          </button>
          <button
            type="button"
            onClick={() => onDecide(true)}
            className="text-sm font-bold text-[var(--accent)]"
          >
            {t.autoPlayPromptAccept}
          </button>
        </div>
      </div>
    </div>
  );
}
