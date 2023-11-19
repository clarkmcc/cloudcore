import { gql } from "@apollo/client";

export const QUERY_HOSTS_LIST = gql`
  query HostsList($projectId: String!) {
    hosts(projectId: $projectId) {
      id
      identifier
      lastHeartbeatTimestamp
      hostname
      online
      createdAt
      updatedAt
      publicIpAddress
      privateIpAddress
      osName
      osFamily
      osVersion
      kernelArchitecture
      kernelVersion
      cpuModel
      cpuCores
    }
  }
`;
