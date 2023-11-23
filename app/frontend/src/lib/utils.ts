import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { AppleIcon } from "@/components/icons/apple-icon.tsx";
import { LinuxIcon } from "@/components/icons/linux-icon.tsx";
import { WindowsIcon } from "@/components/icons/windows-icon.tsx";
import React from "react";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Returns a human-readable name for the given OS.
 * @param os
 */
export function getOsName(os: string): string {
  switch (os) {
    case "darwin":
      return "macOS";
    case "debian":
      return "Debian";
    default:
      return os;
  }
}

/**
 * Takes a GOOS value and returns the appropriate icon.
 * @param goos
 */
export function goosToIcon(goos: string): React.ElementType {
  switch (goos) {
    case "darwin":
      return AppleIcon;
    case "linux":
      return LinuxIcon;
    case "windows":
      return WindowsIcon;
    default:
      throw new Error(`Unknown goos: ${goos}`);
  }
}

/**
 * Returns a human-readable name for the given kernel architecture.
 * @param architecture
 */
export function getArchitecture(architecture?: string): string {
  switch (architecture) {
    case "arm64":
    case "aarch64":
      return "64-bit ARM";
    case "x86_64":
      return "64-bit x86";
    case "i386":
      return "32-bit x86";
    default:
      return "Unknown";
  }
}
