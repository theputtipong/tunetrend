import 'strings.dart';

String _minutesAgo(int n) => '${n}m ago';
String _hoursAgo(int n) => '${n}h ago';
String _daysAgo(int n) => '${n}d ago';
String _weeksAgo(int n) => '${n}w ago';

String _emptyDescription(String countryName) =>
    'No chart data for $countryName in this view yet. New data syncs every 3 hours — check back soon.';

String _closingIn(int n) => 'Closing in ${n}s…';

const enStrings = AppStrings(
  aboutTooltip: 'About TuneTrend',
  themeToggleTooltip: 'Toggle light/dark theme',
  languageToggleTooltip: 'Switch language',
  backToTrends: '← Back to trends',
  tabTrending: 'Trending',
  tabNew: 'New Releases',
  tabMv: 'Music Videos',
  errorTitle: "Couldn't load trending songs",
  errorDescription:
      "We couldn't reach the TuneTrend service. Check your connection and try again.",
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
  aboutStack: ['GO · FIBER', 'NEXT.JS', 'FLUTTER', 'POSTGRESQL', 'CI/CD PIPELINES'],
  onboardingCountryTitle: 'Pick a country',
  onboardingCountryDescription: "Switch between TH, KR, JP, US, and GB to see what's trending there.",
  onboardingTabsTitle: 'Browse by category',
  onboardingTabsDescription: 'Jump between Trending, New Releases, and Music Videos.',
  onboardingLanguageTitle: 'Language',
  onboardingLanguageDescription: 'Switch the whole app between English and Thai.',
  onboardingThemeTitle: 'Light or dark',
  onboardingThemeDescription: 'Toggle the theme — TuneTrend also follows your system setting by default.',
  onboardingAboutTitle: 'About TuneTrend',
  onboardingAboutDescription: 'See the engineering behind the app and how to get in touch.',
  onboardingMenuTitle: 'More options',
  onboardingMenuDescription: 'Tap here for the tour, language, theme, and about settings.',
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
  contactErrorInvalidPhone: 'Please enter a valid Thai phone number (e.g. 0812345678).',
  contactErrorTooLong: 'Message is too long (max 2000 characters).',
  contactErrorRateLimited: 'Too many attempts. Please try again in a few minutes.',
  contactErrorGeneric: 'Something went wrong. Please try again later.',
  supportDevelopment: 'Support our development',
  replayTourTooltip: 'Show me around again',
);
