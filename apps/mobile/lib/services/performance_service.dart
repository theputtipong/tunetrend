import 'package:firebase_performance/firebase_performance.dart';

class PerformanceService {
  factory PerformanceService() => _instance;

  PerformanceService._();

  static final PerformanceService _instance = PerformanceService._();

  final FirebasePerformance _performance = FirebasePerformance.instance;

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
