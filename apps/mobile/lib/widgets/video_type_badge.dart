import 'package:flutter/material.dart';

import '../constants/theme.dart';
import '../constants/video_type.dart';

class VideoTypeBadge extends StatelessWidget {
  final String videoType;

  const VideoTypeBadge({super.key, required this.videoType});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 11, vertical: 5),
      decoration: BoxDecoration(
        color: AppColors.surfaceRaised,
        borderRadius: BorderRadius.circular(999),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 6,
            height: 6,
            decoration: BoxDecoration(
              color: videoTypeDotColor(videoType),
              shape: BoxShape.circle,
            ),
          ),
          const SizedBox(width: 6),
          Text(
            videoType.toUpperCase(),
            style: TextStyle(
              fontSize: 11,
              fontWeight: FontWeight.w600,
              letterSpacing: 0.4,
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }
}
