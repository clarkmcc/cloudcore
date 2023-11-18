import { PageHeader } from "@/components/page-header.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Plus } from "lucide-react";

export function HostsPage() {
  return (
    <>
      <PageHeader
        title="Hosts"
        subtitle="Machines where a cloudcore agent is installed"
      />
      <div className="pl-7">
        <Button>
          <Plus className="mr-2" size={16} />
          Deploy an agent
        </Button>
      </div>
    </>
  );
}
