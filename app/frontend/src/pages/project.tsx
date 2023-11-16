import { LoaderFunction, useLoaderData } from "react-router-dom";
import { useEffect } from "react";

export const loadProject: LoaderFunction = async ({ params }) => {
  // todo: get project from api
  return { projectId: params.projectId };
};

export function ProjectHome() {
  const data = useLoaderData();
  useEffect(() => {
    console.log(data);
  }, [data]);
  return <p>home</p>;
}
