export function hueFromSeed(seed: string) {
  let hash = 0;
  for (const char of seed) {
    hash = (hash * 31 + char.charCodeAt(0)) >>> 0;
  }
  return hash % 360;
}

export function hueStyle(seed: string) {
  return { "--hue": hueFromSeed(seed) } as React.CSSProperties;
}
