import { ProjectMetrics } from "@/types";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import UsageBar from "react-usage-bar";
import { useMemo } from "react";
import { osToColorClasses, osToString } from "@/lib/utils.ts";
import { useTheme } from "@/components/theme-provider.tsx";

type ProjectMetricsCardsProps = {
  metrics?: ProjectMetrics;
  loading?: boolean;
};

export function ProjectMetricsCards({
  metrics,
  loading,
}: ProjectMetricsCardsProps) {
  const { theme } = useTheme();
  const hostsByOsNames = useMemo(
    () =>
      (metrics?.hostsByOsName ?? []).map((v) => ({
        name: osToString(v.osName),
        className: osToColorClasses(v.osName),
        dotClassName: osToColorClasses(v.osName),
        value: v.count,
      })),
    [metrics],
  );

  return (
    <div className="grid grid-cols-4 px-7 lg:space-x-4 lg:space-y-0 space-y-4">
      <Card className="col-span-4 lg:col-span-1">
        <CardHeader>
          <CardTitle>Hosts</CardTitle>
          <CardDescription>Devices where an agent is installed</CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <Skeleton className="h-6" />
          ) : (
            <>
              <div className="text-4xl font-bold">
                {metrics?.totalHosts} hosts
              </div>
              <div className="text-xs text-gray-600">
                {metrics?.totalAgents} agents
              </div>
            </>
          )}
        </CardContent>
      </Card>

      <Card className="col-span-4 lg:col-span-1">
        <CardHeader>
          <CardTitle>Online</CardTitle>
          <CardDescription>How many hosts are reporting</CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <Skeleton className="h-6" />
          ) : (
            <UsageBar
              items={[
                {
                  name: "Online",
                  value: metrics?.onlineHosts ?? 0,
                  className: "bg-green-500 dark:bg-green-500",
                  dotClassName: "bg-green-500 dark:bg-green-500",
                },
                {
                  name: "Offline",
                  value: metrics?.offlineHosts ?? 0,
                  className: "bg-red-500 dark:bg-red-500",
                  dotClassName: "bg-red-500 dark:bg-red-500",
                },
              ]}
              total={metrics?.totalHosts ?? 0}
              showPercentage
              compactLayout
              usageBarContainerClassName="p-0"
              usageBarClassName="m-w-full"
              darkMode={theme === "dark"}
            />
          )}
        </CardContent>
      </Card>

      <Card className="col-span-4 lg:col-span-2">
        <CardHeader>
          <CardTitle>Platforms</CardTitle>
          <CardDescription>Distinct operating systems</CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <Skeleton className="h-6" />
          ) : (
            <UsageBar
              items={hostsByOsNames}
              total={metrics?.totalHosts ?? 0}
              showPercentage
              compactLayout
              darkMode={theme === "dark"}
            />
          )}
        </CardContent>
      </Card>
    </div>
  );
}
