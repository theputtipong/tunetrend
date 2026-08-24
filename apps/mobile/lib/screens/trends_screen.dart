import 'package:flutter/material.dart';
import 'package:showcaseview/showcaseview.dart';
import '../constants/countries.dart';
import '../constants/onboarding.dart';
import '../constants/tabs.dart';
import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import '../i18n/lang.dart';
import '../models/song.dart';
import 'about_screen.dart';
import '../services/api_client.dart';
import '../widgets/country_selector.dart';
import '../widgets/logo_mark.dart';
import '../widgets/song_tile.dart';
import '../widgets/state_message.dart';
import '../widgets/trend_tabs.dart';

enum _MenuAction { replayTour, language, theme, about }

class TrendsScreen extends StatefulWidget {
  final String initialCountry;

  const TrendsScreen({super.key, required this.initialCountry});

  @override
  State<TrendsScreen> createState() => _TrendsScreenState();
}

class _TrendsScreenState extends State<TrendsScreen> {
  final _api = ApiClient();

  final _keyCountry = GlobalKey();
  final _keyTabs = GlobalKey();
  final _keyMenu = GlobalKey();

  late String _country = widget.initialCountry;
  TrendTab _tab = TrendTab.trending;
  late Future<List<Song>> _songsFuture = _load();

  Future<List<Song>> _load() {
    return _api.fetchSongs(_country, _tab);
  }

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!OnboardingController.instance.value && mounted) {
        ShowCaseWidget.of(context).startShowCase([_keyCountry, _keyTabs, _keyMenu]);
        OnboardingController.instance.markSeen();
      }
    });
  }

  void _reload() {
    setState(() {
      _songsFuture = _load();
    });
  }

  void _selectCountry(String code) {
    if (code == _country) return;
    setState(() {
      _country = code;
      _songsFuture = _load();
    });
  }

  void _selectTab(TrendTab tab) {
    if (tab == _tab) return;
    setState(() {
      _tab = tab;
      _songsFuture = _load();
    });
  }

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    return Scaffold(
      body: SafeArea(
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 16, 16, 12),
              child: Row(
                children: [
                  LogoMark(size: 28),
                  const SizedBox(width: 8),
                  Text('TuneTrend', style: displayFont(fontSize: 18)),
                  const Spacer(),
                  Showcase(
                    key: _keyMenu,
                    title: t.onboardingMenuTitle,
                    description: t.onboardingMenuDescription,
                    child: PopupMenuButton<_MenuAction>(
                      tooltip: t.menuTooltip,
                      icon: Icon(Icons.more_vert, color: AppColors.textSecondary, size: 20),
                      onSelected: (action) {
                        switch (action) {
                          case _MenuAction.replayTour:
                            ShowCaseWidget.of(context).startShowCase(
                              [_keyCountry, _keyTabs, _keyMenu],
                            );
                          case _MenuAction.language:
                            LangController.instance.toggle();
                          case _MenuAction.theme:
                            ThemeController.instance.toggle(context);
                          case _MenuAction.about:
                            Navigator.of(context).push(
                              MaterialPageRoute(builder: (_) => const AboutScreen()),
                            );
                        }
                      },
                      itemBuilder: (context) => [
                        PopupMenuItem(
                          value: _MenuAction.replayTour,
                          child: ListTile(
                            leading: Icon(
                              Icons.help_outline,
                              size: 20,
                              color: AppColors.textSecondary,
                            ),
                            title: Text(
                              t.replayTourTooltip,
                              style: TextStyle(color: AppColors.textPrimary, fontSize: 14),
                            ),
                            contentPadding: EdgeInsets.zero,
                          ),
                        ),
                        PopupMenuItem(
                          value: _MenuAction.language,
                          child: ListTile(
                            leading: Icon(
                              Icons.language,
                              size: 20,
                              color: AppColors.textSecondary,
                            ),
                            title: Text(
                              LangController.instance.resolve() == AppLang.en
                                  ? 'ภาษาไทย'
                                  : 'English',
                              style: TextStyle(color: AppColors.textPrimary, fontSize: 14),
                            ),
                            contentPadding: EdgeInsets.zero,
                          ),
                        ),
                        PopupMenuItem(
                          value: _MenuAction.theme,
                          child: ListTile(
                            leading: Icon(
                              ThemeController.instance.resolve(context) == Brightness.dark
                                  ? Icons.light_mode_outlined
                                  : Icons.dark_mode_outlined,
                              size: 20,
                              color: AppColors.textSecondary,
                            ),
                            title: Text(
                              t.themeToggleTooltip,
                              style: TextStyle(color: AppColors.textPrimary, fontSize: 14),
                            ),
                            contentPadding: EdgeInsets.zero,
                          ),
                        ),
                        PopupMenuItem(
                          value: _MenuAction.about,
                          child: ListTile(
                            leading: Icon(
                              Icons.info_outline,
                              size: 20,
                              color: AppColors.textSecondary,
                            ),
                            title: Text(
                              t.aboutTooltip,
                              style: TextStyle(color: AppColors.textPrimary, fontSize: 14),
                            ),
                            contentPadding: EdgeInsets.zero,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Showcase(
                key: _keyCountry,
                title: t.onboardingCountryTitle,
                description: t.onboardingCountryDescription,
                child: CountrySelector(activeCountry: _country, onSelect: _selectCountry),
              ),
            ),
            const SizedBox(height: 14),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Showcase(
                key: _keyTabs,
                title: t.onboardingTabsTitle,
                description: t.onboardingTabsDescription,
                child: TrendTabs(active: _tab, onSelect: _selectTab),
              ),
            ),
            Divider(height: 1, color: AppColors.border),
            Expanded(
              child: RefreshIndicator(
                color: AppColors.accent,
                backgroundColor: AppColors.surface,
                onRefresh: () async => _reload(),
                child: FutureBuilder<List<Song>>(
                  future: _songsFuture,
                  builder: (context, snapshot) {
                    if (snapshot.connectionState != ConnectionState.done) {
                      return const Center(
                        child: CircularProgressIndicator(color: AppColors.accent),
                      );
                    }

                    if (snapshot.hasError) {
                      return ListView(
                        children: [
                          StateMessage(
                            variant: StateMessageVariant.error,
                            title: t.errorTitle,
                            description: t.errorDescription,
                            retryLabel: t.retry,
                            onRetry: _reload,
                          ),
                        ],
                      );
                    }

                    final songs = snapshot.data ?? const [];
                    if (songs.isEmpty) {
                      return ListView(
                        children: [
                          StateMessage(
                            variant: StateMessageVariant.empty,
                            title: t.emptyTitle,
                            description: t.emptyDescription(
                              countryLabel(_country, LangController.instance.resolve()),
                            ),
                          ),
                        ],
                      );
                    }

                    return ListView.builder(
                      padding: const EdgeInsets.symmetric(vertical: 8),
                      itemCount: songs.length,
                      itemBuilder: (context, index) {
                        return SongTile(song: songs[index], rank: index + 1);
                      },
                    );
                  },
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
