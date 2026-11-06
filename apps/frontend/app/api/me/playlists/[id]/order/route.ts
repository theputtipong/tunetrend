import { jsonError, proxyAuthed, readJsonBody, rejectCrossOrigin } from "@/lib/backendProxy";
import { MAX_PLAYLIST_ITEMS, UUID_REGEX, VIDEO_ID_REGEX } from "@/lib/playlist";

const INVALID_ORDER = "videoIds must contain every item in the playlist exactly once";

export async function PUT(request: Request, { params }: { params: Promise<{ id: string }> }) {
  const forbidden = rejectCrossOrigin(request);
  if (forbidden) return forbidden;
  const { id } = await params;
  if (!UUID_REGEX.test(id)) return jsonError(404, "playlist not found");

  const { body, error } = await readJsonBody(request);
  if (error) return error;
  const videoIds = body?.videoIds;
  if (
    !Array.isArray(videoIds) ||
    videoIds.length > MAX_PLAYLIST_ITEMS ||
    !videoIds.every((v) => typeof v === "string" && VIDEO_ID_REGEX.test(v)) ||
    new Set(videoIds).size !== videoIds.length
  ) {
    return jsonError(422, INVALID_ORDER);
  }

  return proxyAuthed(request, `/me/playlists/${id}/order`, {
    method: "PUT",
    body: { videoIds },
  });
}
