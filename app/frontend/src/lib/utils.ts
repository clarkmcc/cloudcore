import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
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
 * Returns a human-readable name for the given kernel architecture.
 * @param architecture
 */
export function getArchitecture(architecture?: string): string {
  switch (architecture) {
    case "arm64":
    case "aarch64":
      return "64-bit ARM"
    case "x86_64":
      return "64-bit x86"
    case "i386":
      return "32-bit x86"
    default:
      return "Unknown"
  }
}
