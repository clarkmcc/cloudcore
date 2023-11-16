import * as React from "react";
import GlobalStyles from "@mui/joy/GlobalStyles";
import Avatar from "@mui/joy/Avatar";
import Box from "@mui/joy/Box";
import Divider from "@mui/joy/Divider";
import IconButton from "@mui/joy/IconButton";
import List from "@mui/joy/List";
import ListItem from "@mui/joy/ListItem";
import ListItemButton, { listItemButtonClasses } from "@mui/joy/ListItemButton";
import ListItemContent from "@mui/joy/ListItemContent";
import Typography from "@mui/joy/Typography";
import Sheet from "@mui/joy/Sheet";
import HomeRoundedIcon from "@mui/icons-material/HomeRounded";
import SettingsRoundedIcon from "@mui/icons-material/SettingsRounded";
import LogoutRoundedIcon from "@mui/icons-material/LogoutRounded";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import { AccountTree, Add, CloudCircle, Computer } from "@mui/icons-material";
import { useAuth0 } from "@auth0/auth0-react";
import Select from "@mui/joy/Select";
import Option from "@mui/joy/Option";
import ColorSchemeToggle from "./ColorSchemeToggle.tsx";
import { useCallback, useState } from "react";
import { ListItemDecorator } from "@mui/joy";
import { useNavigate } from "@tanstack/react-router";

function Toggler({
  renderToggle,
  children,
  rootPage,
}: {
  children: React.ReactNode;
  rootPage: string;
  renderToggle: (params: {
    open: boolean;
    setOpen: React.Dispatch<React.SetStateAction<boolean>>;
  }) => React.ReactNode;
}) {
  const [open, setOpen] = useState(
    window.location.pathname.startsWith(rootPage),
  );
  return (
    <React.Fragment>
      {renderToggle({ open, setOpen })}
      <Box
        sx={{
          display: "grid",
          gridTemplateRows: open ? "1fr" : "0fr",
          transition: "0.2s ease",
          "& > *": {
            overflow: "hidden",
          },
        }}
      >
        {children}
      </Box>
    </React.Fragment>
  );
}

