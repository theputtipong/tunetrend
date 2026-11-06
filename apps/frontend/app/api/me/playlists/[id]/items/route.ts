import { jsonError, proxyAuthed, readJsonBody, rejectCrossOrigin } from "@/lib/backendProxy";
import {
  MAX_ITEM_CHANNEL_LEN,
  MAX_ITEM_TITLE_LEN,
  UUID_REGEX,
  VIDEO_ID_REGEX,
  runeLength,
} from "@/lib/playlist";

function cleanString(value: unknown): string | null {
  if (value === undefined || value === null) return "";
  return typeof value === "string" ? value.trim() : null;
}

export async function POST(request: Request, { params }: { params: Promise<{ id: string }> }) {
  const forbidden = rejectCrossOrigin(request);
  if (forbidden) return forbidden;
  const { id } = await params;
  if (!UUID_REGEX.test(id)) return jsonError(404, "playlist not found");

  const { body, error } = await readJsonBody(request);
  if (error) return error;
  const videoId = body?.videoId;
  const title = cleanString(body?.title);
  const channelTitle = cleanString(body?.channelTitle);
  const thumbnailUrl = cleanString(body?.thumbnailUrl);

  if (typeof videoId !== "string" || !VIDEO_ID_REGEX.test(videoId)) {
    return jsonError(422, "invalid videoId");
  }
  if (!title || runeLength(title) > MAX_ITEM_TITLE_LEN) {
    return jsonError(422, "title must be 1-200 characters");
  }
  if (channelTitle === null || runeLength(channelTitle) > MAX_ITEM_CHANNEL_LEN) {
    return jsonError(422, "channelTitle must be at most 100 characters");
  }
  if (thumbnailUrl === null || (thumbnailUrl && !thumbnailUrl.startsWith("https://i.ytimg.com/"))) {
    return jsonError(422, "thumbnailUrl must start with https://i.ytimg.com/");
  }

  return proxyAuthed(request, `/me/playlists/${id}/items`, {
    method: "POST",
    body: { videoId, title, channelTitle, thumbnailUrl },
  });
}
