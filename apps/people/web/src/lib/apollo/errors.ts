import { CombinedGraphQLErrors } from "@apollo/client";

export const ApiErrorCode = {
  NotFound: "NOT_FOUND",
  BadUserInput: "BAD_USER_INPUT",
  Conflict: "CONFLICT",
  Internal: "INTERNAL",
} as const;

export type ApiErrorCode = (typeof ApiErrorCode)[keyof typeof ApiErrorCode];

export function hasApiErrorCode(error: unknown, code: ApiErrorCode) {
  return (
    CombinedGraphQLErrors.is(error) &&
    error.errors.some((graphQLError) => graphQLError.extensions?.code === code)
  );
}
