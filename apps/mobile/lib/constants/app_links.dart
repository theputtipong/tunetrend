import 'dart:io' show Platform;

import 'package:flutter/foundation.dart' show kIsWeb;

import 'tabs.dart';

const kWebBaseUrl = 'https://tunetrend.pdouvch.com';

const kPrivacyPolicyUrl = '$kWebBaseUrl/privacy';

const kPlayStoreUrl =
    'https://play.google.com/store/apps/details?id=com.tunetrend.tunetrend_mobile';

const kAppStoreUrl = '';

String? get appDownloadUrl {
  if (kIsWeb) return null;
  if (Platform.isAndroid) return kPlayStoreUrl.isEmpty ? null : kPlayStoreUrl;
  if (Platform.isIOS) return kAppStoreUrl.isEmpty ? null : kAppStoreUrl;
  return null;
}

String watchPageUrl(String country, String videoId, TrendTab tab) =>
    '$kWebBaseUrl/${country.toLowerCase()}/watch/$videoId?tab=${tab.key}';
