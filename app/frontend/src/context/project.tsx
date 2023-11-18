import React, {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { Project } from "@/types";
import { useProjectId } from "@/hooks/navigation.ts";

type ProjectContextType = {
  projects: Project[];
  activeProjectId?: string;
  setActiveProject: (project: Project) => void;
  setProjects: (projects: Project[]) => void;
};

const ProjectContext = createContext<ProjectContextType | null>(null);

type Props = {
  children: React.ReactNode;
};

export function ProjectProvider({ children }: Props) {
  const [projects, setProjects] = useState<Project[]>([]);
  const [projectId, setProjectId] = useProjectId();

  function getFromLocalStorage(): Project[] {
    const projects = localStorage.getItem("projects");
    if (projects) return JSON.parse(projects);
    return [];
  }

  const saveToLocalStorage = useCallback((projects: Project[]) => {
    localStorage.setItem("projects", JSON.stringify(projects));
  }, []);

  // On startup, load projects from local storage
  useEffect(() => {
    const projects = getFromLocalStorage();
    if (projects.length > 0) {
      setProjects(projects);
    }
  }, []);

  const handleSetProjects = useCallback(
    (projects: Project[]) => {
      setProjects(projects);
      saveToLocalStorage(projects);

      // On first login, we'll load the list of projects but there won't be
      // any active project, so we'll set the first project as active.
      if (projectId == null && projects.length > 0) {
        setProjectId(projects[0].id);
      }
    },
    [projectId, setProjects, projects, saveToLocalStorage],
  );

  const handleSetActiveProject = useCallback(
    (project: Project) => {
      console.log("setting active project");
      setProjectId(project.id);
    },
    [setProjectId],
  );

  return (
    <ProjectContext.Provider
      value={{
        projects,
        activeProjectId: projectId,
        setActiveProject: handleSetActiveProject,
        setProjects: handleSetProjects,
      }}
    >
      {children}
    </ProjectContext.Provider>
  );
}

export function useProject() {
  const context = useContext(ProjectContext);
  if (!context) {
    throw new Error("useProject must be used within a ProjectProvider");
  }
  return context;
}
