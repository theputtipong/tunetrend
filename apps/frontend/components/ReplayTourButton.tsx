"use client";

import { dictionaries, type Lang } from "@/lib/i18n";
import { HelpIcon } from "./icons";
import { runOnboardingTour } from "./OnboardingTour";

export function ReplayTourButton({
  lang,
  variant = "icon",
  onAction,
}: Readonly<{ lang: Lang; variant?: "icon" | "menu-item"; onAction?: () => void }>) {
  const label = dictionaries[lang].nav.replayTour;

  function handleClick() {
    runOnboardingTour(lang);
    onAction?.();
  }

  if (variant === "menu-item") {
    return (
      <button type="button" onClick={handleClick} className="menu-item">
        <HelpIcon size={16} />
        {label}
      </button>
    );
  }

  return (
    <button type="button" onClick={handleClick} className="theme-toggle" aria-label={label} title={label}>
      <HelpIcon />
    </button>
  );
}
