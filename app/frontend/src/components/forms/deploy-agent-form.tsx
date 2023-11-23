import { z } from "zod";
import { ControllerRenderProps, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group.tsx";
import { Label } from "@radix-ui/react-dropdown-menu";
import { Button } from "@/components/ui/button.tsx";
import { cn, goosToIcon } from "@/lib/utils.ts";
import React, { useCallback, useEffect } from "react";
import { AgentPlatformDownload, DisplayableValue, GOARCH } from "@/types";
import { Switch } from "@/components/ui/switch.tsx";

const FormSchema = z.object({
  os: z.enum(["linux", "windows", "darwin"], {
    required_error: "Please select an operating system",
  }),
  arch: z.enum(["amd64", "arm64", "386"]),
  generate_psk: z.boolean(),
});

type DeployAgentFormProps = {
  downloads: AgentPlatformDownload[];
  onSubmit: (values: z.infer<typeof FormSchema>) => void;
};

export function DeployAgentForm({ onSubmit, downloads }: DeployAgentFormProps) {
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
  });

  const getArchs = useCallback(
    (os: string): DisplayableValue<GOARCH>[] => {
      const download = downloads.find((d) => d.goos.value === os);
      if (!download) return [];
      return download.goarch;
    },
    [downloads],
  );

  const os = form.watch("os");
  const arch = form.watch("arch");

  // Watch the os and arch and make sure that when the OS changes, the arch that
  // is selected is compatible with the OS. If it is not, then select the first
  // compatible arch.
  useEffect(() => {
    const arches = getArchs(os);
    const selectedIncompatibleArch =
      arches.find((a) => a.value === arch) === undefined;
    if (arches.length > 0 && selectedIncompatibleArch) {
      form.setValue("arch", arches[0].value);
    }
  }, [os, arch]);

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="os"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <RadioGroup
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                  className={`grid grid-cols-${downloads.length} gap-4`}
                >
                  {downloads.map((download) => (
                    <OsOption
                      key={download.goos.value}
                      os={download.goos.value}
                      label={download.goos.display}
                      icon={goosToIcon(download.goos.value)}
                      field={field}
                    />
                  ))}
                </RadioGroup>
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        {os && (
          <FormField
            control={form.control}
            name="arch"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <RadioGroup
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                    className={`grid grid-cols-${getArchs(os).length} gap-4`}
                  >
                    {getArchs(os).map((download) => (
                      <ArchOption
                        key={download.value}
                        arch={download.value}
                        label={download.display}
                        field={field}
                      />
                    ))}
                  </RadioGroup>
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        )}

        <FormField
          control={form.control}
          name="generate_psk"
          render={({ field }) => (
            <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
              <div className="space-y-0.5">
                <FormLabel>Generate pre-shared key</FormLabel>
                <FormDescription>
                  Pre-shared keys are used to authenticate agents. If you do not
                  generate one now, you will need to generate one manually
                  later.
                </FormDescription>
              </div>
              <FormControl>
                <Switch
                  checked={field.value}
                  onCheckedChange={field.onChange}
                />
              </FormControl>
            </FormItem>
          )}
        />

        <div className="flex flex-row-reverse space-x-2">
          <Button type="submit">Next</Button>
        </div>
      </form>
    </Form>
  );
}

type OsOptionProps = {
  os: string;
  icon: React.ElementType;
  label: string;
  field: ControllerRenderProps<z.infer<typeof FormSchema>, "os">;
  className?: string;
};

function OsOption(props: OsOptionProps) {
  const Icon = props.icon;
  return (
    <FormItem className={props.className}>
      <FormControl>
        <RadioGroupItem className="sr-only" value={props.os} />
      </FormControl>
      <FormLabel className="font-normal">
        <Label
          className={cn(
            "flex flex-col items-center justify-between rounded-md border-2 border-muted bg-popover p-4 hover:bg-accent hover:text-accent-foreground",
            {
              "border-primary bg-accent": props.field.value === props.os,
            },
          )}
        >
          <Icon className="mb-3 h-6 w-6" />
          {props.label}
        </Label>
      </FormLabel>
    </FormItem>
  );
}

type ArchOptionProps = {
  arch: string;
  // icon: React.ElementType;
  label: string;
  field: ControllerRenderProps<z.infer<typeof FormSchema>, "arch">;
  className?: string;
};

function ArchOption(props: ArchOptionProps) {
  return (
    <FormItem className={props.className}>
      <FormControl>
        <RadioGroupItem className="sr-only" value={props.arch} />
      </FormControl>
      <FormLabel className="font-normal">
        <Label
          className={cn(
            "flex flex-col items-center justify-between rounded-md border-2 border-muted bg-popover p-4 hover:bg-accent hover:text-accent-foreground",
            {
              "border-primary bg-accent": props.field.value === props.arch,
            },
          )}
        >
          {props.label}
        </Label>
      </FormLabel>
    </FormItem>
  );
}
