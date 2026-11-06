"use client";

import { useSearchParams } from "next/navigation";
import { dictionaries, type Lang } from "@/lib/i18n";
import { resolveCategory } from "@/lib/categories";
import { resolvePageTab } from "@/lib/tabs";
import type { Category } from "@/types/category";
import { CategoryFilter } from "./CategoryFilter";
import { Header } from "./Header";
import { OnboardingTour } from "./OnboardingTour";
import { TrendTabs } from "./TrendTabs";

export function NavBar({
  country,
  lang,
  categories,
}: Readonly<{ country: string; lang: Lang; categories: Category[] }>) {
  const searchParams = useSearchParams();
  const tab = resolvePageTab(searchParams.get("tab"));
  const category = resolveCategory(searchParams.get("category"), categories);

  return (
    <>
      <Header activeCountry={country} countrySuffix={`?tab=${tab}`} lang={lang} />
      <CategoryFilter
        country={country.toLowerCase()}
        tab={tab}
        active={category}
        categories={categories}
        lang={lang}
      />
      <div className="tabbar flex items-center justify-between px-4 md:px-8">
        <TrendTabs country={country.toLowerCase()} active={tab} category={category} lang={lang} />
        <span className="sync-note hidden flex-shrink-0 pl-4 lg:inline">
          {tab === "searched"
            ? dictionaries[lang].nav.searchedSyncNote
            : dictionaries[lang].nav.syncNote}
        </span>
      </div>
      <OnboardingTour lang={lang} />
    </>
  );
}
