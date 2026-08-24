import { notFound } from "next/navigation";
import { isValidCountry } from "@/lib/countries";
import { getLang } from "@/lib/i18n/server";
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

  return (
    <div className="mx-auto max-w-6xl">
      <NavBar country={country} lang={lang} />
      {children}
    </div>
  );
}
