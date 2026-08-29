import 'package:flutter/material.dart';
import 'package:showcaseview/showcaseview.dart';

import '../constants/countries.dart';
import '../constants/onboarding.dart';
import '../constants/theme.dart';
import '../i18n/lang.dart';
import '../services/analytics_service.dart';
import 'trends_screen.dart';

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
