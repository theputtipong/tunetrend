import { AlertIcon, EmptyBoxIcon, RetryIcon } from "./icons";

export function StateMessage({
  variant,
  title,
  description,
  onRetry,
  retryLabel = "Retry",
}: Readonly<{
  variant: "error" | "empty";
  title: string;
  description: string;
  onRetry?: () => void;
  retryLabel?: string;
}>) {
  return (
    <div className="flex flex-col items-center gap-4 px-8 py-16 text-center">
      <div className={variant === "error" ? "state-icon state-icon--error" : "state-icon state-icon--empty"}>
        {variant === "error" ? <AlertIcon /> : <EmptyBoxIcon />}
      </div>
      <h2 className="state-title">{title}</h2>
      <p className="state-subtitle">{description}</p>
      {onRetry ? (
        <button type="button" onClick={onRetry} className="retry">
          <RetryIcon />
          {retryLabel}
        </button>
      ) : null}
    </div>
  );
}
