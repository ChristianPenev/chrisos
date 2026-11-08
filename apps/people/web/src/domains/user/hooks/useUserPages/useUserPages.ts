import { useCallback, useEffect, useRef, useState, useTransition } from "react";

import { loadOrganizationUsers } from "@domains/user/actions/loadOrganizationUsers";
import type { UsersPage } from "@domains/user/contracts";

export function useUserPages(organizationId: string, initialPage: UsersPage) {
  const [pages, setPages] = useState<UsersPage[]>([initialPage]);
  const [index, setIndex] = useState(0);
  const [failed, setFailed] = useState(false);
  const [isPending, startTransition] = useTransition();
  const requests = useRef(new Map<string, Promise<UsersPage>>());

  const request = useCallback(
    (cursor: string) => {
      let pending = requests.current.get(cursor);
      if (!pending) {
        pending = loadOrganizationUsers(organizationId, cursor);
        pending.catch(() => requests.current.delete(cursor));
        requests.current.set(cursor, pending);
      }
      return pending;
    },
    [organizationId],
  );

  const current = pages[index];
  const nextCursor = current.nextCursor;
  const nextLoaded = index + 1 < pages.length;

  useEffect(() => {
    if (nextCursor && !nextLoaded) {
      request(nextCursor).catch(() => {});
    }
  }, [nextCursor, nextLoaded, request]);

  function goNext() {
    if (!nextCursor) {
      return;
    }
    setFailed(false);
    if (nextLoaded) {
      setIndex(index + 1);
      return;
    }
    startTransition(async () => {
      try {
        const page = await request(nextCursor);
        startTransition(() => {
          setPages((loaded) => [...loaded.slice(0, index + 1), page]);
          setIndex(index + 1);
        });
      } catch {
        setFailed(true);
      }
    });
  }

  function goPrevious() {
    setFailed(false);
    setIndex((value) => Math.max(0, value - 1));
  }

  const lastPage = pages[pages.length - 1];
  const totalKnown = lastPage.nextCursor === null;
  const start = pages.slice(0, index).reduce((sum, page) => sum + page.users.length, 0) + 1;

  return {
    users: current.users,
    pageNumber: index + 1,
    pageCount: totalKnown ? pages.length : null,
    rangeStart: start,
    rangeEnd: start + current.users.length - 1,
    total: totalKnown ? pages.reduce((sum, page) => sum + page.users.length, 0) : null,
    hasPrevious: index > 0,
    hasNext: nextCursor !== null,
    goNext,
    goPrevious,
    isPending,
    failed,
  };
}
