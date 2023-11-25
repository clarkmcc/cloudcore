import { useState } from "react";
import { ClipboardCopy } from "lucide-react";

export function CommandCopy(props: { command: string }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(props.command);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="relative p-4 bg-neutral-100 dark:bg-neutral-900 dark:text-white font-mono text-sm rounded-md overflow-auto">
      <div className="overflow-x-auto whitespace-nowrap scroll-auto">
        {props.command}
      </div>
      <button
        onClick={handleCopy}
        className="absolute top-2 right-2 flex items-center gap-2 p-1 text-xs bg-neutral-200 hover:bg-neutral-300 dark:bg-neutral-700 dark:hover:bg-neutral-600 rounded"
      >
        <ClipboardCopy className="h-4 w-4" />
        {copied ? "Copied" : "Copy"}
      </button>
    </div>
  );
}
