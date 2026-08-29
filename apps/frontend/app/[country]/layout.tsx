import { notFound } from "next/navigation";
import { isValidCountry } from "@/lib/countries";

export default async function CountryLayout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ country: string }>;
}>) {
  const { country: rawCountry } = await params;
  if (!isValidCountry(rawCountry.toUpperCase())) {
    notFound();
  }

  return children;
}
