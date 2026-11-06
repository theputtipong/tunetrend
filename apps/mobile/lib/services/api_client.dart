import 'package:dio/dio.dart';
import 'package:firebase_performance_dio/firebase_performance_dio.dart';
import 'package:flutter/foundation.dart'
    show VoidCallback, kDebugMode, visibleForTesting;

import '../constants/tabs.dart';
import '../models/category.dart';
import '../models/playlist.dart';
import '../models/search.dart';
import '../models/song.dart';
import '../models/user.dart';
import 'performance_service.dart';

const String _apiBaseUrl = String.fromEnvironment(
  'API_BASE_URL',
  defaultValue: 'http://10.0.2.2:8080',
);

class ApiException implements Exception {
  final String message;
  final int? statusCode;
  ApiException(this.message, {this.statusCode});

  @override
  String toString() => message;
}

class UnauthorizedException extends ApiException {
  UnauthorizedException({super.statusCode = 401}) : super('unauthorized');
}

enum YouTubeBlockReason { userLimit, quotaExhausted, unavailable }

class YouTubeSearchBlockedException extends ApiException {
  final YouTubeBlockReason reason;
  final DateTime? retryAt;
  final int? remaining;

  YouTubeSearchBlockedException(
    super.message, {
    required this.reason,
    this.retryAt,
    this.remaining,
    super.statusCode,
  });
}

class AuthUnavailableException extends ApiException {
  AuthUnavailableException({super.statusCode}) : super('auth_unavailable');
}

class YouTubeErrorException extends ApiException {
  YouTubeErrorException({super.statusCode}) : super('youtube_error');
}

class AlreadyInPlaylistException extends ApiException {
  AlreadyInPlaylistException(super.message) : super(statusCode: 409);
}

class NotFoundException extends ApiException {
  NotFoundException(super.message) : super(statusCode: 404);
}

const _blockReasons = {
  'user_limit': YouTubeBlockReason.userLimit,
  'quota_exhausted': YouTubeBlockReason.quotaExhausted,
  'search_unavailable': YouTubeBlockReason.unavailable,
};

ApiException apiExceptionFromDio(DioException e) {
  if (e.type != DioExceptionType.badResponse) {
    return ApiException("Couldn't reach the TuneTrend service.");
  }
  final status = e.response?.statusCode;
  final data = e.response?.data;
  final body = data is Map<String, dynamic> ? data : const <String, dynamic>{};
  final error = body['error'] as String?;
  final message = error ?? 'TuneTrend API responded with status $status';

  final reason = _blockReasons[error];
  if (reason != null) {
    final retryAt = body['retryAt'];
    return YouTubeSearchBlockedException(
      message,
      reason: reason,
      retryAt: retryAt is String ? DateTime.tryParse(retryAt) : null,
      remaining: (body['remaining'] as num?)?.toInt(),
      statusCode: status,
    );
  }
  if (error == 'auth_unavailable') {
    return AuthUnavailableException(statusCode: status);
  }
  if (error == 'youtube_error') {
    return YouTubeErrorException(statusCode: status);
  }
  if (status == 401 || error == 'unauthorized') {
    return UnauthorizedException(statusCode: status);
  }
  if (status == 409) return AlreadyInPlaylistException(message);
  if (status == 404) return NotFoundException(message);
  return ApiException(message, statusCode: status);
}

class ApiClient {
  factory ApiClient() => _instance;

  ApiClient._()
    : dio = Dio(
        BaseOptions(
          baseUrl: _apiBaseUrl,
          connectTimeout: const Duration(seconds: 10),
          receiveTimeout: const Duration(seconds: 10),
        ),
      ),
      _traced = true {
    dio.interceptors.add(DioFirebasePerformanceInterceptor());
    // Bodies can contain contact-form PII.
    if (kDebugMode) {
      dio.interceptors.add(
        LogInterceptor(
          requestHeader: false,
          requestBody: false,
          responseBody: false,
        ),
      );
    }
  }

