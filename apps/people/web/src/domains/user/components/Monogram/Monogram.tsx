import type { MonogramProps } from "@domains/user/contracts";
import { hueStyle } from "@lib/hue";
import { cn } from "@lib/utils";

export function Monogram({ seed, firstName, lastName, className }: MonogramProps) {
  const initials = `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase();

  return (
    <div
      aria-hidden
      data-slot="monogram"
      style={hueStyle(seed)}
      className={cn(
        "grid size-9 shrink-0 place-items-center rounded-[30%] select-none",
        "text-[0.8125rem] font-semibold tracking-tight",
        "bg-[oklch(0.94_0.045_var(--hue))] text-[oklch(0.42_0.11_var(--hue))]",
        "shadow-[inset_0_0_0_1px_oklch(0.4_0.1_var(--hue)/0.08)]",
        "dark:bg-[oklch(0.3_0.055_var(--hue))] dark:text-[oklch(0.88_0.08_var(--hue))]",
        "dark:shadow-[inset_0_0_0_1px_oklch(0.9_0.08_var(--hue)/0.08)]",
        className,
      )}
    >
      {initials}
    </div>
  );
}
