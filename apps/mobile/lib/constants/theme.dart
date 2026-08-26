import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:shared_preferences/shared_preferences.dart';

class AppColors {
  AppColors._();

  static Color background = _dark.background;
  static Color surface = _dark.surface;
  static Color surfaceRaised = _dark.surfaceRaised;
  static Color border = _dark.border;
  static Color textPrimary = _dark.textPrimary;
  static Color textSecondary = _dark.textSecondary;
  static Color textTertiary = _dark.textTertiary;
  static Color errorBg = _dark.errorBg;
  static Color errorText = _dark.errorText;

  static const accent = Color(0xFFFF7A47);
  static const accentInk = Color(0xFFFFFFFF);

  static void applyBrightness(Brightness brightness) {
    final set = brightness == Brightness.light ? _light : _dark;
    background = set.background;
    surface = set.surface;
    surfaceRaised = set.surfaceRaised;
    border = set.border;
    textPrimary = set.textPrimary;
    textSecondary = set.textSecondary;
    textTertiary = set.textTertiary;
    errorBg = set.errorBg;
    errorText = set.errorText;
  }
}

class _ThemeColorSet {
  const _ThemeColorSet({
    required this.background,
    required this.surface,
    required this.surfaceRaised,
    required this.border,
    required this.textPrimary,
    required this.textSecondary,
    required this.textTertiary,
    required this.errorBg,
    required this.errorText,
  });

  final Color background;
  final Color surface;
  final Color surfaceRaised;
  final Color border;
  final Color textPrimary;
  final Color textSecondary;
  final Color textTertiary;
  final Color errorBg;
  final Color errorText;
}

const _dark = _ThemeColorSet(
  background: Color(0xFF0B0D12),
  surface: Color(0xFF16181F),
  surfaceRaised: Color(0xFF1D2029),
  border: Color(0x663A3F4D),
  textPrimary: Color(0xFFF5F6F8),
  textSecondary: Color(0xFFA7ABB8),
  textTertiary: Color(0xFF7B7F8C),
  errorBg: Color(0xFF3D2416),
  errorText: Color(0xFFE8A874),
);

const _light = _ThemeColorSet(
  background: Color(0xFFFAFAFA),
  surface: Color(0xFFFFFFFF),
  surfaceRaised: Color(0xFFEBECEF),
  border: Color(0xB3D8DAE0),
  textPrimary: Color(0xFF1F222B),
  textSecondary: Color(0xFF565A66),
  textTertiary: Color(0xFF8A8D99),
  errorBg: Color(0xFFFBE4D9),
  errorText: Color(0xFF8A3A17),
);

class ThemeController extends ValueNotifier<Brightness?> {
  ThemeController._() : super(null);

  static final instance = ThemeController._();

  static const _prefKey = 'tunetrend_theme_override';

  Future<void> load() async {
    final prefs = await SharedPreferences.getInstance();
    final stored = prefs.getString(_prefKey);
    if (stored == 'light') {
      value = Brightness.light;
    } else if (stored == 'dark') {
      value = Brightness.dark;
    }
  }

  Brightness resolve(BuildContext context) {
    return value ?? MediaQuery.platformBrightnessOf(context);
  }

  Future<void> toggle(BuildContext context) async {
    final next = resolve(context) == Brightness.light
        ? Brightness.dark
        : Brightness.light;
    value = next;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(
      _prefKey,
      next == Brightness.light ? 'light' : 'dark',
    );
  }
}

ThemeData buildAppTheme(Brightness brightness) {
  final base = brightness == Brightness.light
      ? ThemeData.light(useMaterial3: true)
      : ThemeData.dark(useMaterial3: true);
  final body = GoogleFonts.plusJakartaSansTextTheme(base.textTheme).apply(
    bodyColor: AppColors.textPrimary,
    displayColor: AppColors.textPrimary,
  );

  return base.copyWith(
    brightness: brightness,
    scaffoldBackgroundColor: AppColors.background,
    colorScheme: base.colorScheme.copyWith(
      brightness: brightness,
      surface: AppColors.background,
      primary: AppColors.accent,
      onPrimary: AppColors.accentInk,
    ),
    textTheme: body,
    splashFactory: NoSplash.splashFactory,
    popupMenuTheme: PopupMenuThemeData(
      color: AppColors.surfaceRaised,
      surfaceTintColor: Colors.transparent,
      elevation: 8,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(14),
        side: BorderSide(color: AppColors.border),
      ),
    ),
  );
}

TextStyle displayFont({
  double fontSize = 21,
  FontWeight fontWeight = FontWeight.w700,
  Color? color,
}) {
  return GoogleFonts.spaceGrotesk(
    fontSize: fontSize,
    fontWeight: fontWeight,
    color: color ?? AppColors.textPrimary,
    letterSpacing: -0.2,
  );
}