  @visibleForTesting
  ApiClient.withDio(this.dio) : _traced = false;

  static final ApiClient _instance = ApiClient._();

  final Dio dio;
  final bool _traced;

  String? _authToken;
  VoidCallback? onUnauthorized;

  bool get hasAuthToken => _authToken != null;

  void setAuthToken(String? token) {
    _authToken = (token == null || token.isEmpty) ? null : token;
  }

  Future<T> _trace<T>(String name, Future<T> Function() action) {
    if (!_traced) return action();
    return PerformanceService().traceFunction(name, action);
  }

  Future<List<Song>> fetchSongs(
    String country,
    TrendTab tab, {
    String categoryId = '',
  }) {
    return _trace('fetch_songs', () async {
      try {
        final response = await dio.get(
          tab.endpointPath,
          queryParameters: {
            'country': country,
            if (categoryId.isNotEmpty && tab != TrendTab.searched)
              'category': categoryId,
          },
        );
        return _parseSongs(response.data as Map<String, dynamic>, country);
      } on DioException catch (e) {
        if (e.response?.statusCode == 400 && categoryId.isNotEmpty) {
          return fetchSongs(country, tab);
        }
        throw apiExceptionFromDio(e);
      }
    });
  }

  Future<List<Song>> fetchSearchedTrends(String country) =>
      fetchSongs(country, TrendTab.searched);

  Future<List<Category>> fetchCategories(String country) {
    return _trace('fetch_categories', () async {
      try {
        final response = await dio.get(
          '/categories',
          queryParameters: {'country': country},
        );
        final data = _unwrap(response) as List<dynamic>? ?? const [];
        return data
            .map((e) => Category.fromJson(e as Map<String, dynamic>))
            .toList();
      } on DioException catch (e) {
        throw apiExceptionFromDio(e);
      }
    });
  }

  Future<void> submitContact({
    String? name,
    required String message,
    String? contactEmail,
    String? contactPhone,
  }) {
    return _trace('submit_contact', () async {
      try {
        final response = await dio.post(
          '/contact',
          data: {
            'name': name ?? '',
            'message': message,
            'contactEmail': contactEmail ?? '',
            'contactPhone': contactPhone ?? '',
            'website': '',
          },
        );
        _unwrap(response);
      } on DioException catch (e) {
        throw apiExceptionFromDio(e);
      }
    });
  }

  Future<LocalSearchResult> search(
    String query,
    String country, {
    String category = kMusicCategoryId,
  }) {
    return _trace('search_local', () async {
      final data = await _request(
        'GET',
        '/search',
        query: {'q': query, 'country': country, 'category': category},
      );
      final map = data as Map<String, dynamic>? ?? const {};
      return LocalSearchResult(
        items: _songs(map['items'], country),
        suggestYouTube: map['suggestYouTube'] as bool? ?? false,
      );
    });
  }

  Future<YouTubeSearchResult> searchYouTube(
    String query,
    String country, {
    String category = kMusicCategoryId,
  }) {
    return _trace('search_youtube', () async {
      final data = await _request(
        'GET',
        '/search/youtube',
        query: {'q': query, 'country': country, 'category': category},
        auth: true,
      );
      final map = data as Map<String, dynamic>? ?? const {};
      return YouTubeSearchResult(
        items: _songs(map['items'], country),
        remaining: (map['remaining'] as num?)?.toInt() ?? 0,
        cached: map['cached'] as bool? ?? false,
      );
    });
  }

  Future<bool> recordHit(String videoId, String country, String query) async {
    final data = await _request(
      'POST',
      '/search/hit',
      body: {'videoId': videoId, 'country': country, 'q': query},
    );
    return (data as Map<String, dynamic>?)?['counted'] as bool? ?? false;
  }

