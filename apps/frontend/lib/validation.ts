export const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export const THAI_PHONE_REGEX = /^0(?:[689]\d{8}|[2-57]\d{7})$/;

export const MIN_MESSAGE_LEN = 10;
export const MAX_MESSAGE_LEN = 2000;
export const MAX_NAME_LEN = 100;

export function normalizeThaiPhone(raw: string): string {
  return raw.replace(/[\s()-]/g, "");
}

export function isValidEmail(v: string): boolean {
  return EMAIL_REGEX.test(v.trim());
}

export function isValidThaiPhone(v: string): boolean {
  return THAI_PHONE_REGEX.test(normalizeThaiPhone(v));
}
