import 'package:firebase_performance/firebase_performance.dart';

/// Centralized wrapper around Firebase Performance Monitoring custom traces.
class PerformanceService {
  factory PerformanceService() => _instance;

  PerformanceService._();

  static final PerformanceService _instance = PerformanceService._();

  final FirebasePerformance _performance = FirebasePerformance.instance;

  /// Runs [action] inside a Firebase [Trace] named [name], measuring its
  /// wall-clock duration. `trace.stop()` always runs in a `finally` block —
  /// the trace is never leaked regardless of whether [action] succeeds,
  /// throws, or times out — and [action]'s result or exception is passed
  /// through unchanged.
  Future<T> traceFunction<T>(String name, Future<T> Function() action) async {
    final trace = _performance.newTrace(name);
    await trace.start();
    try {
      final result = await action();
      trace.putAttribute('status', 'success');
      return result;
    } catch (_) {
      trace.putAttribute('status', 'error');
      rethrow;
    } finally {
      await trace.stop();
    }
  }
}
