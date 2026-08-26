import Link from "next/link";
import { COUNTRIES, countryLabel } from "@/lib/countries";
import { dictionaries, type Lang } from "@/lib/i18n";
import { BuyMeCoffeeButton } from "./BuyMeCoffeeButton";
import { GlobeIcon, InfoIcon, LogoMark } from "./icons";
import { LanguageToggle } from "./LanguageToggle";
import { MobileMenu } from "./MobileMenu";
import { ReplayTourButton } from "./ReplayTourButton";
import { ThemeToggle } from "./ThemeToggle";
import type { TabKey } from "@/lib/tabs";

export function Header({
  activeCountry,
  tab,
  lang,
}: Readonly<{ activeCountry: string; tab: TabKey; lang: Lang }>) {
  const t = dictionaries[lang].nav;

  return (
    <div className="flex flex-col gap-3 px-4 py-4 md:px-8 md:py-6">
      <div className="flex items-center justify-between">
        <div className="flex flex-shrink-0 items-center gap-2">
          <LogoMark size={30} />
          <span className="font-display text-[21px] font-bold tracking-[-0.01em]">TuneTrend</span>
        </div>

        <div className="flex flex-shrink-0 items-center gap-2">
          <div className="hidden sm:inline-flex">
            <BuyMeCoffeeButton label={t.buyCoffee} size={32} dataTour="buy-coffee" />
          </div>
          <div className="hidden items-center gap-2 sm:flex">
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
              href={`/${c.code.toLowerCase()}?tab=${tab}`}
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
