import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  getPaginationRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { Host } from "@/types";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { Button } from "@/components/ui/button.tsx";
import { MoreHorizontal } from "lucide-react";
import { Link } from "react-router-dom";
import { useProjectId } from "@/hooks/navigation.ts";
import { useMemo } from "react";
import { Badge } from "@/components/ui/badge.tsx";

type HostWithGroups = Host & {
  groups: string[];
};

type HostsTableProps = {
  hosts: HostWithGroups[];
};

export function HostsTable({ hosts }: HostsTableProps) {
  const [projectId] = useProjectId();
  const columns: ColumnDef<HostWithGroups>[] = useMemo(
    () => [
      {
        accessorKey: "hostname",
        header: "Hostname",
        cell: ({ row }) => {
          const id = row.original.id;
          const hostname = row.getValue<string>("hostname");
          return (
            <Link
              to={`/projects/${projectId}/hosts/${id}`}
              className="text-blue-500 hover:underline"
            >
              {hostname}
            </Link>
          );
        },
      },
      // {
      //   accessorKey: "groups",
      //   header: "Groups",
      //   cell: ({ row }) =>
      //     row.getValue<string[]>("groups")?.map((group) => (
      //       <Badge key={group} variant="secondary" className="mr-1">
      //         {group}
      //       </Badge>
      //     )),
      // },
      {
        accessorKey: "publicIpAddress",
        header: "Public IP",
      },
      {
        accessorKey: "privateIpAddress",
        header: "Private IP",
      },
      {
        accessorKey: "osName",
        header: "OS",
      },
      {
        accessorKey: "osVersion",
        header: "OS version",
      },
      {
        accessorKey: "status",
        header: "Status",
        cell: ({ row }) => {
          switch (row.original.online) {
            case true:
              return (
                <Badge className="text-xs" variant="ghost">
                  <div className="rounded-xl bg-green-500 w-3 h-3 mr-2 inline-block"></div>
                  Online
                </Badge>
              );
            case false:
              return (
                <Badge className="text-xs" variant="ghost">
                  <div className="rounded-xl bg-red-500 w-3 h-3 mr-2 inline-block"></div>
                  Offline
                </Badge>
              );
          }
        },
      },
      {
        id: "actions",
        cell: () => {
          return (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="h-8 w-8 p-0">
                  <span className="sr-only">Open menu</span>
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>Actions</DropdownMenuLabel>
                <DropdownMenuItem>Delete</DropdownMenuItem>
                <DropdownMenuItem>Disable</DropdownMenuItem>
                {/*<DropdownMenuSeparator />*/}
              </DropdownMenuContent>
            </DropdownMenu>
          );
        },
      },
    ],
    [projectId],
  );
  const table = useReactTable({
    data: hosts,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    state: {},
  });
  const rows = useMemo(() => table.getRowModel().rows, [table, hosts]);
  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                return (
                  <TableHead key={header.id}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                  </TableHead>
                );
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {rows.length ? (
            rows.map((row) => (
              <TableRow
                key={row.id}
                data-state={row.getIsSelected() && "selected"}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={columns.length} className="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
      {/*<div className="flex items-center justify-end space-x-2 p-2">*/}
      {/*  <Button*/}
      {/*    variant="outline"*/}
      {/*    size="sm"*/}
      {/*    onClick={() => table.previousPage()}*/}
      {/*    disabled={!table.getCanPreviousPage()}*/}
      {/*  >*/}
      {/*    Previous*/}
      {/*  </Button>*/}
      {/*  <Button*/}
      {/*    variant="outline"*/}
      {/*    size="sm"*/}
      {/*    onClick={() => table.nextPage()}*/}
      {/*    disabled={!table.getCanNextPage()}*/}
      {/*  >*/}
      {/*    Next*/}
      {/*  </Button>*/}
      {/*</div>*/}
    </div>
  );
}
