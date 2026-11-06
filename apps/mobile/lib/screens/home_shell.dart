import 'package:flutter/material.dart';

import '../constants/auth_config.dart';
import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import 'playlists_screen.dart';
import 'search_screen.dart';
import 'trends_screen.dart';

class HomeShell extends StatefulWidget {
  final String initialCountry;

  const HomeShell({super.key, required this.initialCountry});

  @override
  State<HomeShell> createState() => _HomeShellState();
}

class _HomeShellState extends State<HomeShell> {
  int _index = 0;
  late String _country = widget.initialCountry;

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    return Scaffold(
      body: IndexedStack(
        index: _index,
        children: [
          TrendsScreen(
            initialCountry: widget.initialCountry,
            onCountryChanged: (code) => setState(() => _country = code),
          ),
          SearchScreen(country: _country),
          if (kLoginEnabled)
            PlaylistsScreen(country: _country, active: _index == 2),
        ],
      ),
      bottomNavigationBar: NavigationBar(
        selectedIndex: _index,
        onDestinationSelected: (i) => setState(() => _index = i),
        backgroundColor: AppColors.surface,
        surfaceTintColor: Colors.transparent,
        indicatorColor: AppColors.accent.withValues(alpha: 0.18),
        height: 64,
        labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
        destinations: [
          NavigationDestination(
            icon: Icon(Icons.bar_chart_rounded, color: AppColors.textSecondary),
            selectedIcon: const Icon(
              Icons.bar_chart_rounded,
              color: AppColors.accent,
            ),
            label: t.navCharts,
          ),
          NavigationDestination(
            icon: Icon(Icons.search, color: AppColors.textSecondary),
            selectedIcon: const Icon(Icons.search, color: AppColors.accent),
            label: t.navSearch,
          ),
          if (kLoginEnabled)
            NavigationDestination(
              icon: Icon(
                Icons.queue_music_rounded,
                color: AppColors.textSecondary,
              ),
              selectedIcon: const Icon(
                Icons.queue_music_rounded,
                color: AppColors.accent,
              ),
              label: t.navPlaylists,
            ),
        ],
      ),
    );
  }
}
