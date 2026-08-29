import { type CountryCode } from "@/lib/countries";
import { getLang } from "@/lib/i18n/server";
import { fetchCategories } from "@/lib/api";
import { NavBar } from "@/components/NavBar";

export default async function MainLayout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ country: string }>;
}>) {
  const { country: rawCountry } = await params;
  const country = rawCountry.toUpperCase() as CountryCode;
  const lang = await getLang();
  const categories = await fetchCategories(country);

  return (
    <div className="mx-auto max-w-6xl">
      <NavBar country={country} lang={lang} categories={categories} />
      {children}
    </div>
  );
}
