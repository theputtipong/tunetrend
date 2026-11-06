import Link from "next/link";
import { COUNTRIES, countryLabel } from "@/lib/countries";
import { dictionaries, type Lang } from "@/lib/i18n";
import { BuyMeCoffeeButton } from "./BuyMeCoffeeButton";
import { PLAY_STORE_URL } from "@/lib/installPrompt";
import { GetAppIcon, GlobeIcon, InfoIcon, LogoMark } from "./icons";
import { LanguageToggle } from "./LanguageToggle";
import { MobileMenu } from "./MobileMenu";
import { ReplayTourButton } from "./ReplayTourButton";
import { ThemeToggle } from "./ThemeToggle";
import { HeaderAuth } from "./auth/AccountMenu";
import { SearchBox } from "./search/SearchBox";

export function Header({
  activeCountry,
  countrySuffix,
  lang,
  searchQuery,
  searchCategory,
}: Readonly<{
  activeCountry: string;
  countrySuffix: string;
  lang: Lang;
  searchQuery?: string;
  searchCategory?: string;
}>) {
  const t = dictionaries[lang].nav;

  return (
    <div className="flex flex-col gap-3 px-4 py-4 md:px-8 md:py-6">
      <div className="flex flex-wrap items-center justify-between gap-x-4 gap-y-3">
        <Link
          href={`/${activeCountry.toLowerCase()}`}
          className="flex flex-shrink-0 items-center gap-2 text-[var(--text)] hover:text-[var(--text)]"
        >
          <LogoMark size={30} />
          <span className="font-display text-[21px] font-bold tracking-[-0.01em]">TuneTrend</span>
        </Link>

        <div className="order-last w-full md:order-none md:w-auto md:max-w-md md:flex-1">
          <SearchBox
            key={`${searchQuery ?? ""}|${searchCategory ?? ""}`}
            country={activeCountry}
            lang={lang}
            initialQuery={searchQuery}
            category={searchCategory}
          />
        </div>

        <div className="flex flex-shrink-0 items-center gap-2">
          <div className="hidden sm:inline-flex">
            <BuyMeCoffeeButton label={t.buyCoffee} size={32} dataTour="buy-coffee" />
          </div>
          <div className="hidden items-center gap-2 sm:flex">
            <a
              href={PLAY_STORE_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="theme-toggle"
              aria-label={t.getApp}
              title={t.getApp}
            >
              <GetAppIcon />
            </a>
            <Link
              href="/about"
              className="theme-toggle"
              aria-label={t.about}
              title={t.about}
              data-tour="about-link"
            >
              <InfoIcon />
            </Link>
            <LanguageToggle />
            <ThemeToggle lang={lang} />
            <ReplayTourButton lang={lang} />
          </div>
          <MobileMenu lang={lang} />
          <HeaderAuth />
        </div>
      </div>

      <div className="flex items-center gap-2 overflow-x-auto" data-tour="country-selector">
        <span className="mr-1 flex flex-shrink-0 items-center gap-1.5 text-xs font-semibold uppercase tracking-wide text-[var(--text-3)]">
          <GlobeIcon />
          {t.countryLabel}
        </span>
        {COUNTRIES.map((c) => {
          const active = c.code === activeCountry;
          return (
            <Link
              key={c.code}
              href={`/${c.code.toLowerCase()}${countrySuffix}`}
              className={active ? "pill pill--active" : "pill"}
              title={countryLabel(c.code, lang)}
            >
              {c.code}
            </Link>
          );
        })}
      </div>
    </div>
  );
}
