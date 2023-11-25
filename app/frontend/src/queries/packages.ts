import { gql } from "@apollo/client";

export const QUERY_LIST_LATEST_PACKAGES = gql`
  query ListLatestPackages {
    packages {
      goos
      goarch
      goarm
    }
  }
`;
