import { PageContainer } from "@components/layout/PageContainer";
import { PageHeading } from "@components/layout/PageHeading";

export function UserNotFoundView() {
  return (
    <PageContainer slot="user-not-found-view">
      <PageHeading>No one has this user ID</PageHeading>
      <p className="mt-2 max-w-prose text-muted-foreground">
        The link may be incomplete, or the user was removed. Check the ID in the address bar.
      </p>
    </PageContainer>
  );
}
