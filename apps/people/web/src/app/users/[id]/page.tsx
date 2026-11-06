import { UserProfileView } from "@domains/user/views";

export default function UserPage({ params }: PageProps<"/users/[id]">) {
  return <UserProfileView userId={params.then(({ id }) => id)} />;
}
