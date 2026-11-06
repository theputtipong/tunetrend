import { headers } from "next/headers";
import { detectCountryFromAcceptLanguage } from "@/lib/countries";
import { getLang } from "@/lib/i18n/server";
import { Header } from "@/components/Header";

export default async function MeLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  const lang = await getLang();
  const country = detectCountryFromAcceptLanguage((await headers()).get("accept-language"));

  return (
    <div className="mx-auto max-w-6xl">
      <Header activeCountry={country} countrySuffix="" lang={lang} />
      <div className="border-t border-[var(--border)]" />
      {children}
    </div>
  );
}
