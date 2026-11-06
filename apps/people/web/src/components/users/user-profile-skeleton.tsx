import { Skeleton } from "@/components/ui/skeleton";

export function UserProfileSkeleton() {
  return (
    <div className="flex flex-col gap-10" aria-busy aria-label="Loading user">
      <div className="flex items-center gap-5">
        <Skeleton className="size-16 rounded-2xl" />
        <div className="flex flex-col gap-2">
          <Skeleton className="h-10 w-64" />
          <Skeleton className="h-4 w-48" />
        </div>
      </div>
      <div className="divide-y border-y">
        {[0, 1].map((row) => (
          <div key={row} className="flex min-h-12 items-center gap-4 py-2">
            <Skeleton className="h-4 w-16" />
            <Skeleton className="h-4 w-56" />
          </div>
        ))}
      </div>
    </div>
  );
}
