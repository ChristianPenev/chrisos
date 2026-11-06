import { PageContainer } from "@components/layout/PageContainer";
import { PageHeading } from "@components/layout/PageHeading";

export function OrganizationNotFoundView() {
  return (
    <PageContainer slot="organization-not-found-view">
      <PageHeading>That isn&apos;t a valid organization ID</PageHeading>
      <p className="mt-2 max-w-prose text-muted-foreground">
        Organization IDs are UUIDs, like 01a1b32a-91d7-723d-b807-e63b9cbc8331. Check the ID in the
        address bar.
      </p>
    </PageContainer>
  );
}
