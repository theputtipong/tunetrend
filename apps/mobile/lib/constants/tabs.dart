// รหัสหมวดหมู่ YouTube "Music" — เพลงจาก /trends, /trends/new, /trends/mv (ไม่ผ่าน category filter)
// ใช้ค่านี้เท่ากับ MusicCategoryID ฝั่ง backend (apps/backend/internal/domain/song.go)
const kMusicCategoryId = '10';

enum TrendTab {
  trending('trending', '/trends'),
  newReleases('new', '/trends/new'),
  musicVideos('mv', '/trends/mv');

  final String key;
  final String endpointPath;

  const TrendTab(this.key, this.endpointPath);
}
