enum TrendTab {
  trending('trending', '/trends'),
  newReleases('new', '/trends/new'),
  musicVideos('mv', '/trends/mv');

  final String key;
  final String endpointPath;

  const TrendTab(this.key, this.endpointPath);
}
