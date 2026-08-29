import 'dart:async';

import 'package:flutter/material.dart';
import 'package:share_plus/share_plus.dart';
import 'package:youtube_player_iframe/youtube_player_iframe.dart';

import '../constants/app_links.dart';
import '../constants/countries.dart';
import '../constants/tabs.dart';
import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import '../i18n/lang.dart';
import '../models/song.dart';
import '../services/api_client.dart';
import '../widgets/song_list_skeleton.dart';
import '../widgets/song_tile.dart';
import '../widgets/state_message.dart';
import '../widgets/trend_tabs.dart';

const _autoPlayDelay = Duration(seconds: 5);
const _autoPlayPromptCountdownSeconds = 5;

class VideoPlayerScreen extends StatefulWidget {
  final String videoId;
  final String title;
  final String country;
  final String categoryId;
  final TrendTab tab;

  const VideoPlayerScreen({
    super.key,
    required this.videoId,
    required this.title,
    required this.country,
    required this.categoryId,
    required this.tab,
  });

  @override
  State<VideoPlayerScreen> createState() => _VideoPlayerScreenState();
}

class _VideoPlayerScreenState extends State<VideoPlayerScreen> {
  final _api = ApiClient();
  late final YoutubePlayerController _controller;
  StreamSubscription<YoutubePlayerValue>? _playerSub;

  bool? _autoPlayChoice;
  bool _handledEnd = false;

  // Bumped whenever the related tab changes, so a pending _playNextAfterDelay
  // countdown started under the old tab's ordering can detect it's stale and
  // abort instead of jumping to a video from a list the user has since left.
  int _autoPlayGeneration = 0;
  int? _autoAdvanceSecondsLeft;

  // Mutable copies of the widget's initial video, updated in place as
  // continuous play advances — the controller and its WebView are reused
  // across the chain instead of tearing down and recreating a new one per
  // video, which avoids WKWebView teardown races (and the resulting jank).
  late String _videoId = widget.videoId;
  late String _title = widget.title;
  late String _categoryId = widget.categoryId;

  bool get _showMv => _categoryId.isEmpty || _categoryId == kMusicCategoryId;

  late TrendTab _relatedTab = widget.tab;
  final _scrollController = ScrollController();
  final _nowPlayingKey = GlobalKey();
  final _shareButtonKey = GlobalKey();
  List<Song> _relatedSongs = const [];
  late Future<List<Song>> _relatedFuture = _loadRelated();

  Future<List<Song>> _loadRelated() async {
    final songs = await _api.fetchSongs(widget.country, _relatedTab);
    _relatedSongs = songs;
    WidgetsBinding.instance.addPostFrameCallback((_) => _scrollToNowPlaying());
    return songs;
  }

  Future<void> _scrollToNowPlaying() async {
    if (!mounted || !_scrollController.hasClients) return;
    final index = _relatedSongs.indexWhere((song) => song.id == _videoId);
    if (index <= 0) return;

    const estimatedItemExtent = 84.0;
    final estimate = (index * estimatedItemExtent).clamp(
      0.0,
      _scrollController.position.maxScrollExtent,
    );
    _scrollController.jumpTo(estimate);

    // The jump above forces ListView to build items near the target offset;
    // once the now-playing row actually exists, ease to its exact position.
    await Future<void>.delayed(const Duration(milliseconds: 50));
    final targetContext = _nowPlayingKey.currentContext;
    if (targetContext == null || !targetContext.mounted) return;
    await Scrollable.ensureVisible(
      targetContext,
      alignment: 0.3,
      duration: const Duration(milliseconds: 300),
      curve: Curves.easeInOut,
    );
  }

  void _selectRelatedTab(TrendTab tab) {
    if (tab == _relatedTab) return;
    final wasAutoPlaying = _autoPlayChoice == true;
    setState(() {
      _relatedTab = tab;
      _relatedFuture = _loadRelated();
      _autoPlayGeneration++;
      _autoAdvanceSecondsLeft = null;
      if (wasAutoPlaying) _autoPlayChoice = null;
    });
    if (wasAutoPlaying) {
      ScaffoldMessenger.of(context)
        ..hideCurrentSnackBar()
        ..showSnackBar(
          SnackBar(content: Text(currentStrings.autoPlayStoppedMessage)),
        );
    }
  }

