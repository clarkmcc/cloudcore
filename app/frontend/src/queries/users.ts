import { gql } from "@apollo/client";

/**
 * Ensures that the currently logged-in user (through Auth0) exists
 * in the database, and that we have a tenant and project for them.
 * The project and tenant are returned.
 */
export const MUTATION_ENSURE_USER = gql`
  mutation EnsureUser {
    ensureUser {
      id
      name
    }
  }
`;
