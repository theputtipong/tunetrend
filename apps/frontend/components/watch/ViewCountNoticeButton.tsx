"use client";

import { useState } from "react";
import { dictionaries, type Lang } from "@/lib/i18n";
import { InfoIcon } from "../icons";

export function ViewCountNoticeButton({ lang }: Readonly<{ lang: Lang }>) {
  const t = dictionaries[lang].watch;
  const [open, setOpen] = useState(false);

  return (
    <>
      <button
        type="button"
        onClick={() => setOpen(true)}
        title={t.viewCountNoticeTooltip}
        aria-label={t.viewCountNoticeTooltip}
        className="fixed bottom-4 right-4 z-40 flex h-9 w-9 items-center justify-center rounded-full bg-[var(--bg-elev-2)] text-[var(--text-2)] shadow-lg"
      >
        <InfoIcon size={18} />
      </button>

      {open && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4"
          onClick={() => setOpen(false)}
        >
          <div
            className="w-full max-w-sm rounded-2xl bg-[var(--bg-elev)] p-6"
            onClick={(e) => e.stopPropagation()}
          >
            <h2 className="font-display text-lg font-bold">{t.viewCountNoticeTitle}</h2>
            <p className="mt-2 text-sm text-[var(--text-2)]">{t.viewCountNoticeBody}</p>
            <div className="mt-5 flex justify-end">
              <button
                type="button"
                onClick={() => setOpen(false)}
                className="text-sm font-bold text-[var(--accent)]"
              >
                {t.viewCountNoticeDismiss}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
