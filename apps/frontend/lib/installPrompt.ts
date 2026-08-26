export enum AndroidInstallMode {
  PlayStore = "play-store",
  Pwa = "pwa",
  Disabled = "disabled",
}

export enum IosInstallMode {
  AppStore = "app-store",
  Pwa = "pwa",
  Disabled = "disabled",
}

function readAndroidMode(): AndroidInstallMode {
  const raw = process.env.NEXT_PUBLIC_INSTALL_ANDROID_MODE;
  if (
    raw === AndroidInstallMode.PlayStore ||
    raw === AndroidInstallMode.Pwa ||
    raw === AndroidInstallMode.Disabled
  ) {
    return raw;
  }
  return AndroidInstallMode.PlayStore;
}

function readIosMode(): IosInstallMode {
  const raw = process.env.NEXT_PUBLIC_INSTALL_IOS_MODE;
  if (
    raw === IosInstallMode.AppStore ||
    raw === IosInstallMode.Pwa ||
    raw === IosInstallMode.Disabled
  ) {
    return raw;
  }
  return IosInstallMode.Pwa;
}

export const INSTALL_PROMPT_CONFIG: { android: AndroidInstallMode; ios: IosInstallMode } = {
  android: readAndroidMode(),
  ios: readIosMode(),
};

export const PLAY_STORE_URL =
  process.env.NEXT_PUBLIC_PLAY_STORE_URL ??
  "https://play.google.com/store/apps/details?id=com.tunetrend.tunetrend_mobile";

export const APP_STORE_URL = process.env.NEXT_PUBLIC_APP_STORE_URL ?? "";

export enum MobilePlatform {
  Android = "android",
  Ios = "ios",
}

export function detectMobilePlatform(userAgent: string): MobilePlatform | null {
  if (/android/i.test(userAgent)) return MobilePlatform.Android;
  if (/iphone|ipad|ipod/i.test(userAgent)) return MobilePlatform.Ios;
  return null;
}

export interface BeforeInstallPromptEvent extends Event {
  prompt(): Promise<void>;
  userChoice: Promise<{ outcome: "accepted" | "dismissed"; platform: string }>;
}

export function isStandaloneDisplayMode(): boolean {
  const nav = navigator as Navigator & { standalone?: boolean };
  return window.matchMedia("(display-mode: standalone)").matches || nav.standalone === true;
}

const LAST_SHOWN_KEY = "tunetrend-install-prompt-last-shown";
const INSTALLED_KEY = "tunetrend-install-prompt-installed";
const REPROMPT_INTERVAL_MS = 24 * 60 * 60 * 1000;

export function shouldReprompt(): boolean {
  const last = localStorage.getItem(LAST_SHOWN_KEY);
  if (!last) return true;
  return Date.now() - Number(last) >= REPROMPT_INTERVAL_MS;
}

export function markShownNow(): void {
  localStorage.setItem(LAST_SHOWN_KEY, String(Date.now()));
}

export function markInstalled(): void {
  localStorage.setItem(INSTALLED_KEY, "true");
}

export function isMarkedInstalled(): boolean {
  return localStorage.getItem(INSTALLED_KEY) === "true";
}

export type PromptKind = "android-pwa" | "android-store" | "ios-pwa" | "ios-store";

export function computeEligibleKind(
  platform: MobilePlatform,
  hasAndroidInstallEvent: boolean,
): PromptKind | null {
  if (platform === MobilePlatform.Android) {
    if (INSTALL_PROMPT_CONFIG.android === AndroidInstallMode.Pwa) {
      return hasAndroidInstallEvent ? "android-pwa" : null;
    }
    if (INSTALL_PROMPT_CONFIG.android === AndroidInstallMode.PlayStore) {
      return "android-store";
    }
    return null;
  }

  if (INSTALL_PROMPT_CONFIG.ios === IosInstallMode.Pwa) {
    return "ios-pwa";
  }
  if (INSTALL_PROMPT_CONFIG.ios === IosInstallMode.AppStore && APP_STORE_URL) {
    return "ios-store";
  }
  return null;
}
