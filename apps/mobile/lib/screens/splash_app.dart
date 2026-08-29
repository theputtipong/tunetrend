import 'dart:async';

import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/material.dart';
import 'package:lottie/lottie.dart';
import 'package:package_info_plus/package_info_plus.dart';

import '../constants/theme.dart';
import '../services/messaging_service.dart';
import '../services/remote_config_service.dart';
import '../utils/version.dart';
import 'force_update_app.dart';
import 'maintenance_app.dart';
import 'tune_trend_app.dart';

enum _BootPhase { loading, maintenance, forceUpdate, ready }

/// Shown immediately on launch — plays the brand Lottie animation while
/// Firebase Messaging/Remote Config finish loading, then swaps itself for
/// [MaintenanceApp], [ForceUpdateApp], or [TuneTrendApp]. Same
/// self-contained-MaterialApp pattern those two already use.
class SplashApp extends StatefulWidget {
  const SplashApp({super.key});

  @override
  State<SplashApp> createState() => _SplashAppState();
}

class _SplashAppState extends State<SplashApp> {
  // Keeps the animation from flashing by for near-instant boots.
  static const _minDuration = Duration(milliseconds: 1500);

  _BootPhase _phase = _BootPhase.loading;

  @override
  void initState() {
    super.initState();
    _bootstrap();
  }

  Future<void> _bootstrap() async {
    final started = DateTime.now();

    FirebaseMessaging.onBackgroundMessage(firebaseMessagingBackgroundHandler);
    unawaited(MessagingService.initialize());
    await RemoteConfigController.instance.load();

    final snapshot = RemoteConfigController.instance.value;
    final packageInfo = await PackageInfo.fromPlatform();
    final updateRequired =
        !snapshot.isMaintenance &&
        isVersionLower(packageInfo.version, snapshot.minAppVersion);

    final elapsed = DateTime.now().difference(started);
    if (elapsed < _minDuration) {
      await Future.delayed(_minDuration - elapsed);
    }

    if (!mounted) return;
    setState(() {
      _phase = snapshot.isMaintenance
          ? _BootPhase.maintenance
          : updateRequired
          ? _BootPhase.forceUpdate
          : _BootPhase.ready;
    });
  }

  @override
  Widget build(BuildContext context) {
    switch (_phase) {
      case _BootPhase.maintenance:
        return const MaintenanceApp();
      case _BootPhase.forceUpdate:
        return const ForceUpdateApp();
      case _BootPhase.ready:
        return const TuneTrendApp();
      case _BootPhase.loading:
        final brightness =
            WidgetsBinding.instance.platformDispatcher.platformBrightness;
        AppColors.applyBrightness(brightness);
        return MaterialApp(
          debugShowCheckedModeBanner: false,
          theme: buildAppTheme(brightness),
          home: Scaffold(
            backgroundColor: AppColors.background,
            body: Center(
              child: Lottie.asset(
                'assets/lottie/splash.json',
                width: 220,
                height: 220,
              ),
            ),
          ),
        );
    }
  }
}
