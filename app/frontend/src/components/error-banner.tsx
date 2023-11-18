import { ApolloError } from "@apollo/client";

type HostsTableProps = {
  error: ApolloError;
};

export function ErrorBanner(props: HostsTableProps) {
  return (
    <div className="col-span-12 md:col-span-12 lg:col-span-6 border border-red-800 p-4 rounded-lg dark:bg-red-950 bg-red-100">
      <div className="text-lg font-medium dark:text-red-50 text-red-900">
        Oops, that's not right
      </div>
      <div className="text-xs dark:text-red-50 text-red-900">
        There was a problem, please try refreshing this page.
      </div>
      <div className="border border-red-300 bg-red-200 dark:bg-red-900 mt-4 py-2 px-4 rounded-md dark:text-red-50 text-xs text-red-900 font-mono">
        {props.error?.message}
      </div>
    </div>
  );
}
