import { forwardRef, useState } from "react";
import { Plus } from "lucide-react";
import { Button, ButtonProps } from "@/components/ui/button.tsx";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog.tsx";
import { DeployAgentForm } from "@/components/forms/deploy-agent-form.tsx";

// type DeployAnAgentButtonProps = ButtonProps & {};

export const DeployAnAgentButton = forwardRef<HTMLButtonElement, ButtonProps>(
  (props, ref) => {
    const [open, setOpen] = useState(false);
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

            <DeployAgentForm
              onSubmit={(v) => console.log(v)}
              downloads={[
                {
                  goos: {
                    display: "Linux",
                    value: "linux",
                  },
                  goarch: [
                    {
                      display: "64-bit x86",
                      value: "amd64",
                    },
                    {
                      display: "64-bit ARM",
                      value: "arm64",
                    },
                  ],
                },
                {
                  goos: {
                    display: "macOS",
                    value: "darwin",
                  },
                  goarch: [
                    {
                      display: "64-bit x86",
                      value: "amd64",
                    },
                    {
                      display: "64-bit ARM",
                      value: "arm64",
                    },
                  ],
                },
                {
                  goos: {
                    display: "Windows",
                    value: "windows",
                  },
                  goarch: [
                    {
                      display: "64-bit x86",
                      value: "amd64",
                    },
                  ],
                },
              ]}
            />
          </DialogContent>
        </Dialog>
      </>
    );
  },
);
