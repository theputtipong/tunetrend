import 'dart:async';

import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../constants/auth_config.dart';
import '../constants/countries.dart';
import '../constants/tabs.dart';
import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import '../i18n/lang.dart';
import '../i18n/strings.dart';
import '../models/category.dart';
import '../models/search.dart';
import '../models/song.dart';
import '../services/analytics_service.dart';
import '../services/api_client.dart';
import '../services/auth_service.dart';
import '../widgets/category_filter.dart';
import '../widgets/login_prompt.dart';
import '../widgets/song_list_skeleton.dart';
import '../widgets/song_tile.dart';
import '../widgets/state_message.dart';

const _recentKey = 'tunetrend_recent_searches';
const _recentMax = 5;
const _minQueryLength = 2;

enum _YtStatus { idle, loading, done, error }

class SearchScreen extends StatefulWidget {
  final String country;

  const SearchScreen({super.key, required this.country});

  @override
  State<SearchScreen> createState() => _SearchScreenState();
}

class _SearchScreenState extends State<SearchScreen> {
  final _api = ApiClient();
  final _controller = TextEditingController();
  final _focus = FocusNode();

  List<Category> _categories = const [];
  String _category = '';
  List<String> _recent = const [];

  String _query = '';
  int _generation = 0;
  bool _localLoading = false;
  Object? _localError;
  LocalSearchResult? _local;

  _YtStatus _ytStatus = _YtStatus.idle;
  YouTubeSearchResult? _yt;
  ApiException? _ytError;
  int? _remaining;

  String get _apiCategory => _category.isEmpty ? kMusicCategoryId : _category;

  @override
  void initState() {
    super.initState();
    _loadCategories();
    _loadRecent();
    AuthService.instance.addListener(_onAuthChanged);
  }

