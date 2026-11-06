"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import Link from "next/link";
import Script from "next/script";
import { dictionaries, type Lang } from "@/lib/i18n";
import type { User } from "@/types/user";
import { LogoMark } from "../icons";
import { Modal } from "./Modal";

export type LoginReason = "default" | "playlist" | "youtube";

let initializedClientId: string | null = null;
let credentialHandler: ((response: GoogleCredentialResponse) => void) | null = null;

function isDarkTheme(): boolean {
  const explicit = document.documentElement.getAttribute("data-theme");
  if (explicit) return explicit === "dark";
  return window.matchMedia("(prefers-color-scheme: dark)").matches;
}

export function LoginModal({
  reason,
  lang,
  clientId,
  onClose,
  onLoggedIn,
}: Readonly<{
  reason: LoginReason;
  lang: Lang;
  clientId: string;
  onClose: () => void;
  onLoggedIn: (user: User) => void;
}>) {
  const t = dictionaries[lang].auth;
  const buttonRef = useRef<HTMLDivElement>(null);
  const [gsiReady, setGsiReady] = useState(
    () => typeof window !== "undefined" && Boolean(window.google?.accounts?.id),
  );
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const title =
    reason === "playlist"
      ? t.playlistTitle
      : reason === "youtube"
        ? t.youtubeTitle
        : t.defaultTitle;
  const body =
    reason === "playlist" ? t.playlistBody : reason === "youtube" ? t.youtubeBody : t.defaultBody;

  const handleCredential = useCallback(
    async (response: GoogleCredentialResponse) => {
      setSubmitting(true);
      setError(null);
      try {
        const res = await fetch("/api/auth/google", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ credential: response.credential }),
        });
        const payload = await res.json().catch(() => null);
        if (!res.ok || !payload?.success || !payload.data?.user) throw new Error("login failed");
        onLoggedIn(payload.data.user as User);
      } catch {
        setError(t.loginFailed);
      } finally {
        setSubmitting(false);
      }
    },
    [onLoggedIn, t.loginFailed],
  );

  useEffect(() => {
    credentialHandler = handleCredential;
    return () => {
      if (credentialHandler === handleCredential) credentialHandler = null;
    };
  }, [handleCredential]);

  useEffect(() => {
    const gsi = window.google?.accounts?.id;
    if (!gsiReady || !gsi || !buttonRef.current) return;
    if (initializedClientId !== clientId) {
      gsi.initialize({
        client_id: clientId,
        callback: (response) => credentialHandler?.(response),
        ux_mode: "popup",
        cancel_on_tap_outside: true,
      });
      initializedClientId = clientId;
    }
    buttonRef.current.replaceChildren();
    gsi.renderButton(buttonRef.current, {
      type: "standard",
      theme: isDarkTheme() ? "filled_black" : "outline",
      size: "large",
      text: "continue_with",
      shape: "pill",
      logo_alignment: "left",
      width: 280,
      locale: lang,
    });
  }, [gsiReady, clientId, lang]);

  return (
    <Modal labelledBy="login-title" onClose={onClose}>
      <Script
        src="https://accounts.google.com/gsi/client"
        strategy="afterInteractive"
        onReady={() => setGsiReady(true)}
        onError={() => setError(t.loginFailed)}
      />
      <div className="flex flex-col items-center gap-3 text-center">
        <LogoMark size={40} />
        <h2 id="login-title" className="font-display text-xl font-bold">
          {title}
        </h2>
        <p className="text-sm leading-relaxed text-[var(--text-2)]">{body}</p>
        <div className="flex min-h-[44px] items-center justify-center" aria-busy={!gsiReady}>
          <div ref={buttonRef} hidden={!gsiReady || submitting} />
          {!gsiReady || submitting ? (
            <span className="text-sm text-[var(--text-3)]">{t.loadingGoogle}</span>
          ) : null}
        </div>
        {error ? (
          <p role="alert" className="text-sm text-[var(--error-text)]">
            {error}
          </p>
        ) : null}
        <p className="text-xs leading-relaxed text-[var(--text-3)]">
          {t.finePrint} ·{" "}
          <Link href="/privacy" onClick={onClose}>
            {t.privacyLink}
          </Link>
        </p>
        <button type="button" className="btn btn--sm mt-1" onClick={onClose}>
          {t.close}
        </button>
      </div>
    </Modal>
  );
}
