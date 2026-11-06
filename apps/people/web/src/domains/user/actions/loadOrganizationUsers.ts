"use server";

import type { UsersPage } from "@domains/user/contracts";
import { getOrganizationUsers } from "@domains/user/data/getOrganizationUsers";

export async function loadOrganizationUsers(
  organizationId: string,
  after: string,
): Promise<UsersPage> {
  if (typeof organizationId !== "string" || typeof after !== "string") {
    throw new Error("Invalid arguments");
  }
  const page = await getOrganizationUsers(organizationId, after);
  if (!page) {
    throw new Error("Invalid organization ID");
  }
  return page;
}
