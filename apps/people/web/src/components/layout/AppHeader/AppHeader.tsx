export function AppHeader() {
  return (
    <header
      data-slot="app-header"
      className="sticky top-0 z-10 border-b bg-background/85 backdrop-blur-md"
    >
      <div className="mx-auto flex h-12 w-full max-w-5xl items-center gap-2.5 px-4 sm:px-8">
        <svg viewBox="0 0 20 20" aria-hidden className="size-5 text-foreground">
          <circle cx="7.5" cy="10" r="5.5" fill="currentColor" />
          <circle cx="12.5" cy="10" r="5.5" fill="none" stroke="currentColor" strokeWidth="1.5" />
        </svg>
        <span className="text-[0.9375rem] font-semibold tracking-tight">People</span>
      </div>
    </header>
  );
}
