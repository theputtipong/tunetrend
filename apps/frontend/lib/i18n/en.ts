export const en = {
  nav: {
    countryLabel: "Country",
    syncNote: "Refreshes every 3 hours",
    about: "About TuneTrend",
    themeToggle: "Toggle light/dark theme",
    languageToggle: "Switch language",
    buyCoffee: "Buy me a coffee",
    replayTour: "Show me around again",
    menu: "More options",
    musicCategory: "Music",
    privacy: "Privacy Policy",
  },
  tabs: {
    trending: "Trending",
    new: "New Releases",
    mv: "Music Videos",
  },
  countryPageTitle: (countryName: string) => `Trending in ${countryName} · TuneTrend`,
  songList: {
    caption: (count: number, countryName: string) =>
      `Top ${count} · ${countryName} · sorted by views`,
    emptyTitle: "No trending songs yet",
    emptyDescription: (countryName: string) =>
      `No chart data for ${countryName} in this view yet. New data syncs every 3 hours — check back soon.`,
  },
  error: {
    title: "Couldn't load trending songs",
    description: "We couldn't reach the TuneTrend service. Check your connection and try again.",
    retry: "Retry",
  },
  format: {
    minutesAgo: (n: number) => `${n}m ago`,
    hoursAgo: (n: number) => `${n}h ago`,
    daysAgo: (n: number) => `${n}d ago`,
    weeksAgo: (n: number) => `${n}w ago`,
    views: "views",
  },
  watch: {
    back: "Back",
    shareTooltip: "Share",
    shareMessage: (title: string, url: string, appUrl?: string) => {
      const base = `${title}\n${url}\n\nShared via TuneTrend`;
      return appUrl ? `${base}\nGet the app: ${appUrl}` : base;
    },
    shareCopied: "Link copied to clipboard",
    viewCountNoticeTooltip: "About view counts",
    viewCountNoticeTitle: "About view counts",
    viewCountNoticeBody:
      "This video plays through YouTube's official embedded player. Views, " +
      "watch time, and ad revenue are counted normally for the creator — " +
      "exactly as on youtube.com. TuneTrend only aggregates public trending " +
      "data via the YouTube Data API; it never downloads, rehosts, or " +
      "interferes with your content in any way.",
    viewCountNoticeDismiss: "Got it",
    autoPlayPromptTitle: "Keep playing?",
    autoPlayPromptDescription:
      "Automatically play the next video in this list every time one finishes.",
    autoPlayPromptAccept: "Play continuously",
    autoPlayPromptDecline: "No thanks",
    autoPlayPromptCountdown: (secondsLeft: number) => `Auto-continuing in ${secondsLeft}s…`,
    autoAdvanceCountdown: (secondsLeft: number) => `Playing next in ${secondsLeft}s…`,
    autoPlayStoppedMessage: "Continuous play stopped — you switched tabs.",
    discoverHeading: "Explore other categories",
  },
  installPrompt: {
    title: "Get the TuneTrend app",
    description: "For the best experience, get the TuneTrend app.",
    getItButtonAndroid: "Get it on Google Play",
    getItButtonIos: "Download on the App Store",
    installButtonAndroidPwa: "Install app",
    iosAddToHomeScreenDescription: 'Tap the Share button, then choose "Add to Home Screen".',
    gotIt: "Got it",
    continueInBrowser: "Continue in browser",
  },
  about: {
    metaTitle: "About — TuneTrend",
    metaDescription: "The engineering and architecture behind TuneTrend.",
    backToTrends: "← Back to trends",
    eyebrow: "Behind The Code",
    heading: "Bank-grade architecture meets pop-culture data.",
    lead: "TuneTrend didn't just start as a weekend itch to see what the world is listening to. It was built as a sandbox to apply a decade of full-stack experience, microservices architecture, and automated deployments into a highly polished, multi-platform product.",
    stats: [
      { value: "10+", label: "Years engineering scalable systems" },
      { value: "7", label: "Years coding in modern mobile languages" },
      { value: "90%+", label: "Test coverage standard maintained" },
      { value: "100%", label: "Automated CI/CD deployments & monitoring" },
    ],
    bodyP1:
      "By day, my focus is on real financial infrastructure—architecting microservices for complex payment systems, hardening code for strict penetration tests, and ensuring system reliability for one of Thailand's largest banks.",
    bodyP2:
      "By night, that exact same rigor powers TuneTrend. The system is designed to seamlessly pull live YouTube chart data across five countries, refreshed every 3 hours. It routes through a robust Go backend, serving both a Next.js web application and a Flutter mobile app—all sharing a single, cohesive design system and continuous integration pipeline.",
    stackHeading: "The Architecture",
    stackCaption: "Built for scale, speed, and cross-platform consistency.",
    stack: ["GO · FIBER", "NEXT.JS", "FLUTTER", "POSTGRESQL", "CI/CD PIPELINES"],
    contact: {
      openButton: "Contact me",
      nameLabel: "Name (optional)",
      namePlaceholder: "Your name",
      messageLabel: "Message",
      messagePlaceholder: "What would you like to say?",
      methodLabel: "How should I contact you back?",
      methodEmail: "Email",
      methodPhone: "Phone",
      contactEmailPlaceholder: "you@example.com",
      contactPhonePlaceholder: "08xxxxxxxx",
      submit: "Send message",
      submitting: "Sending…",
      successMessage: "Thanks! Your message has been sent.",
      closingIn: (n: number) => `Closing in ${n}s…`,
      errors: {
        messageRequired: "Please write a message.",
        invalidEmail: "Please enter a valid email address.",
        invalidPhone: "Please enter a valid Thai phone number (e.g. 0812345678).",
        tooShort: "Message is too short (min 10 characters).",
        tooLong: "Message is too long (max 2000 characters).",
        rateLimited: "Too many attempts. Please try again in a few minutes.",
        generic: "Something went wrong. Please try again later.",
      },
    },
  },
  privacy: {
    metaTitle: "Privacy Policy — TuneTrend",
    metaDescription: "How TuneTrend collects, uses, and protects your data.",
    backToTrends: "← Back to trends",
    eyebrow: "Legal",
    heading: "Privacy Policy",
    lastUpdated: "Last updated: August 29, 2026",
    intro:
      'TuneTrend (the "app", "we", "us") is a small, independently-run project that shows ' +
      "trending YouTube music and videos by country, available as a website and as Android/iOS " +
      "apps. This page explains what data we collect, why, and how you're in control of it. " +
      "There are no user accounts on TuneTrend — nothing here requires a login.",
    sections: [
      {
        heading: "Information we collect",
        body: [
          "Usage analytics (Firebase Analytics on mobile, Vercel Analytics on web): anonymous, " +
            "aggregated events like which screen you viewed, which video you played, or which " +
            "category you filtered by. These are not tied to your name or contact details.",
          "Crash and performance data (Firebase Crashlytics and Firebase Performance " +
            "Monitoring, mobile only): if the app crashes or runs slowly, we automatically " +
            "receive a stack trace and basic device/OS info so we can fix the bug.",
          "Push notification token (mobile only, if you allow notifications): a device token " +
            "used to send you announcements. We only send broadcast messages to everyone " +
            "subscribed — we never target you individually or know who a token belongs to.",
          "Contact form submissions: if you use the in-app Contact form, we receive whatever " +
            "you type (your message, and an email or Thai phone number so we can reply). This " +
            "is the only place TuneTrend asks for anything resembling personal contact " +
            "information, and it's entirely optional.",
          "IP address (server-side, transient): used only to rate-limit abuse of the contact " +
            "form and our API. We don't log or store IP addresses long-term.",
          "Local preferences: your theme (light/dark), language, and whether you've seen the " +
            "onboarding tour or install prompt are saved in your browser's local storage or " +
            "your device's app storage. This never leaves your device.",
        ],
      },
      {
        heading: "YouTube playback",
        body: [
          "Every video on TuneTrend plays through YouTube's official embedded player " +
            "(iframe on web, the YouTube IFrame Player API on mobile). When you watch a video, " +
            "your interaction with that player is subject to Google/YouTube's own Privacy " +
            "Policy, not just ours. Views, watch time, and ad revenue are counted normally for " +
            "the video's creator — exactly as if you watched it on youtube.com. TuneTrend never " +
            "downloads, rehosts, or modifies anyone's video content.",
        ],
      },
      {
        heading: "Third-party services we use",
        body: [
          "Google Firebase (Analytics, Crashlytics, Cloud Messaging, Remote Config, " +
            "Performance Monitoring) — mobile app infrastructure.",
          "YouTube / Google — video playback and chart data (via the YouTube Data API).",
          "Vercel (Analytics, Speed Insights, hosting) — web app infrastructure.",
          "Resend — delivers the email notification when you submit the Contact form.",
          "Upstash — short-lived rate-limiting counters keyed by IP address.",
          "Google Play Store / Apple App Store — app distribution, subject to their own " +
            "privacy terms for anything handled at the OS/store level.",
        ],
      },
      {
        heading: "No ads, no data sales",
        body: [
          "TuneTrend shows no advertising and includes no ad SDKs or ad trackers. We do not " +
            "sell, rent, or trade your data to anyone, for any reason.",
        ],
      },
      {
        heading: "Data retention",
        body: [
          "Contact form messages are kept only as long as needed to respond to you. Analytics, " +
            "crash, and performance data follow Firebase's and Vercel's own default retention " +
            "periods. Rate-limiting counters expire automatically within minutes.",
        ],
      },
      {
        heading: "Children's privacy",
        body: [
          "TuneTrend is not directed at children under 13, and we do not knowingly collect " +
            "personal information from children. If you believe a child has provided us " +
            "information through the Contact form, please reach out and we'll remove it.",
        ],
      },
      {
        heading: "Your choices",
        body: [
          "You can disable push notifications at any time in your device's system settings. " +
            "You can clear your saved theme/language preferences by clearing your browser's " +
            "site data or the app's storage. Because there's no account system, there's no " +
            "personal profile to delete — the only data we hold that identifies you at all is " +
            "whatever you voluntarily typed into the Contact form, and you can ask us to delete " +
            "it the same way: through that form.",
        ],
      },
      {
        heading: "Changes to this policy",
        body: [
          "If this policy changes, we'll update the date at the top of this page. Continued use " +
            "of TuneTrend after a change means you accept the updated policy.",
        ],
      },
      {
        heading: "Contact us",
        body: [
          "Questions about this policy or your data? Use the Contact button on the About page " +
            "— it reaches the developer directly.",
        ],
      },
    ],
  },
  onboarding: {
    countryTitle: "Pick a country",
    countryDescription: "Switch between TH, KR, JP, US, and GB to see what's trending there.",
    tabsTitle: "Browse by category",
    tabsDescription: "Jump between Trending, New Releases, and Music Videos.",
    languageTitle: "Language",
    languageDescription: "Switch the whole app between English and Thai.",
    themeTitle: "Light or dark",
    themeDescription: "Toggle the theme — TuneTrend also follows your system setting by default.",
    mobileMenuTitle: "More options",
    mobileMenuDescription: "Tap here for the tour, language, theme, and about settings.",
    aboutTitle: "About TuneTrend",
    aboutDescription: "See the engineering behind the app and how to get in touch.",
    coffeeTitle: "Enjoying TuneTrend?",
    coffeeDescription: "You can buy me a coffee if you'd like to support the project.",
    next: "Next",
    previous: "Back",
    done: "Got it",
  },
};

export type Dictionary = typeof en;
