import 'package:flutter/material.dart';

import '../constants/tabs.dart';
import '../constants/theme.dart';
import '../i18n/lang.dart';
import '../models/song.dart';
import '../screens/video_player_screen.dart';
import '../services/analytics_service.dart';
import '../utils/format.dart';
import 'cached_thumbnail.dart';
import 'now_playing_indicator.dart';
import 'video_type_badge.dart';

class SongTile extends StatelessWidget {
  final Song song;
  final int rank;
  final TrendTab tab;
  final bool isPlaying;

  const SongTile({
    super.key,
    required this.song,
    required this.rank,
    required this.tab,
    this.isPlaying = false,
  });

  @override
  Widget build(BuildContext context) {
    final isTop3 = rank <= 3;
    final lang = LangController.instance.resolve();
    final views = formatViewCount(song.viewCount, lang);
    final time = formatRelativeTime(song.publishedAt, lang);

    return InkWell(
      borderRadius: BorderRadius.circular(14),
      onTap: isPlaying
          ? null
          : () {
              AnalyticsService().logVideoPlayed(
                videoId: song.id,
                title: song.title,
                countryCode: song.countryCode,
                categoryId: song.categoryId,
                tab: tab.key,
              );
              Navigator.of(context).push(
                MaterialPageRoute(
                  builder: (_) => VideoPlayerScreen(
                    videoId: song.id,
                    title: song.title,
                    country: song.countryCode,
                    categoryId: song.categoryId,
                    tab: tab,
                  ),
                ),
              );
            },
      child: Stack(
        children: [
          _SongTileContent(
            song: song,
            rank: rank,
            isTop3: isTop3,
            views: views,
            time: time,
          ),
          if (isPlaying)
            Positioned.fill(
              child: Container(
                color:
                    (ThemeController.instance.resolve(context) ==
                                Brightness.dark
                            ? Colors.white
                            : Colors.black)
                        .withValues(alpha: 0.72),
                child: const Center(child: NowPlayingIndicator()),
              ),
            ),
        ],
      ),
    );
  }
}

class _SongTileContent extends StatelessWidget {
  final Song song;
  final int rank;
  final bool isTop3;
  final String views;
  final String time;

  const _SongTileContent({
    required this.song,
    required this.rank,
    required this.isTop3,
    required this.views,
    required this.time,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          SizedBox(
            width: 26,
            child: Text(
              rank.toString().padLeft(2, '0'),
              textAlign: TextAlign.center,
              style: TextStyle(
                fontFamily: displayFont().fontFamily,
                fontWeight: FontWeight.w700,
                fontSize: isTop3 ? 20 : 16,
                color: isTop3 ? AppColors.accent : AppColors.textTertiary,
              ),
            ),
          ),
          const SizedBox(width: 12),
          ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: SizedBox(
              width: 84,
              height: 47,
              child: song.thumbnailUrl.isEmpty
                  ? Container(color: AppColors.surfaceRaised)
                  : CachedThumbnail(imageUrl: song.thumbnailUrl),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  song.title,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: const TextStyle(
                    fontWeight: FontWeight.w600,
                    fontSize: 14.5,
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  song.channelTitle,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: TextStyle(
                    fontSize: 12.5,
                    color: AppColors.textSecondary,
                  ),
                ),
                const SizedBox(height: 6),
                Wrap(
                  spacing: 8,
                  runSpacing: 4,
                  crossAxisAlignment: WrapCrossAlignment.center,
                  children: [
                    VideoTypeBadge(videoType: song.videoType),
                    Text(
                      '$views · $time',
                      style: TextStyle(
                        fontSize: 11.5,
                        color: AppColors.textTertiary,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
