import { Button } from "@/components/ui/button.tsx";
import { PageHeader } from "@/components/page-header.tsx";

export function SettingsPage() {
  return (
    <>
      <PageHeader title="Settings" subtitle="Project-level settings" />
      <div className="pl-7">
        <div className="grid grid-cols-12 h-full">
          <div className="col-span-12 md:col-span-12 lg:col-span-6 border border-red-800 p-4 rounded-lg dark:bg-red-950 bg-red-100">
            <div className="flex">
              <div>
                <div className="text-lg font-medium dark:text-red-50 text-red-900">
                  Delete project
                </div>
                <div className="text-xs dark:text-red-50 text-red-900">
                  This action is irreversible. All hosts and groups will be
                  deleted.
                </div>
              </div>
              <div className="flex-grow"></div>
              <div className="flex flex-column items-center">
                <Button variant="destructive">Delete project</Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
