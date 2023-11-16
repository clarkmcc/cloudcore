import { useAuth0 } from "@auth0/auth0-react";
import { useEffect } from "react";

export function useSyncAuth0User() {
  const { getAccessTokenSilently, getIdTokenClaims } = useAuth0();

  const syncAuth0User = async () => {
    const accessToken = await getAccessTokenSilently();
    const idToken = await getIdTokenClaims();
    // const user = {
    //   accessToken,
    //   idToken: idToken?.__raw,
    // };
  };

  useEffect(() => {
    syncAuth0User().then().catch(console.error);
  }, []);
}
