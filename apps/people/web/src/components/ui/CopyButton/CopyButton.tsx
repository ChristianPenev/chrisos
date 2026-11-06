"use client";

import { Check, Copy } from "lucide-react";
import { useState } from "react";

import { cn } from "@lib/utils";
import { Button } from "@shadcn/button";

import type { CopyButtonProps } from "./contracts";

export function CopyButton({ value, label, compact = false, className }: CopyButtonProps) {
  const [copied, setCopied] = useState(false);

  async function copy() {
    await navigator.clipboard.writeText(value);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  }

  return (
    <Button
      data-slot="copy-button"
      data-copied={copied || undefined}
      variant="ghost"
      size={compact ? "icon-sm" : "sm"}
      onClick={copy}
      aria-label={label}
      title={compact ? label : undefined}
      className={cn("text-muted-foreground hover:text-foreground", className)}
    >
      {copied ? <Check className="text-foreground" /> : <Copy />}
      {compact ? (
        <span aria-live="polite" className="sr-only">
          {copied ? "Copied" : ""}
        </span>
      ) : (
        <span aria-live="polite">{copied ? "Copied" : "Copy"}</span>
      )}
    </Button>
  );
}
