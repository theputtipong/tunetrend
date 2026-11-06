import { NextResponse } from "next/server";
import { cookies } from "next/headers";
import { jsonError, readJsonBody, rejectCrossOrigin } from "@/lib/backendProxy";
import { SESSION_COOKIE, backendFetch, sessionCookieOptions } from "@/lib/session";

const MAX_CREDENTIAL_LEN = 4096;

export async function POST(request: Request) {
  const forbidden = rejectCrossOrigin(request);
  if (forbidden) return forbidden;
  const { body, error } = await readJsonBody(request);
  if (error) return error;
  const credential = body?.credential;
  if (typeof credential !== "string" || !credential || credential.length > MAX_CREDENTIAL_LEN) {
    return jsonError(400, "invalid_input");
  }

  let res: Response;
  try {
    res = await backendFetch("/auth/google", {
      method: "POST",
      body: JSON.stringify({ idToken: credential }),
    });
  } catch {
    return jsonError(502, "network_error");
  }

  const payload = await res.json().catch(() => null);
  const token = payload?.data?.token;
  if (!res.ok || !payload?.success || typeof token !== "string") {
    return jsonError(res.ok ? 502 : res.status, payload?.error ?? "upstream_error");
  }

  (await cookies()).set(SESSION_COOKIE, token, sessionCookieOptions());
  return NextResponse.json({ success: true, data: { user: payload.data.user } });
}
