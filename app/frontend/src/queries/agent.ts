import { gql } from "@apollo/client";

export const MUTATION_BUILD_DEPLOY_AGENT_COMMAND = gql`
  mutation BuildDeployAgentCommand(
    $projectId: String!
    $goos: String!
    $goarch: String!
    $generatePsk: Boolean!
  ) {
    buildDeployAgentCommand(
      projectId: $projectId
      goos: $goos
      goarch: $goarch
      generatePsk: $generatePsk
    )
  }
`;
