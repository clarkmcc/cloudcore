import React, { useCallback } from "react";
import { cn } from "@/lib/utils.ts";
import { Button, ButtonProps } from "@/components/ui/button.tsx";
import { Cloud, Component, Home, Laptop, LogOut, Settings } from "lucide-react";
import { ThemeModeToggle } from "@/components/theme-mode-toggle.tsx";
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/components/ui/avatar.tsx";
import { useAuth0 } from "@auth0/auth0-react";
import { ProjectSelector } from "@/components/project-selector.tsx";
import { useProjectNavigate } from "@/hooks/navigation.ts";
import { matchPath } from "react-router-dom";

interface SidebarProps extends React.HTMLAttributes<HTMLDivElement> {}

export function Sidebar({ className }: SidebarProps) {
  const { user, logout } = useAuth0();
  const navigate = useProjectNavigate();

  const link = useCallback((page: string): Partial<ButtonProps> => {
    // Returns true if the sidebar link is selected. Has a special case for
    // the home page, which is selected if the current path is just the root
    // of a selected project.
    const isSelected =
      page !== "/"
        ? window.location.pathname.endsWith(page)
        : matchPath("/projects/:projectId", window.location.pathname);
    return {
      className: "w-full justify-start",
      variant: isSelected ? "secondary" : "ghost",
      onClick: () => navigate(page),
    };
  }, []);

  function handleLogout() {
    logout({
      logoutParams: {
        returnTo: window.location.origin,
      },
    }).catch(console.error);
  }

  return (
    <div className={cn("pb-2", className, "h-full")}>
      <div className="flex flex-col h-full justify-between">
        {/*Top Sidebar*/}
        <div className="space-y-4">
          <div className="px-3 py-2 space-y-3">
            <div className="flex space-x-2">
              <Cloud className="w-8 h-8" />
              <h1 className="text-2xl font-medium">cloudcore</h1>
            </div>
            <div className="flex space-x-3">
              <ProjectSelector />
              <div>
                <ThemeModeToggle />
              </div>
            </div>
            <div className="space-y-1">
              <Button {...link("/")}>
                <Home className="mr-2 h-4 w-4" />
                Home
              </Button>
              <Button {...link("/hosts")}>
                <Laptop className="mr-2 h-4 w-4" />
                Hosts
              </Button>
              <Button {...link("/hosts/groups")}>
                <Component className="mr-2 h-4 w-4" />
                Groups
              </Button>
            </div>
          </div>
        </div>

        {/* Bottom Sidebar */}
        <div className="px-3 py-2">
          <div className="pb-2">
            <Button {...link("/settings")}>
              <Settings className="mr-2 h-4 w-4" />
              Settings
            </Button>
          </div>

          {/* Profile */}
          <div className="border-t pt-4">
            <div className="flex items-center justify-between">
              {" "}
              {/* Flex container */}
              {/* Avatar and User Info */}
              <div className="flex items-center space-x-4 max-w-[calc(100%-3rem)]">
                {" "}
                {/* Adjust max-width as needed */}
                <Avatar>
                  <AvatarImage src={user?.picture} />
                  <AvatarFallback>CN</AvatarFallback>
                </Avatar>
                <div className="overflow-hidden">
                  <p className="text-sm font-medium leading-none truncate">
                    {user?.given_name} {user?.family_name}
                  </p>
                  <p className="text-sm text-muted-foreground truncate">
                    {user?.email}
                  </p>
                </div>
              </div>
              {/* Logout Button */}
              <div className="flex-shrink-0">
                {" "}
                {/* Prevent button from shrinking */}
                <Button variant="secondary" size="icon" onClick={handleLogout}>
                  <LogOut className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
