import { cn } from "@lib/utils";

export function PageContainer({
  children,
  className,
  slot,
}: {
  children: React.ReactNode;
  className?: string;
  slot: string;
}) {
  return (
    <main
      data-slot={slot}
      className={cn("mx-auto w-full max-w-5xl px-4 pt-10 pb-20 sm:px-8 sm:pt-14", className)}
    >
      {children}
    </main>
  );
}