  Future<AuthSession> loginWithGoogle(String idToken) async {
    final data = await _request(
      'POST',
      '/auth/google',
      body: {'idToken': idToken},
    );
    return AuthSession.fromJson(data as Map<String, dynamic>);
  }

  Future<AppUser> fetchMe() async {
    final data = await _request('GET', '/me', auth: true);
    return AppUser.fromJson(data as Map<String, dynamic>);
  }

  Future<void> deleteMe() => _request('DELETE', '/me', auth: true);

  Future<List<Playlist>> fetchPlaylists() async {
    final data = await _request('GET', '/me/playlists', auth: true);
    return ((data as List<dynamic>?) ?? const [])
        .map((e) => Playlist.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<Playlist> createPlaylist(String name) async {
    final data = await _request(
      'POST',
      '/me/playlists',
      body: {'name': name},
      auth: true,
    );
    return Playlist.fromJson(data as Map<String, dynamic>);
  }

  Future<Playlist> fetchPlaylist(String id) async {
    final data = await _request('GET', '/me/playlists/$id', auth: true);
    return Playlist.fromJson(data as Map<String, dynamic>);
  }

  Future<Playlist> renamePlaylist(String id, String name) async {
    final data = await _request(
      'PATCH',
      '/me/playlists/$id',
      body: {'name': name},
      auth: true,
    );
    return Playlist.fromJson(data as Map<String, dynamic>);
  }

  Future<void> deletePlaylist(String id) =>
      _request('DELETE', '/me/playlists/$id', auth: true);

  Future<void> addPlaylistItem(
    String playlistId, {
    required String videoId,
    required String title,
    required String channelTitle,
    required String thumbnailUrl,
  }) {
    return _request(
      'POST',
      '/me/playlists/$playlistId/items',
      body: {
        'videoId': videoId,
        'title': title,
        'channelTitle': channelTitle,
        'thumbnailUrl': thumbnailUrl,
      },
      auth: true,
    );
  }

  Future<void> removePlaylistItem(String playlistId, String videoId) =>
      _request(
        'DELETE',
        '/me/playlists/$playlistId/items/$videoId',
        auth: true,
      );

  Future<Playlist> reorderPlaylist(String id, List<String> videoIds) async {
    final data = await _request(
      'PUT',
      '/me/playlists/$id/order',
      body: {'videoIds': videoIds},
      auth: true,
    );
    return Playlist.fromJson(data as Map<String, dynamic>);
  }

  Future<Object?> _request(
    String method,
    String path, {
    Map<String, dynamic>? query,
    Object? body,
    bool auth = false,
  }) async {
    final token = _authToken;
    if (auth && token == null) throw UnauthorizedException(statusCode: null);
    try {
      final response = await dio.request<Object?>(
        path,
        data: body,
        queryParameters: query,
        options: Options(
          method: method,
          headers: auth ? {'Authorization': 'Bearer $token'} : null,
        ),
      );
      return _unwrap(response);
    } on DioException catch (e) {
      final error = apiExceptionFromDio(e);
      if (auth && error is UnauthorizedException && token == _authToken) {
        onUnauthorized?.call();
      }
      throw error;
    }
  }

  Object? _unwrap(Response<Object?> response) {
    final body = response.data;
    if (body is! Map<String, dynamic> || body['success'] != true) {
      final error = body is Map<String, dynamic> ? body['error'] : null;
      throw ApiException(
        error as String? ?? 'TuneTrend API returned an error',
        statusCode: response.statusCode,
      );
    }
    return body['data'];
  }

  List<Song> _songs(Object? raw, String country) {
    return ((raw as List<dynamic>?) ?? const [])
        .map((e) => Song.fromJson(e as Map<String, dynamic>, country: country))
        .toList();
  }

  List<Song> _parseSongs(Map<String, dynamic> body, String country) {
    if (body['success'] != true) {
      throw ApiException(
        body['error'] as String? ?? 'TuneTrend API returned an error',
      );
    }
    return _songs(body['data'], country);
  }
}
