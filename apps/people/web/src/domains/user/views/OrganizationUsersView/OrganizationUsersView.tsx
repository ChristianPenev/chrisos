import { notFound } from "next/navigation";
import { Suspense } from "react";

import { PageContainer } from "@components/layout/PageContainer";
import { PageHeading } from "@components/layout/PageHeading";
import { PaginatedUserList, UserListSkeleton } from "@domains/user/components";
import type { OrganizationUsersViewProps } from "@domains/user/contracts";
import { getOrganizationUsers } from "@domains/user/data/getOrganizationUsers";
import { CopyButton } from "@ui/CopyButton";

export function OrganizationUsersView({ organizationId }: OrganizationUsersViewProps) {
  return (
    <PageContainer slot="organization-users-view">
      <PageHeading>Members</PageHeading>
      <Suspense fallback={<UserListSkeleton />}>
        <OrganizationUsersContent organizationId={organizationId} />
      </Suspense>
    </PageContainer>
  );
}

async function OrganizationUsersContent({ organizationId }: OrganizationUsersViewProps) {
  const id = await organizationId;
  const page = await getOrganizationUsers(id);
  if (!page) {
    notFound();
  }

  return (
    <>
      <div className="mt-1 mb-5 flex h-7 min-w-0 items-center gap-1 text-[0.8125rem]">
        <span className="shrink-0 text-muted-foreground">Organization</span>
        <span className="truncate ps-1 font-mono text-xs text-faint">{id}</span>
        <CopyButton compact value={id} label="Copy organization ID" />
      </div>
      {page.users.length > 0 ? (
        <PaginatedUserList organizationId={id} initialPage={page} />
      ) : (
        <div className="rounded-xl border border-dashed px-6 py-16 text-center">
          <p className="text-sm font-medium">No one has joined this organization yet</p>
          <p className="mt-1 text-[0.8125rem] text-muted-foreground">
            People added to this organization will show up here.
          </p>
        </div>
      )}
    </>
  );
}
