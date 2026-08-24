import type { Lang } from "./i18n";

export const COUNTRIES = [
  { code: "TH", label: "Thailand", labelTh: "ไทย" },
  { code: "KR", label: "South Korea", labelTh: "เกาหลีใต้" },
  { code: "JP", label: "Japan", labelTh: "ญี่ปุ่น" },
  { code: "US", label: "United States", labelTh: "สหรัฐอเมริกา" },
  { code: "GB", label: "United Kingdom", labelTh: "สหราชอาณาจักร" },
] as const;

export type CountryCode = (typeof COUNTRIES)[number]["code"];

const COUNTRY_CODE_SET = new Set<string>(COUNTRIES.map((c) => c.code));

export function isValidCountry(code: string): code is CountryCode {
  return COUNTRY_CODE_SET.has(code);
}

export function countryLabel(code: string, lang: Lang = "en"): string {
  const country = COUNTRIES.find((c) => c.code === code);
  if (!country) return code;
  return lang === "th" ? country.labelTh : country.label;
}

export const DEFAULT_COUNTRY: CountryCode = "TH";

export function detectCountryFromAcceptLanguage(acceptLanguage: string | null): CountryCode {
  if (!acceptLanguage) return DEFAULT_COUNTRY;

  const tags = acceptLanguage
    .split(",")
    .map((part) => part.split(";")[0]?.trim().toLowerCase())
    .filter((tag): tag is string => Boolean(tag));

  for (const tag of tags) {
    if (tag.startsWith("th")) return "TH";
    if (tag.startsWith("ko")) return "KR";
    if (tag.startsWith("ja")) return "JP";
    if (tag.startsWith("en-gb") || tag === "en-uk") return "GB";
    if (tag.startsWith("en")) return "US";
  }

  return DEFAULT_COUNTRY;
}
