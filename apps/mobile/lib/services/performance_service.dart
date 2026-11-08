import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_performance/firebase_performance.dart';

class PerformanceService {
  factory PerformanceService() => _instance;

  PerformanceService._();

  static final PerformanceService _instance = PerformanceService._();

  Future<T> traceFunction<T>(String name, Future<T> Function() action) async {
    // No Firebase app (e.g. widget tests, failed init) → run untraced.
    if (Firebase.apps.isEmpty) return action();
    final trace = FirebasePerformance.instance.newTrace(name);
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
