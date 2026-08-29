import 'package:flutter/material.dart';

import '../constants/theme.dart';

enum StateMessageVariant { error, empty }

class StateMessage extends StatelessWidget {
  final StateMessageVariant variant;
  final String title;
  final String description;
  final VoidCallback? onRetry;
  final String retryLabel;

  const StateMessage({
    super.key,
    required this.variant,
    required this.title,
    required this.description,
    this.onRetry,
    this.retryLabel = 'Retry',
  });

  @override
  Widget build(BuildContext context) {
    final isError = variant == StateMessageVariant.error;

    return Center(
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 64),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 88,
              height: 88,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: isError ? AppColors.errorBg : AppColors.surfaceRaised,
              ),
              child: Icon(
                isError
                    ? Icons.warning_amber_rounded
                    : Icons.inventory_2_outlined,
                size: 40,
                color: isError ? AppColors.errorText : AppColors.textTertiary,
              ),
            ),
            const SizedBox(height: 18),
            Text(
              title,
              textAlign: TextAlign.center,
              style: displayFont(fontSize: 21),
            ),
            const SizedBox(height: 8),
            Text(
              description,
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: 14,
                height: 1.5,
                color: AppColors.textSecondary,
              ),
            ),
            if (onRetry != null) ...[
              const SizedBox(height: 18),
              FilledButton.icon(
                onPressed: onRetry,
                style: FilledButton.styleFrom(
                  backgroundColor: AppColors.accent,
                  foregroundColor: AppColors.accentInk,
                  padding: const EdgeInsets.symmetric(
                    horizontal: 20,
                    vertical: 12,
                  ),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(999),
                  ),
                ),
                icon: const Icon(Icons.refresh, size: 18),
                label: Text(
                  retryLabel,
                  style: const TextStyle(fontWeight: FontWeight.w700),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
