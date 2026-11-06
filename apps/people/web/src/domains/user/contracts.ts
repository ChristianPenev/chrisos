export type UserProfileData = {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  createdAt: string;
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

export type UserProfileViewProps = {
  userId: Promise<string>;
};
