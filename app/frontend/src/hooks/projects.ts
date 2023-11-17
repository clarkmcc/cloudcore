import { MutationResult, useMutation } from "@apollo/client";
import { MUTATION_CREATE_PROJECT } from "@/queries/projects.ts";
import { useProject } from "@/context/project.tsx";

type CreateProjectResult = [
  (name: string, description?: string) => Promise<void>,
  MutationResult,
];

export function useCreateProject(): CreateProjectResult {
  const [mutate, status] = useMutation(MUTATION_CREATE_PROJECT);
  const { setProjects, setActiveProject } = useProject();

  return [
    async (name: string, description?: string) => {
      const response = await mutate({
        variables: {
          name,
          description,
        },
      });
      const project = response.data?.projectCreate?.project;
      const allProjects = response.data?.projectCreate?.allProjects;
      setProjects(allProjects);
      setActiveProject(project);
    },
    status,
  ];
}
