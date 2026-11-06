import { NextResponse } from "next/server";
import { cookies } from "next/headers";
import { rejectCrossOrigin } from "@/lib/backendProxy";
import { SESSION_COOKIE } from "@/lib/session";

export async function POST(request: Request) {
  const forbidden = rejectCrossOrigin(request);
  if (forbidden) return forbidden;
  (await cookies()).delete(SESSION_COOKIE);
  return NextResponse.json({ success: true });
}
