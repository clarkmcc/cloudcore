import { gql } from "@apollo/client";

export const MUTATION_CREATE_PROJECT = gql`
  mutation CreateProject($name: String!, $description: String) {
    projectCreate(name: $name, description: $description) {
      project {
        id
        name
      }
      allProjects {
        id
        name
      }
    }
  }
`;

export const QUERY_PROJECT_METRICS = gql`
  query ProjectMetrics($projectId: String!) {
    projectMetrics(projectId: $projectId) {
      offlineHosts
      onlineHosts
      hostsByOsName {
        osName
        count
      }
      totalAgents
      totalHosts
    }
  }
`;
