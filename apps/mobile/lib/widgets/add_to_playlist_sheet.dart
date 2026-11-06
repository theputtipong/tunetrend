import 'package:flutter/material.dart';

import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import '../models/playlist.dart';
import '../models/song.dart';
import '../services/analytics_service.dart';
import '../services/api_client.dart';
import '../services/auth_service.dart';
import 'login_prompt.dart';
import 'playlist_name_dialog.dart';

Future<void> showAddToPlaylistSheet(
  BuildContext context,
  Song song, {
  String source = 'chart',
}) {
  return showModalBottomSheet<void>(
    context: context,
    isScrollControlled: true,
    useSafeArea: true,
    backgroundColor: AppColors.surface,
    shape: const RoundedRectangleBorder(
      borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
    ),
    builder: (context) => _AddToPlaylistSheet(song: song, source: source),
  );
}

class _AddToPlaylistSheet extends StatefulWidget {
  final Song song;
  final String source;

  const _AddToPlaylistSheet({required this.song, required this.source});

  @override
  State<_AddToPlaylistSheet> createState() => _AddToPlaylistSheetState();
}

class _AddToPlaylistSheetState extends State<_AddToPlaylistSheet> {
  final _api = ApiClient();
  Future<List<Playlist>>? _future;
  final Set<String> _saved = {};
  final Set<String> _pending = {};

  @override
  void initState() {
    super.initState();
    AuthService.instance.addListener(_onAuthChanged);
    if (AuthService.instance.isSignedIn) _future = _api.fetchPlaylists();
  }

  @override
  void dispose() {
    AuthService.instance.removeListener(_onAuthChanged);
    super.dispose();
  }

  void _onAuthChanged() {
    if (!mounted) return;
    setState(() {
      _future = AuthService.instance.isSignedIn ? _api.fetchPlaylists() : null;
    });
  }

  void _toast(String message) {
    ScaffoldMessenger.of(context)
      ..hideCurrentSnackBar()
      ..showSnackBar(SnackBar(content: Text(message)));
  }

  Future<void> _add(Playlist playlist) async {
    if (_saved.contains(playlist.id) || _pending.contains(playlist.id)) return;
    final t = currentStrings;
    setState(() => _pending.add(playlist.id));
    try {
      await _api.addPlaylistItem(
        playlist.id,
        videoId: widget.song.id,
        title: widget.song.title,
        channelTitle: widget.song.channelTitle,
        thumbnailUrl: widget.song.safeThumbnailUrl,
      );
      AnalyticsService().logVideoAddedToPlaylist(
        videoId: widget.song.id,
        source: widget.source,
      );
      if (!mounted) return;
      setState(() => _saved.add(playlist.id));
      _toast(t.savedToPlaylist(playlist.name));
    } on AlreadyInPlaylistException {
      if (!mounted) return;
      setState(() => _saved.add(playlist.id));
      _toast(t.alreadyInPlaylist(playlist.name));
    } on ApiException catch (e) {
      if (!mounted) return;
      _toast(e is UnauthorizedException ? t.signInFailed : t.actionFailed);
    } finally {
      if (mounted) setState(() => _pending.remove(playlist.id));
    }
  }

  Future<void> _createAndAdd() async {
    final t = currentStrings;
    final name = await showPlaylistNameDialog(
      context,
      title: t.newPlaylist,
      confirmLabel: t.create,
    );
    if (name == null || !mounted) return;
    try {
      final playlist = await _api.createPlaylist(name);
      AnalyticsService().logPlaylistCreated(source: 'add_sheet');
      if (!mounted) return;
      setState(() => _future = _api.fetchPlaylists());
      await _add(playlist);
    } on ApiException catch (e) {
      if (!mounted) return;
      _toast(e.statusCode == 422 ? e.message : t.actionFailed);
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    final future = _future;
    return Padding(
      padding: EdgeInsets.only(bottom: MediaQuery.viewInsetsOf(context).bottom),
      child: ConstrainedBox(
        constraints: BoxConstraints(
          maxHeight: MediaQuery.sizeOf(context).height * 0.7,
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Center(
              child: Container(
                margin: const EdgeInsets.only(top: 10, bottom: 6),
                width: 36,
                height: 4,
                decoration: BoxDecoration(
                  color: AppColors.border,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
            ),
            if (future == null)
              LoginPrompt(
                title: t.loginPromptTitle,
                description: t.loginPromptDescription,
                buttonLabel: t.continueWithGoogle,
              )
            else ...[
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 8, 20, 8),
                child: Text(
                  t.saveToPlaylistTitle(widget.song.title),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                  style: displayFont(fontSize: 16),
                ),
              ),
              Flexible(
                child: FutureBuilder<List<Playlist>>(
                  future: future,
                  builder: (context, snapshot) {
                    if (snapshot.connectionState != ConnectionState.done) {
                      return const Padding(
                        padding: EdgeInsets.all(24),
                        child: Center(
                          child: CircularProgressIndicator(
                            color: AppColors.accent,
                          ),
                        ),
                      );
                    }
                    if (snapshot.hasError) {
                      return Padding(
                        padding: const EdgeInsets.all(20),
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(
                              t.errorDescription,
                              textAlign: TextAlign.center,
                              style: TextStyle(color: AppColors.textSecondary),
                            ),
                            TextButton(
                              onPressed: () => setState(
                                () => _future = _api.fetchPlaylists(),
                              ),
                              child: Text(
                                t.retry,
                                style: const TextStyle(color: AppColors.accent),
                              ),
                            ),
                          ],
                        ),
                      );
                    }
                    final playlists = snapshot.data ?? const [];
                    return ListView(
                      shrinkWrap: true,
                      children: [
                        for (final p in playlists)
                          ListTile(
                            onTap: () => _add(p),
                            leading: _pending.contains(p.id)
                                ? const SizedBox(
                                    width: 22,
                                    height: 22,
                                    child: CircularProgressIndicator(
                                      strokeWidth: 2,
                                      color: AppColors.accent,
                                    ),
                                  )
                                : Icon(
                                    _saved.contains(p.id)
                                        ? Icons.check_circle
                                        : Icons.queue_music,
                                    color: _saved.contains(p.id)
                                        ? AppColors.accent
                                        : AppColors.textSecondary,
                                  ),
                            title: Text(
                              p.name,
                              style: TextStyle(color: AppColors.textPrimary),
                            ),
                            trailing: Text(
                              '${p.itemCount}',
                              style: TextStyle(color: AppColors.textTertiary),
                            ),
                          ),
                        ListTile(
                          onTap: _createAndAdd,
                          leading: const Icon(
                            Icons.add,
                            color: AppColors.accent,
                          ),
                          title: Text(
                            t.createNewPlaylist,
                            style: const TextStyle(
                              color: AppColors.accent,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ),
                        const SizedBox(height: 12),
                      ],
                    );
                  },
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
