export function baseUrl(path: string): string {
  const { protocol, hostname } = window.location;
  const isLocal =
    hostname.includes("127.0.0.1") || hostname.includes("localhost");
  return `${protocol}//${hostname}${isLocal ? ":10001" : ""}${path}`;
}
