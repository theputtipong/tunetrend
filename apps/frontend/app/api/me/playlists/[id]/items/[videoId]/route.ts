import { jsonError, proxyAuthed, rejectCrossOrigin } from "@/lib/backendProxy";
import { UUID_REGEX, VIDEO_ID_REGEX } from "@/lib/playlist";

export async function DELETE(
  request: Request,
  { params }: { params: Promise<{ id: string; videoId: string }> },
) {
  const forbidden = rejectCrossOrigin(request);
  if (forbidden) return forbidden;
  const { id, videoId } = await params;
  if (!UUID_REGEX.test(id)) return jsonError(404, "playlist not found");
  if (!VIDEO_ID_REGEX.test(videoId)) return jsonError(404, "playlist item not found");
  return proxyAuthed(request, `/me/playlists/${id}/items/${videoId}`, { method: "DELETE" });
}
