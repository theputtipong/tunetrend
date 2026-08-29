import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import '../i18n/lang.dart';
import '../widgets/logo_mark.dart';

const _playStoreUrl =
    'https://play.google.com/store/apps/details?id=com.tunetrend.tunetrend_mobile';

/// Self-contained replacement for [TuneTrendApp], shown instead of the real
/// app when the installed version is below Remote Config's `min_app_version`.
/// Runs before [LangController]/[ThemeController] have loaded their saved
/// preferences, so it reads the device locale and platform brightness
/// directly — same pattern as [MaintenanceApp].
class ForceUpdateApp extends StatelessWidget {
  const ForceUpdateApp({super.key});

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
                    t.updateRequiredTitle,
                    textAlign: TextAlign.center,
                    style: displayFont(fontSize: 20),
                  ),
                  const SizedBox(height: 12),
                  Text(
                    t.updateRequiredDescription,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      color: AppColors.textSecondary,
                      fontSize: 14,
                    ),
                  ),
                  const SizedBox(height: 20),
                  FilledButton(
                    onPressed: () {
                      launchUrl(
                        Uri.parse(_playStoreUrl),
                        mode: LaunchMode.externalApplication,
                      );
                    },
                    style: FilledButton.styleFrom(
                      backgroundColor: AppColors.accent,
                      foregroundColor: AppColors.accentInk,
                    ),
                    child: Text(t.updateNowButton),
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
