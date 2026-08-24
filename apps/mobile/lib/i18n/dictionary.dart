import 'en.dart';
import 'lang.dart';
import 'strings.dart';
import 'th.dart';

const Map<AppLang, AppStrings> dictionaries = {
  AppLang.en: enStrings,
  AppLang.th: thStrings,
};

AppStrings get currentStrings => dictionaries[LangController.instance.resolve()]!;
