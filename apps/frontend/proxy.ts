import { NextResponse, type NextRequest } from "next/server";
import { detectCountryFromAcceptLanguage } from "@/lib/countries";
import { LANG_COOKIE, detectLangFromAcceptLanguage } from "@/lib/i18n";
import { getClientIp, ratelimit } from "@/lib/ratelimit";

export async function proxy(request: NextRequest) {
  const ip = getClientIp(request);
  if (ip) {
    try {
      const { success, limit, remaining, reset } = await ratelimit.limit(ip);
      if (!success) {
        return new NextResponse("Too many requests. Please slow down and try again shortly.", {
          status: 429,
          headers: {
            "Retry-After": Math.max(0, Math.ceil((reset - Date.now()) / 1000)).toString(),
            "X-RateLimit-Limit": limit.toString(),
            "X-RateLimit-Remaining": remaining.toString(),
          },
        });
      }
    } catch {}
  }

  const response =
    request.nextUrl.pathname === "/"
      ? NextResponse.redirect(
          new URL(
            `/${detectCountryFromAcceptLanguage(request.headers.get("accept-language")).toLowerCase()}`,
            request.url,
          ),
        )
      : NextResponse.next();

  if (!request.cookies.has(LANG_COOKIE)) {
    const lang = detectLangFromAcceptLanguage(request.headers.get("accept-language"));
    response.cookies.set(LANG_COOKIE, lang, { path: "/", maxAge: 60 * 60 * 24 * 365 });
  }

  return response;
}

export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico|icon.svg|apple-icon.png).*)"],
};
