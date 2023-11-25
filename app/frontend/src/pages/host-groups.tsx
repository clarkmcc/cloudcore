import { PageHeader } from "@/components/page-header.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Plus } from "lucide-react";
import { HostGroupCreateDialog } from "@/components/host-group-create-dialog.tsx";
import { useState } from "react";

export function HostGroupsPage() {
  const [open, setOpen] = useState(false);

  return (
    <>
      <PageHeader
        title="Host groups"
        subtitle="Groups of hosts that can be interacted with as a single unit"
      />
      <div className="pl-7">
        <Button onClick={() => setOpen(true)}>
          <Plus className="mr-2" size={16} />
          New group
        </Button>
      </div>
      <HostGroupCreateDialog open={open} setOpen={setOpen} />
    </>
  );
}
