import 'package:flutter/material.dart';

import '../constants/auth_config.dart';
import '../constants/tabs.dart';
import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import '../i18n/lang.dart';
import '../models/song.dart';
import '../screens/video_player_screen.dart';
import '../services/analytics_service.dart';
import '../utils/format.dart';
import 'add_to_playlist_sheet.dart';
import 'cached_thumbnail.dart';
import 'now_playing_indicator.dart';
import 'video_type_badge.dart';

class SongTile extends StatelessWidget {
  final Song song;
  final int? rank;
  final TrendTab tab;
  final bool isPlaying;
  final String? badge;
  final VoidCallback? onPlay;
  final List<Song>? queue;
  final Widget? trailing;
  final bool showAddToPlaylist;
  final String analyticsSource;
  final bool dimmed;

  const SongTile({
    super.key,
    required this.song,
    required this.rank,
    required this.tab,
    this.isPlaying = false,
    this.badge,
    this.onPlay,
    this.queue,
    this.trailing,
    this.showAddToPlaylist = true,
    this.analyticsSource = 'chart',
    this.dimmed = false,
  });

  @override
  Widget build(BuildContext context) {
    final rank = this.rank;
    final isTop3 = rank != null && rank <= 3;
    final lang = LangController.instance.resolve();
    final t = currentStrings;
    final views = song.hits != null
        ? t.searchHits(compactNumber(song.hits!))
        : song.viewCount.isEmpty
        ? ''
        : formatViewCount(song.viewCount, lang);
    final time = formatRelativeTime(song.publishedAt, lang);
    final trailing =
        this.trailing ??
        (showAddToPlaylist && kLoginEnabled && !isPlaying
            ? IconButton(
                tooltip: t.addToPlaylistTooltip,
                onPressed: () => showAddToPlaylistSheet(
                  context,
                  song,
                  source: analyticsSource,
                ),
                icon: Icon(Icons.add, color: AppColors.textSecondary, size: 22),
              )
            : null);

    return InkWell(
      borderRadius: BorderRadius.circular(14),
      onTap: isPlaying
          ? null
          : () {
              onPlay?.call();
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
                    queue: queue,
                  ),
                ),
              );
            },
      child: Stack(
        children: [
          _SongTileContent(
            dimmed: dimmed,
            song: song,
            rank: rank,
            isTop3: isTop3,
            views: views,
            time: time,
            badge: badge,
            trailing: trailing,
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
  final int? rank;
  final bool isTop3;
  final String views;
  final String time;
  final String? badge;
  final Widget? trailing;
  final bool dimmed;

  const _SongTileContent({
    this.dimmed = false,
    required this.song,
    required this.rank,
    required this.isTop3,
    required this.views,
    required this.time,
    this.badge,
    this.trailing,
  });

  @override
  Widget build(BuildContext context) {
    final rank = this.rank;
    final meta = [views, time].where((s) => s.isNotEmpty).join(' · ');
    return Padding(
      padding: EdgeInsets.fromLTRB(12, 10, trailing == null ? 12 : 0, 10),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          if (rank != null) ...[
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
          ],
          Opacity(
            opacity: dimmed ? 0.4 : 1,
            child: ClipRRect(
              borderRadius: BorderRadius.circular(8),
              child: SizedBox(
                width: 84,
                height: 47,
                child: song.thumbnailUrl.isEmpty
                    ? Container(color: AppColors.surfaceRaised)
                    : CachedThumbnail(imageUrl: song.thumbnailUrl),
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Opacity(
              opacity: dimmed ? 0.55 : 1,
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
                      if (badge != null) _MetaChip(label: badge!),
                      if (song.isNew)
                        _MetaChip(label: currentStrings.newBadge, accent: true)
                      else if (song.rankChange != null && song.rankChange != 0)
                        _RankChange(change: song.rankChange!),
                      if (meta.isNotEmpty)
                        Text(
                          meta,
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
          ),
          ?trailing,
        ],
      ),
    );
  }
}

class _MetaChip extends StatelessWidget {
  final String label;
  final bool accent;

  const _MetaChip({required this.label, this.accent = false});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
      decoration: BoxDecoration(
        color: accent
            ? AppColors.accent.withValues(alpha: 0.15)
            : AppColors.surfaceRaised,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        label,
        style: TextStyle(
          fontSize: 11,
          fontWeight: FontWeight.w600,
          color: accent ? AppColors.accent : AppColors.textSecondary,
        ),
      ),
    );
  }
}

class _RankChange extends StatelessWidget {
  final int change;

  const _RankChange({required this.change});

  @override
  Widget build(BuildContext context) {
    final up = change > 0;
    return Text(
      '${up ? '▲' : '▼'} ${change.abs()}',
      style: TextStyle(
        fontSize: 11.5,
        fontWeight: FontWeight.w700,
        color: up ? const Color(0xFF3DBE7A) : AppColors.textTertiary,
      ),
    );
  }
}
