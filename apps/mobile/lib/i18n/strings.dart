class AboutStat {
  final String value;
  final String label;
  const AboutStat(this.value, this.label);
}

class AppStrings {
  final String aboutTooltip;
  final String themeToggleTooltip;
  final String languageToggleTooltip;
  final String backToTrends;

  final String tabTrending;
  final String tabNew;
  final String tabNewGeneric;
  final String tabMv;
  final String musicCategory;

  final String autoPlayPromptTitle;
  final String autoPlayPromptDescription;
  final String autoPlayPromptAccept;
  final String autoPlayPromptDecline;
  final String Function(int secondsLeft) autoPlayPromptCountdown;
  final String Function(int secondsLeft) autoAdvanceCountdown;
  final String autoPlayStoppedMessage;

  final String viewCountNoticeTooltip;
  final String viewCountNoticeTitle;
  final String viewCountNoticeBody;
  final String viewCountNoticeDismiss;

  final String maintenanceTitle;
  final String maintenanceDescription;

  final String updateRequiredTitle;
  final String updateRequiredDescription;
  final String updateNowButton;

  final String shareTooltip;
  final String Function(String title, String url, String? appUrl) shareMessage;

  final String errorTitle;
  final String errorDescription;
  final String retry;

  final String emptyTitle;
  final String Function(String countryName) emptyDescription;

  final String viewsSuffix;
  final String Function(int n) minutesAgo;
  final String Function(int n) hoursAgo;
  final String Function(int n) daysAgo;
  final String Function(int n) weeksAgo;

  final String aboutEyebrow;
  final String aboutHeading;
  final String aboutLead;
  final List<AboutStat> aboutStats;
  final String aboutBodyP1;
  final String aboutBodyP2;
  final String aboutStackHeading;
  final String aboutStackCaption;
  final List<String> aboutStack;

  final String onboardingCountryTitle;
  final String onboardingCountryDescription;
  final String onboardingTabsTitle;
  final String onboardingTabsDescription;
  final String onboardingCategoryTitle;
  final String onboardingCategoryDescription;
  final String onboardingLanguageTitle;
  final String onboardingLanguageDescription;
  final String onboardingThemeTitle;
  final String onboardingThemeDescription;
  final String onboardingAboutTitle;
  final String onboardingAboutDescription;
  final String onboardingMenuTitle;
  final String onboardingMenuDescription;

  final String menuTooltip;

  final String contactOpenButton;
  final String contactNameLabel;
  final String contactNamePlaceholder;
  final String contactMessageLabel;
  final String contactMessagePlaceholder;
  final String contactMethodLabel;
  final String contactMethodEmail;
  final String contactMethodPhone;
  final String contactEmailPlaceholder;
  final String contactPhonePlaceholder;
  final String contactSubmit;
  final String contactSubmitting;
  final String contactSuccessMessage;
  final String Function(int n) closingIn;
  final String contactErrorMessageRequired;
  final String contactErrorInvalidEmail;
  final String contactErrorInvalidPhone;
  final String contactErrorTooLong;
  final String contactErrorRateLimited;
  final String contactErrorGeneric;

  final String supportDevelopment;
  final String replayTourTooltip;
  final String privacyPolicy;

  const AppStrings({
    required this.aboutTooltip,
    required this.themeToggleTooltip,
    required this.languageToggleTooltip,
    required this.backToTrends,
    required this.tabTrending,
    required this.tabNew,
    required this.tabNewGeneric,
    required this.tabMv,
    required this.musicCategory,
    required this.autoPlayPromptTitle,
    required this.autoPlayPromptDescription,
    required this.autoPlayPromptAccept,
    required this.autoPlayPromptDecline,
    required this.autoPlayPromptCountdown,
    required this.autoAdvanceCountdown,
    required this.autoPlayStoppedMessage,
    required this.viewCountNoticeTooltip,
    required this.viewCountNoticeTitle,
    required this.viewCountNoticeBody,
    required this.viewCountNoticeDismiss,
    required this.maintenanceTitle,
    required this.maintenanceDescription,
    required this.updateRequiredTitle,
    required this.updateRequiredDescription,
    required this.updateNowButton,
    required this.shareTooltip,
    required this.shareMessage,
    required this.errorTitle,
    required this.errorDescription,
    required this.retry,
    required this.emptyTitle,
    required this.emptyDescription,
    required this.viewsSuffix,
    required this.minutesAgo,
    required this.hoursAgo,
    required this.daysAgo,
    required this.weeksAgo,
    required this.aboutEyebrow,
    required this.aboutHeading,
    required this.aboutLead,
    required this.aboutStats,
    required this.aboutBodyP1,
    required this.aboutBodyP2,
    required this.aboutStackHeading,
    required this.aboutStackCaption,
    required this.aboutStack,
    required this.onboardingCountryTitle,
    required this.onboardingCountryDescription,
    required this.onboardingTabsTitle,
    required this.onboardingTabsDescription,
    required this.onboardingCategoryTitle,
    required this.onboardingCategoryDescription,
    required this.onboardingLanguageTitle,
    required this.onboardingLanguageDescription,
    required this.onboardingThemeTitle,
    required this.onboardingThemeDescription,
    required this.onboardingAboutTitle,
    required this.onboardingAboutDescription,
    required this.onboardingMenuTitle,
    required this.onboardingMenuDescription,
    required this.menuTooltip,
    required this.contactOpenButton,
    required this.contactNameLabel,
    required this.contactNamePlaceholder,
    required this.contactMessageLabel,
    required this.contactMessagePlaceholder,
    required this.contactMethodLabel,
    required this.contactMethodEmail,
    required this.contactMethodPhone,
    required this.contactEmailPlaceholder,
    required this.contactPhonePlaceholder,
    required this.contactSubmit,
    required this.contactSubmitting,
    required this.contactSuccessMessage,
    required this.closingIn,
    required this.contactErrorMessageRequired,
    required this.contactErrorInvalidEmail,
    required this.contactErrorInvalidPhone,
    required this.contactErrorTooLong,
    required this.contactErrorRateLimited,
    required this.contactErrorGeneric,
    required this.supportDevelopment,
    required this.replayTourTooltip,
    required this.privacyPolicy,
  });
}
