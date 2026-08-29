import Link from "next/link";
import { dictionaries } from "@/lib/i18n";
import { getLang } from "@/lib/i18n/server";
import { LogoMark } from "@/components/icons";
import { LanguageToggle } from "@/components/LanguageToggle";
import { MobileMenu } from "@/components/MobileMenu";
import { ThemeToggle } from "@/components/ThemeToggle";

export async function generateMetadata() {
  const lang = await getLang();
  const t = dictionaries[lang].privacy;
  return { title: t.metaTitle, description: t.metaDescription };
}

export default async function PrivacyPage() {
  const lang = await getLang();
  const t = dictionaries[lang].privacy;

  return (
    <div className="mx-auto max-w-3xl px-4 py-4 md:px-8 md:py-6">
      <div className="flex items-center justify-between">
        <Link href="/" className="flex items-center gap-2">
          <LogoMark size={30} />
          <span className="font-display text-[21px] font-bold tracking-[-0.01em]">TuneTrend</span>
        </Link>
        <div className="flex items-center gap-2">
          <Link href="/" className="pill">
            {t.backToTrends}
          </Link>
          <div className="hidden items-center gap-2 sm:flex">
            <LanguageToggle />
            <ThemeToggle lang={lang} />
          </div>
          <MobileMenu lang={lang} showReplayTour={false} showPrivacy={false} />
        </div>
      </div>

      <div className="mt-12 flex flex-col gap-4">
        <span className="text-xs font-semibold uppercase tracking-wide text-[var(--accent)]">
          {t.eyebrow}
        </span>
        <h1 className="font-display text-3xl font-bold tracking-[-0.01em] md:text-4xl">
          {t.heading}
        </h1>
        <p className="text-sm text-[var(--text-3)]">{t.lastUpdated}</p>
        <p className="max-w-xl text-[15px] leading-relaxed text-[var(--text-2)]">{t.intro}</p>
      </div>

      <div className="mt-10 flex flex-col gap-8">
        {t.sections.map((section) => (
          <div key={section.heading}>
            <h2 className="font-display text-lg font-bold">{section.heading}</h2>
            <div className="mt-2 flex flex-col gap-2 text-[15px] leading-relaxed text-[var(--text-2)]">
              {section.body.map((paragraph, i) => (
                <p key={i}>{paragraph}</p>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
