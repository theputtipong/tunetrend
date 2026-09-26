import 'package:flutter_test/flutter_test.dart';

import 'package:tunetrend_mobile/constants/app_links.dart';
import 'package:tunetrend_mobile/constants/tabs.dart';
import 'package:tunetrend_mobile/i18n/en.dart';
import 'package:tunetrend_mobile/i18n/th.dart';

void main() {
  test('watchPageUrl points at the TuneTrend web player, not youtube.com', () {
    expect(
      watchPageUrl('TH', 'abc123', TrendTab.trending),
      'https://tunetrend.pdouvch.com/th/watch/abc123?tab=trending',
    );
    expect(
      watchPageUrl('kr', 'xyz789', TrendTab.musicVideos),
      'https://tunetrend.pdouvch.com/kr/watch/xyz789?tab=mv',
    );
  });

  test('the store link is filled in, so shares carry a download link', () {
    expect(kPlayStoreUrl, isNotEmpty);

    final url = watchPageUrl('TH', 'abc123', TrendTab.trending);
    for (final strings in [enStrings, thStrings]) {
      final text = strings.shareMessage('Song', url, kPlayStoreUrl);
      expect(text, contains(url));
      expect(text, contains(kPlayStoreUrl));
    }
  });

  test('without a store link the web link still goes out (iOS today)', () {
    final url = watchPageUrl('TH', 'abc123', TrendTab.trending);
    final text = enStrings.shareMessage('Song', url, null);

    expect(text, contains(url));
    expect(text, isNot(contains('play.google.com')));
  });
}
