"use client";

import { dictionaries } from "@/lib/i18n";
import { useAuth } from "../auth/AuthProvider";
import { EmptyBoxIcon } from "../icons";

export function LoginRequired() {
  const { lang, loginEnabled, openLogin } = useAuth();
  const t = dictionaries[lang];

  return (
    <div className="flex flex-col items-center gap-4 px-8 py-16 text-center">
      <div className="state-icon state-icon--empty">
        <EmptyBoxIcon />
      </div>
      <h1 className="state-title">{t.playlist.loginTitle}</h1>
      <p className="state-subtitle">{t.playlist.loginBody}</p>
      {loginEnabled ? (
        <button type="button" className="btn" onClick={() => openLogin("playlist")}>
          <span className="g-dot" aria-hidden="true" />
          {t.search.loginButton}
        </button>
      ) : null}
    </div>
  );
}
