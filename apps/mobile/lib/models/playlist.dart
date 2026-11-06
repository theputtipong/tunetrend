import 'song.dart';

class PlaylistItem {
  final String videoId;
  final int position;
  final String title;
  final String channelTitle;
  final String thumbnailUrl;
  final String addedAt;

  const PlaylistItem({
    required this.videoId,
    required this.position,
    required this.title,
    required this.channelTitle,
    required this.thumbnailUrl,
    required this.addedAt,
  });

  factory PlaylistItem.fromJson(Map<String, dynamic> json) {
    return PlaylistItem(
      videoId: json['videoId'] as String? ?? '',
      position: (json['position'] as num?)?.toInt() ?? 0,
      title: json['title'] as String? ?? '',
      channelTitle: json['channelTitle'] as String? ?? '',
      thumbnailUrl: json['thumbnailUrl'] as String? ?? '',
      addedAt: json['addedAt'] as String? ?? '',
    );
  }

  Song toSong(String country) => Song(
    id: videoId,
    title: title,
    channelTitle: channelTitle,
    thumbnailUrl: thumbnailUrl,
    viewCount: '',
    countryCode: country,
    categoryId: '',
    publishedAt: addedAt,
    videoType: 'General',
  );
}

class Playlist {
  final String id;
  final String name;
  final String updatedAt;
  final int itemCount;
  final List<PlaylistItem> items;

  const Playlist({
    required this.id,
    required this.name,
    required this.updatedAt,
    required this.itemCount,
    this.items = const [],
  });

  factory Playlist.fromJson(Map<String, dynamic> json) {
    final items = ((json['items'] as List<dynamic>?) ?? const [])
        .map((e) => PlaylistItem.fromJson(e as Map<String, dynamic>))
        .toList();
    return Playlist(
      id: json['id'] as String? ?? '',
      name: json['name'] as String? ?? '',
      updatedAt: json['updatedAt'] as String? ?? '',
      itemCount: items.isNotEmpty
          ? items.length
          : (json['itemCount'] as num?)?.toInt() ?? 0,
      items: items,
    );
  }
}
