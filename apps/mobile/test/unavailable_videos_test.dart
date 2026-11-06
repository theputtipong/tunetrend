import 'package:flutter_test/flutter_test.dart';
import 'package:tunetrend_mobile/services/unavailable_videos.dart';

void main() {
  test('marks a video once and notifies listeners', () {
    final set = UnavailableVideos.instance;
    var notified = 0;
    void listener() => notified++;
    set.addListener(listener);

    set.mark('dead1');
    set.mark('dead1');
    set.mark('');

    expect(set.contains('dead1'), isTrue);
    expect(set.contains('alive'), isFalse);
    expect(notified, 1);
    set.removeListener(listener);
  });
}
