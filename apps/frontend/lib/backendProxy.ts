import { NextResponse } from "next/server";
import { cookies } from "next/headers";
import { getClientIp } from "./ratelimit";
import { SESSION_COOKIE, backendFetch, forwardedClientHeaders, getSessionToken } from "./session";

export function jsonError(status: number, error: string) {
  return NextResponse.json({ success: false, error }, { status });
}

function expectedOrigin(request: Request): string | null {
  const url = new URL(request.url);
  const host =
    request.headers.get("x-forwarded-host")?.split(",")[0]?.trim() ||
    request.headers.get("host") ||
    url.host;
  const proto =
    request.headers.get("x-forwarded-proto")?.split(",")[0]?.trim() || url.protocol.slice(0, -1);
  if (!host || !proto) return null;
  return `${proto}://${host}`;
}

export function rejectCrossOrigin(request: Request): NextResponse | null {
  if (request.headers.get("sec-fetch-site") === "same-origin") return null;
  const origin = request.headers.get("origin");
  if (origin && origin === expectedOrigin(request)) return null;
  return jsonError(403, "forbidden_origin");
}

type JsonBody =
  | { body: Record<string, unknown> | null; error?: undefined }
  | { body?: undefined; error: NextResponse };

export async function readJsonBody(request: Request): Promise<JsonBody> {
  const contentType = request.headers.get("content-type")?.split(";")[0]?.trim().toLowerCase();
  if (contentType !== "application/json") {
    return { error: jsonError(415, "unsupported_media_type") };
  }
  try {
    const body = await request.json();
    return { body: body && typeof body === "object" && !Array.isArray(body) ? body : null };
  } catch {
    return { body: null };
  }
}

async function forward(
  request: Request,
  path: string,
  init: { method: string; body?: unknown; token?: string | null },
): Promise<NextResponse> {
  const headers = forwardedClientHeaders(getClientIp(request));

  let res: Response;
  try {
    res = await backendFetch(path, {
      method: init.method,
      token: init.token,
      headers,
      body: init.body === undefined ? undefined : JSON.stringify(init.body),
    });
  } catch {
    return jsonError(502, "network_error");
  }

  const payload = await res.json().catch(() => null);
  if (res.status === 401 && init.token) {
    (await cookies()).delete(SESSION_COOKIE);
  }
  return NextResponse.json(payload ?? { success: false, error: "upstream_error" }, {
    status: res.status,
  });
}

export async function proxyPublic(
  request: Request,
  path: string,
  init: { method: string; body?: unknown },
) {
  return forward(request, path, init);
}

export async function proxyAuthed(
  request: Request,
  path: string,
  init: { method: string; body?: unknown },
) {
  const token = await getSessionToken();
  if (!token) return jsonError(401, "unauthorized");
  return forward(request, path, { ...init, token });
}
