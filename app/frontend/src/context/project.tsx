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
  const [activeProjectId, setActiveProjectId] = useState<string | undefined>();
  const [projectId, setProjectId] = useProjectId();

  function getFromLocalStorage(): Project[] {
    const projects = localStorage.getItem("projects");
    if (projects) return JSON.parse(projects);
    return [];
  }

  function saveToLocalStorage(projects: Project[]) {
    localStorage.setItem("projects", JSON.stringify(projects));
  }

  // On startup, load projects from local storage
  useEffect(() => {
    const projects = getFromLocalStorage();
    if (projects.length > 0) {
      setProjects(projects);
    }
  }, []);

  useEffect(() => {
    if (activeProjectId == null && projects.length > 0 && projectId != null) {
      console.log("active project is null");
      const project = projects.find((project) => project.id === projectId);
      if (project) {
        console.log("setting active project", project.id);
        setActiveProjectId(project.id);
      }
    }
  }, [activeProjectId, projectId, projects]);

  // // If we have projects, but no active project, select one.
  // useEffect(() => {
  //   if (projects.length > 0 && activeProject == null) {
  //     setActiveProject(projects[0]);
  //   }
  // }, [activeProject, projects]);

  const handleSetProjects = useCallback(
    (projects: Project[]) => {
      console.log("setting projects", projects);
      setProjects(projects);
      saveToLocalStorage(projects);
      // if (activeProjectId == null && projects.length > 0) {
      //   setActiveProjectId(projects[0].id);
      // }
    },
    [setProjects, saveToLocalStorage],
  );

  const handleSetActiveProject = useCallback(
    (project: Project) => {
      console.log("setting active project");
      setActiveProjectId(project.id);
      setProjectId(project.id);
    },
    [setActiveProjectId, setProjectId],
  );

  return (
    <ProjectContext.Provider
      value={{
        projects,
        activeProjectId,
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
