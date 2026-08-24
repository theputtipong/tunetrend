import { redirect } from "next/navigation";
import { DEFAULT_COUNTRY } from "@/lib/countries";

export default function RootPage() {
  redirect(`/${DEFAULT_COUNTRY.toLowerCase()}`);
}
