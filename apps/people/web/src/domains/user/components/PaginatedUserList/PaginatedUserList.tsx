"use client";

import { ChevronLeft, ChevronRight } from "lucide-react";
import { useEffect, useRef } from "react";

import { UserList } from "@domains/user/components/UserList";
import type { PaginatedUserListProps } from "@domains/user/contracts";
import { useUserPages } from "@domains/user/hooks/useUserPages";
import { cn } from "@lib/utils";
import { Button } from "@shadcn/button";

export function PaginatedUserList({ organizationId, initialPage }: PaginatedUserListProps) {
  const pager = useUserPages(organizationId, initialPage);
  const panel = useRef<HTMLDivElement>(null);
  const shownPage = useRef(pager.pageNumber);

  useEffect(() => {
    if (shownPage.current === pager.pageNumber) {
      return;
    }
    shownPage.current = pager.pageNumber;
    const element = panel.current;
    if (element && element.getBoundingClientRect().top < 0) {
      const reduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
      element.scrollIntoView({ block: "start", behavior: reduced ? "auto" : "smooth" });
    }
  }, [pager.pageNumber]);

  const singlePage = !pager.hasPrevious && !pager.hasNext;
  const noun = pager.total === 1 ? "member" : "members";
  const summary = singlePage
    ? `${pager.total} ${noun}`
    : `Showing ${pager.rangeStart}–${pager.rangeEnd}${pager.total ? ` of ${pager.total}` : ""}`;

  return (
    <div data-slot="paginated-user-list" className="flex flex-col">
      <div ref={panel} className="scroll-mt-16 overflow-hidden rounded-xl border">
        <div
          aria-busy={pager.isPending}
          className={cn(
            "transition-opacity duration-150",
            pager.isPending && "opacity-50 motion-reduce:opacity-70",
          )}
        >
          <UserList users={pager.users} />
        </div>
        <div
          className={cn(
            "flex min-h-12 items-center justify-between gap-4",
            "border-t bg-sunken py-2 ps-4 pe-2 sm:ps-5",
          )}
        >
          <p className="text-[0.8125rem] text-muted-foreground tabular-nums" aria-live="polite">
            {summary}
          </p>
          {singlePage ? null : (
            <nav aria-label="Pagination" className="flex items-center gap-1">
              <Button
                variant="ghost"
                size="icon-sm"
                onClick={pager.goPrevious}
                disabled={!pager.hasPrevious || pager.isPending}
                aria-label="Previous page"
              >
                <ChevronLeft />
              </Button>
              <span className="min-w-20 text-center text-[0.8125rem] tabular-nums">
                Page {pager.pageNumber}
                {pager.pageCount ? (
                  <span className="text-muted-foreground"> of {pager.pageCount}</span>
                ) : null}
              </span>
              <Button
                variant="ghost"
                size="icon-sm"
                onClick={pager.goNext}
                disabled={!pager.hasNext || pager.isPending}
                aria-label="Next page"
              >
                <ChevronRight />
              </Button>
            </nav>
          )}
        </div>
      </div>
      {pager.failed ? (
        <p role="alert" className="mt-3 text-[0.8125rem] text-destructive">
          Couldn&apos;t load the next page. Check that the API is reachable, then try again.
        </p>
      ) : null}
    </div>
  );
}
