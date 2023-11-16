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
import { Breadcrumbs, Link } from "@mui/joy";
import HomeRoundedIcon from "@mui/icons-material/HomeRounded";
import Typography from "@mui/joy/Typography";
import { ChevronRightRounded } from "@mui/icons-material";
import { useSyncAuth0User } from "./hooks/auth.ts";

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
  useSyncAuth0User();

  return (
    <Box sx={{ display: "flex", minHeight: "100dvh" }}>
      <Sidebar />
      <Box
        component="main"
        className="MainContent"
        sx={{
          px: {
            xs: 2,
            md: 6,
          },
          pt: {
            xs: "calc(12px + var(--Header-height))",
            sm: "calc(12px + var(--Header-height))",
            md: 3,
          },
          pb: {
            xs: 2,
            sm: 2,
            md: 3,
          },
          flex: 1,
          display: "flex",
          flexDirection: "column",
          minWidth: 0,
          height: "100dvh",
          gap: 1,
        }}
      >
        <Box sx={{ display: "flex", alignItems: "center" }}>
          <Breadcrumbs
            size="sm"
            aria-label="breadcrumbs"
            separator={<ChevronRightRounded />}
            sx={{ pl: 0 }}
          >
            <Link
              underline="none"
              color="neutral"
              href="#some-link"
              aria-label="Home"
            >
              <HomeRoundedIcon />
            </Link>
            <Link
              underline="hover"
              color="neutral"
              href="#some-link"
              fontSize={12}
              fontWeight={500}
            >
              Project 1
            </Link>
            <Typography color="primary" fontWeight={500} fontSize={12}>
              Home
            </Typography>
          </Breadcrumbs>
        </Box>
        <Box
          sx={{
            display: "flex",
            my: 1,
            gap: 1,
            flexDirection: { xs: "column", sm: "row" },
            alignItems: { xs: "start", sm: "center" },
            flexWrap: "wrap",
            justifyContent: "space-between",
          }}
        >
          <Typography level="h2">Home</Typography>
        </Box>
        <Outlet />
      </Box>
    </Box>
  );
}

function App() {
  const auth0 = useAuth0();
  return <RouterProvider router={router} context={{ auth0 }} />;
}

export default App;
