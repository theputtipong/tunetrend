import 'strings.dart';

String _minutesAgo(int n) => '${n}m ago';
String _hoursAgo(int n) => '${n}h ago';
String _daysAgo(int n) => '${n}d ago';
String _weeksAgo(int n) => '${n}w ago';

String _emptyDescription(String countryName) =>
    'No chart data for $countryName in this view yet. New data syncs every 3 hours — check back soon.';

String _closingIn(int n) => 'Closing in ${n}s…';
String _autoPlayPromptCountdown(int n) => 'Auto-continuing in ${n}s…';
String _autoAdvanceCountdown(int n) => 'Playing next in ${n}s…';
String _shareMessage(String title, String url, String? appUrl) {
  final base = '$title\n$url\n\nShared via TuneTrend';
  return appUrl == null ? base : '$base\nGet the app: $appUrl';
}

String _searchHits(String n) => '$n searches';
String _youtubeRemaining(int remaining, int limit) =>
    '$remaining/$limit left today';
String _youtubeQuotaExhausted(String? time) => time == null
    ? 'Lots of searches today used up the quota. Chart search still works as usual.'
    : 'Lots of searches today used up the quota. It should be back around $time. Chart search still works as usual.';
String _youtubeUnavailable(String? time) => time == null
    ? 'Please try again later. Chart search still works as usual.'
    : 'Please try again around $time. Chart search still works as usual.';
String _youtubeUserLimit(int limit) =>
    "You've searched YouTube $limit times today. Try again tomorrow, or search our charts first.";
String _playlistSongCount(int n) => n == 1 ? '1 song' : '$n songs';
String _playlistNameInvalid(int max) => 'Name must be 1–$max characters';
String _deletePlaylistBody(String name) =>
    'Delete "$name" and all its songs. This can\'t be undone.';
String _saveToPlaylistTitle(String title) => 'Save "$title" to…';
String _savedToPlaylist(String name) => 'Saved to $name';
String _alreadyInPlaylist(String name) => 'Already in $name';

