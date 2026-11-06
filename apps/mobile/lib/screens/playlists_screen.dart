import 'package:flutter/material.dart';

import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import '../models/playlist.dart';
import '../services/analytics_service.dart';
import '../services/api_client.dart';
import '../services/auth_service.dart';
import '../widgets/cached_thumbnail.dart';
import '../widgets/login_prompt.dart';
import '../widgets/playlist_name_dialog.dart';
import '../widgets/song_list_skeleton.dart';
import '../widgets/state_message.dart';
import 'playlist_detail_screen.dart';

enum _AccountAction { signOut, deleteAccount }

class PlaylistsScreen extends StatefulWidget {
  final String country;
  final bool active;

  const PlaylistsScreen({
    super.key,
    required this.country,
    required this.active,
  });

  @override
  State<PlaylistsScreen> createState() => _PlaylistsScreenState();
}

class _PlaylistsScreenState extends State<PlaylistsScreen> {
  final _api = ApiClient();
  Future<List<Playlist>>? _future;
  bool _wasSignedIn = false;

  @override
  void initState() {
    super.initState();
    AuthService.instance.addListener(_onAuthChanged);
    _wasSignedIn = AuthService.instance.isSignedIn;
    if (_wasSignedIn) _future = _api.fetchPlaylists();
  }

