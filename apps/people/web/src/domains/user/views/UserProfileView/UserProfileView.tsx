import { Suspense } from "react";

import { UserProfile, UserProfileSkeleton } from "@domains/user/components";
import type { UserProfileViewProps } from "@domains/user/contracts";
import { placeholderUser } from "@domains/user/mocks/placeholderUser";

export function UserProfileView({ userId }: UserProfileViewProps) {
  return (
    <main data-slot="user-profile-view" className="mx-auto w-full max-w-2xl px-4 py-12 sm:py-20">
      <Suspense fallback={<UserProfileSkeleton />}>
        <UserProfileContent userId={userId} />
      </Suspense>
    </main>
  );
}

async function UserProfileContent({ userId }: UserProfileViewProps) {
  const user = placeholderUser(await userId);
  return <UserProfile user={user} />;
}
