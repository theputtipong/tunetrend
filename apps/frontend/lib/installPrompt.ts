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
  if (raw === IosInstallMode.AppStore || raw === IosInstallMode.Pwa || raw === IosInstallMode.Disabled) {
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
