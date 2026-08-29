import 'package:flutter/material.dart';
import 'package:shimmer/shimmer.dart';

import '../constants/theme.dart';

/// Placeholder list shown instead of a spinner while songs are loading, so
/// the trends list layout is visible immediately even on a slow connection.
class SongListSkeleton extends StatelessWidget {
  final int itemCount;

  const SongListSkeleton({super.key, this.itemCount = 8});

  @override
  Widget build(BuildContext context) {
    return Shimmer.fromColors(
      baseColor: AppColors.surfaceRaised,
      highlightColor: AppColors.border,
      child: ListView.builder(
        padding: const EdgeInsets.symmetric(vertical: 8),
        itemCount: itemCount,
        itemBuilder: (context, index) => const _SongTileSkeleton(),
      ),
    );
  }
}

class _SongTileSkeleton extends StatelessWidget {
  const _SongTileSkeleton();

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          const SizedBox(width: 26),
          const SizedBox(width: 12),
          ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: Container(
              width: 84,
              height: 47,
              color: AppColors.surfaceRaised,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                _bar(width: double.infinity, height: 14.5),
                const SizedBox(height: 6),
                _bar(width: 120, height: 12.5),
                const SizedBox(height: 8),
                _bar(width: 90, height: 11.5),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _bar({required double width, required double height}) {
    return Container(
      width: width,
      height: height,
      decoration: BoxDecoration(
        color: AppColors.surfaceRaised,
        borderRadius: BorderRadius.circular(4),
      ),
    );
  }
}
