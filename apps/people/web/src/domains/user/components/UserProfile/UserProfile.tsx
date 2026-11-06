import { Monogram } from "@domains/user/components/Monogram";
import type { UserProfileProps } from "@domains/user/contracts";
import { formatDate } from "@lib/formatDate";
import { hueStyle } from "@lib/hue";
import { cn } from "@lib/utils";
import { CopyButton } from "@ui/CopyButton";

export function UserProfile({ user }: UserProfileProps) {
  const fullName = `${user.firstName} ${user.lastName}`;
  const initials = `${user.firstName.charAt(0)}${user.lastName.charAt(0)}`.toUpperCase();

  return (
    <article data-slot="user-profile" style={hueStyle(user.id)} className="max-w-2xl">
      <div
        aria-hidden
        className={cn(
          "relative h-36 overflow-hidden rounded-2xl sm:h-44",
          "bg-[oklch(0.955_0.035_var(--hue))] dark:bg-[oklch(0.235_0.04_var(--hue))]",
          "animate-in duration-500 fade-in-0 motion-reduce:animate-none",
        )}
      >
        <span
          className={cn(
            "absolute end-[-0.03em] bottom-[-0.2em] leading-none select-none",
            "text-[11rem] font-semibold tracking-[-0.06em] font-stretch-condensed sm:text-[15rem]",
            "text-[oklch(0.91_0.055_var(--hue))] dark:text-[oklch(0.29_0.06_var(--hue))]",
          )}
        >
          {initials}
        </span>
      </div>

      <Monogram
        seed={user.id}
        firstName={user.firstName}
        lastName={user.lastName}
        className="relative -mt-9 ms-5 size-[4.5rem] text-2xl ring-4 ring-background"
      />

      <h1
        className={cn(
          "mt-4 text-[2rem] leading-[1.1] font-semibold tracking-[-0.025em] sm:text-[2.5rem]",
          "font-stretch-semi-condensed text-balance",
        )}
      >
        {fullName}
      </h1>
      <a
        href={`mailto:${user.email}`}
        className="mt-1.5 inline-block text-muted-foreground hover:text-foreground"
      >
        {user.email}
      </a>

      <dl className="mt-8 divide-y overflow-hidden rounded-xl border text-sm">
        <Row term="Joined">
          <time dateTime={user.createdAt}>{formatDate(user.createdAt)}</time>
        </Row>
        <Row term="User ID" action={<CopyButton compact value={user.id} label="Copy user ID" />}>
          <span className="font-mono text-xs break-all text-muted-foreground">{user.id}</span>
        </Row>
      </dl>
    </article>
  );
}

function Row({
  term,
  action,
  children,
}: {
  term: string;
  action?: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <div className="grid min-h-12 grid-cols-[6rem_1fr_auto] items-center gap-x-4 py-2 ps-4 pe-2">
      <dt className="text-muted-foreground">{term}</dt>
      <dd className="min-w-0">{children}</dd>
      <div className="justify-self-end">{action}</div>
    </div>
  );
}
