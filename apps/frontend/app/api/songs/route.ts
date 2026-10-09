import { NextResponse } from "next/server";
import { fetchSongs } from "@/lib/api";
import { isValidCountry } from "@/lib/countries";
import { resolveTab } from "@/lib/tabs";

// Sites allowed to read the chart cross-origin (pdouvch.com shows a live top 5).
const CORS_ORIGINS = new Set(["https://pdouvch.com", "http://localhost:3000"]);

function corsHeaders(request: Request): HeadersInit {
  const origin = request.headers.get("origin");
  if (!origin || !CORS_ORIGINS.has(origin)) return { Vary: "Origin" };
  return { "Access-Control-Allow-Origin": origin, Vary: "Origin" };
}

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const country = (searchParams.get("country") ?? "").toUpperCase();
  const tab = resolveTab(searchParams.get("tab"));
  const category = searchParams.get("category") ?? undefined;
  const headers = corsHeaders(request);

  if (!isValidCountry(country)) {
    return NextResponse.json(
      { success: false, error: "invalid_country" },
      { status: 400, headers },
    );
  }

  try {
    const songs = await fetchSongs(country, tab, category);
    return NextResponse.json({ success: true, data: songs }, { headers });
  } catch {
    return NextResponse.json({ success: false, error: "upstream_error" }, { status: 502, headers });
  }
}
