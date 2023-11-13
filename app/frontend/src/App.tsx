import "./App.css";
import { useAuth0, withAuthenticationRequired } from "@auth0/auth0-react";
import Sidebar from "./components/Sidebar.tsx";
import Box from "@mui/joy/Box";
import {
  Outlet,
  RootRoute,
  Route,
  Router,
  RouterProvider,
} from "@tanstack/react-router";
import { Home } from "./routes";
import { Projects } from "./routes/Projects.tsx";
import { Agents, AllHosts, Deploy, Groups } from "./routes/hosts";

const rootRoute = new RootRoute({
  component: withAuthenticationRequired(Root),
});

const homeRoute = new Route({
  getParentRoute: () => rootRoute,
  component: Home,
  path: "/",
});

const projectsRoute = new Route({
  getParentRoute: () => rootRoute,
  component: Projects,
  path: "/projects",
});

const hostsRoute = new Route({
  getParentRoute: () => rootRoute,
  path: "/hosts",
});

const hostsAllRoute = new Route({
  getParentRoute: () => hostsRoute,
  path: "/",
  component: AllHosts,
});

const hostsGroupsRoute = new Route({
  getParentRoute: () => hostsRoute,
  path: "/groups",
  component: Groups,
});

const hostsAgentsRoute = new Route({
  getParentRoute: () => hostsRoute,
  path: "/agents",
  component: Agents,
});

const hostsDeployRoute = new Route({
  getParentRoute: () => hostsRoute,
  path: "/deploy",
  component: Deploy,
});

const settingsRoute = new Route({
  getParentRoute: () => rootRoute,
  path: "/settings",
  component: () => <p>Settings route</p>,
});

const routeTree = rootRoute.addChildren([
  homeRoute,
  projectsRoute,
  hostsRoute,
  hostsAllRoute,
  hostsGroupsRoute,
  hostsAgentsRoute,
  hostsDeployRoute,
  settingsRoute,
]);
const router = new Router({ routeTree });

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

function Root() {
  return (
    <Box sx={{ display: "flex", minHeight: "100dvh" }}>
      <Sidebar />
      <Outlet />
    </Box>
  );
}

function App() {
  const auth0 = useAuth0();
  return <RouterProvider router={router} context={{ auth0 }} />;
}

export default App;
