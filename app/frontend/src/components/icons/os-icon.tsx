import React, { forwardRef } from "react";
import { AppleIcon } from "@/components/icons/apple-icon.tsx";
import { DebianIcon } from "@/components/icons/debian-icon.tsx";

type OsIconProps = React.SVGProps<SVGSVGElement> & {
  osName: string;
}

export const OsIcon = forwardRef<SVGSVGElement, OsIconProps>(({osName, ...rest}, ref) => {
  switch (osName) {
    case "darwin":
      // @ts-expect-error not sure what the issue is here
      return <AppleIcon ref={ref} {...rest}/>
    case "debian":
      // @ts-expect-error not sure what the issue is here
      return <DebianIcon ref={ref} {...rest}/>
    default:
      return <></>
  }
});
