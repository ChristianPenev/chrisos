---
name: component
description: Conventions for components, hooks and views in people/web — where a new component goes (components/shadcn vs components/ui vs domains/<domain>), folder-per-component anatomy, index.ts barrels and the self-import cycle trap, contracts.ts types, data-slot attributes, cn() and the 100-char line rule, when 'use client' is earned, import aliases and ordering. Use when creating, moving or splitting a component/hook/view, adding a shadcn primitive, or hitting import/order, max-len or no-restricted-imports lint errors.
---

# Component

## Where does it go?

| It is…                                                        | Put it in                                     | Model                                  |
| ------------------------------------------------------------- | --------------------------------------------- | -------------------------------------- |
| A primitive from the shadcn registry (Base UI)                | `src/components/shadcn/<name>.tsx` (CLI only) | `shadcn/button.tsx`                    |
| Our own generic component that knows nothing about any domain | `src/components/ui/<Name>/`                   | `ui/CopyButton/`                       |
| A component that knows about a business entity (user, org, …) | `src/domains/<domain>/components/<Name>/`     | `domains/user/components/UserProfile/` |
| A full page composition that a route renders                  | `src/domains/<domain>/views/<Name>View/`      | `domains/user/views/UserProfileView/`  |
| A hook used by one domain                                     | `src/domains/<domain>/hooks/use<Name>/`       | —                                      |
| A hook used across domains                                    | `src/hooks/use<Name>/`                        | —                                      |
| A pure helper with no React                                   | `src/lib/`                                    | `lib/utils.ts`                         |

Rule of thumb: if you could reuse it in a different app without renaming props, it's `ui/`.
`CopyButton` takes `value` + `label`, so it's `ui/`. `Monogram` takes `firstName` + `lastName`,
so it's `domains/user/`.

Add a primitive with `npx shadcn@latest add <name>`. `components.json` points `ui` at
`@shadcn`, so it lands in `components/shadcn/`. Run `npm run format` afterwards. Don't edit the
file by hand beyond formatting; wrap it in a `ui/` component.

## Anatomy of a component

```
domains/user/components/UserProfile/
├── UserProfile.tsx      # the component (PascalCase file = component name)
└── index.ts             # export { UserProfile } from "./UserProfile";
```

1. **Types**: put `<Name>Props` in the domain's `contracts.ts` (ui atoms get a local
   `contracts.ts` next to the component, e.g. `ui/CopyButton/contracts.ts`). Private
   sub-components used in one file (like `Row` in `UserProfile.tsx`) can keep inline types.
2. **Barrels**: add the component to `domains/<domain>/components/index.ts` (and views to
   `views/index.ts`). Consumers outside the domain import from the barrel:
   `import { UserProfile } from "@domains/user/components"`.
3. **No self-import through the barrel.** Inside the same `components/` folder, import a
   sibling by its own folder: `@domains/user/components/Monogram`, never
   `@domains/user/components`. The barrel re-exports the importer, so that would be a cycle.
4. **`data-slot="<kebab-name>"`** on the root element (`data-slot="user-profile"`). It's the
   stable hook for tests and CSS. Never select on Tailwind classes.
5. **Styling**: Tailwind, merged with `cn` from `@lib/utils` (a tailwind-merge drop-in). Put the
   caller's `className` last in the same `cn()` call so it wins.
6. **Line length 100**: Prettier wraps code but not strings. Split long class lists into
   several `cn()` arguments, grouped by concern (layout / type / color / states). See `Monogram`.

## Server vs client

- A component is a Server Component by default. Add `"use client"` only if **this file** uses
  state/effects/hooks, attaches event handlers, or touches browser APIs (`navigator`,
  `window`). `CopyButton` earns it; `UserProfile` renders `CopyButton` but doesn't.
- Keep client islands small and leaf-level. Pass data down as serializable props.
- Don't branch on `typeof window` in render. It causes hydration mismatches.
- Server-only code (anything that reads `cookies()`/`headers()` or uses the server Apollo
  client) must never be re-exported from a barrel that a client file can reach. Import it by
  direct path. Types it shares go in `contracts.ts`.

## Views

A view is the page body: layout container, Suspense boundaries, and data → component wiring.
It takes plain props (or promises of them) from the route; it never reads `params` itself.
The model is `UserProfileView`: a `<main>` shell renders statically, and an inner async
`UserProfileContent` awaits the data under `<Suspense fallback={<UserProfileSkeleton />}>`.
Every async view section ships a matching `<Name>Skeleton` component.

## Imports

Enforced by ESLint (`npm run lint:fix` sorts them):

```tsx
import { Check } from "lucide-react"; // 1. packages (alphabetized)
import { useState } from "react";

import { Monogram } from "@domains/user/components/Monogram"; // 2. our aliases
import { cn } from "@lib/utils";
import { CopyButton } from "@ui/CopyButton";

import type { CopyButtonProps } from "./contracts"; // 3. same-folder
```

`@/*` and `../` are banned. Aliases: `@shadcn/*`, `@ui/*`, `@components/*`, `@domains/*`,
`@lib/*`, `@hooks/*`, `@generated/*` (defined in `tsconfig.json`).

## Verify

```bash
npm run format && npm run lint && npm run typecheck
```
