import 'dart:io' show Platform;

import 'package:flutter/foundation.dart' show kIsWeb;

const kPrivacyPolicyUrl = 'https://tunetrend.vercel.app/privacy';

// Not published yet (applicationId: com.tunetrend.tunetrend_mobile).
// Fill in with the real listing URL once TuneTrend ships on the Play Store —
// e.g. https://play.google.com/store/apps/details?id=com.tunetrend.tunetrend_mobile
const kPlayStoreUrl = '';

// Not published yet (bundle id: com.tunetrend.tunetrendMobile).
// Fill in with https://apps.apple.com/app/idXXXXXXXXXX once TuneTrend ships
// on the App Store.
const kAppStoreUrl = '';

/// The store link to promote alongside shared content, for the current
/// platform — Play Store on Android, App Store on iOS, once each is
/// actually published. Null until the matching constant above is filled in.
String? get appDownloadUrl {
  if (kIsWeb) return null;
  if (Platform.isAndroid) return kPlayStoreUrl.isEmpty ? null : kPlayStoreUrl;
  if (Platform.isIOS) return kAppStoreUrl.isEmpty ? null : kAppStoreUrl;
  return null;
}
