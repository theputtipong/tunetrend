"use client";

import { useEffect } from "react";
import "driver.js/dist/driver.css";
import { dictionaries, type Lang } from "@/lib/i18n";

const STORAGE_KEY = "tunetrend-onboarding-seen";

function isMobileViewport() {
  return typeof window !== "undefined" && window.matchMedia("(max-width: 639px)").matches;
}

function buildSteps(t: (typeof dictionaries)[Lang]["onboarding"]) {
  const mobile = isMobileViewport();

  const desktopOnlySteps = mobile
    ? [
        {
          element: '[data-tour="mobile-menu"]',
          popover: { title: t.mobileMenuTitle, description: t.mobileMenuDescription },
        },
      ]
    : [
        {
          element: '[data-tour="language-toggle"]',
          popover: { title: t.languageTitle, description: t.languageDescription },
        },
        {
          element: '[data-tour="theme-toggle"]',
          popover: { title: t.themeTitle, description: t.themeDescription },
        },
        {
          element: '[data-tour="about-link"]',
          popover: { title: t.aboutTitle, description: t.aboutDescription },
        },
        {
          element: '[data-tour="buy-coffee"]',
          popover: { title: t.coffeeTitle, description: t.coffeeDescription },
        },
      ];

  return [
    {
      element: '[data-tour="country-selector"]',
      popover: { title: t.countryTitle, description: t.countryDescription },
    },
    {
      element: '[data-tour="tabs"]',
      popover: { title: t.tabsTitle, description: t.tabsDescription },
    },
    ...desktopOnlySteps,
  ];
}

let isTourRunning = false;

export function runOnboardingTour(lang: Lang) {
  if (isTourRunning) return;
  isTourRunning = true;

  const t = dictionaries[lang].onboarding;

  import("driver.js")
    .then(({ driver }) => {
      driver({
        showProgress: true,
        popoverClass: "tunetrend-tour",
        nextBtnText: t.next,
        prevBtnText: t.previous,
        doneBtnText: t.done,
        onDestroyed: () => {
          isTourRunning = false;
          localStorage.setItem(STORAGE_KEY, "true");
        },
        steps: buildSteps(t),
      }).drive();
    })
    .catch(() => {
      isTourRunning = false;
    });
}

export function OnboardingTour({ lang }: Readonly<{ lang: Lang }>) {
  useEffect(() => {
    if (localStorage.getItem(STORAGE_KEY) === "true") return;
    runOnboardingTour(lang);
  }, [lang]);

  return null;
}
