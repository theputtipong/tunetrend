import 'dart:math' as math;

import 'package:flutter/material.dart';

import '../constants/theme.dart';

const _barCount = 10;

/// A small looping equalizer-bar animation, shown over the song row that
/// matches the video currently playing.
class NowPlayingIndicator extends StatefulWidget {
  const NowPlayingIndicator({super.key});

  @override
  State<NowPlayingIndicator> createState() => _NowPlayingIndicatorState();
}

class _NowPlayingIndicatorState extends State<NowPlayingIndicator>
    with SingleTickerProviderStateMixin {
  late final AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 900),
    )..repeat();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _controller,
      builder: (context, _) {
        return Row(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.end,
          children: List.generate(_barCount, (i) {
            final phase = (_controller.value + i / _barCount) % 1.0;
            final wave = (math.sin(phase * 2 * math.pi) + 1) / 2;
            return Container(
              width: 2.2,
              height: 4 + 10 * wave,
              margin: const EdgeInsets.symmetric(horizontal: 0.8),
              decoration: BoxDecoration(
                color: AppColors.accent,
                borderRadius: BorderRadius.circular(2),
              ),
            );
          }),
        );
      },
    );
  }
}
