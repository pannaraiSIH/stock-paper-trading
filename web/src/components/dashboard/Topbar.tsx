import { memo } from "react";
import { Menu } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

interface TopbarProps {
  connected?: boolean;
  syncedTime?: string;
  onOpenMobileNav: () => void;
  avatarInitials: string;
}

export const Topbar = memo(function Topbar({
  connected,
  syncedTime,
  onOpenMobileNav,
  avatarInitials,
}: TopbarProps) {
  return (
    <header className="topbar">
      <Button
        variant="outline"
        size="icon"
        onClick={onOpenMobileNav}
        aria-label="Open navigation menu"
        className="hidden shrink-0 max-[60rem]:inline-flex"
      >
        <Menu className="size-4.5" />
      </Button>

      {connected !== undefined && (
        <div className="topbar__status">
          <Badge
            variant="outline"
            className="gap-1.5 font-mono text-[0.68rem] tracking-wide uppercase"
          >
            <i
              aria-hidden="true"
              className={cn(
                "size-1.5 shrink-0 rounded-full",
                connected ? "bg-positive" : "bg-negative",
              )}
            />
            {connected ? "Live prices" : "Reconnecting…"}
          </Badge>
          <span className="topbar__synced">
            Synced <time className="mono">{syncedTime || "—"}</time>
          </span>
        </div>
      )}

      <Button
        variant="outline"
        size="icon"
        aria-label="Account menu"
        className="ml-auto rounded-md font-mono text-xs max-sm:order-2"
      >
        {avatarInitials}
      </Button>
    </header>
  );
});
