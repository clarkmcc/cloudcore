import { PageHeader } from "@/components/page-header.tsx";
import moment from "moment";
import "react-usage-bar/build/index.css";
import { useQuery } from "@apollo/client";
import { QUERY_PROJECT_METRICS } from "@/queries/projects.ts";
import { ErrorBanner } from "@/components/error-banner.tsx";
import { useProjectId } from "@/hooks/navigation.ts";
import { ProjectMetrics } from "@/types";
import { useEffect } from "react";
import { ProjectMetricsCards } from "@/components/project-metrics-cards.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import { DeployAnAgentButton } from "@/components/deploy-an-agent-button.tsx";

export function HomePage() {
  const [projectId] = useProjectId();
  const { data, loading, error, refetch } = useQuery<{
    projectMetrics: ProjectMetrics;
  }>(QUERY_PROJECT_METRICS, {
    variables: { projectId },
  });
  const metrics = data?.projectMetrics;
  const hasHosts = (metrics?.totalHosts ?? 0) > 0;

  // Refetch on page load
  useEffect(() => {
    refetch().catch(console.error);
  }, [refetch]);

  if (error) {
    return (
      <div className="p-7">
        <ErrorBanner error={error} />
      </div>
    );
  }

  return (
    <>
      <PageHeader title="Home" subtitle={moment().format("dddd, MMMM D")} />
      {!hasHosts && loading && (
        <div className="px-7">
          <Skeleton className="w-full h-6" />
        </div>
      )}
      {hasHosts && <ProjectMetricsCards metrics={metrics} loading={loading} />}
      {!hasHosts && !loading && (
        <div className="px-7">
          <div className="bg-gray-100 dark:bg-neutral-900 rounded-xl p-24 flex flex-row items-center justify-center">
            <div className="w-1/2">
              <div className="text-2xl text-gray-500 dark:text-gray-300 font-bold text-center">
                No hosts yet
              </div>
              <div className="text-gray-500 dark:text-gray-300 text-xs text-center">
                Start by installing the agent on a host
              </div>
              <div className="flex justify-center items-center py-4">
                <DeployAnAgentButton />
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
