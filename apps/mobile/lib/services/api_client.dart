import 'dart:convert';
import 'package:http/http.dart' as http;
import '../constants/tabs.dart';
import '../models/song.dart';

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
  final http.Client _client;

  ApiClient({http.Client? client}) : _client = client ?? http.Client();

  Future<List<Song>> fetchSongs(String country, TrendTab tab) async {
    final uri = Uri.parse(
      '$_apiBaseUrl${tab.endpointPath}',
    ).replace(queryParameters: {'country': country});

    final http.Response response;
    try {
      response = await _client.get(uri).timeout(const Duration(seconds: 10));
    } catch (_) {
      throw ApiException("Couldn't reach the TuneTrend service.");
    }

    if (response.statusCode != 200) {
      throw ApiException('TuneTrend API responded with status ${response.statusCode}');
    }

    final body = jsonDecode(response.body) as Map<String, dynamic>;
    if (body['success'] != true) {
      throw ApiException(body['error'] as String? ?? 'TuneTrend API returned an error');
    }

    final data = (body['data'] as List<dynamic>?) ?? const [];
    return data.map((e) => Song.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<void> submitContact({
    String? name,
    required String message,
    String? contactEmail,
    String? contactPhone,
  }) async {
    final uri = Uri.parse('$_apiBaseUrl/contact');

    final http.Response response;
    try {
      response = await _client
          .post(
            uri,
            headers: {'Content-Type': 'application/json'},
            body: jsonEncode({
              'name': name ?? '',
              'message': message,
              'contactEmail': contactEmail ?? '',
              'contactPhone': contactPhone ?? '',
              'website': '',
            }),
          )
          .timeout(const Duration(seconds: 10));
    } catch (_) {
      throw ApiException("Couldn't reach the TuneTrend service.");
    }

    final body = jsonDecode(response.body) as Map<String, dynamic>;
    if (body['success'] != true) {
      throw ApiException(
        body['error'] as String? ?? 'TuneTrend API returned an error',
        statusCode: response.statusCode,
      );
    }
  }
}
