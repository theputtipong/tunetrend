import 'dart:ui';

import '../i18n/lang.dart';

class Country {
  final String code;
  final String label;
  final String labelTh;

  const Country(this.code, this.label, this.labelTh);
}

const List<Country> kCountries = [
  Country('TH', 'Thailand', 'ไทย'),
  Country('KR', 'South Korea', 'เกาหลีใต้'),
  Country('JP', 'Japan', 'ญี่ปุ่น'),
  Country('US', 'United States', 'สหรัฐอเมริกา'),
  Country('GB', 'United Kingdom', 'สหราชอาณาจักร'),
];

const String kDefaultCountry = 'TH';

bool isValidCountry(String code) {
  return kCountries.any((c) => c.code == code);
}

String countryLabel(String code, AppLang lang) {
  for (final country in kCountries) {
    if (country.code == code) {
      return lang == AppLang.th ? country.labelTh : country.label;
    }
  }
  return code;
}

String detectCountryFromLocale(Locale locale) {
  final language = locale.languageCode.toLowerCase();
  final region = (locale.countryCode ?? '').toUpperCase();

  if (language == 'th') return 'TH';
  if (language == 'ko') return 'KR';
  if (language == 'ja') return 'JP';
  if (language == 'en' && region == 'GB') return 'GB';
  if (language == 'en') return 'US';

  return kDefaultCountry;
}
