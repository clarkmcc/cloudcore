import { useAuth0 } from "@auth0/auth0-react";
import { useMutation } from "@apollo/client";
import { MUTATION_ENSURE_USER } from "@/queries/users.ts";
import { useEffect } from "react";
import { useProject } from "@/context/project.tsx";

export function useEnsureUser() {
  const { isAuthenticated, getAccessTokenSilently } = useAuth0();
  const [mutate, { data, loading, error }] = useMutation(MUTATION_ENSURE_USER);
  const { setProjects } = useProject();

  useEffect(() => {
    if (isAuthenticated) {
      (async () => {
        const token = await getAccessTokenSilently();

        // Upsert the user and load all the user's projects into the project
        // context. If the user was just created, then there would only be one.
        const response = await mutate({ variables: { accessToken: token } });
        const allProjects = response?.data?.ensureUser ?? [];
        if (allProjects.length > 0) {
          setProjects(allProjects);
        }
      })();
    }
  }, [isAuthenticated, getAccessTokenSilently]);

  return { data, loading, error };
}
