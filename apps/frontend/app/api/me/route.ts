import { cookies } from "next/headers";
import { proxyAuthed, rejectCrossOrigin } from "@/lib/backendProxy";
import { SESSION_COOKIE } from "@/lib/session";

export async function GET(request: Request) {
  return proxyAuthed(request, "/me", { method: "GET" });
}

export async function DELETE(request: Request) {
  const forbidden = rejectCrossOrigin(request);
  if (forbidden) return forbidden;
  const res = await proxyAuthed(request, "/me", { method: "DELETE" });
  if (res.ok) (await cookies()).delete(SESSION_COOKIE);
  return res;
}