  Future<void> _onPlayerValueChanged(YoutubePlayerValue value) async {
    if (value.playerState != PlayerState.ended || _handledEnd) return;
    _handledEnd = true;
    if (!mounted) return;

    if (_autoPlayChoice == null) {
      final accepted = await _askAutoPlay();
      if (!mounted) return;
      setState(() => _autoPlayChoice = accepted);
    }

    if (_autoPlayChoice != true) return;
    await _playNextAfterDelay();
  }

  Future<bool> _askAutoPlay() async {
    final accepted = await showDialog<bool>(
      context: context,
      barrierDismissible: false,
      builder: (context) => const _AutoPlayPromptDialog(),
    );
    return accepted ?? false;
  }

  void _shareSong() {
    final box =
        _shareButtonKey.currentContext?.findRenderObject() as RenderBox?;
    final origin = box != null
        ? (box.localToGlobal(Offset.zero) & box.size)
        : null;
    final url = 'https://www.youtube.com/watch?v=$_videoId';

    SharePlus.instance.share(
      ShareParams(
        text: currentStrings.shareMessage(_title, url, appDownloadUrl),
        subject: _title,
        sharePositionOrigin: origin,
      ),
    );
  }

  void _showViewCountNotice() {
    final t = currentStrings;
    showDialog<void>(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: AppColors.surfaceRaised,
        title: Text(
          t.viewCountNoticeTitle,
          style: TextStyle(color: AppColors.textPrimary),
        ),
        content: Text(
          t.viewCountNoticeBody,
          style: TextStyle(color: AppColors.textSecondary),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: Text(
              t.viewCountNoticeDismiss,
              style: const TextStyle(
                color: AppColors.accent,
                fontWeight: FontWeight.w700,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _playNextAfterDelay() async {
    final songs = _relatedSongs;
    if (songs.isEmpty) return;

    final currentIndex = songs.indexWhere((song) => song.id == _videoId);
    final Song? next = currentIndex == -1
        ? songs.first
        : (currentIndex + 1 < songs.length ? songs[currentIndex + 1] : null);
    if (next == null) return;

    final generation = _autoPlayGeneration;
    for (
      var secondsLeft = _autoPlayDelay.inSeconds;
      secondsLeft > 0;
      secondsLeft--
    ) {
      if (!mounted || generation != _autoPlayGeneration) return;
      setState(() => _autoAdvanceSecondsLeft = secondsLeft);
      await Future<void>.delayed(const Duration(seconds: 1));
    }
    if (!mounted || generation != _autoPlayGeneration) return;

    setState(() {
      _videoId = next.id;
      _title = next.title;
      _categoryId = next.categoryId;
      _handledEnd = false;
      _autoAdvanceSecondsLeft = null;
    });
    await _controller.loadVideoById(videoId: next.id);
    WidgetsBinding.instance.addPostFrameCallback((_) => _scrollToNowPlaying());
  }

  @override
  void initState() {
    super.initState();
    _controller = YoutubePlayerController.fromVideoId(
      videoId: widget.videoId,
      autoPlay: true,
      params: const YoutubePlayerParams(showFullscreenButton: true),
    );
    _playerSub = _controller.stream.listen(_onPlayerValueChanged);
  }

  @override
  void dispose() {
    _playerSub?.cancel();
    _controller.close();
    _scrollController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    final lang = LangController.instance.resolve();

    return Scaffold(
      backgroundColor: Colors.black,
      floatingActionButton: FloatingActionButton.small(
        heroTag: null,
        backgroundColor: AppColors.surfaceRaised,
        tooltip: t.viewCountNoticeTooltip,
        onPressed: _showViewCountNotice,
        child: Icon(
          Icons.verified_outlined,
          color: AppColors.textSecondary,
          size: 20,
        ),
      ),
      body: SafeArea(
        child: Column(
          children: [
            Row(
              children: [
                IconButton(
                  onPressed: () => Navigator.of(context).pop(),
                  icon: const Icon(
                    Icons.arrow_back,
                    color: Colors.white,
                    size: 20,
                  ),
                ),
                Expanded(
                  child: Text(
                    _title,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.w600,
                      fontSize: 14.5,
                    ),
                  ),
                ),
                IconButton(
                  key: _shareButtonKey,
                  onPressed: _shareSong,
                  tooltip: t.shareTooltip,
                  icon: const Icon(Icons.share, color: Colors.white, size: 20),
                ),
              ],
            ),
            SizedBox(
              width: double.infinity,
              child: AspectRatio(
                aspectRatio: 16 / 9,
                child: YoutubePlayer(controller: _controller),
              ),
            ),
            if (_autoAdvanceSecondsLeft != null)
              Container(
                width: double.infinity,
                color: Colors.black,
                padding: const EdgeInsets.fromLTRB(16, 10, 16, 10),
                child: Align(
                  alignment: Alignment.centerLeft,
                  child: Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 12,
                      vertical: 8,
                    ),
                    decoration: BoxDecoration(
                      color: const Color(0xFF2A2D37),
                      borderRadius: BorderRadius.circular(10),
                    ),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        const Icon(
                          Icons.skip_next,
                          color: AppColors.accent,
                          size: 16,
                        ),
                        const SizedBox(width: 6),
                        Text(
                          t.autoAdvanceCountdown(_autoAdvanceSecondsLeft!),
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 12.5,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            Expanded(
              child: Container(
                width: double.infinity,
                color: AppColors.background,
                child: Column(
                  children: [
                    Padding(
                      padding: const EdgeInsets.fromLTRB(16, 14, 16, 0),
                      child: TrendTabs(
                        active: _relatedTab,
                        showMusicVideos: _showMv,
                        onSelect: _selectRelatedTab,
                      ),
                    ),
                    Divider(height: 1, color: AppColors.border),
                    Expanded(
                      child: FutureBuilder<List<Song>>(
                        future: _relatedFuture,
                        builder: (context, snapshot) {
                          if (snapshot.connectionState !=
                              ConnectionState.done) {
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
                                  onRetry: () => setState(
                                    () => _relatedFuture = _loadRelated(),
                                  ),
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
                                    countryLabel(widget.country, lang),
                                  ),
                                ),
                              ],
                            );
                          }

                          return ListView.builder(
                            controller: _scrollController,
                            padding: const EdgeInsets.symmetric(vertical: 8),
                            itemCount: songs.length,
                            itemBuilder: (context, index) {
                              final song = songs[index];
                              final isPlaying = song.id == _videoId;
                              return SongTile(
                                key: isPlaying ? _nowPlayingKey : null,
                                song: song,
                                rank: index + 1,
                                tab: _relatedTab,
                                isPlaying: isPlaying,
                              );
                            },
                          );
                        },
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _AutoPlayPromptDialog extends StatefulWidget {
  const _AutoPlayPromptDialog();

  @override
  State<_AutoPlayPromptDialog> createState() => _AutoPlayPromptDialogState();
}

class _AutoPlayPromptDialogState extends State<_AutoPlayPromptDialog> {
  int _secondsLeft = _autoPlayPromptCountdownSeconds;
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _timer = Timer.periodic(const Duration(seconds: 1), _onTick);
  }

  void _onTick(Timer timer) {
    if (_secondsLeft <= 1) {
      timer.cancel();
      Navigator.of(context).pop(true);
      return;
    }
    setState(() => _secondsLeft--);
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    return AlertDialog(
      backgroundColor: AppColors.surfaceRaised,
      title: Text(
        t.autoPlayPromptTitle,
        style: TextStyle(color: AppColors.textPrimary),
      ),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            t.autoPlayPromptDescription,
            style: TextStyle(color: AppColors.textSecondary),
          ),
          const SizedBox(height: 8),
          Text(
            t.autoPlayPromptCountdown(_secondsLeft),
            style: TextStyle(color: AppColors.textTertiary, fontSize: 12.5),
          ),
        ],
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(false),
          child: Text(
            t.autoPlayPromptDecline,
            style: TextStyle(color: AppColors.textSecondary),
          ),
        ),
        TextButton(
          onPressed: () => Navigator.of(context).pop(true),
          child: Text(
            t.autoPlayPromptAccept,
            style: const TextStyle(
              color: AppColors.accent,
              fontWeight: FontWeight.w700,
            ),
          ),
        ),
      ],
    );
  }
}
