const longDate = new Intl.DateTimeFormat("en-GB", { dateStyle: "long", timeZone: "UTC" });

export function formatDate(iso: string) {
  return longDate.format(new Date(iso));
}
