import 'package:flutter/material.dart';

import '../constants/tabs.dart';
import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import '../models/playlist.dart';
import '../models/song.dart';
import '../services/api_client.dart';
import '../services/unavailable_videos.dart';
import '../widgets/playlist_name_dialog.dart';
import '../widgets/song_list_skeleton.dart';
import '../widgets/song_tile.dart';
import '../widgets/state_message.dart';
import 'video_player_screen.dart';

enum _PlaylistAction { rename, delete }

class PlaylistDetailScreen extends StatefulWidget {
  final String playlistId;
  final String initialName;
  final String country;

  const PlaylistDetailScreen({
    super.key,
    required this.playlistId,
    required this.initialName,
    required this.country,
  });

  @override
  State<PlaylistDetailScreen> createState() => _PlaylistDetailScreenState();
}

class _PlaylistDetailScreenState extends State<PlaylistDetailScreen> {
  final _api = ApiClient();

  late String _name = widget.initialName;
  List<PlaylistItem> _items = const [];
  bool _loading = true;
  bool _failed = false;

  List<Song> get _songs => _items.map((i) => i.toSong(widget.country)).toList();

  @override
  void initState() {
    super.initState();
    UnavailableVideos.instance.addListener(_onUnavailableChanged);
    _load();
  }

  @override
  void dispose() {
    UnavailableVideos.instance.removeListener(_onUnavailableChanged);
    super.dispose();
  }