  @override
  void didUpdateWidget(SearchScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.country != widget.country) {
      _category = '';
      _categories = const [];
      _loadCategories();
      if (_query.isNotEmpty) _runSearch(_query);
    }
  }

  @override
  void dispose() {
    AuthService.instance.removeListener(_onAuthChanged);
    _controller.dispose();
    _focus.dispose();
    super.dispose();
  }

  void _onAuthChanged() {
    if (!mounted) return;
    setState(() {
      if (!AuthService.instance.isSignedIn) {
        _ytStatus = _YtStatus.idle;
        _yt = null;
        _ytError = null;
        _remaining = null;
      }
    });
  }

  Future<void> _loadCategories() async {
    final country = widget.country;
    try {
      final categories = await _api.fetchCategories(country);
      if (mounted && country == widget.country) {
        setState(() => _categories = categories);
      }
    } catch (_) {
      return;
    }
  }

  Future<void> _loadRecent() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final recent = prefs.getStringList(_recentKey) ?? const [];
      if (mounted) setState(() => _recent = recent);
    } catch (_) {
      return;
    }
  }

  Future<void> _remember(String query) async {
    final next = [
      query,
      ..._recent.where((q) => q.toLowerCase() != query.toLowerCase()),
    ].take(_recentMax).toList();
    setState(() => _recent = next);
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setStringList(_recentKey, next);
    } catch (_) {
      return;
    }
  }

  void _submit(String raw) {
    final query = raw.trim();
    if (query.runes.length < _minQueryLength) {
      ScaffoldMessenger.of(context)
        ..hideCurrentSnackBar()
        ..showSnackBar(
          SnackBar(content: Text(currentStrings.searchQueryTooShort)),
        );
      return;
    }
    _focus.unfocus();
    _remember(query);
    _runSearch(query);
  }

  void _selectRecent(String query) {
    _controller.text = query;
    _submit(query);
  }

  void _clear() {
    _controller.clear();
    _generation++;
    setState(() {
      _query = '';
      _local = null;
      _localError = null;
      _localLoading = false;
      _ytStatus = _YtStatus.idle;
      _yt = null;
      _ytError = null;
    });
    _focus.requestFocus();
  }

  Future<void> _runSearch(String query) async {
    final generation = ++_generation;
    final country = widget.country;
    final category = _apiCategory;
    setState(() {
      _query = query;
      _localLoading = true;
      _localError = null;
      _local = null;
      _ytStatus = _YtStatus.idle;
      _yt = null;
      _ytError = null;
    });
    try {
      final result = await _api.search(query, country, category: category);
      if (!mounted || generation != _generation) return;
      setState(() {
        _local = result;
        _localLoading = false;
      });
      AnalyticsService().logSearchPerformed(
        query: query,
        countryCode: country,
        categoryId: category,
        resultCount: result.items.length,
      );
      if (result.suggestYouTube &&
          kLoginEnabled &&
          AuthService.instance.isSignedIn) {
        _searchYouTube();
      }
    } catch (e) {
      if (!mounted || generation != _generation) return;
      setState(() {
        _localError = e;
        _localLoading = false;
      });
    }
  }

  Future<void> _searchYouTube() async {
    final query = _query;
    if (query.isEmpty || !AuthService.instance.isSignedIn) return;
    final generation = _generation;
    final country = widget.country;
    final category = _apiCategory;
    setState(() {
      _ytStatus = _YtStatus.loading;
      _ytError = null;
    });
    void log(String outcome, {int count = 0, bool cached = false}) {
      AnalyticsService().logYouTubeSearch(
        query: query,
        countryCode: country,
        categoryId: category,
        outcome: outcome,
        resultCount: count,
        cached: cached,
      );
    }

    try {
      final result = await _api.searchYouTube(
        query,
        country,
        category: category,
      );
      log('ok', count: result.items.length, cached: result.cached);
      if (!mounted || generation != _generation) return;
      setState(() {
        _yt = result;
        _remaining = result.remaining;
        _ytStatus = _YtStatus.done;
      });
    } on UnauthorizedException {
      log('unauthorized');
      if (!mounted) return;
      setState(() => _ytStatus = _YtStatus.idle);
    } on YouTubeSearchBlockedException catch (e) {
      log(e.message);
      if (!mounted || generation != _generation) return;
      setState(() {
        _ytError = e;
        if (e.remaining != null) _remaining = e.remaining;
        _ytStatus = _YtStatus.error;
      });
    } on ApiException catch (e) {
      log(e is YouTubeErrorException ? 'youtube_error' : 'error');
      if (!mounted || generation != _generation) return;
      setState(() {
        _ytError = e;
        _ytStatus = _YtStatus.error;
      });
    }
  }

  void _selectCategory(String categoryId) {
    if (categoryId == _category) return;
    setState(() => _category = categoryId);
    AnalyticsService().logCategorySelected(categoryId);
    if (_query.isNotEmpty) _runSearch(_query);
  }

  void _recordHit(Song song) {
    unawaited(
      _api
          .recordHit(song.id, widget.country, _query)
          .then<void>((_) {}, onError: (_) {}),
    );
  }

  String? _badgeFor(Song song, AppStrings t) {
    final rank = song.chartRank;
    if (rank != null && rank > 0) {
      final tab = switch (song.chartTab) {
        'new' => t.tabNew,
        'trending' => t.tabTrending,
        _ => '',
      };
      return tab.isEmpty ? '#$rank' : '#$rank $tab';
    }
    if (song.source == 'searched') return t.searchedBefore;
    return null;
  }

  Widget _resultTile(Song song, AppStrings t, String source) {
    return SongTile(
      song: song,
      rank: null,
      tab: TrendTab.trending,
      badge: _badgeFor(song, t),
      onPlay: () => _recordHit(song),
      analyticsSource: source,
    );
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
              child: ValueListenableBuilder<TextEditingValue>(
                valueListenable: _controller,
                builder: (context, value, _) => TextField(
                  controller: _controller,
                  focusNode: _focus,
                  textInputAction: TextInputAction.search,
                  onSubmitted: _submit,
                  style: TextStyle(color: AppColors.textPrimary, fontSize: 15),
                  decoration: InputDecoration(
                    hintText: t.searchHint,
                    hintStyle: TextStyle(color: AppColors.textTertiary),
                    filled: true,
                    fillColor: AppColors.surface,
                    prefixIcon: Icon(
                      Icons.search,
                      color: AppColors.textSecondary,
                    ),
                    suffixIcon: value.text.isEmpty
                        ? null
                        : IconButton(
                            tooltip: t.clearSearchTooltip,
                            onPressed: _clear,
                            icon: Icon(
                              Icons.close,
                              color: AppColors.textSecondary,
                              size: 20,
                            ),
                          ),
                    contentPadding: const EdgeInsets.symmetric(vertical: 12),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(999),
                      borderSide: BorderSide(color: AppColors.border),
                    ),
                    enabledBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(999),
                      borderSide: BorderSide(color: AppColors.border),
                    ),
                    focusedBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(999),
                      borderSide: const BorderSide(color: AppColors.accent),
                    ),
                  ),
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: CategoryFilter(
                active: _category,
                categories: _categories,
                musicLabel: t.musicCategory,
                onSelect: _selectCategory,
              ),
            ),
            if (_categories.isNotEmpty) const SizedBox(height: 10),
            Divider(height: 1, color: AppColors.border),
            Expanded(child: _buildBody(t)),
          ],
        ),
      ),
    );
  }

  Widget _buildBody(AppStrings t) {
    if (_query.isEmpty) return _buildIdle(t);
    if (_localLoading) return const SongListSkeleton();
    if (_localError != null) {
      final error = _localError;
      return ListView(
        children: [
          StateMessage(
            variant: StateMessageVariant.error,
            title: t.errorTitle,
            description: error is ApiException && error.statusCode == 400
                ? error.message
                : t.errorDescription,
            retryLabel: t.retry,
            onRetry: () => _runSearch(_query),
          ),
        ],
      );
    }

    final localItems = _local?.items ?? const <Song>[];
    final localIds = localItems.map((s) => s.id).toSet();
    final ytItems = (_yt?.items ?? const <Song>[])
        .where((s) => !localIds.contains(s.id))
        .toList();
    final lang = LangController.instance.resolve();

    return ListView(
      padding: const EdgeInsets.only(bottom: 24),
      keyboardDismissBehavior: ScrollViewKeyboardDismissBehavior.onDrag,
      children: [
        _SectionHeader(
          title: t.localSectionTitle,
          count: localItems.length,
          trailing: countryLabel(widget.country, lang),
        ),
        if (localItems.isEmpty)
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 4, 16, 12),
            child: Text(
              t.localEmpty,
              style: TextStyle(color: AppColors.textTertiary, fontSize: 13),
            ),
          ),
        for (final song in localItems) _resultTile(song, t, 'search_local'),
        if (kLoginEnabled) ..._buildYouTubeSection(t, ytItems),
      ],
    );
  }

  Widget _buildIdle(AppStrings t) {
    return ListView(
      padding: const EdgeInsets.all(16),
      keyboardDismissBehavior: ScrollViewKeyboardDismissBehavior.onDrag,
      children: [
        if (_recent.isNotEmpty) ...[
          Wrap(
            spacing: 8,
            runSpacing: 8,
            crossAxisAlignment: WrapCrossAlignment.center,
            children: [
              Text(
                t.recentSearches,
                style: TextStyle(color: AppColors.textTertiary, fontSize: 13),
              ),
              for (final q in _recent)
                ActionChip(
                  label: Text(q),
                  onPressed: () => _selectRecent(q),
                  backgroundColor: AppColors.surface,
                  side: BorderSide(color: AppColors.border),
                  labelStyle: TextStyle(
                    color: AppColors.textSecondary,
                    fontSize: 13,
                  ),
                ),
            ],
          ),
        ],
        StateMessage(
          variant: StateMessageVariant.empty,
          title: t.searchPromptTitle,
          description: t.searchPromptDescription,
        ),
      ],
    );
  }

  List<Widget> _buildYouTubeSection(AppStrings t, List<Song> items) {
    final signedIn = AuthService.instance.isSignedIn;
    final remaining = _remaining;
    return [
      const SizedBox(height: 8),
      _SectionHeader(
        title: t.youtubeSectionTitle,
        count: _ytStatus == _YtStatus.done ? items.length : null,
        trailing: signedIn && remaining != null
            ? t.youtubeRemaining(remaining, kYouTubeDailySearchLimit)
            : null,
      ),
      if (!signedIn)
        LoginPrompt(
          title: t.youtubeLoginTitle,
          description: t.youtubeLoginDescription,
          buttonLabel: t.signInWithGoogle,
          showPrivacyNote: false,
          icon: Icons.travel_explore,
          onSignedIn: _searchYouTube,
        )
      else
        switch (_ytStatus) {
          _YtStatus.idle => _YouTubeCta(onSearch: _searchYouTube),
          _YtStatus.loading => Padding(
            padding: const EdgeInsets.all(20),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const SizedBox(
                  width: 18,
                  height: 18,
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    color: AppColors.accent,
                  ),
                ),
                const SizedBox(width: 10),
                Text(
                  t.youtubeLoading,
                  style: TextStyle(color: AppColors.textSecondary),
                ),
              ],
            ),
          ),
          _YtStatus.error => _YouTubeNotice(
            error: _ytError,
            onRetry: _searchYouTube,
          ),
          _YtStatus.done when items.isEmpty => Padding(
            padding: const EdgeInsets.fromLTRB(16, 4, 16, 12),
            child: Text(
              t.youtubeEmpty,
              style: TextStyle(color: AppColors.textTertiary, fontSize: 13),
            ),
          ),
          _YtStatus.done => Column(
            children: [
              for (final song in items) _resultTile(song, t, 'search_youtube'),
            ],
          ),
        },
    ];
  }
}

