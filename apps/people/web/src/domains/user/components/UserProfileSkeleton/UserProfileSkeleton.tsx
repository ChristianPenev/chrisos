import { cn } from "@lib/utils";
import { Skeleton } from "@shadcn/skeleton";

export function UserProfileSkeleton() {
  return (
    <div
      data-slot="user-profile-skeleton"
      className="max-w-2xl"
      aria-busy
      aria-label="Loading profile"
    >
      <Skeleton className="h-36 rounded-2xl sm:h-44" />
      <Skeleton
        className={cn("relative -mt-9 ms-5 size-[4.5rem]", "rounded-[30%] ring-4 ring-background")}
      />
      <Skeleton className="mt-4 h-9 w-72 max-w-full sm:h-11" />
      <Skeleton className="mt-2.5 h-4 w-52" />
      <div className="mt-8 divide-y rounded-xl border">
        {["w-28", "w-64"].map((width) => (
          <div key={width} className="flex min-h-12 items-center gap-4 ps-4">
            <Skeleton className="h-3.5 w-16" />
            <Skeleton className={`h-3.5 ${width}`} />
          </div>
        ))}
      </div>
    </div>
  );
}
