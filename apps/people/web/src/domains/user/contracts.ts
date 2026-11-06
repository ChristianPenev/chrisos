import type { UserListItemFragment } from "@domains/user/graphql/organizationUsers.generated";

export type UserProfileData = {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  createdAt: string;
};

export type UserListItemData = UserListItemFragment;

export type UsersPage = {
  users: UserListItemData[];
  nextCursor: string | null;
};

export type MonogramProps = {
  seed: string;
  firstName: string;
  lastName: string;
  className?: string;
};

export type UserProfileProps = {
  user: UserProfileData;
};

export type UserListItemProps = {
  user: UserListItemData;
};

export type UserListProps = {
  users: UserListItemData[];
};

export type PaginatedUserListProps = {
  organizationId: string;
  initialPage: UsersPage;
};

export type UserProfileViewProps = {
  userId: Promise<string>;
};

export type OrganizationUsersViewProps = {
  organizationId: Promise<string>;
};
