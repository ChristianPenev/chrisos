import { Monogram } from "@domains/user/components/Monogram";
import type { UserListItemProps } from "@domains/user/contracts";
import { formatDate } from "@lib/formatDate";
import { hueStyle } from "@lib/hue";
import { cn } from "@lib/utils";
import { CopyButton } from "@ui/CopyButton";

export function UserListItem({ user }: UserListItemProps) {
  const fullName = `${user.firstName} ${user.lastName}`;

  return (
    <tr
      data-slot="user-list-item"
      style={hueStyle(user.id)}
      className={cn(
        "group border-b transition-colors duration-150 last:border-b-0",
        "hover:bg-[oklch(0.975_0.015_var(--hue))] dark:hover:bg-[oklch(0.21_0.02_var(--hue))]",
        "animate-in duration-300 fade-in-0 motion-reduce:animate-none",
      )}
    >
      <td className="py-3 ps-4 pe-3 sm:ps-5">
        <div className="flex min-w-0 items-center gap-3.5">
          <Monogram seed={user.id} firstName={user.firstName} lastName={user.lastName} />
          <div className="min-w-0">
            <p className="truncate text-sm font-medium">{fullName}</p>
            <a
              href={`mailto:${user.email}`}
              className={cn(
                "block truncate text-[0.8125rem] text-muted-foreground",
                "hover:text-foreground",
              )}
            >
              {user.email}
            </a>
          </div>
        </div>
      </td>
      <td className="px-3 py-3 text-[0.8125rem]">
        <span
          title={user.role.name}
          className={cn(
            "block truncate",
            user.role.isOwner ? "font-medium text-foreground" : "text-muted-foreground",
          )}
        >
          {user.role.name}
        </span>
      </td>
      <td
        className={cn(
          "hidden px-3 py-3 whitespace-nowrap sm:table-cell",
          "text-[0.8125rem] text-muted-foreground tabular-nums",
        )}
      >
        <time dateTime={user.createdAt}>{formatDate(user.createdAt)}</time>
      </td>
      <td className="w-px py-3 ps-1 pe-2 sm:pe-3">
        <CopyButton
          compact
          value={user.id}
          label={`Copy ${fullName}'s user ID`}
          className={cn(
            "[@media(hover:hover)]:opacity-0 [@media(hover:hover)]:group-hover:opacity-100",
            "focus-visible:opacity-100 data-copied:opacity-100",
          )}
        />
      </td>
    </tr>
  );
}