const enStrings = AppStrings(
  aboutTooltip: 'About TuneTrend',
  themeToggleTooltip: 'Toggle light/dark theme',
  languageToggleTooltip: 'Switch language',
  backToTrends: '← Back to trends',
  tabTrending: 'Trending',
  tabNew: 'New Releases',
  tabNewGeneric: 'New Videos',
  tabMv: 'Music Videos',
  musicCategory: 'Music',
  autoPlayPromptTitle: 'Keep playing?',
  autoPlayPromptDescription:
      'Automatically play the next video in this list every time one finishes.',
  autoPlayPromptAccept: 'Play continuously',
  autoPlayPromptDecline: 'No thanks',
  autoPlayPromptCountdown: _autoPlayPromptCountdown,
  autoAdvanceCountdown: _autoAdvanceCountdown,
  autoPlayStoppedMessage: 'Continuous play stopped — you switched tabs.',
  viewCountNoticeTooltip: 'About view counts',
  viewCountNoticeTitle: 'About view counts',
  viewCountNoticeBody:
      "This video plays through YouTube's official embedded player. Views, "
      "watch time, and ad revenue are counted normally for the creator — "
      "exactly as on youtube.com. TuneTrend only aggregates public trending "
      "data via the YouTube Data API; it never downloads, rehosts, or "
      "interferes with your content in any way.",
  viewCountNoticeDismiss: 'Got it',
  maintenanceTitle: 'TuneTrend is under maintenance',
  maintenanceDescription:
      "We're making some improvements and will be back shortly. "
      'Thanks for your patience!',
  updateRequiredTitle: 'Update required',
  updateRequiredDescription: 'A new version of TuneTrend is available. Please update to keep using the app.',
  updateNowButton: 'Update now',
  shareTooltip: 'Share',
  shareMessage: _shareMessage,
  errorTitle: "Couldn't load trending songs",
  errorDescription: "We couldn't reach the TuneTrend service. Check your connection and try again.",
  retry: 'Retry',
  emptyTitle: 'No trending songs yet',
  emptyDescription: _emptyDescription,
  viewsSuffix: 'views',
  minutesAgo: _minutesAgo,
  hoursAgo: _hoursAgo,
  daysAgo: _daysAgo,
  weeksAgo: _weeksAgo,
  aboutEyebrow: 'Behind The Code',
  aboutHeading: 'Bank-grade architecture meets pop-culture data.',
  aboutLead:
      "TuneTrend didn't just start as a weekend itch to see what the world is "
      'listening to. It was built as a sandbox to apply a decade of full-stack '
      'experience, microservices architecture, and automated deployments into a '
      'highly polished, multi-platform product.',
  aboutStats: [
    AboutStat('10+', 'Years engineering scalable systems'),
    AboutStat('7', 'Years coding in modern mobile languages'),
    AboutStat('90%+', 'Test coverage standard maintained'),
    AboutStat('100%', 'Automated CI/CD deployments & monitoring'),
  ],
  aboutBodyP1:
      'By day, my focus is on real financial infrastructure—architecting '
      'microservices for complex payment systems, hardening code for strict '
      "penetration tests, and ensuring system reliability for one of Thailand's "
      'largest banks.',
  aboutBodyP2:
      'By night, that exact same rigor powers TuneTrend. The system is designed '
      'to seamlessly pull live YouTube chart data across five countries, '
      'refreshed every 3 hours. It routes through a robust Go backend, serving '
      'both a Next.js web application and this Flutter mobile app—all sharing a '
      'single, cohesive design system and continuous integration pipeline.',
  aboutStackHeading: 'The Architecture',
  aboutStackCaption: 'Built for scale, speed, and cross-platform consistency.',
  aboutStack: [
    'GO · FIBER',
    'NEXT.JS',
    'FLUTTER',
    'POSTGRESQL',
    'CI/CD PIPELINES',
  ],
  onboardingCountryTitle: 'Pick a country',
  onboardingCountryDescription:
      "Switch between TH, KR, JP, US, and GB to see what's trending there.",
  onboardingTabsTitle: 'Switch views',
  onboardingTabsDescription:
      'Jump between Trending, New Releases, and Music Videos.',
  onboardingCategoryTitle: 'Filter by category',
  onboardingCategoryDescription:
      'Narrow the chart down to a specific YouTube category, like Gaming or Sports. '
      'Swipe sideways to see more.',
  onboardingLanguageTitle: 'Language',
  onboardingLanguageDescription:
      'Switch the whole app between English and Thai.',
  onboardingThemeTitle: 'Light or dark',
  onboardingThemeDescription: 'Toggle the theme — TuneTrend also follows your system setting by default.',
  onboardingAboutTitle: 'About TuneTrend',
  onboardingAboutDescription:
      'See the engineering behind the app and how to get in touch.',
  onboardingMenuTitle: 'More options',
  onboardingMenuDescription:
      'Tap here for the tour, language, theme, and about settings.',
  menuTooltip: 'More options',
  contactOpenButton: 'Contact me',
  contactNameLabel: 'Name (optional)',
  contactNamePlaceholder: 'Your name',
  contactMessageLabel: 'Message',
  contactMessagePlaceholder: 'What would you like to say?',
  contactMethodLabel: 'How should I contact you back?',
  contactMethodEmail: 'Email',
  contactMethodPhone: 'Phone',
  contactEmailPlaceholder: 'you@example.com',
  contactPhonePlaceholder: '08xxxxxxxx',
  contactSubmit: 'Send message',
  contactSubmitting: 'Sending…',
  contactSuccessMessage: 'Thanks! Your message has been sent.',
  closingIn: _closingIn,
  contactErrorMessageRequired: 'Please write a message.',
  contactErrorInvalidEmail: 'Please enter a valid email address.',
  contactErrorInvalidPhone:
      'Please enter a valid Thai phone number (e.g. 0812345678).',
  contactErrorTooLong: 'Message is too long (max 2000 characters).',
  contactErrorRateLimited:
      'Too many attempts. Please try again in a few minutes.',
  contactErrorGeneric: 'Something went wrong. Please try again later.',
  supportDevelopment: 'Support our development',
  replayTourTooltip: 'Show me around again',
  privacyPolicy: 'Privacy Policy',
  navCharts: 'Charts',
  navSearch: 'Search',
  navPlaylists: 'Playlists',
  tabSearched: 'Most Searched',
  searchedTabCaption: 'Songs people in this country searched for and played most in the last 7 days',
  newBadge: 'New',
  searchHits: _searchHits,
  searchHint: 'Search songs, artists…',
  recentSearches: 'Recent:',
  clearSearchTooltip: 'Clear',
  searchPromptTitle: 'Find a song',
  searchPromptDescription:
      'Search TuneTrend charts instantly, no sign-in needed.',
  searchQueryTooShort: 'Type at least 2 characters',
  localSectionTitle: 'In charts',
  localEmpty: 'Not found in TuneTrend charts',
  searchedBefore: 'Searched before',
  youtubeSectionTitle: 'From YouTube',
  youtubeMoreHint: 'Can\'t find it? Search YouTube too. Results are saved so others find them faster.',
  youtubeSearchButton: 'Search more on YouTube',
  youtubeRemaining: _youtubeRemaining,
  youtubeLoading: 'Loading…',
  youtubeEmpty: 'No results on YouTube',
  youtubeLoginTitle: 'Sign in to search YouTube',
  youtubeLoginDescription:
      'Chart search still works as usual. YouTube search is for members only.',
  signInWithGoogle: 'Sign in with Google',
  loginUnavailable: 'Login is temporarily unavailable. Please try again.',
  videoUnavailable: 'Video no longer available',
  signInFailed: 'Couldn\'t sign in. Please try again.',
  youtubeUnavailableTitle: 'YouTube search is temporarily unavailable',
  youtubeQuotaExhausted: _youtubeQuotaExhausted,
  youtubeUnavailable: _youtubeUnavailable,
  youtubeUserLimitTitle: 'You\'ve used today\'s YouTube searches',
  youtubeUserLimit: _youtubeUserLimit,
  youtubeErrorTitle: 'Couldn\'t search YouTube',
  youtubeErrorDescription: 'Something went wrong. Please try again.',
  playlistsTitle: 'My playlists',
  playlistSongCount: _playlistSongCount,
  newPlaylist: 'New playlist',
  createNewPlaylist: 'New playlist',
  playlistNameLabel: 'Playlist name',
  playlistNameInvalid: _playlistNameInvalid,
  create: 'Create',
  save: 'Save',
  cancel: 'Cancel',
  delete: 'Delete',
  playlistsEmptyTitle: 'No playlists yet',
  playlistsEmptyDescription: 'Tap + on any song to save it for later.',
  playlistItemsEmpty: 'No songs in this playlist yet',
  playlistPrivate: 'Private playlist',
  playAll: 'Play all',
  rename: 'Rename',
  deletePlaylistTitle: 'Delete this playlist?',
  deletePlaylistBody: _deletePlaylistBody,
  removeFromPlaylist: 'Remove',
  addToPlaylistTooltip: 'Save to playlist',
  saveToPlaylistTitle: _saveToPlaylistTitle,
  savedToPlaylist: _savedToPlaylist,
  alreadyInPlaylist: _alreadyInPlaylist,
  actionFailed: 'Something went wrong. Please try again.',
  loginPromptTitle: 'Save songs for later',
  loginPromptDescription: 'Sign in with Google to create your own playlists. Works on both web and app.',
  continueWithGoogle: 'Continue with Google',
  loginPrivacyNote: 'We only keep your name, email and profile photo. Nothing is posted to your YouTube.',
  accountMenuTooltip: 'Account',
  signOut: 'Sign out',
  deleteAccount: 'Delete account',
  deleteAccountTitle: 'Delete your account?',
  deleteAccountBody:
      'Your account and all your playlists will be permanently deleted.',
  accountDeleted: 'Account deleted',
);
