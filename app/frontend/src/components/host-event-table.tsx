import { AgentEvent } from "@/types";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx";
import moment from "moment";
import { cn } from "@/lib/utils.ts";
import { AgentEventType } from "@/types/enums.ts";

type HostEventTableProps = {
  events: AgentEvent[];
};

export function HostEventTable({ events }: HostEventTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[250px] hidden lg:block">Timestamp</TableHead>
          <TableHead className="w-[150px]">Type</TableHead>
          <TableHead>Message</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {events.map((event) => (
          <TableRow key={event.id}>
            <TableCell className="font-medium hidden lg:block">
              {moment(event.createdAt).format("llll")}
            </TableCell>
            <TableCell>
              <code
                className={cn(
                  "relative rounded px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold",
                  AgentEventType.toBackgroundClass(event.type),
                )}
              >
                {event.type}
              </code>
            </TableCell>
            <TableCell>
              {event.message}{" "}
              <span className="text-gray-500">
                {moment(event.createdAt).fromNow()}
              </span>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
