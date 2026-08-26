"use client";

import { useRouter } from "next/navigation";
import type { Lang } from "@/lib/i18n";
import { setLang, useLang } from "@/lib/i18n/useLang";
import { GlobeIcon } from "./icons";

export function LanguageToggle({
  variant = "icon",
  onAction,
}: Readonly<{ variant?: "icon" | "menu-item"; onAction?: () => void }>) {
  const lang = useLang();
  const router = useRouter();

  function toggle() {
    const next: Lang = lang === "en" ? "th" : "en";
    setLang(next);
    router.refresh();
    onAction?.();
  }

  if (variant === "menu-item") {
    return (
      <button type="button" onClick={toggle} className="menu-item">
        <GlobeIcon size={16} />
        {lang === "en" ? "ภาษาไทย" : "English"}
      </button>
    );
  }

  return (
    <button
      type="button"
      onClick={toggle}
      className="theme-toggle lang-toggle"
      aria-label="Switch language"
      title="Switch language"
      data-tour="language-toggle"
    >
      {lang === "en" ? "TH" : "EN"}
    </button>
  );
}
