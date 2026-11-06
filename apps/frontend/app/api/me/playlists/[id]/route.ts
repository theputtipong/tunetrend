import { jsonError, proxyAuthed, readJsonBody, rejectCrossOrigin } from "@/lib/backendProxy";
import { UUID_REGEX, isValidPlaylistName } from "@/lib/playlist";

type Ctx = { params: Promise<{ id: string }> };

export async function GET(request: Request, { params }: Ctx) {
  const { id } = await params;
  if (!UUID_REGEX.test(id)) return jsonError(404, "playlist not found");
  return proxyAuthed(request, `/me/playlists/${id}`, { method: "GET" });
}

export async function PATCH(request: Request, { params }: Ctx) {
  const forbidden = rejectCrossOrigin(request);
  if (forbidden) return forbidden;
  const { id } = await params;
  if (!UUID_REGEX.test(id)) return jsonError(404, "playlist not found");
  const { body, error } = await readJsonBody(request);
  if (error) return error;
  if (!isValidPlaylistName(body?.name)) return jsonError(422, "name must be 1-60 characters");
  return proxyAuthed(request, `/me/playlists/${id}`, {
    method: "PATCH",
    body: { name: (body.name as string).trim() },
  });
}

export async function DELETE(request: Request, { params }: Ctx) {
  const forbidden = rejectCrossOrigin(request);
  if (forbidden) return forbidden;
  const { id } = await params;
  if (!UUID_REGEX.test(id)) return jsonError(404, "playlist not found");
  return proxyAuthed(request, `/me/playlists/${id}`, { method: "DELETE" });
}
