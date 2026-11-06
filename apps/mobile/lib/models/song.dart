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
  final String source;
  final String chartTab;
  final int? chartRank;
  final int? rankChange;
  final bool isNew;
  final int? hits;

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
    this.source = '',
    this.chartTab = '',
    this.chartRank,
    this.rankChange,
    this.isNew = false,
    this.hits,
  });

  factory Song.fromJson(Map<String, dynamic> json, {String country = ''}) {
    final countryCode = json['countryCode'] as String? ?? '';
    return Song(
      id: json['id'] as String? ?? '',
      title: json['title'] as String? ?? '',
      channelTitle: json['channelTitle'] as String? ?? '',
      thumbnailUrl: json['thumbnailUrl'] as String? ?? '',
      viewCount: json['viewCount'] as String? ?? '0',
      countryCode: countryCode.isEmpty ? country : countryCode,
      categoryId: json['categoryId'] as String? ?? '',
      publishedAt: json['publishedAt'] as String? ?? '',
      videoType: json['videoType'] as String? ?? 'General',
      source: json['source'] as String? ?? '',
      chartTab: json['chartTab'] as String? ?? '',
      chartRank: (json['chartRank'] as num?)?.toInt(),
      rankChange: (json['rankChange'] as num?)?.toInt(),
      isNew: json['isNew'] as bool? ?? false,
      hits: (json['hits'] as num?)?.toInt(),
    );
  }

  String get safeThumbnailUrl => thumbnailUrl.startsWith(kYtImageHost)
      ? thumbnailUrl
      : '${kYtImageHost}vi/$id/hqdefault.jpg';
}

const kYtImageHost = 'https://i.ytimg.com/';
