export function UserNotFoundView() {
  return (
    <main data-slot="user-not-found-view" className="mx-auto w-full max-w-2xl px-4 py-12 sm:py-20">
      <h1 className="font-heading text-3xl font-semibold tracking-tight">No user with this ID</h1>
      <p className="mt-2 max-w-prose text-muted-foreground">
        The link may be incomplete, or the user was removed. Check the ID in the address bar and try
        again.
      </p>
    </main>
  );
}
