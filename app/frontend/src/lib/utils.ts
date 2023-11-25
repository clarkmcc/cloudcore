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
export function goosToString(os: string): string {
  switch (os) {
    case "darwin":
      return "macOS";
    case "linux":
      return "Linux";
    case "windows":
      return "Windows";
    default:
      return os;
  }
}

export function osToColorClasses(os: string): string {
  switch (os) {
    case "darwin":
      return "bg-blue-500";
    case "ubuntu":
      return "bg-[#e95420]";
    case "debian":
      return "bg-[#d70a53]";
    case "redhat":
      return "bg-[#ee0000]";
    default:
      return "";
  }
}

export function osToString(os: string): string {
  switch (os) {
    case "debian":
      return "Debian";
    case "ubuntu":
      return "Ubuntu";
    case "darwin":
      return "macOS";
    case "redhat":
      return "Red Hat";
    default:
      return goosToString(os);
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
 */
export function goarchToString(goarch: string, goarm?: string): string {
  switch (goarch) {
    case "arm64":
    case "aarch64":
      return "64-bit ARM";
    case "x86_64":
    case "amd64":
      return "64-bit x86";
    case "i386":
    case "386":
      return "32-bit x86";
    case "arm":
      switch (goarm) {
        case "5":
          return "ARMv5";
        default:
          return "Unknown ARM";
      }
    default:
      return "Unknown";
  }
}
