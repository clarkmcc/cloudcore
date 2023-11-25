import { forwardRef, useMemo, useState } from "react";
import { Plus } from "lucide-react";
import { Button, ButtonProps } from "@/components/ui/button.tsx";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog.tsx";
import {
  DeployAgentForm,
  DeployAgentFormProps,
  DeployAgentFormSchema,
} from "@/components/forms/deploy-agent-form.tsx";
import { QUERY_LIST_LATEST_PACKAGES } from "@/queries/packages.ts";
import { useMutation, useQuery } from "@apollo/client";
import { AgentPlatformDownload, GOARCH, GOARM, GOOS, Package } from "@/types";
import { goarchToString, goosToString } from "@/lib/utils.ts";
import { ErrorBanner } from "@/components/error-banner.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import { MUTATION_BUILD_DEPLOY_AGENT_COMMAND } from "@/queries/agent.ts";
import { z } from "zod";
import { useProjectId } from "@/hooks/navigation.ts";
import { CommandCopy } from "@/components/command-copy.tsx";

// type DeployAnAgentButtonProps = ButtonProps & {};

export const DeployAnAgentButton = forwardRef<HTMLButtonElement, ButtonProps>(
  (props, ref) => {
    const [open, setOpen] = useState(false);
    const [projectId] = useProjectId();

    const [mutate, { data, loading, error, reset }] = useMutation(
      MUTATION_BUILD_DEPLOY_AGENT_COMMAND,
    );
    const command = data?.buildDeployAgentCommand ?? null;

    const handleSubmit = async (
      values: z.infer<typeof DeployAgentFormSchema>,
    ) => {
      await mutate({ variables: { ...values, projectId } });
    };

    function handleDone() {
      setOpen(false);
      reset();
    }

    return (
      <>
        <Button ref={ref} {...props} onClick={() => setOpen(true)}>
          <Plus className="mr-2" size={16} />
          Deploy an agent
        </Button>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Deploy an agent</DialogTitle>
              <DialogDescription>
                Agents are used to manage your hosts. In order for cloudcore to
                see your host, you'll need to install the cloudcore agent.
              </DialogDescription>
            </DialogHeader>

            {error && <ErrorBanner error={error} />}

            {!command && (
              <DeployAgentFormLoader
                loading={loading}
                onSubmit={handleSubmit}
              />
            )}

            {command && (
              <>
                <CommandCopy command={command} />
                <div className="flex flex-row-reverse ">
                  <div className="space-x-2">
                    <Button variant="ghost" onClick={reset}>
                      Deploy another agent
                    </Button>
                    <Button onClick={handleDone}>Done</Button>
                  </div>
                </div>
              </>
            )}
          </DialogContent>
        </Dialog>
      </>
    );
  },
);

function DeployAgentFormLoader(props: Omit<DeployAgentFormProps, "downloads">) {
  const { data, loading, error } = useQuery<{ packages: Array<Package> }>(
    QUERY_LIST_LATEST_PACKAGES,
  );
  const downloads = useMemo((): AgentPlatformDownload[] => {
    if (!data) return [];

    // data.packages returns an array of {goos: string, goarch: string}
    // with multiple entries for each goos. We want to group them by goos
    // and then map them to the format that the DeployAgentForm expects.
    const out: { [key: string]: { goarch: GOARCH; goarm: GOARM }[] } = {};
    data.packages.forEach((p) => {
      if (!out[p.goos]) out[p.goos] = [];
      if (out[p.goos].find((a) => a.goarch === p.goarch && a.goarm == p.goarm))
        return;
      out[p.goos].push({ goarch: p.goarch, goarm: p.goarm });
    });
    return Object.entries(out).map(([goos, goarch]) => ({
      goos: {
        display: goosToString(goos),
        value: goos as GOOS,
      },
      goarch: goarch.map(({ goarch, goarm }) => ({
        display: goarchToString(goarch, goarm),
        value: goarch,
      })),
    }));
  }, [data]);

  if (error) {
    return <ErrorBanner error={error} />;
  }

  if (loading) {
    return <Skeleton className="h-64" />;
  }

  return <DeployAgentForm {...props} downloads={downloads} />;
}
