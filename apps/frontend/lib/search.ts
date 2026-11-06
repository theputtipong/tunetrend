export const MIN_QUERY_LEN = 2;
export const MAX_QUERY_LEN = 100;

export function isValidSearchQuery(q: string): boolean {
  const len = Array.from(q.trim()).length;
  return len >= MIN_QUERY_LEN && len <= MAX_QUERY_LEN;
}

export function isValidCategoryId(id: string): boolean {
  return /^\d{1,3}$/.test(id);
}
