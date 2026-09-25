import type { MetadataRoute } from "next";
import { ANDROID_PACKAGE_NAME, PLAY_STORE_URL } from "@/lib/installPrompt";

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: "TuneTrend",
    short_name: "TuneTrend",
    description: "Trending YouTube Music videos by country",
    start_url: "/",
    display: "standalone",
    background_color: "#fafafa",
    theme_color: "#ff7a47",
    icons: [
      { src: "/icon.svg", sizes: "any", type: "image/svg+xml" },
      { src: "/apple-icon.png", sizes: "180x180", type: "image/png" },
    ],
    // The Android app is live, so point browsers at it instead of offering a PWA install.
    related_applications: [{ platform: "play", id: ANDROID_PACKAGE_NAME, url: PLAY_STORE_URL }],
    prefer_related_applications: true,
  };
}
