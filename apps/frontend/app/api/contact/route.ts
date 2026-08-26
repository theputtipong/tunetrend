import { NextResponse } from "next/server";
import { contactRatelimit, getClientIp } from "@/lib/ratelimit";
import { isValidEmail, isValidThaiPhone, MAX_MESSAGE_LEN, MAX_NAME_LEN } from "@/lib/validation";

const API_BASE_URL = process.env.API_BASE_URL ?? "http://localhost:8080";

type ContactBody = {
  name?: string;
  message?: string;
  contactEmail?: string;
  contactPhone?: string;
  website?: string;
};

function upstreamErrorResponse(status: number) {
  if (status === 429) {
    return NextResponse.json({ success: false, error: "rate_limited" }, { status: 429 });
  }
  if (status === 400) {
    return NextResponse.json({ success: false, error: "invalid_input" }, { status: 400 });
  }
  return NextResponse.json({ success: false, error: "upstream_error" }, { status: 502 });
}

type CleanContactBody = {
  name: string;
  message: string;
  contactEmail: string;
  contactPhone: string;
};

function sanitizeAndValidateBody(body: ContactBody): CleanContactBody | null {
  const message = (body.message ?? "").trim();
  const name = (body.name ?? "").trim();
  const contactEmail = (body.contactEmail ?? "").trim();
  const contactPhone = (body.contactPhone ?? "").trim();

  const isValid =
    message.length > 0 &&
    message.length <= MAX_MESSAGE_LEN &&
    name.length <= MAX_NAME_LEN &&
    (!!contactEmail || !!contactPhone) &&
    (!contactEmail || isValidEmail(contactEmail)) &&
    (!contactPhone || isValidThaiPhone(contactPhone));

  return isValid ? { name, message, contactEmail, contactPhone } : null;
}

export async function POST(request: Request) {
  const ip = getClientIp(request);
  if (ip) {
    try {
      const { success } = await contactRatelimit.limit(ip);
      if (!success) {
        return NextResponse.json({ success: false, error: "rate_limited" }, { status: 429 });
      }
    } catch {}
  }

  let body: ContactBody;
  try {
    body = await request.json();
  } catch {
    return NextResponse.json({ success: false, error: "invalid_input" }, { status: 400 });
  }

  if (body.website) {
    return NextResponse.json({ success: true });
  }

  const clean = sanitizeAndValidateBody(body);
  if (!clean) {
    return NextResponse.json({ success: false, error: "invalid_input" }, { status: 400 });
  }
  const { message, name, contactEmail, contactPhone } = clean;

  try {
    const res = await fetch(`${API_BASE_URL}/contact`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, message, contactEmail, contactPhone, website: "" }),
      signal: AbortSignal.timeout(10_000),
    });

    const upstream = await res.json().catch(() => null);
    if (!res.ok || !upstream?.success) {
      return upstreamErrorResponse(res.status);
    }

    return NextResponse.json({ success: true });
  } catch {
    return NextResponse.json({ success: false, error: "network_error" }, { status: 502 });
  }
}
