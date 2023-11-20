import { PageHeader } from "@/components/page-header.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Plus } from "lucide-react";
import { HostsTable } from "@/components/hosts-table.tsx";
import { useQuery } from "@apollo/client";
import { QUERY_HOSTS_LIST } from "@/queries/hosts.ts";
import { useProjectId } from "@/hooks/navigation.ts";
import { ErrorBanner } from "@/components/error-banner.tsx";
import CircularProgress from "@mui/material/CircularProgress";
import { useEffect } from "react";

export function HostsPage() {
  const [projectId] = useProjectId();
  const { data, loading, error, refetch } = useQuery(QUERY_HOSTS_LIST, {
    variables: { projectId },
    refetchWritePolicy: "merge",
    pollInterval: 10000,
  });

  // Refetch on page load
  useEffect(() => {
    refetch().catch(console.error);
  }, [refetch]);

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
      <div className="px-7 pt-7">
        {error && <ErrorBanner error={error} />}
        {loading && <CircularProgress color="primary" />}
        {!error && !loading && <HostsTable hosts={data?.hosts ?? []} />}
      </div>
    </>
  );
}
