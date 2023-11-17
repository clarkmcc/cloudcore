import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx";
import { Plus } from "lucide-react";
import { useProject } from "@/context/project.tsx";
import { useState } from "react";
import { ProjectCreateDialog } from "@/components/project-create-dialog.tsx";

export function ProjectSelector() {
  const { projects, activeProject, setActiveProject } = useProject();
  const [open, setOpen] = useState(false);

  function handleCreateProject(value: string) {
    switch (value) {
      case "new":
        setOpen(true);
        console.log("opening");
        break;
      default:
        // eslint-disable-next-line no-case-declarations
        const project = projects.find((project) => project.id === value);
        if (project) {
          setActiveProject(project);
        }
        break;
    }
  }

  return (
    <>
      <Select value={activeProject?.id} onValueChange={handleCreateProject}>
        <SelectTrigger>
          <SelectValue placeholder="Project" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="new" className="cursor-pointer">
            <div className="flex">
              <Plus className="w-4 h-5 mr-2" />
              Create a project
            </div>
          </SelectItem>
          {projects.map((project) => (
            <SelectItem key={project.id} value={project.id}>
              {project.name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <ProjectCreateDialog open={open} setOpen={setOpen} />
    </>
  );
}
