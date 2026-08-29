import type { Metadata, Viewport } from "next";
import { Space_Grotesk, Plus_Jakarta_Sans } from "next/font/google";
import { Analytics } from "@vercel/analytics/next";
import { SpeedInsights } from "@vercel/speed-insights/next";
import { getLang } from "@/lib/i18n/server";
import { InstallPrompt } from "@/components/InstallPrompt";
import { ServiceWorkerRegister } from "@/components/ServiceWorkerRegister";
import "./globals.css";

// display: "optional" instead of the next/font default ("swap") — this site's
// Thai copy mixes in a lot of English/technical terms (e.g. "full-stack",
// "microservices"), and those Latin words are the only part of the text that
// actually uses this self-hosted font. Swapping it in after first paint can
// change those words' width just enough to shift a paragraph's line wrap,
// pushing every section below it down/up by a full line (a real, measured
// CLS regression). "optional" means the browser commits to the fallback font
// for the whole page life if the webfont isn't ready almost immediately,
// instead of swapping mid-render.
const spaceGrotesk = Space_Grotesk({
  variable: "--font-space-grotesk",
  subsets: ["latin"],
  weight: ["500", "700"],
  display: "optional",
});

const plusJakarta = Plus_Jakarta_Sans({
  variable: "--font-plus-jakarta",
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
  display: "optional",
});

export const metadata: Metadata = {
  title: "TuneTrend",
  description: "Trending YouTube Music videos by country",
};

export const viewport: Viewport = {
  themeColor: [
    { media: "(prefers-color-scheme: light)", color: "#fafafa" },
    { media: "(prefers-color-scheme: dark)", color: "oklch(15% 0.014 260)" },
  ],
};

const THEME_INIT_SCRIPT = `
try {
  var t = localStorage.getItem("tunetrend-theme");
  if (t === "light" || t === "dark") {
    document.documentElement.setAttribute("data-theme", t);
  }
} catch (e) {}
`;

export default async function RootLayout({ children }: LayoutProps<"/">) {
  const lang = await getLang();
  return (
    <html
      lang={lang}
      className={`${spaceGrotesk.variable} ${plusJakarta.variable} h-full antialiased`}
      suppressHydrationWarning
    >
      <head>
        <script dangerouslySetInnerHTML={{ __html: THEME_INIT_SCRIPT }} />
      </head>
      <body className="min-h-full">
        {children}
        <Analytics />
        <SpeedInsights />
        <ServiceWorkerRegister />
        <InstallPrompt lang={lang} />
      </body>
    </html>
  );
}
