import { OrganizationUsersView } from "@domains/user/views";

export default function OrganizationUsersPage({
  params,
}: PageProps<"/organizations/[organizationId]/users">) {
  return (
    <OrganizationUsersView organizationId={params.then(({ organizationId }) => organizationId)} />
  );
}
