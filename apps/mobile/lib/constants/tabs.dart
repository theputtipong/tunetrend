// Must match MusicCategoryID in apps/backend/internal/domain/song.go.
const kMusicCategoryId = '10';

enum TrendTab {
  trending('trending', '/trends'),
  newReleases('new', '/trends/new'),
  musicVideos('mv', '/trends/mv'),
  searched('searched', '/trends/searched');

  final String key;
  final String endpointPath;

  const TrendTab(this.key, this.endpointPath);
}
