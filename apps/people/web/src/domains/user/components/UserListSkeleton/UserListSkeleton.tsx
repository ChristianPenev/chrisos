import { cn } from "@lib/utils";
import { Skeleton } from "@shadcn/skeleton";

const rows = [
  ["w-32", "w-48"],
  ["w-28", "w-56"],
  ["w-36", "w-44"],
  ["w-24", "w-52"],
  ["w-32", "w-40"],
];

export function UserListSkeleton() {
  return (
    <div data-slot="user-list-skeleton" aria-busy aria-label="Loading members">
      <div className="mt-1 mb-5 flex h-7 items-center">
        <Skeleton className="h-3.5 w-72 max-w-full" />
      </div>
      <div className="overflow-hidden rounded-xl border">
        <div className="h-10 border-b bg-sunken" />
        {rows.map(([name, email], index) => (
          <div
            key={index}
            className={cn(
              "flex h-16 items-center gap-3.5",
              "border-b ps-4 pe-5 last:border-b-0 sm:ps-5",
            )}
          >
            <Skeleton className="size-9 rounded-[30%]" />
            <div className="flex flex-1 flex-col gap-1.5">
              <Skeleton className={`h-3.5 ${name}`} />
              <Skeleton className={`h-3 ${email}`} />
            </div>
            <Skeleton className="h-3 w-16 sm:me-20 sm:w-24" />
            <Skeleton className="hidden h-3 w-24 sm:block" />
          </div>
        ))}
      </div>
    </div>
  );
}
