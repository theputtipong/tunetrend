import { cookies } from "next/headers";
import { DEFAULT_LANG, LANG_COOKIE, isValidLang, type Lang } from "./index";

export async function getLang(): Promise<Lang> {
  const store = await cookies();
  const value = store.get(LANG_COOKIE)?.value;
  return isValidLang(value) ? value : DEFAULT_LANG;
}
