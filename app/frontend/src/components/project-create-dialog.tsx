import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog.tsx";
import { Label } from "@/components/ui/label.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useForm } from "react-hook-form";
import { useCreateProject } from "@/hooks/projects.ts";

type ProjectCreateDialogProps = {
  open: boolean;
  setOpen: (open: boolean) => void;
};

interface Form {
  name: string;
  description?: string;
}

export function ProjectCreateDialog({
  open,
  setOpen,
}: ProjectCreateDialogProps) {
  const { register, handleSubmit } = useForm<Form>();
  const [create, status] = useCreateProject();

  async function onSubmit(data: Form) {
    await create(data.name, data.description);
    setOpen(false);
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create a new project</DialogTitle>
          <DialogDescription>
            Projects are used to organize your hosts and groups. They act is
            independent environments for you to manage your fleet.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="space-y-4">
            <div className="grid w-full items-center gap-1.5">
              <Label htmlFor="name">Name</Label>
              <Input
                required
                type="name"
                id="name"
                placeholder="My project"
                {...register("name", { required: true })}
              />
            </div>
            <div className="grid w-full items-center gap-1.5">
              <Label htmlFor="description">Description</Label>
              <Input
                type="description"
                id="description"
                placeholder="..."
                {...register("description")}
              />
            </div>
            <Button type="submit" disabled={status.loading} className="w-full">
              Create
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
