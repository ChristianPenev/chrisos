"use client";

import { Check, Copy } from "lucide-react";
import { useState } from "react";

import { Button } from "@shadcn/button";

import type { CopyButtonProps } from "./contracts";

export function CopyButton({ value, label }: CopyButtonProps) {
  const [copied, setCopied] = useState(false);

  async function copy() {
    await navigator.clipboard.writeText(value);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  }

  return (
    <Button data-slot="copy-button" variant="ghost" size="sm" onClick={copy} aria-label={label}>
      {copied ? <Check /> : <Copy />}
      <span aria-live="polite">{copied ? "Copied" : "Copy"}</span>
    </Button>
  );
}
