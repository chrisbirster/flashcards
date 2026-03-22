export interface AppNavigationItem {
  to: string;
  label: string;
  description: string;
}

export const appNavigation: AppNavigationItem[] = [
  { to: "/", label: "Home", description: "Overview and usage" },
  { to: "/stats", label: "Stats", description: "Study analytics and trends" },
  { to: "/notes/view", label: "Notes", description: "Browse and edit notes" },
  {
    to: "/marketplace",
    label: "Marketplace",
    description: "Browse and publish listings",
  },
  {
    to: "/templates",
    label: "Templates",
    description: "Manage card templates",
  },
  { to: "/decks", label: "Decks", description: "Organize study decks" },
  {
    to: "/study-groups",
    label: "Study Groups",
    description: "Source decks and member installs",
  },
];

export function pageTitleForPath(pathname: string): string {
  if (pathname === "/") return "Home";
  if (pathname.startsWith("/stats")) return "Stats";
  if (pathname.startsWith("/notes/view")) return "Notes";
  if (pathname.startsWith("/notes/add")) return "Add Note";
  if (pathname.startsWith("/marketplace")) return "Marketplace";
  if (pathname.startsWith("/templates")) return "Templates";
  if (pathname.startsWith("/decks")) return "Decks";
  if (pathname.startsWith("/study-groups")) return "Study Groups";
  if (pathname.startsWith("/study/")) return "Study";
  if (pathname.startsWith("/tools/empty-cards")) return "Empty Cards";
  return "Vutadex";
}
