import { AgentEventType } from "@/types/enums.ts";

export type Project = {
  id: string;
  name: string;
};

export type Host = {
  id: string;
  createdAt: string;
  updatedAt: string;
  lastHeartbeatTimestamp: string;
  identifier: string;
  online: boolean;
  hostname?: string;
  publicIpAddress?: string;
  privateIpAddress?: string;
  osName?: string;
  osFamily?: string;
  osVersion?: string;
  kernelArchitecture?: string;
  kernelVersion?: string;
  cpuModel?: string;
  cpuCores?: number;
};

export type AgentEvent = {
  id: string;
  createdAt: string;
  type: AgentEventType;
  message: string;
};
