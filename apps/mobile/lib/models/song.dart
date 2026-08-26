class Song {
  final String id;
  final String title;
  final String channelTitle;
  final String thumbnailUrl;
  final String viewCount;
  final String countryCode;
  final String categoryId;
  final String publishedAt;
  final String videoType;

  const Song({
    required this.id,
    required this.title,
    required this.channelTitle,
    required this.thumbnailUrl,
    required this.viewCount,
    required this.countryCode,
    required this.categoryId,
    required this.publishedAt,
    required this.videoType,
  });

  factory Song.fromJson(Map<String, dynamic> json) {
    return Song(
      id: json['id'] as String? ?? '',
      title: json['title'] as String? ?? '',
      channelTitle: json['channelTitle'] as String? ?? '',
      thumbnailUrl: json['thumbnailUrl'] as String? ?? '',
      viewCount: json['viewCount'] as String? ?? '0',
      countryCode: json['countryCode'] as String? ?? '',
      categoryId: json['categoryId'] as String? ?? '',
      publishedAt: json['publishedAt'] as String? ?? '',
      videoType: json['videoType'] as String? ?? 'General',
    );
  }
}
