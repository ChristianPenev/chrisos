import { cn } from "@lib/utils";

export function PageHeading({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <h1
      data-slot="page-heading"
      className={cn(
        "text-[1.75rem] leading-tight font-semibold tracking-[-0.02em]",
        "font-stretch-semi-condensed text-balance",
        className,
      )}
    >
      {children}
    </h1>
  );
}
