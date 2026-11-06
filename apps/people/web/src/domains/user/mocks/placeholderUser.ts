import type { UserProfileData } from "@domains/user/contracts";

export function placeholderUser(id: string): UserProfileData {
  return {
    id,
    email: "jane.doe@example.com",
    firstName: "Jane",
    lastName: "Doe",
    createdAt: "2026-11-06T13:58:00.250Z",
  };
}
