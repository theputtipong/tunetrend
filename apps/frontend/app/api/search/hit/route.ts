import { jsonError, proxyPublic, readJsonBody, rejectCrossOrigin } from "@/lib/backendProxy";
import { isValidCountry } from "@/lib/countries";
import { VIDEO_ID_REGEX } from "@/lib/playlist";
import { isValidSearchQuery } from "@/lib/search";

export async function POST(request: Request) {
  const forbidden = rejectCrossOrigin(request);
  if (forbidden) return forbidden;
  const { body, error } = await readJsonBody(request);
  if (error) return error;
  const videoId = body?.videoId;
  const country = typeof body?.country === "string" ? body.country.toUpperCase() : "";
  const q = typeof body?.q === "string" ? body.q.trim() : "";

  if (typeof videoId !== "string" || !VIDEO_ID_REGEX.test(videoId)) {
    return jsonError(400, "invalid videoId");
  }
  if (!isValidCountry(country)) return jsonError(400, "invalid country");
  if (!isValidSearchQuery(q)) return jsonError(400, "q must be 2-100 characters");

  return proxyPublic(request, "/search/hit", { method: "POST", body: { videoId, country, q } });
}
