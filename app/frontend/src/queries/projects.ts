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
