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

export const QUERY_HOST_DETAILS = gql`
  query Host($projectId: String!, $hostId: String!) {
    host(projectId: $projectId, hostId: $hostId) {
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
      events {
        id
        createdAt
        type
        message
      }
    }
  }
`;
