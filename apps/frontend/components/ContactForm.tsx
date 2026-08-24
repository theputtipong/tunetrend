"use client";

import { useEffect, useState } from "react";
import type { Dictionary } from "@/lib/i18n/en";
import { CheckIcon } from "./icons";
import { isValidEmail, isValidThaiPhone, MAX_MESSAGE_LEN, MIN_MESSAGE_LEN } from "@/lib/validation";

type Method = "email" | "phone";
type Status = "idle" | "loading" | "success" | "error";

const AUTO_CLOSE_SECONDS = 3;

export function ContactForm({
  t,
  onClose,
}: Readonly<{ t: Dictionary["about"]["contact"]; onClose?: () => void }>) {
  const [name, setName] = useState("");
  const [message, setMessage] = useState("");
  const [method, setMethod] = useState<Method>("email");
  const [contactValue, setContactValue] = useState("");
  const [website, setWebsite] = useState("");
  const [status, setStatus] = useState<Status>("idle");
  const [errorCode, setErrorCode] = useState<string | null>(null);
  const [fieldError, setFieldError] = useState<{ message?: string; contact?: string }>({});
  const [secondsLeft, setSecondsLeft] = useState(AUTO_CLOSE_SECONDS);

  useEffect(() => {
    if (status !== "success") return;
    if (secondsLeft <= 0) {
      onClose?.();
      return;
    }
    const timer = setTimeout(() => setSecondsLeft((s) => s - 1), 1000);
    return () => clearTimeout(timer);
  }, [status, secondsLeft, onClose]);

  function validate(): boolean {
    const errors: { message?: string; contact?: string } = {};

    if (!message.trim()) {
      errors.message = t.errors.messageRequired;
    } else if (message.trim().length < MIN_MESSAGE_LEN) {
      errors.message = t.errors.tooShort;
    } else if (message.length > MAX_MESSAGE_LEN) {
      errors.message = t.errors.tooLong;
    }

    if (method === "email" && !isValidEmail(contactValue)) {
      errors.contact = t.errors.invalidEmail;
    } else if (method === "phone" && !isValidThaiPhone(contactValue)) {
      errors.contact = t.errors.invalidPhone;
    }

    setFieldError(errors);
    return Object.keys(errors).length === 0;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!validate()) return;

    setStatus("loading");
    setErrorCode(null);

    try {
      const res = await fetch("/api/contact", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name,
          message,
          contactEmail: method === "email" ? contactValue : "",
          contactPhone: method === "phone" ? contactValue : "",
          website,
        }),
      });

      const body = await res.json();
      if (!body.success) {
        setStatus("error");
        setErrorCode(body.error ?? "generic");
        return;
      }

      setStatus("success");
      setName("");
      setMessage("");
      setContactValue("");
    } catch {
      setStatus("error");
      setErrorCode("network_error");
    }
  }

  const errorMessage = errorCode === "rate_limited" ? t.errors.rateLimited : t.errors.generic;

  if (status === "success") {
    return (
      <div className="flex w-full flex-col items-center gap-3 py-6 text-center">
        <div
          className="flex h-14 w-14 items-center justify-center rounded-full"
          style={{ backgroundColor: "var(--accent)", color: "var(--accent-ink)" }}
        >
          <CheckIcon size={28} />
        </div>
        <p className="text-[15px] font-semibold">{t.successMessage}</p>
        <p className="text-sm text-[var(--text-3)]">{t.closingIn(secondsLeft)}</p>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="flex w-full flex-col gap-3">
      <div
        style={{ position: "absolute", left: "-9999px", width: 1, height: 1, overflow: "hidden" }}
        aria-hidden="true"
      >
        <input
          tabIndex={-1}
          autoComplete="off"
          value={website}
          onChange={(e) => setWebsite(e.target.value)}
        />
      </div>

      <label className="flex flex-col gap-1.5 text-sm">
        <span className="text-[var(--text-2)]">{t.nameLabel}</span>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder={t.namePlaceholder}
          className="rounded-lg border border-[var(--border)] bg-[var(--bg-elev)] px-3 py-2 text-[15px] outline-none focus:border-[var(--accent)]"
        />
      </label>

      <label className="flex flex-col gap-1.5 text-sm">
        <span className="text-[var(--text-2)]">{t.messageLabel}</span>
        <textarea
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder={t.messagePlaceholder}
          rows={4}
          className="rounded-lg border border-[var(--border)] bg-[var(--bg-elev)] px-3 py-2 text-[15px] outline-none focus:border-[var(--accent)]"
        />
        {fieldError.message && (
          <span className="text-xs text-red-500">{fieldError.message}</span>
        )}
      </label>

      <div className="flex flex-col gap-1.5 text-sm">
        <span className="text-[var(--text-2)]">{t.methodLabel}</span>
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => setMethod("email")}
            className={method === "email" ? "pill pill--active" : "pill"}
          >
            {t.methodEmail}
          </button>
          <button
            type="button"
            onClick={() => setMethod("phone")}
            className={method === "phone" ? "pill pill--active" : "pill"}
          >
            {t.methodPhone}
          </button>
        </div>
        <input
          type={method === "email" ? "email" : "tel"}
          value={contactValue}
          onChange={(e) => setContactValue(e.target.value)}
          placeholder={method === "email" ? t.contactEmailPlaceholder : t.contactPhonePlaceholder}
          className="rounded-lg border border-[var(--border)] bg-[var(--bg-elev)] px-3 py-2 text-[15px] outline-none focus:border-[var(--accent)]"
        />
        {fieldError.contact && (
          <span className="text-xs text-red-500">{fieldError.contact}</span>
        )}
      </div>

      <button
        type="submit"
        disabled={status === "loading"}
        className="retry w-fit"
        style={{ marginTop: 0 }}
      >
        {status === "loading" ? t.submitting : t.submit}
      </button>

      {status === "error" && <p className="text-sm text-red-500">{errorMessage}</p>}
    </form>
  );
}
