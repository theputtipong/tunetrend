"use client";

import { useEffect } from "react";
import { StateMessage } from "@/components/StateMessage";
import { dictionaries } from "@/lib/i18n";
import { useLang } from "@/lib/i18n/useLang";

export default function CountryTrendsError({
  error,
  reset,
}: Readonly<{ error: Error & { digest?: string }; reset: () => void }>) {
  const lang = useLang();
  const t = dictionaries[lang].error;

  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <StateMessage
      variant="error"
      title={t.title}
      description={t.description}
      retryLabel={t.retry}
      onRetry={reset}
    />
  );
}
