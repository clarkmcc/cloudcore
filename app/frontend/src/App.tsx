import "./App.css";
import {
  createBrowserRouter,
  createRoutesFromElements,
  Outlet,
  Route,
  RouterProvider,
} from "react-router-dom";
import { Sidebar } from "@/components/sidebar.tsx";
import { loadProject, ProjectHome } from "@/pages/project.tsx";
import { withAuthenticationRequired } from "@auth0/auth0-react";
import { useEnsureUser } from "@/hooks/user.ts";
import { AuthenticatedApolloProvider } from "@/queries/client.tsx";
import { ProjectProvider } from "@/context/project.tsx";
import { SettingsPage } from "@/pages/settings.tsx";
import { HostsPage } from "@/pages/hosts.tsx";
import { HostGroupsPage } from "@/pages/host-groups.tsx";
import { HostDetails } from "@/pages/host-details.tsx";

function Container() {
  useEnsureUser();
  return (
    <div className="flex flex-row h-full">
      <Sidebar className="max-w-[200px]"/>
      <div className="border-l w-full">
        <Outlet />
      </div>
    </div>
  );
}

const AuthenticatedContainer = withAuthenticationRequired(() => {
  return (
    <AuthenticatedApolloProvider>
      <ProjectProvider>
        <Container />
      </ProjectProvider>
    </AuthenticatedApolloProvider>
  );
});

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<AuthenticatedContainer />}>
      <Route path="/projects/:projectId">
        <Route index element={<ProjectHome />} loader={loadProject} />
        <Route path="/projects/:projectId/hosts" element={<HostsPage />} />
        <Route
          path="/projects/:projectId/hosts/:hostId"
          element={<HostDetails />}
        />
        <Route
          path="/projects/:projectId/hosts/groups"
          element={<HostGroupsPage />}
          loader={loadProject}
        />
        <Route
          path="/projects/:projectId/settings"
          element={<SettingsPage />}
        />
      </Route>
    </Route>,
  ),
);

function App() {
  return <RouterProvider router={router} />;
}

export default App;
