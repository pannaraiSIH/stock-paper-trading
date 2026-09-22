"use client";

import { memo } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { BookmarkIcon, BriefcaseIcon, GearIcon, GridIcon, ListIcon } from "../icons";

interface SidebarProps {
  email: string | null;
  onSignOut: () => void;
  mobileOpen: boolean;
  onMobileOpenChange: (open: boolean) => void;
}

const NAV_ITEMS = [
  { key: "dashboard", label: "Dashboard", icon: GridIcon, href: "/dashboard" },
  { key: "watchlist", label: "Watchlist", icon: BookmarkIcon, href: "/watchlist" },
  { key: "portfolio", label: "Portfolio", icon: BriefcaseIcon, href: "/portfolio" },
  { key: "orders", label: "Orders", icon: ListIcon, href: "/orders" },
  { key: "settings", label: "Settings", icon: GearIcon, href: null },
];

interface SidebarNavProps {
  email: string | null;
  onSignOut: () => void;
  onNavigate?: () => void;
}

function SidebarNav({ email, onSignOut, onNavigate }: SidebarNavProps) {
  const pathname = usePathname();

  return (
    <>
      <div className="sidebar__brand">
        <span className="brand__mark">P</span>
        <span className="brand__name">Paperline</span>
      </div>

      <nav className="sidebar__nav" aria-label="Primary">
        {NAV_ITEMS.map(({ key, label, icon: Icon, href }) => {
          const active = href !== null && pathname.startsWith(href);
          const itemClassName = cn(
            "h-auto w-full justify-start gap-2 px-2.5 py-2 font-normal",
            active ? "text-primary" : "text-muted-foreground"
          );

          if (href === null) {
            return (
              <Button
                key={key}
                variant="ghost"
                className={itemClassName}
                disabled
                title="Coming soon"
              >
                <Icon className="size-[1.05rem]" />
                <span>{label}</span>
              </Button>
            );
          }

          return (
            <Button
              key={key}
              asChild
              variant={active ? "secondary" : "ghost"}
              className={itemClassName}
            >
              <Link
                href={href}
                aria-current={active ? "page" : undefined}
                onClick={() => onNavigate?.()}
              >
                <Icon className="size-[1.05rem]" />
                <span>{label}</span>
              </Link>
            </Button>
          );
        })}
      </nav>

      <div className="sidebar__account">
        {email && <span className="watchlist__name">{email}</span>}
        <Button
          type="button"
          variant="link"
          className="h-auto justify-start p-0 text-xs text-muted-foreground"
          onClick={onSignOut}
        >
          Sign out
        </Button>
      </div>
    </>
  );
}

export const Sidebar = memo(function Sidebar({
  email,
  onSignOut,
  mobileOpen,
  onMobileOpenChange,
}: SidebarProps) {
  return (
    <>
      <aside className="sidebar">
        <SidebarNav email={email} onSignOut={onSignOut} />
      </aside>

      <Sheet open={mobileOpen} onOpenChange={onMobileOpenChange}>
        <SheetContent side="left" className="w-72 border-none bg-(--color-paper-2) p-0">
          <SheetHeader className="sr-only">
            <SheetTitle>Navigation</SheetTitle>
            <SheetDescription>Dashboard navigation and account menu</SheetDescription>
          </SheetHeader>
          <div className="flex h-full flex-col gap-10 px-4 py-6">
            <SidebarNav
              email={email}
              onSignOut={onSignOut}
              onNavigate={() => onMobileOpenChange(false)}
            />
          </div>
        </SheetContent>
      </Sheet>
    </>
  );
});
