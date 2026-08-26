import 'package:flutter_test/flutter_test.dart';

import 'package:tunetrend_mobile/main.dart';

void main() {
  testWidgets('renders the TuneTrend header and tabs', (WidgetTester tester) async {
    await tester.pumpWidget(const TuneTrendApp());
    await tester.pump();

    expect(find.text('TuneTrend'), findsOneWidget);
    expect(find.text('Trending'), findsOneWidget);
    expect(find.text('New Releases'), findsOneWidget);
    expect(find.text('Music Videos'), findsOneWidget);
  });
}
