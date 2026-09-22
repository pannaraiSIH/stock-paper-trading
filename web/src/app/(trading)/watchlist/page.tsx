"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Plus } from "lucide-react";
import { api } from "@/lib/api";
import { errorMessage } from "@/lib/errors";
import type { StockSearchResult, WatchlistItem } from "@/types";
import { Button } from "@/components/ui/button";
import {
  WatchlistPanel,
  type WatchlistRow,
} from "@/components/dashboard/WatchlistPanel";
import { CommandPalette } from "@/components/dashboard/CommandPalette";
import { PageLoading } from "@/components/PageLoading";

export default function WatchlistPage() {
  const router = useRouter();

  const [loading, setLoading] = useState(true);
  const [items, setItems] = useState<WatchlistItem[]>([]);
  const [names, setNames] = useState<Record<string, string>>({});
  const [addOpen, setAddOpen] = useState(false);

  useEffect(() => {
    (async () => {
      setLoading(true);
      try {
        const watchlistItems = await api.getWatchlistItems();
        setItems(watchlistItems);

        const detailEntries = await Promise.all(
          watchlistItems.map(async (item) => {
            try {
              const details = await api.getStockDetails(item.symbol);
              return [item.symbol, details.name] as const;
            } catch {
              return null;
            }
          }),
        );
        setNames(Object.fromEntries(detailEntries.filter((e) => e !== null)));
      } catch (err) {
        toast.error(errorMessage(err, "Failed to load your watchlist."));
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const removeItem = useCallback(async (item: WatchlistItem) => {
    setItems((prev) => prev.filter((i) => i.id !== item.id));
    try {
      await api.deleteWatchlistItem(item.id);
      toast(`Removed ${item.symbol} from watchlist.`, {
        action: {
          label: "Undo",
          onClick: async () => {
            try {
              const restored = await api.addWatchlistItem(item.symbol);
              setItems((prev) => [...prev, restored]);
            } catch (err) {
              toast.error(
                errorMessage(err, `Could not restore ${item.symbol}.`),
              );
            }
          },
        },
      });
    } catch (err) {
      setItems((prev) => [...prev, item]);
      toast.error(errorMessage(err, `Could not remove ${item.symbol}.`));
    }
  }, []);

  async function addSymbol(result: StockSearchResult) {
    if (items.some((i) => i.symbol === result.symbol)) {
      toast.error(`${result.symbol} is already on your watchlist.`);
      return;
    }
    try {
      const item = await api.addWatchlistItem(result.symbol);
      setItems((prev) => [...prev, item]);
      setNames((prev) => ({ ...prev, [result.symbol]: result.name }));
    } catch (err) {
      toast.error(errorMessage(err, `Could not add ${result.symbol}.`));
    }
  }

  const rows: WatchlistRow[] = useMemo(
    () =>
      items.map((item) => ({
        item,
        name: names[item.symbol] ?? null,
        livePrice: null,
        prevClose: null,
      })),
    [items, names],
  );

  const handleSelect = useCallback(
    (symbol: string) =>
      router.push(`/dashboard?symbol=${encodeURIComponent(symbol)}`),
    [router],
  );

  if (loading) {
    return <PageLoading label="watchlist" />;
  }

  return (
    <>
      <main id="main" className="main">
        <section
          className="row reveal"
          style={{ "--i": 0 } as React.CSSProperties}
        >
          <div className="flex items-center justify-between">
            <h1 className="text-lg font-display font-semibold text-foreground">
              Watchlist
            </h1>
            <Button
              type="button"
              onClick={() => setAddOpen(true)}
              className="gap-1.5"
            >
              <Plus className="size-4" />
              Add symbol
            </Button>
          </div>
        </section>

        <section
          className="row reveal"
          style={{ "--i": 1 } as React.CSSProperties}
        >
          <WatchlistPanel
            rows={rows}
            activeSymbol=""
            onSelect={handleSelect}
            onRemove={removeItem}
          />
        </section>
      </main>

      <footer className="statusbar">
        <span>Paperline &mdash; paper trading sandbox</span>
      </footer>

      <CommandPalette
        open={addOpen}
        onClose={() => setAddOpen(false)}
        onSelect={addSymbol}
        livePrices={{}}
      />
    </>
  );
}
