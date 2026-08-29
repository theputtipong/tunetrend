import Link from "next/link";
import { dictionaries } from "@/lib/i18n";
import { getLang } from "@/lib/i18n/server";
import { LogoMark } from "@/components/icons";
import { BuyMeCoffeeButton } from "@/components/BuyMeCoffeeButton";
import { ContactFab } from "@/components/ContactFab";
import { LanguageToggle } from "@/components/LanguageToggle";
import { MobileMenu } from "@/components/MobileMenu";
import { ThemeToggle } from "@/components/ThemeToggle";

export async function generateMetadata() {
  const lang = await getLang();
  const t = dictionaries[lang].about;
  return { title: t.metaTitle, description: t.metaDescription };
}

export default async function AboutPage() {
  const lang = await getLang();
  const t = dictionaries[lang].about;

  return (
    <div className="mx-auto max-w-3xl px-4 py-4 md:px-8 md:py-6">
      <div className="flex items-center justify-between">
        <Link href="/" className="flex items-center gap-2">
          <LogoMark size={30} />
          <span className="font-display text-[21px] font-bold tracking-[-0.01em]">TuneTrend</span>
        </Link>
        <div className="flex items-center gap-2">
          <div className="hidden sm:inline-flex">
            <BuyMeCoffeeButton label={dictionaries[lang].nav.buyCoffee} size={32} />
          </div>
          <Link href="/" className="pill">
            {t.backToTrends}
          </Link>
          <div className="hidden items-center gap-2 sm:flex">
            <LanguageToggle />
            <ThemeToggle lang={lang} />
          </div>
          <MobileMenu lang={lang} showAbout={false} showReplayTour={false} />
        </div>
      </div>

      <div className="mt-12 flex flex-col gap-4">
        <span className="text-xs font-semibold uppercase tracking-wide text-[var(--accent)]">
          {t.eyebrow}
        </span>
        <h1 className="font-display text-3xl font-bold tracking-[-0.01em] md:text-4xl">
          {t.heading}
        </h1>
        <p className="max-w-xl text-[15px] leading-relaxed text-[var(--text-2)]">{t.lead}</p>
      </div>

      <div className="mt-10 grid grid-cols-2 gap-3 md:grid-cols-4">
        {t.stats.map((stat) => (
          <div key={stat.label} className="stat-card">
            <div className="stat-value">{stat.value}</div>
            <div className="stat-label">{stat.label}</div>
          </div>
        ))}
      </div>

      <div className="mt-10 flex flex-col gap-4 text-[15px] leading-relaxed text-[var(--text-2)]">
        <p>{t.bodyP1}</p>
        <p>{t.bodyP2}</p>
      </div>

      <div className="mt-10">
        <h2 className="font-display text-lg font-bold">{t.stackHeading}</h2>
        <p className="mt-1 text-sm text-[var(--text-3)]">{t.stackCaption}</p>
        <div className="mt-3 flex flex-wrap gap-2">
          {t.stack.map((item) => (
            <span key={item} className="badge">
              {item}
            </span>
          ))}
        </div>
      </div>

      <footer className="mt-14 flex flex-col items-start gap-3 border-t border-[var(--border)] py-6 pb-24">
        <div className="sm:hidden">
          <BuyMeCoffeeButton label={dictionaries[lang].nav.buyCoffee} size={26} showLabel />
        </div>
        <Link href="/privacy" className="text-sm text-[var(--text-3)] underline">
          {dictionaries[lang].nav.privacy}
        </Link>
      </footer>

      <ContactFab lang={lang} />
    </div>
  );
}
