import 'package:firebase_analytics/firebase_analytics.dart';
import 'package:flutter/foundation.dart';

class AnalyticsService {
  factory AnalyticsService() => _instance;

  AnalyticsService._();

  static final AnalyticsService _instance = AnalyticsService._();

  final FirebaseAnalytics _analytics = FirebaseAnalytics.instance;

  late final FirebaseAnalyticsObserver observer = FirebaseAnalyticsObserver(
    analytics: _analytics,
  );

  Future<void> logCountryChanged(String countryCode) {
    return _log('country_changed', {'country_code': countryCode});
  }

  Future<void> logCategorySelected(String categoryId) {
    return _log('category_selected', {
      'category_id': categoryId.isEmpty ? 'all' : categoryId,
    });
  }

  Future<void> logVideoPlayed({
    required String videoId,
    required String title,
    required String countryCode,
    required String categoryId,
    required String tab,
  }) {
    return _log('video_played', {
      'video_id': videoId,
      'title': title,
      'country_code': countryCode,
      'category_id': categoryId.isEmpty ? 'none' : categoryId,
      'tab': tab,
    });
  }

  Future<void> _log(String name, Map<String, Object> parameters) async {
    try {
      await _analytics.logEvent(name: name, parameters: parameters);
    } catch (error, stack) {
      debugPrint('AnalyticsService: failed to log "$name": $error\n$stack');
    }
  }
}
