import 'package:flutter/material.dart';

import 'theme.dart';

const Map<String, Color> _videoTypeDotColor = {
  'MV': AppColors.accent,
  'Lyric': Color(0xFF4FB8C9),
  'Audio Track': Color(0xFF93A8E0),
  'Cover': Color(0xFFDE8AC9),
  'Live Performance': Color(0xFF6FC98C),
  'General': Color(0xFF8B8E99),
};

const Color _fallbackDotColor = Color(0xFF8B8E99);

Color videoTypeDotColor(String videoType) {
  return _videoTypeDotColor[videoType] ?? _fallbackDotColor;
}
