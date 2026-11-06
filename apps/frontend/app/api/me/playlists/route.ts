import { jsonError, proxyAuthed, readJsonBody, rejectCrossOrigin } from "@/lib/backendProxy";
import { isValidPlaylistName } from "@/lib/playlist";

export async function GET(request: Request) {
  return proxyAuthed(request, "/me/playlists", { method: "GET" });
}

export async function POST(request: Request) {
  const forbidden = rejectCrossOrigin(request);
  if (forbidden) return forbidden;
  const { body, error } = await readJsonBody(request);
  if (error) return error;
  if (!isValidPlaylistName(body?.name)) return jsonError(422, "name must be 1-60 characters");
  return proxyAuthed(request, "/me/playlists", {
    method: "POST",
    body: { name: (body.name as string).trim() },
  });
}
