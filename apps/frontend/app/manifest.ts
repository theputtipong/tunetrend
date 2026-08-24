import type { MetadataRoute } from "next";

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
  };
}
