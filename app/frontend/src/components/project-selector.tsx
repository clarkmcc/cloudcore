import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx";
import { Plus } from "lucide-react";
import { useState } from "react";
import { Skeleton } from "@/components/ui/skeleton.tsx";

export function ProjectSelector() {
  const [loading] = useState(true);
  return (
    <Select>
      <SelectTrigger>
        {loading ? (
          <Skeleton className="h-4 w-full mr-2" />
        ) : (
          <SelectValue placeholder="Project" />
        )}
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="new">
          <div className="flex">
            <Plus className="w-4 h-5 mr-2" />
            Create a project
          </div>
        </SelectItem>
        <SelectItem value="1">Project 1</SelectItem>
        <SelectItem value="2">Project 2</SelectItem>
      </SelectContent>
    </Select>
  );
}
