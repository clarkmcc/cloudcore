import { PageHeader } from "@/components/page-header.tsx";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card.tsx";
import { useParams } from "react-router-dom";
import { useQuery } from "@apollo/client";
import { QUERY_HOST_DETAILS } from "@/queries/hosts.ts";
import { ErrorBanner } from "@/components/error-banner.tsx";
import { ShieldCheck, TerminalSquare, Unlink } from "lucide-react";
import { HostEventTable } from "@/components/host-event-table.tsx";
import { Button } from "@/components/ui/button.tsx";

export function HostDetails() {
  const { hostId, projectId } = useParams();
  const missingParams = hostId == null || projectId == null;
  const { data, loading, error } = useQuery(QUERY_HOST_DETAILS, {
    variables: { hostId, projectId },
    skip: missingParams,
    pollInterval: 10000,
  });

  if (missingParams) {
    return <p>No host</p>;
  }
  if (error) {
    return (
      <div className="p-7">
        <ErrorBanner error={error} />
      </div>
    );
  }

  return (
    <>
      <PageHeader
        loading={loading}
        title={data?.host?.hostname}
        subtitle={data?.host?.hostId ?? data?.host?.id}
        backButton
      />

      <div className="px-7">
        <div className="mb-4 space-x-4">
          <Button>
            <TerminalSquare className="mr-2 h-4 w-4" />
            Connect
          </Button>
          <Button variant="destructive" disabled>
            <Unlink className="mr-2 h-4 w-4" />
            Revoke Access
          </Button>
        </div>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <Card className="col-span-4 2xl:col-span-2">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Events</CardTitle>
              <ShieldCheck className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <HostEventTable events={data?.host?.events ?? []} />
            </CardContent>
          </Card>
        </div>
      </div>
    </>
  );
}
