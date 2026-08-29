import 'package:flutter/material.dart';

import '../constants/theme.dart';

/// A gently bobbing chevron over a trailing fade, hinting that the row
/// underneath it scrolls horizontally. Purely decorative — taps pass through.
class ScrollHintChevron extends StatefulWidget {
  const ScrollHintChevron({super.key});

  @override
  State<ScrollHintChevron> createState() => _ScrollHintChevronState();
}

class _ScrollHintChevronState extends State<ScrollHintChevron>
    with SingleTickerProviderStateMixin {
  late final AnimationController _controller;
  late final Animation<double> _offset;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 800),
    )..repeat(reverse: true);
    _offset = Tween<double>(
      begin: 0,
      end: 4,
    ).animate(CurvedAnimation(parent: _controller, curve: Curves.easeInOut));
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return IgnorePointer(
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 24,
            height: 36,
            decoration: BoxDecoration(
              gradient: LinearGradient(
                colors: [
                  AppColors.background.withValues(alpha: 0),
                  AppColors.background,
                ],
              ),
            ),
          ),
          Container(
            height: 36,
            color: AppColors.background,
            alignment: Alignment.center,
            padding: const EdgeInsets.only(right: 2),
            child: AnimatedBuilder(
              animation: _offset,
              builder: (context, child) => Transform.translate(
                offset: Offset(_offset.value, 0),
                child: child,
              ),
              child: Icon(
                Icons.chevron_right,
                size: 18,
                color: AppColors.textTertiary,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
