import { useAuth0 } from "@auth0/auth0-react";
import { useMutation } from "@apollo/client";
import { MUTATION_ENSURE_USER } from "@/queries/users.ts";
import { useEffect } from "react";
import { useProject } from "@/context/project.tsx";
import { Project } from "@/types";

export function useEnsureUser() {
  const { isAuthenticated, getAccessTokenSilently } = useAuth0();
  const [mutate, { data, loading, error }] = useMutation(MUTATION_ENSURE_USER);
  const { setProjects, setActiveProject } = useProject();

  useEffect(() => {
    if (isAuthenticated) {
      (async () => {
        const token = await getAccessTokenSilently();

        const storedProjects = getFromLocalStorage();
        if (storedProjects.length > 0) {
          setProjects(storedProjects);
          setActiveProject(storedProjects[0]);
        }

        // Upsert the user and load all the user's projects into the project
        // context. If the user was just created, then there would only be one.
        const response = await mutate({ variables: { accessToken: token } });
        const projects = response?.data?.ensureUser ?? [];
        if (projects.length > 0) {
          // If we already loaded the projects from local storage, then just
          // update them here, but don't overwrite the project ID in the URL.
          if (storedProjects.length > 0) {
            setProjects(projects);
          } else {
            setActiveProject(projects[0]);
          }
          saveToLocalStorage(projects);
        }
      })();
    }
  }, [isAuthenticated, getAccessTokenSilently]);

  return { data, loading, error };
}

function getFromLocalStorage(): Project[] {
  const projects = localStorage.getItem("projects");
  if (projects) return JSON.parse(projects);
  return [];
}

function saveToLocalStorage(projects: Project[]) {
  localStorage.setItem("projects", JSON.stringify(projects));
}
