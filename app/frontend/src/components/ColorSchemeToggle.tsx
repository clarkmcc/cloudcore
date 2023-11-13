import DarkModeRoundedIcon from "@mui/icons-material/DarkModeRounded";
import LightModeIcon from "@mui/icons-material/LightMode";
import { IconButtonProps, useColorScheme } from "@mui/joy";
import { useEffect, useState } from "react";
import IconButton from "@mui/joy/IconButton";

export default function ColorSchemeToggle({
  onClick,
  sx,
  ...props
}: IconButtonProps) {
  const { mode, setMode } = useColorScheme();
  const [mounted, setMounted] = useState(false);
  useEffect(() => {
    setMounted(true);
  }, []);
  if (!mounted) {
    return (
      <IconButton
        size="sm"
        variant="outlined"
        color="neutral"
        {...props}
        sx={sx}
        disabled
      />
    );
  }
  return (
    <IconButton
      id="toggle-mode"
      size="sm"
      variant="outlined"
      color="neutral"
      {...props}
      onClick={(event) => {
        if (mode === "light") {
          setMode("dark");
        } else {
          setMode("light");
        }
        onClick?.(event);
      }}
      sx={[
        {
          "& > *:first-child": {
            display: mode === "dark" ? "none" : "initial",
          },
          "& > *:last-child": {
            display: mode === "light" ? "none" : "initial",
          },
        },
        ...(Array.isArray(sx) ? sx : [sx]),
      ]}
    >
      <DarkModeRoundedIcon />
      <LightModeIcon />
    </IconButton>
  );
}
