import { en, type Dictionary } from "./en";
import { th } from "./th";

export type Lang = "en" | "th";
export type { Dictionary };

export const LANG_COOKIE = "tunetrend-lang";
export const DEFAULT_LANG: Lang = "en";

export const dictionaries: Record<Lang, Dictionary> = { en, th };

export function isValidLang(value: string | null | undefined): value is Lang {
  return value === "en" || value === "th";
}

export function detectLangFromAcceptLanguage(acceptLanguage: string | null): Lang {
  if (!acceptLanguage) return DEFAULT_LANG;

  const tags = acceptLanguage
    .split(",")
    .map((part) => part.split(";")[0]?.trim().toLowerCase())
    .filter((tag): tag is string => Boolean(tag));

  for (const tag of tags) {
    if (tag.startsWith("th")) return "th";
  }

  return DEFAULT_LANG;
}
