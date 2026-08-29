import 'package:flutter/material.dart';

import '../constants/theme.dart';

class PillButton extends StatelessWidget {
  final String label;
  final bool active;
  final VoidCallback onTap;

  const PillButton({
    super.key,
    required this.label,
    required this.active,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      borderRadius: BorderRadius.circular(999),
      onTap: onTap,
      child: Container(
        alignment: Alignment.center,
        padding: const EdgeInsets.symmetric(horizontal: 14),
        decoration: BoxDecoration(
          color: active ? AppColors.accent : AppColors.surface,
          borderRadius: BorderRadius.circular(999),
          border: active ? null : Border.all(color: AppColors.border),
        ),
        child: Text(
          label,
          style: TextStyle(
            fontWeight: FontWeight.w600,
            fontSize: 13,
            color: active ? AppColors.accentInk : AppColors.textSecondary,
          ),
        ),
      ),
    );
  }
}
