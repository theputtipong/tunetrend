"use client";

import { useEffect, useRef, type ReactNode } from "react";

export function Modal({
  labelledBy,
  onClose,
  children,
}: Readonly<{ labelledBy: string; onClose: () => void; children: ReactNode }>) {
  const cardRef = useRef<HTMLDivElement>(null);
  const onCloseRef = useRef(onClose);

  useEffect(() => {
    onCloseRef.current = onClose;
  }, [onClose]);

  useEffect(() => {
    const previous = document.activeElement as HTMLElement | null;
    const focusable = cardRef.current?.querySelector<HTMLElement>(
      "button, [href], input, [tabindex]:not([tabindex='-1'])",
    );
    focusable?.focus();

    function onKeyDown(e: KeyboardEvent) {
      if (e.key === "Escape") onCloseRef.current();
    }
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("keydown", onKeyDown);
      previous?.focus?.();
    };
  }, []);

  return (
    <div
      className="modal-backdrop"
      onMouseDown={(e) => {
        if (e.target === e.currentTarget) onClose();
      }}
    >
      <div
        ref={cardRef}
        role="dialog"
        aria-modal="true"
        aria-labelledby={labelledBy}
        className="modal-card"
      >
        {children}
      </div>
    </div>
  );
}

export function ConfirmDialog({
  title,
  body,
  confirmLabel,
  cancelLabel,
  error,
  busy,
  onConfirm,
  onCancel,
}: Readonly<{
  title: string;
  body: string;
  confirmLabel: string;
  cancelLabel: string;
  error?: string | null;
  busy?: boolean;
  onConfirm: () => void;
  onCancel: () => void;
}>) {
  return (
    <Modal labelledBy="confirm-title" onClose={onCancel}>
      <h2 id="confirm-title" className="font-display text-lg font-bold">
        {title}
      </h2>
      <p className="mt-2 text-sm text-[var(--text-2)]">{body}</p>
      {error ? <p className="mt-3 text-sm text-[var(--error-text)]">{error}</p> : null}
      <div className="mt-5 flex justify-end gap-3">
        <button type="button" className="btn" onClick={onCancel} disabled={busy}>
          {cancelLabel}
        </button>
        <button type="button" className="btn btn--accent" onClick={onConfirm} disabled={busy}>
          {confirmLabel}
        </button>
      </div>
    </Modal>
  );
}