  @override
  void didUpdateWidget(PlaylistsScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.active && !oldWidget.active && AuthService.instance.isSignedIn) {
      _future = _api.fetchPlaylists();
    }
  }

  @override
  void dispose() {
    AuthService.instance.removeListener(_onAuthChanged);
    super.dispose();
  }

  void _onAuthChanged() {
    final signedIn = AuthService.instance.isSignedIn;
    if (!mounted) return;
    setState(() {
      if (signedIn && !_wasSignedIn) _future = _api.fetchPlaylists();
      if (!signedIn) _future = null;
      _wasSignedIn = signedIn;
    });
  }

  Future<void> _reload() async {
    final future = _api.fetchPlaylists();
    setState(() => _future = future);
    try {
      await future;
    } catch (_) {
      return;
    }
  }

  void _toast(String message) {
    ScaffoldMessenger.of(context)
      ..hideCurrentSnackBar()
      ..showSnackBar(SnackBar(content: Text(message)));
  }

  Future<void> _create() async {
    final t = currentStrings;
    final name = await showPlaylistNameDialog(
      context,
      title: t.newPlaylist,
      confirmLabel: t.create,
    );
    if (name == null || !mounted) return;
    try {
      await _api.createPlaylist(name);
      AnalyticsService().logPlaylistCreated(source: 'playlists_screen');
      if (mounted) _reload();
    } on ApiException catch (e) {
      if (mounted) _toast(e.statusCode == 422 ? e.message : t.actionFailed);
    }
  }

  Future<void> _open(Playlist playlist) async {
    await Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => PlaylistDetailScreen(
          playlistId: playlist.id,
          initialName: playlist.name,
          country: widget.country,
        ),
      ),
    );
    if (mounted && AuthService.instance.isSignedIn) _reload();
  }

  Future<void> _onAccountAction(_AccountAction action) async {
    final t = currentStrings;
    switch (action) {
      case _AccountAction.signOut:
        await AuthService.instance.signOut();
      case _AccountAction.deleteAccount:
        final ok = await showConfirmDialog(
          context,
          title: t.deleteAccountTitle,
          body: t.deleteAccountBody,
          confirmLabel: t.deleteAccount,
        );
        if (!ok || !mounted) return;
        try {
          await AuthService.instance.deleteAccount();
          if (mounted) _toast(t.accountDeleted);
        } on ApiException {
          if (mounted) _toast(t.actionFailed);
        }
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    final future = _future;
    final signedIn = AuthService.instance.isSignedIn;
    return Scaffold(
      floatingActionButton: signedIn
          ? FloatingActionButton.extended(
              heroTag: null,
              onPressed: _create,
              backgroundColor: AppColors.accent,
              foregroundColor: AppColors.accentInk,
              icon: const Icon(Icons.add),
              label: Text(
                t.newPlaylist,
                style: const TextStyle(fontWeight: FontWeight.w700),
              ),
            )
          : null,
      body: SafeArea(
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 16, 8, 8),
              child: Row(
                children: [
                  Expanded(
                    child: Text(
                      t.playlistsTitle,
                      style: displayFont(fontSize: 20),
                    ),
                  ),
                  if (signedIn) _AccountMenu(onSelected: _onAccountAction),
                ],
              ),
            ),
            Expanded(
              child: future == null
                  ? Center(
                      child: SingleChildScrollView(
                        child: LoginPrompt(
                          title: t.loginPromptTitle,
                          description: t.loginPromptDescription,
                          buttonLabel: t.continueWithGoogle,
                        ),
                      ),
                    )
                  : RefreshIndicator(
                      color: AppColors.accent,
                      backgroundColor: AppColors.surface,
                      onRefresh: _reload,
                      child: FutureBuilder<List<Playlist>>(
                        future: future,
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
                                  onRetry: _reload,
                                ),
                              ],
                            );
                          }
                          final playlists = snapshot.data ?? const [];
                          if (playlists.isEmpty) {
                            return ListView(
                              children: [
                                StateMessage(
                                  variant: StateMessageVariant.empty,
                                  title: t.playlistsEmptyTitle,
                                  description: t.playlistsEmptyDescription,
                                ),
                              ],
                            );
                          }
                          return ListView.separated(
                            padding: const EdgeInsets.fromLTRB(8, 4, 8, 96),
                            itemCount: playlists.length,
                            separatorBuilder: (_, _) =>
                                const SizedBox(height: 2),
                            itemBuilder: (context, index) {
                              final p = playlists[index];
                              return _PlaylistRow(
                                playlist: p,
                                onTap: () => _open(p),
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

class _AccountMenu extends StatelessWidget {
  final ValueChanged<_AccountAction> onSelected;

  const _AccountMenu({required this.onSelected});

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    final user = AuthService.instance.user;
    final initial = (user?.name.isNotEmpty ?? false)
        ? user!.name.characters.first.toUpperCase()
        : '?';
    return PopupMenuButton<_AccountAction>(
      tooltip: t.accountMenuTooltip,
      onSelected: onSelected,
      icon: CircleAvatar(
        radius: 15,
        backgroundColor: AppColors.accent,
        child: Text(
          initial,
          style: const TextStyle(
            color: AppColors.accentInk,
            fontWeight: FontWeight.w700,
            fontSize: 13,
          ),
        ),
      ),
      itemBuilder: (context) => [
        if (user != null)
          PopupMenuItem<_AccountAction>(
            enabled: false,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  user.name,
                  style: TextStyle(
                    color: AppColors.textPrimary,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                Text(
                  user.email,
                  style: TextStyle(color: AppColors.textTertiary, fontSize: 12),
                ),
              ],
            ),
          ),
        PopupMenuItem(
          value: _AccountAction.signOut,
          child: ListTile(
            leading: Icon(
              Icons.logout,
              size: 20,
              color: AppColors.textSecondary,
            ),
            title: Text(
              t.signOut,
              style: TextStyle(color: AppColors.textPrimary, fontSize: 14),
            ),
            contentPadding: EdgeInsets.zero,
          ),
        ),
        PopupMenuItem(
          value: _AccountAction.deleteAccount,
          child: ListTile(
            leading: Icon(
              Icons.delete_forever_outlined,
              size: 20,
              color: AppColors.errorText,
            ),
            title: Text(
              t.deleteAccount,
              style: TextStyle(color: AppColors.errorText, fontSize: 14),
            ),
            contentPadding: EdgeInsets.zero,
          ),
        ),
      ],
    );
  }
}

class _PlaylistRow extends StatelessWidget {
  final Playlist playlist;
  final VoidCallback onTap;

  const _PlaylistRow({required this.playlist, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(14),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 10),
        child: Row(
          children: [
            ClipRRect(
              borderRadius: BorderRadius.circular(10),
              child: SizedBox(
                width: 52,
                height: 52,
                child:
                    playlist.items.isNotEmpty &&
                        playlist.items.first.thumbnailUrl.isNotEmpty
                    ? CachedThumbnail(
                        imageUrl: playlist.items.first.thumbnailUrl,
                      )
                    : Container(
                        color: AppColors.surfaceRaised,
                        child: Icon(
                          Icons.queue_music,
                          color: AppColors.textTertiary,
                        ),
                      ),
              ),
            ),
            const SizedBox(width: 14),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    playlist.name,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 15,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    t.playlistSongCount(playlist.itemCount),
                    style: TextStyle(
                      fontSize: 12.5,
                      color: AppColors.textSecondary,
                    ),
                  ),
                ],
              ),
            ),
            Icon(Icons.chevron_right, color: AppColors.textTertiary),
          ],
        ),
      ),
    );
  }
}
