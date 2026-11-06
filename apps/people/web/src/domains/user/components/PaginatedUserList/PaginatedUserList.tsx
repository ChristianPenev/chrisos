"use client";

import { useState, useTransition } from "react";

import { loadOrganizationUsers } from "@domains/user/actions/loadOrganizationUsers";
import { UserList } from "@domains/user/components/UserList";
import type { PaginatedUserListProps } from "@domains/user/contracts";
import { cn } from "@lib/utils";
import { Button } from "@shadcn/button";

function countLabel(count: number, complete: boolean) {
  const noun = count === 1 ? "member" : "members";
  return complete ? `${count} ${noun}` : `Showing ${count} ${noun}`;
}

export function PaginatedUserList({ organizationId, initialPage }: PaginatedUserListProps) {
  const [users, setUsers] = useState(initialPage.users);
  const [nextCursor, setNextCursor] = useState(initialPage.nextCursor);
  const [failed, setFailed] = useState(false);
  const [isPending, startTransition] = useTransition();

  function loadMore() {
    if (!nextCursor) {
      return;
    }
    setFailed(false);
    startTransition(async () => {
      try {
        const page = await loadOrganizationUsers(organizationId, nextCursor);
        setUsers((current) => [...current, ...page.users]);
        setNextCursor(page.nextCursor);
      } catch {
        setFailed(true);
      }
    });
  }

  return (
    <div data-slot="paginated-user-list" className="flex flex-col">
      <div className="overflow-hidden rounded-xl border">
        <UserList users={users} />
        {nextCursor ? (
          <div
            className={cn(
              "flex items-center justify-between gap-4",
              "border-t bg-sunken px-4 py-2 sm:ps-5",
            )}
          >
            <p className="text-[0.8125rem] text-muted-foreground tabular-nums" aria-live="polite">
              {countLabel(users.length, false)}
            </p>
            <Button variant="outline" size="sm" onClick={loadMore} disabled={isPending}>
              {isPending ? "Loading…" : "Load more"}
            </Button>
          </div>
        ) : null}
      </div>
      {nextCursor ? null : (
        <p className="mt-3 text-[0.8125rem] text-faint tabular-nums">
          {countLabel(users.length, true)}
        </p>
      )}
      {failed ? (
        <p role="alert" className="mt-3 text-[0.8125rem] text-destructive">
          Couldn&apos;t load more members. Check that the API is reachable, then try again.
        </p>
      ) : null}
    </div>
  );
}
