import { Sidebar } from "@/components/sidebar.tsx";
import { createBrowserRouter, Outlet } from "react-router-dom";

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

// const AuthenticatedContainer = withAuthenticationRequired(Container);

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Container />,
  },
  {
    path: "/projects/:projectId",
    element: <p>project</p>,
  },
]);

// const rootRoute = new RootRoute({
//   component: withAuthenticationRequired(Container),
// });
//
// const projectsRoute = new Route({
//   getParentRoute: () => rootRoute,
//   path: "projects",
// });
//
// const projectInstanceRoute = new Route({
//   getParentRoute: () => projectsRoute,
//   path: "$projectId",
//   component: () => <p>home</p>,
// });
//
// // const homeRoute = new Route({
// //   getParentRoute: () => projectInstanceRoute,
// //   component: () => <p>home</p>,
// //   path: "/",
// // });
//
// const hostsRoute = new Route({
//   getParentRoute: () => projectInstanceRoute,
//   path: "/hosts",
// });
//
// const allHostsRoute = new Route({
//   getParentRoute: () => projectInstanceRoute,
//   component: () => <div className="bg-red-500">hosts</div>,
//   path: "/",
// });
//
// const groupsRoute = new Route({
//   getParentRoute: () => projectInstanceRoute,
//   component: () => <div className="bg-red-500">host groups</div>,
//   path: "/groups",
// });
//
// const settingsRoute = new Route({
//   getParentRoute: () => projectInstanceRoute,
//   component: () => <div className="bg-red-500">settings</div>,
//   path: "/settings",
// });
//
// const routeTree = rootRoute.addChildren([
//   projectInstanceRoute,
//   hostsRoute,
//   allHostsRoute,
//   groupsRoute,
//   settingsRoute,
// ]);
// export const router = new Router({
//   routeTree,
// });

// declare module "@tanstack/react-router" {
//   interface Register {
//     // This infers the type of our router and registers it across your entire project
//     router: typeof router;
//   }
// }
