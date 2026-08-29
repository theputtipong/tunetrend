import 'package:firebase_crashlytics/firebase_crashlytics.dart';
import 'package:firebase_remote_config/firebase_remote_config.dart';
import 'package:flutter/foundation.dart';

class RemoteConfigSnapshot {
  final bool isMaintenance;
  final String minAppVersion;

  const RemoteConfigSnapshot({
    required this.isMaintenance,
    required this.minAppVersion,
  });

  static const initial = RemoteConfigSnapshot(
    isMaintenance: false,
    minAppVersion: '0.0.0',
  );
}

class RemoteConfigController extends ValueNotifier<RemoteConfigSnapshot> {
  RemoteConfigController._() : super(RemoteConfigSnapshot.initial);

  static final instance = RemoteConfigController._();

  Future<void> load() async {
    final remoteConfig = FirebaseRemoteConfig.instance;

    await remoteConfig.setConfigSettings(
      RemoteConfigSettings(
        fetchTimeout: const Duration(seconds: 10),
        minimumFetchInterval: kDebugMode
            ? Duration.zero
            : const Duration(minutes: 15),
      ),
    );
    await remoteConfig.setDefaults(const {
      'is_maintenance': false,
      'min_app_version': '0.0.0',
    });

    try {
      final activated = await remoteConfig.fetchAndActivate();
      debugPrint('🔧 [RemoteConfig] fetchAndActivate: activated=$activated');
    } catch (e, stackTrace) {
      // Offline or throttled — fall back to whatever's already cached (or
      // the defaults set above) instead of blocking app launch on this.
      debugPrint('⚠️ [RemoteConfig] fetchAndActivate ล้มเหลว: $e');
      FirebaseCrashlytics.instance.recordError(
        e,
        stackTrace,
        reason: 'RemoteConfig fetchAndActivate failed',
        fatal: false,
      );
    }

    value = RemoteConfigSnapshot(
      isMaintenance: remoteConfig.getBool('is_maintenance'),
      minAppVersion: remoteConfig.getString('min_app_version'),
    );
  }
}
