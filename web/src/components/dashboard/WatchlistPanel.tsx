import { memo } from "react";
import { fmtPct, fmtUSD } from "@/lib/format";
import type { WatchlistItem } from "@/types";
import { Button } from "@/components/ui/button";
import { XIcon } from "../icons";

export interface WatchlistRow {
  item: WatchlistItem;
  name: string | null;
  livePrice: number | null;
  prevClose: number | null;
}

interface WatchlistPanelProps {
  rows: WatchlistRow[];
  activeSymbol: string;
  onSelect: (symbol: string) => void;
  onRemove: (item: WatchlistItem) => void;
}

export const WatchlistPanel = memo(function WatchlistPanel({
  rows,
  activeSymbol,
  onSelect,
  onRemove,
}: WatchlistPanelProps) {
  return (
    <div className="panel panel--list">
      <h2 className="panel__title">Watchlist</h2>
      {rows.length === 0 ? (
        <p className="watchlist__empty">No symbols yet. Search above to add one.</p>
      ) : (
        <ul className="watchlist">
          {rows.map(({ item, name, livePrice, prevClose }) => {
            const price = livePrice ?? item.price;
            const change = price !== null && prevClose !== null ? price - prevClose : null;
            const changePct = change !== null && prevClose ? (change / prevClose) * 100 : null;
            const up = (change ?? 0) >= 0;

            return (
              <li key={item.id}>
                <button
                  type="button"
                  className={`watchlist__symbol-btn${item.symbol === activeSymbol ? " is-active" : ""}`}
                  aria-label={`Open ${item.symbol}`}
                  onClick={() => onSelect(item.symbol)}
                >
                  <span className="watchlist__symbol">
                    <span className="watchlist__ticker">{item.symbol}</span>
                    <span className="watchlist__name">{name ?? " "}</span>
                  </span>
                </button>
                <span className="watchlist__price">{price === null ? "—" : fmtUSD(price)}</span>
                <span className={`watchlist__change ${up ? "change--up" : "change--down"}`}>
                  {changePct === null ? " " : `${up ? "▲" : "▼"} ${fmtPct(changePct)}`}
                </span>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="size-7 shrink-0 text-muted-foreground hover:text-[var(--color-negative)]"
                  aria-label={`Remove ${item.symbol} from watchlist`}
                  onClick={() => onRemove(item)}
                >
                  <XIcon className="size-3.5" />
                </Button>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
});
