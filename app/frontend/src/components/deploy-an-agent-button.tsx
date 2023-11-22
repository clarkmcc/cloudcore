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
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group.tsx";
import { Label } from "@radix-ui/react-dropdown-menu";
import { AppleIcon } from "@/components/icons/apple-icon.tsx";
import { LinuxIcon } from "@/components/icons/linux-icon.tsx";
import { WindowsIcon } from "@/components/icons/windows-icon.tsx";

// type DeployAnAgentButtonProps = ButtonProps & {};

export const DeployAnAgentButton = forwardRef<HTMLButtonElement, ButtonProps>(
  (props, ref) => {
    const [open, setOpen] = useState(false);
    const [os] = useState("macos");
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

            <div>
              <RadioGroup value={os} className="grid grid-cols-3 gap-4">
                <div>
                  <RadioGroupItem
                    value="macos"
                    id="macos"
                    className="peer sr-only"
                  />
                  <Label className="flex flex-col items-center justify-between rounded-md border-2 border-muted bg-popover p-4 hover:bg-accent hover:text-accent-foreground peer-data-[state=checked]:border-primary [&:has([data-state=checked])]:border-primary">
                    <AppleIcon className="mb-3 h-6 w-6" />
                    macOS
                  </Label>
                </div>
                <div>
                  <RadioGroupItem
                    value="linux"
                    id="linux"
                    className="peer sr-only"
                  />
                  <Label className="flex flex-col items-center justify-between rounded-md border-2 border-muted bg-popover p-4 hover:bg-accent hover:text-accent-foreground peer-data-[state=checked]:border-primary [&:has([data-state=checked])]:border-primary">
                    <LinuxIcon className="mb-3 h-6 w-6" />
                    Linux
                  </Label>
                </div>
                <div>
                  <RadioGroupItem
                    disabled
                    value="windows"
                    id="windows"
                    className="peer sr-only"
                  />
                  <Label className="flex flex-col items-center justify-between rounded-md border-2 border-muted bg-popover p-4 hover:bg-accent hover:text-accent-foreground peer-data-[state=checked]:border-primary [&:has([data-state=checked])]:border-primary">
                    <WindowsIcon className="mb-3 h-6 w-6" />
                    Windows
                  </Label>
                </div>
              </RadioGroup>
            </div>

            <div>
              <RadioGroup value={os} className="grid grid-cols-3 gap-4">
                <div>
                  <RadioGroupItem
                    value="macos"
                    id="macos"
                    className="peer sr-only"
                  />
                  <Label className="flex flex-col items-center justify-between rounded-md border-2 border-muted bg-popover p-4 hover:bg-accent hover:text-accent-foreground peer-data-[state=checked]:border-primary [&:has([data-state=checked])]:border-primary">
                    64-bit ARM
                  </Label>
                </div>
                <div>
                  <RadioGroupItem
                    value="linux"
                    id="linux"
                    className="peer sr-only"
                  />
                  <Label className="flex flex-col items-center justify-between rounded-md border-2 border-muted bg-popover p-4 hover:bg-accent hover:text-accent-foreground peer-data-[state=checked]:border-primary [&:has([data-state=checked])]:border-primary">
                    32-bit ARM
                  </Label>
                </div>
                <div>
                  <RadioGroupItem
                    value="linux"
                    id="linux"
                    className="peer sr-only"
                  />
                  <Label className="flex flex-col items-center justify-between rounded-md border-2 border-muted bg-popover p-4 hover:bg-accent hover:text-accent-foreground peer-data-[state=checked]:border-primary [&:has([data-state=checked])]:border-primary">
                    64-bit x86
                  </Label>
                </div>
              </RadioGroup>
            </div>
          </DialogContent>
        </Dialog>
      </>
    );
  },
);
