import type { NextRequest } from "next/server";
import { jsonError, proxyAuthed } from "@/lib/backendProxy";
import { isValidCountry } from "@/lib/countries";
import { isValidSearchQuery, isValidCategoryId } from "@/lib/search";

export async function GET(request: NextRequest) {
  const params = request.nextUrl.searchParams;
  const q = (params.get("q") ?? "").trim();
  const country = (params.get("country") ?? "").toUpperCase();
  const category = params.get("category") ?? "";

  if (!isValidSearchQuery(q)) return jsonError(400, "q must be 2-100 characters");
  if (!isValidCountry(country)) return jsonError(400, "invalid country");
  if (category && !isValidCategoryId(category)) return jsonError(400, "invalid category");

  const qs = new URLSearchParams({ q, country });
  if (category) qs.set("category", category);
  return proxyAuthed(request, `/search/youtube?${qs.toString()}`, { method: "GET" });
}
