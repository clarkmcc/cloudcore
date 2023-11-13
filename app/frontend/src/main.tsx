import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import "@fontsource/inter";
import { Auth0Provider } from "@auth0/auth0-react";
import { CssBaseline, CssVarsProvider } from "@mui/joy";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Auth0Provider
      domain={import.meta.env.VITE_AUTH0_DOMAIN}
      clientId={import.meta.env.VITE_AUTH0_CLIENT_ID}
      authorizationParams={{
        redirect_uri: window.location.origin,
      }}
      useRefreshTokens={true}
      cacheLocation="localstorage"
    >
      <CssVarsProvider>
        <CssBaseline />
        <App />
      </CssVarsProvider>
    </Auth0Provider>
  </React.StrictMode>,
);
