export type Project = {
  id: string;
  name: string;
};

export type Host = {
  id: string;
  createdAt: string;
  updatedAt: string;
  identifier: string;
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
