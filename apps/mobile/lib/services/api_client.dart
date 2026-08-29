import 'package:dio/dio.dart';
import 'package:firebase_performance_dio/firebase_performance_dio.dart';
import 'package:flutter/foundation.dart' show kDebugMode;

import '../constants/tabs.dart';
import '../models/category.dart';
import '../models/song.dart';
import 'performance_service.dart';

const String _apiBaseUrl = String.fromEnvironment(
  'API_BASE_URL',
  defaultValue: 'http://localhost:8080',
);

class ApiException implements Exception {
  final String message;
  final int? statusCode;
  ApiException(this.message, {this.statusCode});

  @override
  String toString() => message;
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
      ) {
    dio.interceptors.add(DioFirebasePerformanceInterceptor());
    // Request/response bodies can carry user PII (contact form details) —
    // only log them in debug builds, never in release.
    if (kDebugMode) {
      dio.interceptors.add(
        LogInterceptor(requestBody: true, responseBody: true),
      );
    }
  }

  static final ApiClient _instance = ApiClient._();

  /// The shared Dio client, wired with Firebase Performance and (debug-only)
  /// logging interceptors. Prefer the typed methods below over calling this
  /// directly.
  final Dio dio;

  Future<List<Song>> fetchSongs(
    String country,
    TrendTab tab, {
    String categoryId = '',
  }) {
    return PerformanceService().traceFunction('fetch_songs', () async {
      try {
        final response = await dio.get(
          tab.endpointPath,
          queryParameters: {
            'country': country,
            if (categoryId.isNotEmpty) 'category': categoryId,
          },
        );
        return _parseSongs(response.data as Map<String, dynamic>);
      } on DioException catch (e) {
        if (e.response?.statusCode == 400 && categoryId.isNotEmpty) {
          return fetchSongs(country, tab);
        }
        throw _toApiException(e);
      }
    });
  }

  Future<List<Category>> fetchCategories(String country) {
    return PerformanceService().traceFunction('fetch_categories', () async {
      try {
        final response = await dio.get(
          '/categories',
          queryParameters: {'country': country},
        );
        final body = response.data as Map<String, dynamic>;
        if (body['success'] != true) {
          throw ApiException(
            body['error'] as String? ?? 'TuneTrend API returned an error',
          );
        }
        final data = (body['data'] as List<dynamic>?) ?? const [];
        return data
            .map((e) => Category.fromJson(e as Map<String, dynamic>))
            .toList();
      } on DioException catch (e) {
        throw _toApiException(e);
      }
    });
  }

  Future<void> submitContact({
    String? name,
    required String message,
    String? contactEmail,
    String? contactPhone,
  }) {
    return PerformanceService().traceFunction('submit_contact', () async {
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
        final body = response.data as Map<String, dynamic>;
        if (body['success'] != true) {
          throw ApiException(
            body['error'] as String? ?? 'TuneTrend API returned an error',
            statusCode: response.statusCode,
          );
        }
      } on DioException catch (e) {
        throw _toApiException(e);
      }
    });
  }

  List<Song> _parseSongs(Map<String, dynamic> body) {
    if (body['success'] != true) {
      throw ApiException(
        body['error'] as String? ?? 'TuneTrend API returned an error',
      );
    }
    final data = (body['data'] as List<dynamic>?) ?? const [];
    return data.map((e) => Song.fromJson(e as Map<String, dynamic>)).toList();
  }

  ApiException _toApiException(DioException e) {
    if (e.type == DioExceptionType.badResponse) {
      final data = e.response?.data;
      final serverError = data is Map<String, dynamic>
          ? data['error'] as String?
          : null;
      return ApiException(
        serverError ??
            'TuneTrend API responded with status ${e.response?.statusCode}',
        statusCode: e.response?.statusCode,
      );
    }
    return ApiException("Couldn't reach the TuneTrend service.");
  }
}
