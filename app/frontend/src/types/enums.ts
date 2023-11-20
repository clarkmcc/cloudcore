export enum AgentEventType {
  AgentStartup = "AGENT_STARTUP",
  AgentShutdown = "AGENT_SHUTDOWN",
}

// eslint-disable-next-line @typescript-eslint/no-namespace
export namespace AgentEventType {
  export function toBackgroundClass(type: AgentEventType): string {
    switch (type) {
      case AgentEventType.AgentStartup:
        return "dark:bg-green-900 bg-green-300";
      case AgentEventType.AgentShutdown:
        return "dark:bg-red-900 bg-red-300";
      default:
        return "bg-muted";
    }
  }
}
