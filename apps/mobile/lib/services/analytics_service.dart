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

  Future<void> logSearchPerformed({
    required String query,
    required String countryCode,
    required String categoryId,
    required int resultCount,
  }) {
    return _log('search_performed', {
      'query_length': query.trim().length,
      'country_code': countryCode,
      'category_id': categoryId,
      'result_count': resultCount,
    });
  }

  Future<void> logYouTubeSearch({
    required String query,
    required String countryCode,
    required String categoryId,
    required String outcome,
    int resultCount = 0,
    bool cached = false,
  }) {
    return _log('youtube_search', {
      'query_length': query.trim().length,
      'country_code': countryCode,
      'category_id': categoryId,
      'outcome': outcome,
      'result_count': resultCount,
      'cached': cached ? 1 : 0,
    });
  }

  Future<void> logPlaylistCreated({required String source}) {
    return _log('playlist_created', {'source': source});
  }

  Future<void> logVideoAddedToPlaylist({
    required String videoId,
    required String source,
  }) {
    return _log('video_added_to_playlist', {
      'video_id': videoId,
      'source': source,
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
