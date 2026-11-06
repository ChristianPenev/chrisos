import { Suspense } from "react";

import { UserProfile, type UserProfileData } from "@/components/users/user-profile";
import { UserProfileSkeleton } from "@/components/users/user-profile-skeleton";

function placeholderUser(id: string): UserProfileData {
  return {
    id,
    email: "jane.doe@example.com",
    firstName: "Jane",
    lastName: "Doe",
    createdAt: "2026-11-06T13:58:00.250Z",
  };
}

export default function UserPage({ params }: { params: Promise<{ id: string }> }) {
  return (
    <main className="mx-auto w-full max-w-2xl px-4 py-12 sm:py-20">
      <Suspense fallback={<UserProfileSkeleton />}>
        {params.then(({ id }) => (
          <UserProfile user={placeholderUser(id)} />
        ))}
      </Suspense>
    </main>
  );
}
