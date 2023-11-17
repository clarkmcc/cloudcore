import React, { createContext, useContext, useState } from "react";
import { Project } from "@/types";
import { useProjectId } from "@/hooks/navigation.ts";

type ProjectContextType = {
  projects: Project[];
  activeProject: Project | null;
  setActiveProject: (project: Project) => void;
  setProjects: (projects: Project[]) => void;
};

const ProjectContext = createContext<ProjectContextType | null>(null);

type Props = {
  children: React.ReactNode;
};

export function ProjectProvider({ children }: Props) {
  const [projects, setProjects] = useState<Project[]>([]);
  const [activeProject, setActiveProject] = useState<Project | null>(null);
  const [_, setProjectId] = useProjectId();

  const handleSetActiveProject = (project: Project) => {
    setActiveProject(project);
    setProjectId(project.id);
  };

  return (
    <ProjectContext.Provider
      value={{
        projects,
        activeProject,
        setActiveProject: handleSetActiveProject,
        setProjects,
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
