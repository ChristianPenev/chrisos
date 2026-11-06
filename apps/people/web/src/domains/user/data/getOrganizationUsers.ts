import type { UsersPage } from "@domains/user/contracts";
import { OrganizationUsersDocument } from "@domains/user/graphql/organizationUsers.generated";
import { ApiErrorCode, hasApiErrorCode } from "@lib/apollo/errors";
import { getClient } from "@lib/apollo/serverClient";

export async function getOrganizationUsers(
  organizationId: string,
  after?: string,
): Promise<UsersPage | null> {
  try {
    const { data } = await getClient().query({
      query: OrganizationUsersDocument,
      variables: { organizationId, after },
    });
    if (!data) {
      throw new Error("OrganizationUsers returned no data");
    }
    const { nodes, pageInfo } = data.users;
    return { users: nodes, nextCursor: pageInfo.hasNextPage ? (pageInfo.endCursor ?? null) : null };
  } catch (error) {
    if (hasApiErrorCode(error, ApiErrorCode.BadUserInput)) {
      return null;
    }
    throw error;
  }
}
