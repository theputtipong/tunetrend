import 'package:flutter/material.dart';
import 'package:showcaseview/showcaseview.dart';

import '../constants/countries.dart';
import '../constants/onboarding.dart';
import '../constants/tabs.dart';
import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import '../i18n/lang.dart';
import '../models/category.dart';
import '../models/song.dart';
import 'about_screen.dart';
import '../services/analytics_service.dart';
import '../services/api_client.dart';
import '../widgets/category_filter.dart';
import '../widgets/country_selector.dart';
import '../widgets/logo_mark.dart';
import '../widgets/song_list_skeleton.dart';
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
  final _keyCategory = GlobalKey();
  final _keyTabs = GlobalKey();
  final _keyMenu = GlobalKey();

  List<GlobalKey> get _tourKeys => [
    _keyCountry,
    if (_categories.isNotEmpty) _keyCategory,
    _keyTabs,
    _keyMenu,
  ];

  late String _country = widget.initialCountry;
  TrendTab _tab = TrendTab.trending;
  String _category = '';
  List<Category> _categories = const [];
  late Future<List<Song>> _songsFuture = _load();

  Future<List<Song>> _load() {
    final categoryId = _tab == TrendTab.musicVideos ? '' : _category;
    return _api.fetchSongs(_country, _tab, categoryId: categoryId);
  }

  Future<void> _loadCategories() async {
    try {
      final categories = await _api.fetchCategories(_country);
      if (mounted) setState(() => _categories = categories);
    } catch (_) {
      // Category filter is a bonus affordance — fall back to no filter on failure.
    }
  }

  @override
  void initState() {
    super.initState();
    _loadCategories().then((_) {
      if (!mounted) return;
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (!OnboardingController.instance.value && mounted) {
          ShowCaseWidget.of(context).startShowCase(_tourKeys);
          OnboardingController.instance.markSeen();
        }
      });
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
      _category = '';
      _categories = const [];
      _songsFuture = _load();
    });
    _loadCategories();
    AnalyticsService().logCountryChanged(code);
  }

  void _selectTab(TrendTab tab) {
    if (tab == _tab) return;
    setState(() {
      _tab = tab;
      _songsFuture = _load();
    });
  }

  void _selectCategory(String categoryId) {
    if (categoryId == _category) return;
    setState(() {
      _category = categoryId;
      if (categoryId.isNotEmpty && _tab == TrendTab.musicVideos) {
        _tab = TrendTab.trending;
      }
      _songsFuture = _load();
    });
    AnalyticsService().logCategorySelected(categoryId);
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
                      icon: Icon(
                        Icons.more_vert,
                        color: AppColors.textSecondary,
                        size: 20,
                      ),
                      onSelected: (action) {
                        switch (action) {
                          case _MenuAction.replayTour:
                            ShowCaseWidget.of(context).startShowCase(_tourKeys);
                          case _MenuAction.language:
                            LangController.instance.toggle();
                          case _MenuAction.theme:
                            ThemeController.instance.toggle(context);
                          case _MenuAction.about:
                            Navigator.of(context).push(
                              MaterialPageRoute(
                                builder: (_) => const AboutScreen(),
                              ),
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
                              style: TextStyle(
                                color: AppColors.textPrimary,
                                fontSize: 14,
                              ),
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
                              style: TextStyle(
                                color: AppColors.textPrimary,
                                fontSize: 14,
                              ),
                            ),
                            contentPadding: EdgeInsets.zero,
                          ),
                        ),
                        PopupMenuItem(
                          value: _MenuAction.theme,
                          child: ListTile(
                            leading: Icon(
                              ThemeController.instance.resolve(context) ==
                                      Brightness.dark
                                  ? Icons.light_mode_outlined
                                  : Icons.dark_mode_outlined,
                              size: 20,
                              color: AppColors.textSecondary,
                            ),
                            title: Text(
                              t.themeToggleTooltip,
                              style: TextStyle(
                                color: AppColors.textPrimary,
                                fontSize: 14,
                              ),
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
                              style: TextStyle(
                                color: AppColors.textPrimary,
                                fontSize: 14,
                              ),
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
                child: CountrySelector(
                  activeCountry: _country,
                  onSelect: _selectCountry,
                ),
              ),
            ),
            const SizedBox(height: 14),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Showcase(
                key: _keyCategory,
                title: t.onboardingCategoryTitle,
                description: t.onboardingCategoryDescription,
                child: CategoryFilter(
                  active: _category,
                  categories: _categories,
                  musicLabel: t.musicCategory,
                  onSelect: _selectCategory,
                ),
              ),
            ),
            const SizedBox(height: 14),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Showcase(
                key: _keyTabs,
                title: t.onboardingTabsTitle,
                description: t.onboardingTabsDescription,
                child: TrendTabs(
                  active: _tab,
                  showMusicVideos: _category.isEmpty,
                  onSelect: _selectTab,
                ),
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
                      return const SongListSkeleton();
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
                              countryLabel(
                                _country,
                                LangController.instance.resolve(),
                              ),
                            ),
                          ),
                        ],
                      );
                    }

                    return ListView.builder(
                      padding: const EdgeInsets.symmetric(vertical: 8),
                      itemCount: songs.length,
                      itemBuilder: (context, index) {
                        return SongTile(
                          song: songs[index],
                          rank: index + 1,
                          tab: _tab,
                        );
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
