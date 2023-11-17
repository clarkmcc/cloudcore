import { useAuth0 } from "@auth0/auth0-react";
import { useMutation } from "@apollo/client";
import { MUTATION_ENSURE_USER } from "@/queries/users.ts";
import { useEffect } from "react";

export function useEnsureUser() {
  const { isAuthenticated, getAccessTokenSilently } = useAuth0();
  const [mutate, { data, loading, error }] = useMutation(MUTATION_ENSURE_USER);

  useEffect(() => {
    if (isAuthenticated) {
      (async () => {
        const token = await getAccessTokenSilently();
        console.log(token);
        await mutate({ variables: { accessToken: token } });
      })();
    }
  }, [isAuthenticated, getAccessTokenSilently]);

  return { data, loading, error };
}
