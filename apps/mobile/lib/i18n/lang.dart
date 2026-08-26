import 'dart:ui';

import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

enum AppLang { en, th }

AppLang detectLangFromLocale(Locale locale) {
  return locale.languageCode.toLowerCase() == 'th' ? AppLang.th : AppLang.en;
}

class LangController extends ValueNotifier<AppLang?> {
  LangController._() : super(null);

  static final instance = LangController._();

  static const _prefKey = 'tunetrend_lang_override';

  Future<void> load() async {
    final prefs = await SharedPreferences.getInstance();
    final stored = prefs.getString(_prefKey);
    if (stored == 'en') {
      value = AppLang.en;
    } else if (stored == 'th') {
      value = AppLang.th;
    }
  }

  AppLang resolve() {
    return value ?? detectLangFromLocale(PlatformDispatcher.instance.locale);
  }

  Future<void> toggle() async {
    final next = resolve() == AppLang.en ? AppLang.th : AppLang.en;
    value = next;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_prefKey, next == AppLang.en ? 'en' : 'th');
  }
}
