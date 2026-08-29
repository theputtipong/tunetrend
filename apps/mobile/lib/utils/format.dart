import 'package:intl/intl.dart';

import '../i18n/dictionary.dart';
import '../i18n/lang.dart';

final NumberFormat _compactViews = NumberFormat.compact(locale: 'en_US');

String formatViewCount(String viewCount, AppLang lang) {
  final n = num.tryParse(viewCount);
  final suffix = dictionaries[lang]!.viewsSuffix;
  if (n == null) return '— $suffix';
  return '${_compactViews.format(n)} $suffix';
}

String formatRelativeTime(String iso, AppLang lang) {
  final date = DateTime.tryParse(iso);
  if (date == null) return '';

  final t = dictionaries[lang]!;
  final diffMinutes = DateTime.now().difference(date).inMinutes;
  if (diffMinutes < 60) return t.minutesAgo(diffMinutes < 0 ? 0 : diffMinutes);

  final diffHours = (diffMinutes / 60).round();
  if (diffHours < 24) return t.hoursAgo(diffHours);

  final diffDays = (diffHours / 24).round();
  if (diffDays < 7) return t.daysAgo(diffDays);

  final diffWeeks = (diffDays / 7).round();
  return t.weeksAgo(diffWeeks);
}
