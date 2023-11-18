import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
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
import { Badge } from "@/components/ui/badge.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { Button } from "@/components/ui/button.tsx";
import { MoreHorizontal } from "lucide-react";

type HostWithGroups = Host & {
  groups: string[];
};

type HostsTableProps = {
  hosts: HostWithGroups[];
};

export function HostsTable(props: HostsTableProps) {
  const columns: ColumnDef<HostWithGroups>[] = [
    {
      accessorKey: "hostname",
      header: "Hostname",
    },
    {
      accessorKey: "groups",
      header: "Groups",
      cell: ({ row }) =>
        row.getValue<string[]>("groups").map((group) => (
          <Badge key={group} variant="secondary" className="mr-1">
            {group}
          </Badge>
        )),
    },
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
        const status = row.getValue<string>("status");
        return (
          <div>
            <div className="rounded-xl bg-green-500 w-3 h-3 mr-2 inline-block"></div>
            {status}
          </div>
        );
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const payment = row.original;

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
              <DropdownMenuItem
                onClick={() => navigator.clipboard.writeText(payment.id)}
              >
                Copy payment ID
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>View customer</DropdownMenuItem>
              <DropdownMenuItem>View payment details</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];
  const table = useReactTable({
    data: [
      {
        hostname: "Clark's MacBook Pro",
        groups: ["macOS", "laptop"],
        publicIpAddress: "35.133.10.18",
        privateIpAddress: "10.0.0.10",
        osName: "darwin",
        osVersion: "14.0",
        status: "Online",
      },
    ],
    columns,
    getCoreRowModel: getCoreRowModel(),
  });
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
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
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
    </div>
  );
}
