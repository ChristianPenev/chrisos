import { ApolloClient, HttpLink, InMemoryCache } from "@apollo/client";
import { cache } from "react";

const apiUrl = process.env.API_URL ?? "http://localhost:8080";

export const getClient = cache(
  () =>
    new ApolloClient({
      link: new HttpLink({ uri: `${apiUrl}/graphql` }),
      cache: new InMemoryCache(),
    }),
);
