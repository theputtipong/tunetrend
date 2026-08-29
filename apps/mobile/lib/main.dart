import 'dart:async';

import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_crashlytics/firebase_crashlytics.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:package_info_plus/package_info_plus.dart';
import 'package:showcaseview/showcaseview.dart';

import 'constants/countries.dart';
import 'constants/onboarding.dart';
import 'constants/theme.dart';
import 'firebase_options.dart';
import 'i18n/lang.dart';
import 'screens/force_update_app.dart';
import 'screens/maintenance_app.dart';
import 'screens/trends_screen.dart';
import 'services/analytics_service.dart';
import 'services/messaging_service.dart';
import 'services/remote_config_service.dart';
import 'utils/version.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);

  // Route uncaught framework and platform/async errors to Crashlytics.
  FlutterError.onError = FirebaseCrashlytics.instance.recordFlutterFatalError;
  PlatformDispatcher.instance.onError = (error, stack) {
    FirebaseCrashlytics.instance.recordError(error, stack, fatal: true);
    return true;
  };

  FirebaseMessaging.onBackgroundMessage(firebaseMessagingBackgroundHandler);
  unawaited(MessagingService.initialize());

  await RemoteConfigController.instance.load();

  final snapshot = RemoteConfigController.instance.value;
  final packageInfo = await PackageInfo.fromPlatform();
  final updateRequired =
      !snapshot.isMaintenance &&
      isVersionLower(packageInfo.version, snapshot.minAppVersion);

  runApp(
    snapshot.isMaintenance
        ? const MaintenanceApp()
        : updateRequired
        ? const ForceUpdateApp()
        : const TuneTrendApp(),
  );
}

class TuneTrendApp extends StatefulWidget {
  const TuneTrendApp({super.key});

  @override
  State<TuneTrendApp> createState() => _TuneTrendAppState();
}

class _TuneTrendAppState extends State<TuneTrendApp>
    with WidgetsBindingObserver {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
    ThemeController.instance.addListener(_onThemeChanged);
    ThemeController.instance.load();
    LangController.instance.addListener(_onLangChanged);
    LangController.instance.load();
    OnboardingController.instance.load();
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    ThemeController.instance.removeListener(_onThemeChanged);
    LangController.instance.removeListener(_onLangChanged);
    super.dispose();
  }

  void _onThemeChanged() => setState(() {});
  void _onLangChanged() => setState(() {});

  @override
  void didChangePlatformBrightness() {
    if (ThemeController.instance.value == null) {
      setState(() {});
    }
  }

  @override
  Widget build(BuildContext context) {
    final brightness =
        ThemeController.instance.value ??
        WidgetsBinding.instance.platformDispatcher.platformBrightness;
    AppColors.applyBrightness(brightness);

    final locale = WidgetsBinding.instance.platformDispatcher.locale;
    final initialCountry = detectCountryFromLocale(locale);

    return MaterialApp(
      title: 'TuneTrend',
      debugShowCheckedModeBanner: false,
      theme: buildAppTheme(brightness),
      navigatorObservers: [AnalyticsService().observer],
      builder: (context, child) => ShowCaseWidget(builder: (context) => child!),
      home: TrendsScreen(initialCountry: initialCountry),
    );
  }
}
