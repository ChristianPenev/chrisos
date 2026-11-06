import { CopyButton } from "@/components/users/copy-button";
import { Monogram } from "@/components/users/monogram";

export type UserProfileData = {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  createdAt: string;
};

const addedOn = new Intl.DateTimeFormat("en-GB", { dateStyle: "long" });

export function UserProfile({ user }: { user: UserProfileData }) {
  const fullName = `${user.firstName} ${user.lastName}`;

  return (
    <article className="flex flex-col gap-10">
      <header className="flex items-center gap-5">
        <Monogram seed={user.id} firstName={user.firstName} lastName={user.lastName} />
        <div className="min-w-0">
          <h1 className="font-heading text-4xl font-semibold tracking-tight text-balance sm:text-5xl">
            {fullName}
          </h1>
          <a
            href={`mailto:${user.email}`}
            className="mt-1 block truncate text-muted-foreground underline-offset-4 hover:text-foreground hover:underline"
          >
            {user.email}
          </a>
        </div>
      </header>

      <dl className="divide-y border-y text-sm">
        <Row term="Added">
          <time dateTime={user.createdAt}>{addedOn.format(new Date(user.createdAt))}</time>
        </Row>
        <Row term="User ID" action={<CopyButton value={user.id} label="Copy user ID" />}>
          <span className="font-mono text-xs break-all">{user.id}</span>
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
    <div className="grid min-h-12 grid-cols-[7rem_1fr_auto] items-center gap-x-4 py-2">
      <dt className="text-muted-foreground">{term}</dt>
      <dd className="min-w-0">{children}</dd>
      <div className="justify-self-end">{action}</div>
    </div>
  );
}
