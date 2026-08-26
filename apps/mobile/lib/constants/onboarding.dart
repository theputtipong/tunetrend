import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

class OnboardingController extends ValueNotifier<bool> {
  OnboardingController._() : super(false);

  static final instance = OnboardingController._();

  static const _prefKey = 'tunetrend_onboarding_seen';

  Future<void> load() async {
    final prefs = await SharedPreferences.getInstance();
    value = prefs.getBool(_prefKey) ?? false;
  }

  Future<void> markSeen() async {
    value = true;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool(_prefKey, true);
  }
}
