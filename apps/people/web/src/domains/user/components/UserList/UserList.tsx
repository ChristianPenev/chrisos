import { UserListItem } from "@domains/user/components/UserListItem";
import type { UserListProps } from "@domains/user/contracts";

export function UserList({ users }: UserListProps) {
  return (
    <table data-slot="user-list" className="w-full table-fixed border-collapse text-start">
      <thead className="border-b bg-sunken text-[0.8125rem] text-muted-foreground">
        <tr>
          <th scope="col" className="py-2.5 ps-4 pe-3 text-start font-normal sm:ps-5">
            Name
          </th>
          <th scope="col" className="w-28 px-3 py-2.5 text-start font-normal sm:w-44">
            Role
          </th>
          <th scope="col" className="hidden w-36 px-3 py-2.5 text-start font-normal sm:table-cell">
            Joined
          </th>
          <th scope="col" className="w-12 py-2.5">
            <span className="sr-only">Actions</span>
          </th>
        </tr>
      </thead>
      <tbody>
        {users.map((user) => (
          <UserListItem key={user.id} user={user} />
        ))}
      </tbody>
    </table>
  );
}
