export const en = {
  nav: {
    countryLabel: "Country",
    syncNote: "Refreshes every 3 hours",
    searchedSyncNote: "Updates every 5 minutes",
    about: "About TuneTrend",
    themeToggle: "Toggle light/dark theme",
    languageToggle: "Switch language",
    buyCoffee: "Buy me a coffee",
    replayTour: "Show me around again",
    menu: "More options",
    musicCategory: "Music",
    privacy: "Privacy Policy",
    getApp: "Get the TuneTrend app",
    searchPlaceholder: "Search songs, artists…",
    searchLabel: "Search",
    login: "Sign in",
    accountMenu: "Account menu",
    myPlaylists: "My playlists",
    logout: "Sign out",
    deleteAccount: "Delete account",
    newBadge: "New",
  },
  tabs: {
    trending: "Trending",
    new: "New Releases",
    mv: "Music Videos",
    searched: "Most Searched",
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
  searched: {
    hint: "What people in this country searched for most in the last 7 days, with the latest results for each search",
    searches: (n: string) => `${n} searches`,
    viewAll: "See all results",
    emptyTitle: "Nothing searched yet",
    emptyDescription:
      "A search shows up here once at least 2 people have searched for it. Try searching for a song you love.",
  },
  search: {
    metaTitle: (q: string) => `"${q}" · TuneTrend search`,
    promptTitle: "Search TuneTrend",
    promptDescription: "Type at least 2 characters to search songs and artists.",
    chartSection: "In TuneTrend charts",
    chartCount: (n: number, countryName: string) => `${n} results · ${countryName}`,
    chartRankBadge: (rank: number, tabLabel: string) => `#${rank} ${tabLabel}`,
    searchedBefore: "Searched before",
    noLocalTitle: "Not in our charts yet",
    noLocalDescription: "Nothing in the TuneTrend charts matches this search.",
    errorTitle: "Search failed",
    errorDescription: "We couldn't reach the TuneTrend service. Please try again.",
    youtubeSection: "From YouTube",
    loading: "Loading…",
    ctaTitle: "Can't find what you're looking for?",
    ctaBody: "Keep searching on YouTube. What you find is saved so others can find it faster too.",
    ctaButton: "Search more on YouTube",
    remaining: (n: number) => `${n} left today`,
    loginTitle: "Sign in to keep searching on YouTube",
    loginBody: "Chart search still works as usual. YouTube search is for members only.",
    loginButton: "Sign in with Google",
    quotaTitle: "YouTube search is temporarily unavailable",
    quotaBody: (time: string) =>
      `So many people searched today that the quota ran out. It'll be back around ${time}. You can still search the TuneTrend charts meanwhile.`,
    userLimitTitle: "You've used all your YouTube searches today",
    userLimitBody: "You can search again tomorrow, or try our charts first.",
    unavailableTitle: "YouTube search isn't available right now",
    unavailableBody: "Please try again later. Chart search still works as usual.",
    youtubeErrorTitle: "YouTube search failed",
    youtubeErrorBody: "Something went wrong while contacting YouTube. Please try again.",
    retry: "Try again",
    youtubeEmpty: "No results on YouTube.",
  },
  auth: {
    defaultTitle: "Sign in to TuneTrend",
    defaultBody: "Create your own playlists and search more on YouTube. Works on both web and app.",
    playlistTitle: "Save songs for later",
    playlistBody: "Sign in with Google to create your own playlists. Works on both web and app.",
    youtubeTitle: "Sign in to keep searching on YouTube",
    youtubeBody: "Chart search still works as usual. YouTube search is for members only.",
    finePrint:
      "We only keep your name, email and profile picture. Nothing is ever posted to your YouTube.",
    privacyLink: "Privacy Policy",
    loadingGoogle: "Loading Google sign-in…",
    loginFailed: "Sign-in failed. Please try again.",
    close: "Close",
    cancel: "Cancel",
    deleteTitle: "Delete your account?",
    deleteBody:
      "Your account and all of your playlists will be permanently deleted. This can't be undone.",
    deleteConfirm: "Delete account",
    deleteFailed: "Couldn't delete your account. Please try again.",
    deleteSessionExpired: "Your session has expired. Please sign in again to delete your account.",
  },
  playlist: {
    addButton: "Save to playlist",
    popoverTitle: "Save to playlist",
    newPlaceholder: "New playlist name",
    create: "Create",
    loading: "Loading…",
    noPlaylists: "No playlists yet — create one below.",
    errorLimit: "You've reached the limit of 50 playlists.",
    errorItemLimit: "This playlist is full (200 songs).",
    errorName: "Name must be 1–60 characters.",
    errorGeneric: "Couldn't save. Please try again.",
    metaTitle: "My playlists · TuneTrend",
    heading: "My playlists",
    newPlaylist: "+ New playlist",
    songCount: (n: number) => `${n} ${n === 1 ? "song" : "songs"}`,
    privateBadge: "Private playlist",
    updated: (rel: string) => `Last edited ${rel}`,
    playAll: "▶ Play all",
    shuffle: "Shuffle",
    rename: "Rename",
    delete: "Delete",
    save: "Save",
    cancel: "Cancel",
    nameLabel: "Playlist name",
    deleteTitle: "Delete this playlist?",
    deleteBody: (name: string) => `"${name}" and all songs in it will be deleted.`,
    removeItem: "Remove from playlist",
    moveUp: "Move up",
    moveDown: "Move down",
    dragHandle: "Drag to reorder",
    unavailableTitle: "This video is no longer available",
    unavailableBody: "Removed or made private on YouTube",
    skipOnPlay: "Skipped when playing",
    emptyTitle: "No songs in this playlist yet",
    emptyBody: "Tap + on any song to add it here.",
    listEmptyTitle: "No playlists yet",
    listEmptyBody: "Create your first playlist, then tap + on songs you love.",
    loginTitle: "Sign in to see your playlists",
    loginBody: "Playlists are saved to your Google-linked TuneTrend account.",
    notFoundTitle: "Playlist not found",
    notFoundBody: "It may have been deleted.",
    backToList: "← My playlists",
    backToCharts: "← Back to charts",
    playingFrom: (name: string) => `Playing from ${name}`,
    actionFailed: "Something went wrong. Please try again.",
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
    lastUpdated: "Last updated: November 7, 2026",
    intro:
      'TuneTrend (the "app", "we", "us") is a small, independently-run project that shows ' +
      "trending YouTube music and videos by country, available as a website and as Android/iOS " +
      "apps. This page explains what data we collect, why, and how you're in control of it. " +
      "Browsing charts and searching them never requires a login. Signing in with Google is " +
      "optional and only needed to create playlists or to search further on YouTube.",
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
          "Search terms (no sign-in needed): when you search for a song, we save the search term, " +
            'the country and the latest results for that term to show on the public "Most Searched" ' +
            "tab. A term only appears there once at least 2 people have searched for it within 7 " +
            "days, so private searches stay private. To count distinct people, we keep a one-way " +
            "hash of your IP address (it can't be turned back into the IP) with each search for 7 " +
            "days, then delete it. Search terms are never linked to your account.",
          "Local preferences: your theme (light/dark), language, and whether you've seen the " +
            "onboarding tour or install prompt are saved in your browser's local storage or " +
            "your device's app storage. This never leaves your device.",
        ],
      },
      {
        heading: "Your account and playlists (only if you sign in)",
        body: [
          "If you choose to sign in with Google, we receive and store your Google account " +
            "name, email address and profile picture URL, plus an internal account ID and the " +
            "time you last signed in. We never get your Google password and never ask for " +
            "access to your YouTube account — nothing is ever posted to YouTube on your behalf.",
          "Playlists you create are stored on TuneTrend's own server: each playlist's name and, " +
            "for each saved song, its YouTube video ID, title, channel name and thumbnail link. " +
            "Playlists are private to your account and are not shared with anyone.",
          "On the web, your sign-in session is kept in a secure, HTTP-only cookie that expires " +
            "after 30 days. To limit costs, we also count how many YouTube searches each account " +
            "makes per day.",
          "When you play a song from the search results page, we count that play (video ID, " +
            "country and search term — not who you are) to build the Most Searched chart.",
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
          "Google Sign-In — optional sign-in; shares your name, email and profile picture with us.",
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
            "periods. Rate-limiting counters expire automatically within minutes. Your account " +
            "details and playlists are kept until you delete your account.",
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
            "site data or the app's storage.",
          "Deleting your account: sign in, open the profile menu (your picture at the top " +
            'right) and choose "Delete account". Your account and all of your playlists are ' +
            "permanently deleted right away. Signing out only ends the session on that device. " +
            "If you can't sign in any more, ask us through the Contact form and we'll delete it " +
            "for you, along with anything you typed into that form.",
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