  void _onUnavailableChanged() {
    if (mounted) setState(() {});
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _failed = false;
    });
    try {
      final playlist = await _api.fetchPlaylist(widget.playlistId);
      if (!mounted) return;
      setState(() {
        _name = playlist.name;
        _items = playlist.items;
        _loading = false;
      });
    } on NotFoundException {
      if (mounted) Navigator.of(context).pop();
    } on ApiException {
      if (!mounted) return;
      setState(() {
        _loading = false;
        _failed = true;
      });
    }
  }

  void _toast(String message) {
    ScaffoldMessenger.of(context)
      ..hideCurrentSnackBar()
      ..showSnackBar(SnackBar(content: Text(message)));
  }

  void _play(int index) {
    final songs = _songs;
    if (songs.isEmpty) return;
    final song = songs[index];
    Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => VideoPlayerScreen(
          videoId: song.id,
          title: song.title,
          country: widget.country,
          categoryId: '',
          tab: TrendTab.trending,
          queue: songs,
        ),
      ),
    );
  }

  Future<void> _rename() async {
    final t = currentStrings;
    final name = await showPlaylistNameDialog(
      context,
      title: t.rename,
      confirmLabel: t.save,
      initial: _name,
    );
    if (name == null || name == _name || !mounted) return;
    try {
      final updated = await _api.renamePlaylist(widget.playlistId, name);
      if (mounted) setState(() => _name = updated.name);
    } on ApiException catch (e) {
      if (mounted) _toast(e.statusCode == 422 ? e.message : t.actionFailed);
    }
  }

  Future<void> _delete() async {
    final t = currentStrings;
    final ok = await showConfirmDialog(
      context,
      title: t.deletePlaylistTitle,
      body: t.deletePlaylistBody(_name),
      confirmLabel: t.delete,
    );
    if (!ok || !mounted) return;
    try {
      await _api.deletePlaylist(widget.playlistId);
      if (mounted) Navigator.of(context).pop();
    } on ApiException {
      if (mounted) _toast(t.actionFailed);
    }
  }

  Future<void> _remove(PlaylistItem item) async {
    final previous = _items;
    setState(
      () => _items = [
        for (final i in _items)
          if (i.videoId != item.videoId) i,
      ],
    );
    try {
      await _api.removePlaylistItem(widget.playlistId, item.videoId);
    } on ApiException {
      if (!mounted) return;
      setState(() => _items = previous);
      _toast(currentStrings.actionFailed);
    }
  }

  Future<void> _reorder(int oldIndex, int newIndex) async {
    final previous = _items;
    final next = [..._items];
    final moved = next.removeAt(oldIndex);
    next.insert(newIndex, moved);
    setState(() => _items = next);
    try {
      final updated = await _api.reorderPlaylist(
        widget.playlistId,
        next.map((i) => i.videoId).toList(),
      );
      if (mounted) setState(() => _items = updated.items);
    } on ApiException {
      if (!mounted) return;
      setState(() => _items = previous);
      _toast(currentStrings.actionFailed);
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    return Scaffold(
      appBar: AppBar(
        backgroundColor: AppColors.background,
        surfaceTintColor: Colors.transparent,
        title: Text(
          _name,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
          style: displayFont(fontSize: 18),
        ),
        actions: [
          PopupMenuButton<_PlaylistAction>(
            icon: Icon(Icons.more_vert, color: AppColors.textSecondary),
            onSelected: (action) {
              switch (action) {
                case _PlaylistAction.rename:
                  _rename();
                case _PlaylistAction.delete:
                  _delete();
              }
            },
            itemBuilder: (context) => [
              PopupMenuItem(
                value: _PlaylistAction.rename,
                child: ListTile(
                  leading: Icon(
                    Icons.edit_outlined,
                    size: 20,
                    color: AppColors.textSecondary,
                  ),
                  title: Text(
                    t.rename,
                    style: TextStyle(
                      color: AppColors.textPrimary,
                      fontSize: 14,
                    ),
                  ),
                  contentPadding: EdgeInsets.zero,
                ),
              ),
              PopupMenuItem(
                value: _PlaylistAction.delete,
                child: ListTile(
                  leading: Icon(
                    Icons.delete_outline,
                    size: 20,
                    color: AppColors.errorText,
                  ),
                  title: Text(
                    t.delete,
                    style: TextStyle(color: AppColors.errorText, fontSize: 14),
                  ),
                  contentPadding: EdgeInsets.zero,
                ),
              ),
            ],
          ),
        ],
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    final t = currentStrings;
    if (_loading) return const SongListSkeleton();
    if (_failed) {
      return ListView(
        children: [
          StateMessage(
            variant: StateMessageVariant.error,
            title: t.errorTitle,
            description: t.errorDescription,
            retryLabel: t.retry,
            onRetry: _load,
          ),
        ],
      );
    }

    final songs = _songs;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 4, 16, 12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                '${t.playlistPrivate} · ${t.playlistSongCount(_items.length)}',
                style: TextStyle(color: AppColors.textTertiary, fontSize: 12.5),
              ),
              const SizedBox(height: 10),
              FilledButton.icon(
                onPressed: songs.isEmpty ? null : () => _play(0),
                style: FilledButton.styleFrom(
                  backgroundColor: AppColors.accent,
                  foregroundColor: AppColors.accentInk,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(999),
                  ),
                ),
                icon: const Icon(Icons.play_arrow_rounded),
                label: Text(
                  t.playAll,
                  style: const TextStyle(fontWeight: FontWeight.w700),
                ),
              ),
            ],
          ),
        ),
        Divider(height: 1, color: AppColors.border),
        Expanded(
          child: _items.isEmpty
              ? ListView(
                  children: [
                    StateMessage(
                      variant: StateMessageVariant.empty,
                      title: t.playlistItemsEmpty,
                      description: t.playlistsEmptyDescription,
                    ),
                  ],
                )
              : ReorderableListView.builder(
                  padding: const EdgeInsets.symmetric(vertical: 8),
                  buildDefaultDragHandles: false,
                  itemCount: _items.length,
                  onReorderItem: _reorder,
                  itemBuilder: (context, index) {
                    final item = _items[index];
                    final unavailable = UnavailableVideos.instance.contains(
                      item.videoId,
                    );
                    return Material(
                      key: ValueKey(item.videoId),
                      color: Colors.transparent,
                      child: SongTile(
                        song: songs[index],
                        rank: null,
                        dimmed: unavailable,
                        badge: unavailable ? t.videoUnavailable : null,
                        tab: TrendTab.trending,
                        queue: songs,
                        analyticsSource: 'playlist',
                        trailing: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            IconButton(
                              tooltip: t.removeFromPlaylist,
                              onPressed: () => _remove(item),
                              icon: Icon(
                                Icons.close,
                                size: 20,
                                color: AppColors.textTertiary,
                              ),
                            ),
                            ReorderableDragStartListener(
                              index: index,
                              child: Padding(
                                padding: const EdgeInsets.fromLTRB(4, 8, 12, 8),
                                child: Icon(
                                  Icons.drag_handle,
                                  color: AppColors.textTertiary,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                    );
                  },
                ),
        ),
      ],
    );
  }
}
