import React, { forwardRef } from "react";
import { cn } from "@/lib/utils.ts";

export const WindowsIcon = forwardRef<
  SVGSVGElement,
  React.SVGProps<SVGSVGElement>
>((props, ref) => (
  <svg
    ref={ref}
    {...props}
    className={cn("fill-black dark:fill-white", props.className)}
    overflow="visible"
    viewBox="0 0 21 21"
  >
    <path fill="#f35325" d="M0 0h10v10H0z" />
    <path fill="#81bc06" d="M11 0h10v10H11z" />
    <path fill="#05a6f0" d="M0 11h10v10H0z" />
    <path fill="#ffba08" d="M11 11h10v10H11z" />
  </svg>
));
