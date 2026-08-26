"use client";

import { useSyncExternalStore } from "react";
import { dictionaries, type Lang } from "@/lib/i18n";
import { MoonIcon, SunIcon } from "./icons";

const STORAGE_KEY = "tunetrend-theme";

function subscribe(callback: () => void) {
  const mql = window.matchMedia("(prefers-color-scheme: dark)");
  mql.addEventListener("change", callback);
  window.addEventListener("storage", callback);
  return () => {
    mql.removeEventListener("change", callback);
    window.removeEventListener("storage", callback);
  };
}

function getSnapshot(): boolean {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored === "dark") return true;
  if (stored === "light") return false;
  return window.matchMedia("(prefers-color-scheme: dark)").matches;
}

function getServerSnapshot(): boolean {
  return false;
}

export function ThemeToggle({
  lang = "en",
  variant = "icon",
  onAction,
}: Readonly<{ lang?: Lang; variant?: "icon" | "menu-item"; onAction?: () => void }>) {
  const isDark = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);
  const label = dictionaries[lang].nav.themeToggle;

  function toggle() {
    const next = !isDark;
    document.documentElement.setAttribute("data-theme", next ? "dark" : "light");
    localStorage.setItem(STORAGE_KEY, next ? "dark" : "light");
    window.dispatchEvent(new Event("storage"));
    onAction?.();
  }

  if (variant === "menu-item") {
    return (
      <button type="button" onClick={toggle} className="menu-item">
        {isDark ? <SunIcon size={16} /> : <MoonIcon size={16} />}
        {label}
      </button>
    );
  }

  return (
    <button
      type="button"
      onClick={toggle}
      className="theme-toggle"
      aria-label={label}
      title={label}
      data-tour="theme-toggle"
    >
      {isDark ? <SunIcon /> : <MoonIcon />}
    </button>
  );
}
