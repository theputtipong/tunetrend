export function resolveCategory(raw: string | undefined | null, available: ReadonlyArray<{ id: string }>): string {
  if (!raw) return "";
  return available.some((c) => c.id === raw) ? raw : "";
}
