import 'package:flutter/material.dart';
import '../constants/theme.dart';

class LogoMark extends StatelessWidget {
  final double size;

  // ignore: prefer_const_constructors_in_immutables
  LogoMark({super.key, this.size = 30});

  @override
  Widget build(BuildContext context) {
    final scale = size / 48;
    return SizedBox(
      width: size,
      height: size * 52 / 48,
      child: Stack(
        clipBehavior: Clip.none,
        children: [
          _bar(left: 6 * scale, height: 14 * scale, scale: scale, opacity: 0.5),
          _bar(left: 15 * scale, height: 20 * scale, scale: scale, opacity: 0.68),
          _bar(left: 24 * scale, height: 27 * scale, scale: scale, opacity: 0.84),
          _bar(left: 33 * scale, height: 34 * scale, scale: scale, opacity: 1),
          Positioned(
            left: 29 * scale,
            top: 1 * scale,
            child: Container(
              width: 14 * scale,
              height: 14 * scale,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: AppColors.background,
                border: Border.all(color: AppColors.accent, width: 1.6 * scale),
              ),
              child: Text(
                '1',
                style: TextStyle(
                  fontSize: 8.5 * scale,
                  fontWeight: FontWeight.w700,
                  height: 1,
                  color: AppColors.accent,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _bar({
    required double left,
    required double height,
    required double scale,
    required double opacity,
  }) {
    return Positioned(
      left: left,
      bottom: 8 * scale,
      child: Container(
        width: 6 * scale,
        height: height,
        decoration: BoxDecoration(
          color: AppColors.accent.withValues(alpha: opacity),
          borderRadius: BorderRadius.circular(3 * scale),
        ),
      ),
    );
  }
}
