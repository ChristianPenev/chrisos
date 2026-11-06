# CLAUDE.md

People web — the Next.js 16 frontend for the people service (`../api`, Go + gqlgen GraphQL).
Read `AGENTS.md` first: this Next version differs from training data, and its docs live in
`node_modules/next/dist/docs/`.

## Skills (how-to router)

How-to detail lives in skills (`.claude/skills/<name>/SKILL.md`). This file holds only the
rules that must always hold. When a task touches one of these areas, load the skill first:

| Skill           | Use when                                                                                                                                                             |
| --------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `component`     | Creating/moving a component, hook or view; deciding shadcn vs ui vs domain; folder anatomy, barrels, contracts.ts, data-slot, `'use client'`, import aliases         |
| `page`          | Adding/changing a route: thin page.tsx → view, `PageProps<'/route'>`, Suspense under Cache Components, not-found, metadata                                           |
| `data-fetching` | GraphQL operations + codegen, Apollo on server and client, `'use cache'` vs Suspense streaming, server-only modules; "blocking-route" errors, `next/headers` imports |

Add a row when you add a skill. Skill names are singular, kebab-case, domain-oriented.

## Commands

```bash
npm run dev            # next dev (proxies /graphql to API_URL, default http://localhost:8080)
npm run build
npm run lint           # eslint (max-len 100, import/order, alias rules)
npm run lint:fix
npm run format         # prettier --write .
npm run typecheck      # next typegen && tsc --noEmit
npm run codegen        # GraphQL types from ../api/schema/*.graphqls
```

Before finishing any change: `npm run format && npm run lint && npm run typecheck`.

## Structure

```
src/
├── app/                         # routes only: page.tsx / layout.tsx / not-found.tsx stay thin
├── components/
│   ├── shadcn/                  # primitives pulled via the shadcn CLI (lowercase files)
│   └── ui/<Name>/               # our own generic components, no domain knowledge
├── domains/<domain>/
│   ├── contracts.ts             # all props/data types for the domain
│   ├── components/<Name>/       # <Name>.tsx + index.ts, barrelled in components/index.ts
│   ├── views/<Name>View/        # page compositions that routes render
│   ├── hooks/use<Name>/         # domain hooks (when needed)
│   ├── graphql/                 # *.graphql operations + generated *.generated.ts
│   └── mocks/                   # placeholder data until the API is wired
├── lib/                         # cn() and other framework-free helpers
└── generated/types.ts           # codegen base schema types (never edit)
```

Aliases: `@shadcn/*`, `@ui/*`, `@components/*`, `@domains/*`, `@lib/*`, `@hooks/*`,
`@generated/*`. There is no `@/*`, and `../` imports are banned; `./` is for the same folder.

## Invariants

- **Routes are thin.** A `page.tsx` turns route params into view props and renders one view
  from `@domains/<domain>/views`. Markup lives in the view.
- **Types live in `contracts.ts`**, not in component files. Server-only files can then share
  types with client files without the client importing them.
- **Never re-export a server-only module from a barrel that client code can reach.** Server-only
  means anything that touches `next/headers`, `cookies()`, or the server Apollo client.
  Turbopack follows every import, including `export type`, and the build fails with
  "importing a module that depends on next/headers". Import such files by direct path.
- **Cache Components is on** (`cacheComponents` + `partialPrefetching`). Uncached async work and
  runtime APIs (`params`, `cookies`, `headers`, `searchParams`) must sit under `<Suspense>` with a
  skeleton fallback, or be cached with `'use cache'` + `cacheLife`. This is the opposite of the
  "no Suspense above content" rule in older App Router codebases. See `data-fetching`.
- **`'use client'` only when earned**: the file itself uses hooks, event handlers or browser
  APIs. Client-ness flows down; a needless directive keeps the subtree out of RSC.
- **No Apollo streaming.** Don't add `@apollo/client-integration-nextjs`, `PreloadQuery` or
  `useSuspenseQuery`. Server data reaches client components as props.
- **Max line length is 100**, code included. Prettier wraps code; split long Tailwind strings
  into several arguments of one `cn()` call. Only `components/shadcn/**` is exempt.
- **Don't hand-edit `components/shadcn/**`** beyond what the shadcn CLI writes. Wrap it in a
  `ui/` component instead.
- **No code comments.** Name things so they don't need one. Put reasoning in the PR or chat.
