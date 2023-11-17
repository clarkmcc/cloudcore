import { NavigateOptions } from "react-router/dist/lib/context";

import {
  useNavigate,
  useParams,
  NavigateFunction,
  To,
  useLocation,
} from "react-router-dom";

export function useProjectNavigate(): NavigateFunction {
  const navigate = useNavigate();
  const { projectId } = useParams();

  return (to: To | number, options?: NavigateOptions): void => {
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
  };
}

// Custom hook to set or update the project ID in the URL
export function useProjectId(): [
  string | undefined,
  (projectId: string) => void,
] {
  const navigate = useNavigate();
  const location = useLocation();
  const { projectId } = useParams();

  const setProjectId = (newProjectId: string) => {
    const pathParts = location.pathname.split("/").filter(Boolean);

    // Check if the current URL already has a project ID
    if (pathParts[0] === "projects" && pathParts[1]) {
      // Replace the existing project ID
      pathParts[1] = newProjectId;
    } else {
      // Add the project ID
      pathParts.unshift("projects", newProjectId);
    }

    // Navigate to the new URL
    navigate(`/${pathParts.join("/")}`);
  };

  return [projectId, setProjectId];
}
