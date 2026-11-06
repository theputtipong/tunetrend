"use client";

import { useEffect, useRef, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { dictionaries } from "@/lib/i18n";
import { useAuth } from "./AuthProvider";
import { ConfirmDialog } from "./Modal";

export function HeaderAuth() {
  const { user, lang, loginEnabled, openLogin } = useAuth();
  const t = dictionaries[lang].nav;

  if (user) return <AccountMenu />;
  if (!loginEnabled) return null;

  return (
    <button type="button" className="btn btn--sm" onClick={() => openLogin("default")}>
      <span className="g-dot" aria-hidden="true" />
      <span className="hidden sm:inline">{t.login}</span>
      <span className="sr-only sm:hidden">{t.login}</span>
    </button>
  );
}

function AccountMenu() {
  const { user, lang, logout, deleteAccount, openLogin } = useAuth();
  const t = dictionaries[lang].nav;
  const ta = dictionaries[lang].auth;
  const [isOpen, setIsOpen] = useState(false);
  const [confirming, setConfirming] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!isOpen) return;
    function onClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setIsOpen(false);
    }
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === "Escape") setIsOpen(false);
    }
    document.addEventListener("mousedown", onClickOutside);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("mousedown", onClickOutside);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [isOpen]);

  if (!user) return null;
  const initial = (user.name || user.email || "?").trim().charAt(0).toUpperCase();

  async function onConfirmDelete() {
    setBusy(true);
    setError(null);
    const result = await deleteAccount();
    setBusy(false);
    if (result === "deleted") {
      setConfirming(false);
    } else if (result === "unauthorized") {
      setError(ta.deleteSessionExpired);
      openLogin("default");
    } else {
      setError(ta.deleteFailed);
    }
  }

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        className="avatar"
        aria-label={t.accountMenu}
        aria-haspopup="menu"
        aria-expanded={isOpen}
        onClick={() => setIsOpen((v) => !v)}
      >
        {user.avatarUrl.startsWith("https://") ? (
          <Image
            src={user.avatarUrl}
            alt=""
            fill
            sizes="32px"
            unoptimized
            referrerPolicy="no-referrer"
          />
        ) : (
          initial
        )}
      </button>

      {isOpen && (
        <div
          role="menu"
          className="absolute right-0 top-full z-50 mt-2 min-w-[220px] rounded-xl border border-[var(--border)] bg-[var(--bg-elev)] p-1.5 shadow-lg"
        >
          <div className="px-3 py-2">
            <div className="truncate text-sm font-semibold">{user.name}</div>
            <div className="truncate text-xs text-[var(--text-3)]">{user.email}</div>
          </div>
          <div className="my-1 border-t border-[var(--border)]" />
          <Link
            role="menuitem"
            href="/me/playlists"
            className="menu-item"
            onClick={() => setIsOpen(false)}
          >
            {t.myPlaylists}
          </Link>
          <button
            role="menuitem"
            type="button"
            className="menu-item"
            onClick={() => {
              setIsOpen(false);
              void logout();
            }}
          >
            {t.logout}
          </button>
          <button
            role="menuitem"
            type="button"
            className="menu-item text-[var(--error-text)]"
            onClick={() => {
              setIsOpen(false);
              setError(null);
              setConfirming(true);
            }}
          >
            {t.deleteAccount}
          </button>
        </div>
      )}

      {confirming && (
        <ConfirmDialog
          title={ta.deleteTitle}
          body={ta.deleteBody}
          confirmLabel={ta.deleteConfirm}
          cancelLabel={ta.cancel}
          error={error}
          busy={busy}
          onConfirm={() => void onConfirmDelete()}
          onCancel={() => setConfirming(false)}
        />
      )}
    </div>
  );
}
