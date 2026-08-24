"use client";

import { useState } from "react";
import { dictionaries, type Lang } from "@/lib/i18n";
import { ContactModal } from "./ContactModal";

export function ContactFab({ lang }: Readonly<{ lang: Lang }>) {
  const [isOpen, setIsOpen] = useState(false);
  const t = dictionaries[lang].about.contact;

  return (
    <>
      <button
        type="button"
        onClick={() => setIsOpen(true)}
        className="fixed bottom-5 right-5 z-40 inline-flex items-center gap-2 rounded-full px-5 py-3 text-sm font-bold shadow-lg"
        style={{ backgroundColor: "var(--accent)", color: "var(--accent-ink)" }}
      >
        {t.openButton}
      </button>
      {isOpen && <ContactModal t={t} onClose={() => setIsOpen(false)} />}
    </>
  );
}
