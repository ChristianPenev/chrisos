---
name: data-fetching
description: GraphQL data in people/web — writing .graphql operations next to the domain, running codegen (near-operation-file → *.generated.ts typed documents), the planned Apollo setup (per-request server client via React cache(), client provider for mutations), choosing 'use cache' + cacheLife vs Suspense streaming under Cache Components, notFound vs error, server-only import rules, and the Apollo streaming ban. Use when adding a query or mutation, wiring a view to the API, regenerating types, or hitting blocking-route, "depends on next/headers", or codegen errors.
---

# Data fetching

> **Status:** the API is not wired yet. `UserProfileView` reads `domains/user/mocks/`. The
> codegen config is live; the Apollo client files described under "Apollo" don't exist yet.
> Create them in that shape when wiring the first query, then update this note.

## GraphQL operations + codegen

- The schema source is `../api/schema/*.graphqls`, owned by the Go service (gqlgen). Never
  hand-write API types.
- Put operations next to the domain: `src/domains/<domain>/graphql/<name>.graphql`. Name
  operations after the entity and action (`query UserProfile`, `mutation CreateOrganization`).
  Prefer a fragment per component that renders the data (`fragment UserProfile on User`).
- Run `npm run codegen` (or `codegen:watch`). It writes:
  - `src/generated/types.ts`: base schema types.
  - `<name>.generated.ts` next to each `.graphql`: operation result/variable types plus a
    `TypedDocumentNode` (`UserProfileDocument`) to pass straight to Apollo.
- Generated files are committed, never edited, and excluded from lint/prettier.
  `npm run codegen:check` fails if they're stale.
- Once real data exists, derive `contracts.ts` data types from generated fragment types instead
  of hand-written shapes (`export type UserProfileData = UserProfileFragment`), then delete the
  mock.

## Cache Components: cached or streamed?

`cacheComponents` is on, so every server read is one of these:

| Data                                                        | Do                                                                                      |
| ----------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| Same for every viewer, can be a little stale (lookup lists) | `'use cache'` + `cacheLife("minutes" \| "hours")` (+ `cacheTag` to revalidate on write) |
| Depends on the viewer or must be fresh (auth'd data)        | Plain async component under `<Suspense>` with a skeleton                                |
| Keyed by a route param                                      | Either. The param is a function argument, so it becomes part of the cache key           |

Rules from the Next docs (`node_modules/next/dist/docs/01-app/01-getting-started/08-caching.md`,
`03-api-reference/01-directives/use-cache.md`):

- A `'use cache'` function can't call `cookies()`/`headers()`. Read them outside, then pass the
  value in as an argument.
- Arguments and return values must be serializable. Pass ids, not client instances. Create the
  Apollo client inside the cached function.
- Always pair `'use cache'` with `cacheLife`. Mutations call `revalidateTag`/`updateTag` for the
  tags they affect.

## Apollo (planned shape)

- **Server**: `src/lib/apollo/serverClient.ts` exports `getClient`, wrapped in React `cache()`
  so one request shares one `ApolloClient` (`HttpLink` to `API_URL`/graphql, `InMemoryCache`).
  It's server-only: import it by direct path, never from a barrel.
- **Data functions** live in the domain, e.g. `domains/user/data/getUser.ts`:
  ```ts
  export async function getUser(id: string) {
    const client = await getClient();
    const { data } = await client.query({ query: UserProfileDocument, variables: { id } });
    return data?.user ?? null;
  }
  ```
  The view's async content component calls it, renders `notFound()` on `null`, and passes the
  result to presentational components as props. `generateMetadata` calls the same function.
- **Client**: only add an `ApolloProvider` (in a `"use client"` provider in `src/app/`) when a
  client component needs `useMutation`/`useQuery` from `@apollo/client/react`. Reads that the
  server can do stay on the server.
- In dev, `/graphql` is proxied to `API_URL` by `next.config.ts` rewrites (dev only).

## Bans

- **No Apollo streaming**: no `@apollo/client-integration-nextjs`,
  `@apollo/client-react-streaming`, `PreloadQuery`, `useSuspenseQuery`. A sister codebase
  removed them. The open stream kept the tab spinning, and every SSR query got serialized into
  the HTML. Next's own Suspense streaming covers this.
- **No barrel re-exports of server-only files**, not even `export type`. Turbopack follows them
  into client bundles and the build fails with "importing a module that depends on
  next/headers". Shared types go in `contracts.ts`.
- **No bare `fetch` to the API.** Go through the generated documents and Apollo, so types
  stay in sync with the schema.
