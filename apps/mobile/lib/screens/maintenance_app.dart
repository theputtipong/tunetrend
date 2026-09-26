import 'package:flutter/material.dart';

import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import '../i18n/lang.dart';
import '../widgets/logo_mark.dart';

class MaintenanceApp extends StatelessWidget {
  const MaintenanceApp({super.key});

  @override
  Widget build(BuildContext context) {
    final brightness =
        WidgetsBinding.instance.platformDispatcher.platformBrightness;
    AppColors.applyBrightness(brightness);

    final lang = detectLangFromLocale(
      WidgetsBinding.instance.platformDispatcher.locale,
    );
    final t = dictionaries[lang]!;

    return MaterialApp(
      title: 'TuneTrend',
      debugShowCheckedModeBanner: false,
      theme: buildAppTheme(brightness),
      home: Scaffold(
        backgroundColor: AppColors.background,
        body: SafeArea(
          child: Center(
            child: Padding(
              padding: const EdgeInsets.all(24),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  LogoMark(size: 48),
                  const SizedBox(height: 20),
                  Text(
                    t.maintenanceTitle,
                    textAlign: TextAlign.center,
                    style: displayFont(fontSize: 20),
                  ),
                  const SizedBox(height: 12),
                  Text(
                    t.maintenanceDescription,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      color: AppColors.textSecondary,
                      fontSize: 14,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
