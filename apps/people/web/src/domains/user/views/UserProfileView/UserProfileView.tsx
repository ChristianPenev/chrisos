import { Suspense } from "react";

import { PageContainer } from "@components/layout/PageContainer";
import { UserProfile, UserProfileSkeleton } from "@domains/user/components";
import type { UserProfileViewProps } from "@domains/user/contracts";
import { placeholderUser } from "@domains/user/mocks/placeholderUser";

export function UserProfileView({ userId }: UserProfileViewProps) {
  return (
    <PageContainer slot="user-profile-view" className="sm:pt-10">
      <Suspense fallback={<UserProfileSkeleton />}>
        <UserProfileContent userId={userId} />
      </Suspense>
    </PageContainer>
  );
}

async function UserProfileContent({ userId }: UserProfileViewProps) {
  const user = placeholderUser(await userId);
  return <UserProfile user={user} />;
}
