import { NextResponse } from "next/server";
import { fetchSongs } from "@/lib/api";
import { isValidCountry } from "@/lib/countries";
import { resolveTab } from "@/lib/tabs";

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const country = (searchParams.get("country") ?? "").toUpperCase();
  const tab = resolveTab(searchParams.get("tab"));
  const category = searchParams.get("category") ?? undefined;

  if (!isValidCountry(country)) {
    return NextResponse.json({ success: false, error: "invalid_country" }, { status: 400 });
  }

  try {
    const songs = await fetchSongs(country, tab, category);
    return NextResponse.json({ success: true, data: songs });
  } catch {
    return NextResponse.json({ success: false, error: "upstream_error" }, { status: 502 });
  }
}
