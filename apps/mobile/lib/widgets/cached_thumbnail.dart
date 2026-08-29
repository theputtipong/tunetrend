import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:shimmer/shimmer.dart';

import '../constants/theme.dart';

/// Network image with disk/memory caching (so the same thumbnail isn't
/// re-downloaded on every rebuild or scroll) and a shimmering placeholder
/// while it loads.
class CachedThumbnail extends StatelessWidget {
  final String imageUrl;
  final double? width;
  final double? height;
  final BoxFit fit;

  const CachedThumbnail({
    super.key,
    required this.imageUrl,
    this.width,
    this.height,
    this.fit = BoxFit.cover,
  });

  @override
  Widget build(BuildContext context) {
    return CachedNetworkImage(
      imageUrl: imageUrl,
      width: width,
      height: height,
      fit: fit,
      placeholder: (context, url) => Shimmer.fromColors(
        baseColor: AppColors.surfaceRaised,
        highlightColor: AppColors.border,
        child: Container(
          width: width,
          height: height,
          color: AppColors.surfaceRaised,
        ),
      ),
      errorWidget: (context, url, error) => Container(
        width: width,
        height: height,
        color: AppColors.surfaceRaised,
      ),
    );
  }
}
