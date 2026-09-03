import 'package:flutter/material.dart';

import '../constants/tabs.dart';
import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import '../i18n/strings.dart';

class TrendTabs extends StatelessWidget {
  final TrendTab active;
  final bool showMusicVideos;
  final ValueChanged<TrendTab> onSelect;

  const TrendTabs({
    super.key,
    required this.active,
    required this.onSelect,
    this.showMusicVideos = true,
  });

  String _labelFor(TrendTab tab, AppStrings t) {
    switch (tab) {
      case TrendTab.trending:
        return t.tabTrending;
      case TrendTab.newReleases:
        return showMusicVideos ? t.tabNew : t.tabNewGeneric;
      case TrendTab.musicVideos:
        return t.tabMv;
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    final visibleTabs = showMusicVideos
        ? TrendTab.values
        : TrendTab.values.where((tab) => tab != TrendTab.musicVideos);
    return Row(
      children: visibleTabs.map((tab) {
        final isActive = tab == active;
        return Padding(
          padding: const EdgeInsets.only(right: 22),
          child: InkWell(
            onTap: () => onSelect(tab),
            child: Container(
              padding: const EdgeInsets.symmetric(vertical: 12),
              decoration: BoxDecoration(
                border: Border(
                  bottom: BorderSide(
                    color: isActive ? AppColors.accent : Colors.transparent,
                    width: 2,
                  ),
                ),
              ),
              child: Text(
                _labelFor(tab, t),
                style: TextStyle(
                  fontWeight: FontWeight.w600,
                  fontSize: 13.5,
                  color: isActive
                      ? AppColors.textPrimary
                      : AppColors.textTertiary,
                ),
              ),
            ),
          ),
        );
      }).toList(),
    );
  }
}
