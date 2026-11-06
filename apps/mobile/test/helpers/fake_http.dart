import 'dart:convert';
import 'dart:typed_data';

import 'package:dio/dio.dart';

class FakeResponse {
  final int status;
  final Object? body;
  const FakeResponse(this.status, this.body);
}

class RecordedRequest {
  final String method;
  final String path;
  final Map<String, dynamic> query;
  final Map<String, dynamic> headers;
  final Object? data;

  RecordedRequest(this.method, this.path, this.query, this.headers, this.data);
}

class FakeAdapter implements HttpClientAdapter {
  final Map<String, FakeResponse> routes = {};
  final List<RecordedRequest> requests = [];

  void on(String method, String path, int status, Object? body) {
    routes['$method $path'] = FakeResponse(status, body);
  }

  @override
  Future<ResponseBody> fetch(
    RequestOptions options,
    Stream<Uint8List>? requestStream,
    Future<void>? cancelFuture,
  ) async {
    requests.add(
      RecordedRequest(
        options.method,
        options.path,
        options.queryParameters,
        options.headers,
        options.data,
      ),
    );
    final route = routes['${options.method} ${options.path}'];
    if (route == null) {
      throw DioException.connectionError(
        requestOptions: options,
        reason: 'no route',
      );
    }
    return ResponseBody.fromString(
      jsonEncode(route.body),
      route.status,
      headers: {
        Headers.contentTypeHeader: [Headers.jsonContentType],
      },
    );
  }

  @override
  void close({bool force = false}) {}
}

Dio fakeDio(FakeAdapter adapter) {
  return Dio(BaseOptions(baseUrl: 'http://test'))..httpClientAdapter = adapter;
}
