"use client";

import { useSearchParams } from "next/navigation";
import { dictionaries, type Lang } from "@/lib/i18n";
import { resolveTab } from "@/lib/tabs";
import { Header } from "./Header";
import { OnboardingTour } from "./OnboardingTour";
import { TrendTabs } from "./TrendTabs";

export function NavBar({ country, lang }: Readonly<{ country: string; lang: Lang }>) {
  const searchParams = useSearchParams();
  const tab = resolveTab(searchParams.get("tab"));

  return (
    <>
      <Header activeCountry={country} tab={tab} lang={lang} />
      <div className="tabbar flex items-center justify-between px-4 md:px-8">
        <TrendTabs country={country.toLowerCase()} active={tab} lang={lang} />
        <span className="sync-note hidden sm:inline">{dictionaries[lang].nav.syncNote}</span>
      </div>
      <OnboardingTour lang={lang} />
    </>
  );
}
