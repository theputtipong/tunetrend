import { notFound } from "next/navigation";
import { isValidCountry, type CountryCode } from "@/lib/countries";
import { getLang } from "@/lib/i18n/server";
import { fetchCategories } from "@/lib/api";
import { NavBar } from "@/components/NavBar";

export default async function CountryLayout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ country: string }>;
}>) {
  const { country: rawCountry } = await params;
  const country = rawCountry.toUpperCase();
  if (!isValidCountry(country)) {
    notFound();
  }
  const lang = await getLang();
  const categories = await fetchCategories(country as CountryCode);

  return (
    <div className="mx-auto max-w-6xl">
      <NavBar country={country} lang={lang} categories={categories} />
      {children}
    </div>
  );
}