class _SectionHeader extends StatelessWidget {
  final String title;
  final int? count;
  final String? trailing;

  const _SectionHeader({required this.title, this.count, this.trailing});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 14, 16, 6),
      child: Row(
        children: [
          Text(title, style: displayFont(fontSize: 15)),
          if (count != null) ...[
            const SizedBox(width: 6),
            Text(
              '$count',
              style: TextStyle(
                color: AppColors.textTertiary,
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
          const Spacer(),
          if (trailing != null)
            Flexible(
              child: Text(
                trailing!,
                overflow: TextOverflow.ellipsis,
                style: TextStyle(color: AppColors.textTertiary, fontSize: 12),
              ),
            ),
        ],
      ),
    );
  }
}

class _YouTubeCta extends StatelessWidget {
  final VoidCallback onSearch;

  const _YouTubeCta({required this.onSearch});

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            t.youtubeMoreHint,
            style: TextStyle(
              color: AppColors.textSecondary,
              fontSize: 13,
              height: 1.5,
            ),
          ),
          const SizedBox(height: 10),
          OutlinedButton.icon(
            onPressed: onSearch,
            style: OutlinedButton.styleFrom(
              foregroundColor: AppColors.accent,
              side: const BorderSide(color: AppColors.accent),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(999),
              ),
            ),
            icon: const Icon(Icons.travel_explore, size: 18),
            label: Text(
              t.youtubeSearchButton,
              style: const TextStyle(fontWeight: FontWeight.w700),
            ),
          ),
        ],
      ),
    );
  }
}

