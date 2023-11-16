import { NavigateOptions } from "react-router/dist/lib/context";

import { useNavigate, useParams, NavigateFunction, To } from "react-router-dom";

export function useProjectNavigate(): NavigateFunction {
  const navigate = useNavigate();
  const { projectId } = useParams();

  return (to: To | number, options?: NavigateOptions): void => {
    if (typeof to === "string" || typeof to === "object") {
      navigate(
        `/projects/${projectId}${typeof to === "string" ? to : to.pathname}`,
        options,
      );
    } else {
      // If 'to' is a number, use it as is for navigate
      // navigate(to, options);
    }
  };
}
