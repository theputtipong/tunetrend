const BMAC_URL = "https://buymeacoffee.com/theputtipong";

const BMAC_GIF_URL =
  "https://media2.giphy.com/media/v1.Y2lkPTc5MGI3NjExOHJ1aHdnNjlrcXowa3h5YWp2YzljbHpkNzN1ZHQ0Ym5rdWNkcmsyNSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9cw/TDQOtnWgsBx99cNoyH/giphy.gif";

export function BuyMeCoffeeButton({
  label,
  size = 40,
  dataTour,
  showLabel = false,
}: Readonly<{
  label: string;
  size?: number;
  dataTour?: string;
  showLabel?: boolean;
}>) {
  return (
    <a
      href={BMAC_URL}
      target="_blank"
      rel="noopener noreferrer"
      aria-label={label}
      title={label}
      data-tour={dataTour}
      className="inline-flex items-center gap-2 no-underline"
    >
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={BMAC_GIF_URL}
        alt={label}
        width={size}
        height={size}
        className="rounded-full"
      />
      {showLabel && (
        <span className="text-[13.5px] font-semibold text-[var(--text)] sm:hidden">{label}</span>
      )}
    </a>
  );
}
