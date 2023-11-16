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

function Container() {
  return (
    <div className="grid grid-cols-12 h-full">
      <div className="col-span-6 sm:col-span-3 md:col-span-2 border-r">
        <Sidebar />
      </div>
      <div className="col-span-6 sm:col-span-9 md:col-span-10">
        <Outlet />
      </div>
    </div>
  );
}

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route element={<Container />}>
      <Route path="/projects/:projectId">
        <Route index element={<ProjectHome />} loader={loadProject} />
        <Route path="/projects/:projectId/hosts" element={<p>hosts</p>} />
        <Route
          path="/projects/:projectId/hosts/groups"
          element={<p>host groups</p>}
        />
        <Route path="/projects/:projectId/settings" element={<p>settings</p>} />
      </Route>
    </Route>,
  ),
);

function App() {
  return <RouterProvider router={router} />;
}

export default App;
