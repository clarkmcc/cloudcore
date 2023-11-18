import {
  ApolloProvider,
  ApolloClient,
  InMemoryCache,
  createHttpLink,
  NormalizedCacheObject,
} from "@apollo/client";
import { useAuth0 } from "@auth0/auth0-react";
import React, { useEffect, useState } from "react";
import { baseUrl } from "@/queries/base.ts";
import { setContext } from "@apollo/client/link/context";

// Function to create Apollo client
const createApolloClient = (getAccessTokenSilently: () => Promise<string>) => {
  const httpLink = createHttpLink({
    uri: baseUrl("/graphql"),
  });

  const authLink = setContext(async () => {
    const token = await getAccessTokenSilently().catch(console.error);
    return {
      headers: {
        authorization: token ? `Bearer ${token}` : "",
      },
    };
  });

  return new ApolloClient({
    link: authLink.concat(httpLink),
    cache: new InMemoryCache(),
    connectToDevTools: true,
  });
};

type Props = {
  children: React.ReactNode;
};

// Component
export const AuthenticatedApolloProvider = ({ children }: Props) => {
  const { isAuthenticated, getAccessTokenSilently } = useAuth0();
  const [client, setClient] =
    useState<ApolloClient<NormalizedCacheObject> | null>(null);

  useEffect(() => {
    if (isAuthenticated) {
      const apolloClient = createApolloClient(getAccessTokenSilently);
      setClient(apolloClient);
    }
  }, [isAuthenticated, getAccessTokenSilently]);

  if (!client || !isAuthenticated) {
    // loginWithRedirect().then().catch(console.error);
    return <div>Loading...</div>; // or any other loading indicator
  }

  return <ApolloProvider client={client}>{children}</ApolloProvider>;
};
