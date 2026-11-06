import { cn } from "@/lib/utils";

function hueFrom(seed: string) {
  let hash = 0;
  for (const char of seed) {
    hash = (hash * 31 + char.charCodeAt(0)) >>> 0;
  }
  return hash % 360;
}

export function Monogram({
  seed,
  firstName,
  lastName,
  className,
}: {
  seed: string;
  firstName: string;
  lastName: string;
  className?: string;
}) {
  const initials = `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase();

  return (
    <div
      aria-hidden
      style={{ "--hue": hueFrom(seed) } as React.CSSProperties}
      className={cn(
        "grid size-16 shrink-0 place-items-center rounded-2xl font-heading text-2xl font-semibold tracking-tight",
        "bg-[oklch(0.93_0.05_var(--hue))] text-[oklch(0.4_0.1_var(--hue))]",
        "dark:bg-[oklch(0.32_0.06_var(--hue))] dark:text-[oklch(0.9_0.07_var(--hue))]",
        className,
      )}
    >
      {initials}
    </div>
  );
}
