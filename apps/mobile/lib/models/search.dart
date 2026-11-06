import 'song.dart';

class LocalSearchResult {
  final List<Song> items;
  final bool suggestYouTube;

  const LocalSearchResult({required this.items, required this.suggestYouTube});
}

class YouTubeSearchResult {
  final List<Song> items;
  final int remaining;
  final bool cached;

  const YouTubeSearchResult({
    required this.items,
    required this.remaining,
    required this.cached,
  });
}
