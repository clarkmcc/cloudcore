import { NavigateOptions } from "react-router/dist/lib/context";

import { useNavigate, useParams, To, useLocation } from "react-router-dom";
import { useCallback } from "react";

export function useProjectNavigate() {
  const navigate = useNavigate();
  const { projectId } = useParams();

  return useCallback(
    (to: To | number, options?: NavigateOptions): void => {
      if (typeof to === "string" || typeof to === "object") {
        navigate(
          `/projects/${projectId ?? ""}${
            typeof to === "string" ? to : to.pathname
          }`,
          options,
        );
      } else {
        // If 'to' is a number, use it as is for navigate
        // navigate(to, options);
      }
    },
    [projectId],
  );
}

// Custom hook to set or update the project ID in the URL
export function useProjectId(): [string | undefined, (id: string) => void] {
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const { projectId } = useParams();

  const setProjectId = useCallback(
    (newProjectId: string) => {
      if (!newProjectId || newProjectId === projectId) return;

      // Extract the path after the current project ID
      const pathAfterProjectId =
        pathname.split(`/projects/${projectId}`)[1] || "";

      // Navigate to the new project ID with the same subpath
      navigate(`/projects/${newProjectId}${pathAfterProjectId}`, {
        replace: true,
      });
    },
    [projectId, pathname, navigate],
  );

  return [projectId, setProjectId];
}