class _YouTubeNotice extends StatelessWidget {
  final ApiException? error;
  final VoidCallback onRetry;

  const _YouTubeNotice({required this.error, required this.onRetry});

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    final error = this.error;
    final String title;
    final String body;
    var canRetry = false;
    if (error is YouTubeSearchBlockedException) {
      final retryAt = error.retryAt;
      final time = retryAt == null
          ? null
          : DateFormat.Hm().format(retryAt.toLocal());
      switch (error.reason) {
        case YouTubeBlockReason.userLimit:
          title = t.youtubeUserLimitTitle;
          body = t.youtubeUserLimit(kYouTubeDailySearchLimit);
        case YouTubeBlockReason.quotaExhausted:
          title = t.youtubeUnavailableTitle;
          body = t.youtubeQuotaExhausted(time);
        case YouTubeBlockReason.unavailable:
          title = t.youtubeUnavailableTitle;
          body = t.youtubeUnavailable(time);
      }
    } else {
      title = t.youtubeErrorTitle;
      body = t.youtubeErrorDescription;
      canRetry = true;
    }

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: AppColors.errorBg,
        borderRadius: BorderRadius.circular(14),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(Icons.info_outline, color: AppColors.errorText, size: 20),
          const SizedBox(width: 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: TextStyle(
                    color: AppColors.errorText,
                    fontWeight: FontWeight.w700,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  body,
                  style: TextStyle(
                    color: AppColors.errorText,
                    fontSize: 13,
                    height: 1.5,
                  ),
                ),
                if (canRetry)
                  TextButton(
                    onPressed: onRetry,
                    style: TextButton.styleFrom(padding: EdgeInsets.zero),
                    child: Text(
                      t.retry,
                      style: TextStyle(
                        color: AppColors.errorText,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                  ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