export default function Sidebar() {
  const { user, logout, loginWithRedirect } = useAuth0();
  const navigate = useNavigate();

  const buildSidebarButton = useCallback((page: string) => {
    return {
      selected: window.location.pathname === page,
      onClick: () => navigate({ params: {}, to: page }),
    };
  }, []);

  return (
    <Sheet
      className="Sidebar"
      sx={{
        position: {
          xs: "fixed",
          md: "sticky",
        },
        transform: {
          xs: "translateX(calc(100% * (var(--SideNavigation-slideIn, 0) - 1)))",
          md: "none",
        },
        transition: "transform 0.4s, width 0.4s",
        zIndex: 1,
        height: "100dvh",
        width: "var(--Sidebar-width)",
        top: 0,
        p: 2,
        flexShrink: 0,
        display: "flex",
        flexDirection: "column",
        gap: 2,
        borderRight: "1px solid",
        borderColor: "divider",
      }}
    >
      <GlobalStyles
        styles={(theme) => ({
          ":root": {
            "--Sidebar-width": "220px",
            [theme.breakpoints.up("lg")]: {
              "--Sidebar-width": "240px",
            },
          },
        })}
      />
      <Box
        className="Sidebar-overlay"
        sx={{
          position: "fixed",
          zIndex: 9998,
          top: 0,
          left: 0,
          width: "100vw",
          height: "100vh",
          opacity: "var(--SideNavigation-slideIn)",
          backgroundColor: "var(--joy-palette-background-backdrop)",
          transition: "opacity 0.4s",
          transform: {
            xs: "translateX(calc(100% * (var(--SideNavigation-slideIn, 0) - 1) + var(--SideNavigation-slideIn, 0) * var(--Sidebar-width, 0px)))",
            lg: "translateX(-100%)",
          },
        }}
        // onClick={() => closeSidebar()}
      />
      <Box sx={{ display: "flex", gap: 1, alignItems: "center" }}>
        <IconButton variant="soft" color="primary" size="sm">
          <CloudCircle />
        </IconButton>
        <Typography level="title-lg">cloudcore</Typography>
        <ColorSchemeToggle sx={{ ml: "auto" }} />
      </Box>
      {/*<Input*/}
      {/*  size="sm"*/}
      {/*  startDecorator={<SearchRoundedIcon />}*/}
      {/*  placeholder="Search"*/}
      {/*/>*/}
      <Select placeholder="Select a project">
        <Option value="">
          <ListItemDecorator>
            <Add />
          </ListItemDecorator>
          New project
        </Option>
        <Option value="1">Project 1</Option>
        <Option value="2">Project 2</Option>
      </Select>
      <Box
        sx={{
          minHeight: 0,
          overflow: "hidden auto",
          flexGrow: 1,
          display: "flex",
          flexDirection: "column",
          [`& .${listItemButtonClasses.root}`]: {
            gap: 1.5,
          },
        }}
      >
        <List
          size="sm"
          sx={{
            gap: 1,
            "--List-nestedInsetStart": "30px",
            "--ListItem-radius": (theme) => theme.vars.radius.sm,
          }}
        >
          <ListItem>
            <ListItemButton {...buildSidebarButton("/")}>
              <HomeRoundedIcon />
              <ListItemContent>
                <Typography level="title-sm">Home</Typography>
              </ListItemContent>
            </ListItemButton>
          </ListItem>

          <ListItem>
            <ListItemButton {...buildSidebarButton("/projects")}>
              <AccountTree />
              <ListItemContent>
                <Typography level="title-sm">Projects</Typography>
              </ListItemContent>
            </ListItemButton>
          </ListItem>

          {/*<ListItem>*/}
          {/*  <ListItemButton selected>*/}
          {/*    <ShoppingCartRoundedIcon />*/}
          {/*    <ListItemContent>*/}
          {/*      <Typography level="title-sm">Orders</Typography>*/}
          {/*    </ListItemContent>*/}
          {/*  </ListItemButton>*/}
          {/*</ListItem>*/}

          <ListItem nested>
            <Toggler
              rootPage={"/hosts"}
              renderToggle={({ open, setOpen }) => (
                <ListItemButton onClick={() => setOpen(!open)}>
                  <Computer />
                  <ListItemContent>
                    <Typography level="title-sm">Hosts</Typography>
                  </ListItemContent>
                  <KeyboardArrowDownIcon
                    sx={{ transform: open ? "rotate(180deg)" : "none" }}
                  />
                </ListItemButton>
              )}
            >
              <List sx={{ gap: 0.5 }}>
                <ListItem sx={{ mt: 0.5 }}>
                  <ListItemButton {...buildSidebarButton("/hosts")}>
                    All hosts
                  </ListItemButton>
                </ListItem>
                <ListItem>
                  <ListItemButton {...buildSidebarButton("/hosts/groups")}>
                    Groups
                  </ListItemButton>
                </ListItem>
                <ListItem>
                  <ListItemButton {...buildSidebarButton("/hosts/agents")}>
                    Agents
                  </ListItemButton>
                </ListItem>
                <ListItem>
                  <ListItemButton {...buildSidebarButton("/hosts/deploy")}>
                    Deploy
                  </ListItemButton>
                </ListItem>
              </List>
            </Toggler>
          </ListItem>

          {/*<ListItem>*/}
          {/*  <ListItemButton*/}
          {/*    role="menuitem"*/}
          {/*    component="a"*/}
          {/*    href="/joy-ui/getting-started/templates/messages/"*/}
          {/*  >*/}
          {/*    <QuestionAnswerRoundedIcon />*/}
          {/*    <ListItemContent>*/}
          {/*      <Typography level="title-sm">Messages</Typography>*/}
          {/*    </ListItemContent>*/}
          {/*    <Chip size="sm" color="primary" variant="solid">*/}
          {/*      4*/}
          {/*    </Chip>*/}
          {/*  </ListItemButton>*/}
          {/*</ListItem>*/}

          {/*<ListItem nested>*/}
          {/*  <Toggler*/}
          {/*    renderToggle={({ open, setOpen }) => (*/}
          {/*      <ListItemButton onClick={() => setOpen(!open)}>*/}
          {/*        <GroupRoundedIcon />*/}
          {/*        <ListItemContent>*/}
          {/*          <Typography level="title-sm">Users</Typography>*/}
          {/*        </ListItemContent>*/}
          {/*        <KeyboardArrowDownIcon*/}
          {/*          sx={{ transform: open ? "rotate(180deg)" : "none" }}*/}
          {/*        />*/}
          {/*      </ListItemButton>*/}
          {/*    )}*/}
          {/*  >*/}
          {/*    <List sx={{ gap: 0.5 }}>*/}
          {/*      <ListItem sx={{ mt: 0.5 }}>*/}
          {/*        <ListItemButton*/}
          {/*          role="menuitem"*/}
          {/*          component="a"*/}
          {/*          href="/joy-ui/getting-started/templates/profile-dashboard/"*/}
          {/*        >*/}
          {/*          My profile*/}
          {/*        </ListItemButton>*/}
          {/*      </ListItem>*/}
          {/*      <ListItem>*/}
          {/*        <ListItemButton>Create a new user</ListItemButton>*/}
          {/*      </ListItem>*/}
          {/*      <ListItem>*/}
          {/*        <ListItemButton>Roles & permission</ListItemButton>*/}
          {/*      </ListItem>*/}
          {/*    </List>*/}
          {/*  </Toggler>*/}
          {/*</ListItem>*/}
        </List>

        <List
          size="sm"
          sx={{
            mt: "auto",
            flexGrow: 0,
            "--ListItem-radius": (theme) => theme.vars.radius.sm,
            "--List-gap": "8px",
            mb: 2,
          }}
        >
          {/*<ListItem>*/}
          {/*  <ListItemButton>*/}
          {/*    <SupportRoundedIcon />*/}
          {/*    Support*/}
          {/*  </ListItemButton>*/}
          {/*</ListItem>*/}
          <ListItem>
            <ListItemButton {...buildSidebarButton("/settings")}>
              <SettingsRoundedIcon />
              Settings
            </ListItemButton>
          </ListItem>
        </List>
        {/*<Card*/}
        {/*  invertedColors*/}
        {/*  variant="soft"*/}
        {/*  color="warning"*/}
        {/*  size="sm"*/}
        {/*  sx={{ boxShadow: "none" }}*/}
        {/*>*/}
        {/*  <Stack*/}
        {/*    direction="row"*/}
        {/*    justifyContent="space-between"*/}
        {/*    alignItems="center"*/}
        {/*  >*/}
        {/*    <Typography level="title-sm">Used space</Typography>*/}
        {/*    <IconButton size="sm">*/}
        {/*      <CloseRoundedIcon />*/}
        {/*    </IconButton>*/}
        {/*  </Stack>*/}
        {/*  <Typography level="body-xs">*/}
        {/*    Your team has used 80% of your available space. Need more?*/}
        {/*  </Typography>*/}
        {/*  <LinearProgress*/}
        {/*    variant="outlined"*/}
        {/*    value={80}*/}
        {/*    determinate*/}
        {/*    sx={{ my: 1 }}*/}
        {/*  />*/}
        {/*  <Button size="sm" variant="solid">*/}
        {/*    Upgrade plan*/}
        {/*  </Button>*/}
        {/*</Card>*/}
      </Box>
      <Divider />
      <Box sx={{ display: "flex", gap: 1, alignItems: "center" }}>
        <Avatar
          variant="outlined"
          size="sm"
          src={user?.picture}
          onClick={() => loginWithRedirect()}
        />
        <Box sx={{ minWidth: 0, flex: 1 }}>
          <Typography level="title-sm">
            {user?.given_name} {user?.family_name}
          </Typography>
          <Typography level="body-xs" sx={{ overflow: "hidden" }}>
            {user?.email}
          </Typography>
        </Box>
        <IconButton
          size="sm"
          variant="plain"
          color="neutral"
          onClick={() => logout()}
        >
          <LogoutRoundedIcon />
        </IconButton>
      </Box>
    </Sheet>
  );
}
