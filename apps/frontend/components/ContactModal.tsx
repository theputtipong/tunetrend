"use client";

import { useEffect } from "react";
import type { Dictionary } from "@/lib/i18n/en";
import { ContactForm } from "./ContactForm";

export function ContactModal({
  t,
  onClose,
}: Readonly<{ t: Dictionary["about"]["contact"]; onClose: () => void }>) {
  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    document.addEventListener("keydown", onKeyDown);
    return () => document.removeEventListener("keydown", onKeyDown);
  }, [onClose]);

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4"
      onClick={onClose}
    >
      <div
        className="relative w-full max-w-md rounded-2xl bg-[var(--bg-elev)] p-6"
        onClick={(e) => e.stopPropagation()}
      >
        <button
          type="button"
          onClick={onClose}
          aria-label="Close"
          className="absolute right-4 top-4 text-xl leading-none text-[var(--text-3)] hover:text-[var(--text)]"
        >
          ✕
        </button>
        <h2 className="font-display mb-4 pr-6 text-lg font-bold">{t.openButton}</h2>
        <ContactForm t={t} onClose={onClose} />
      </div>
    </div>
  );
}
