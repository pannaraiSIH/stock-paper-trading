"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Plus } from "lucide-react";
import { api, ApiError, PAGE_SIZE } from "@/lib/api";
import { errorMessage } from "@/lib/errors";
import type { StockSearchResult, WatchlistItem } from "@/types";
import { Button } from "@/components/ui/button";
import {
  WatchlistPanel,
  type WatchlistRow,
} from "@/components/dashboard/WatchlistPanel";
import { CommandPalette } from "@/components/dashboard/CommandPalette";
import { PageLoading } from "@/components/PageLoading";

async function fetchWatchlistPage(offset: number) {
  const page = await api.getWatchlistItemsPage(PAGE_SIZE, offset);
  const nameEntries = await Promise.all(
    page.items.map(async (item) => {
      try {
        const details = await api.getStockDetails(item.symbol);
        return [item.symbol, details.name] as const;
      } catch {
        return null;
      }
    }),
  );
  return {
    items: page.items,
    hasMore: page.hasMore,
    names: Object.fromEntries(nameEntries.filter((e) => e !== null)),
  };
}

export default function WatchlistPage() {
  const router = useRouter();

  const [loading, setLoading] = useState(true);
  const [items, setItems] = useState<WatchlistItem[]>([]);
  const [names, setNames] = useState<Record<string, string>>({});
  const [page, setPage] = useState(0);
  const [hasMore, setHasMore] = useState(false);
  const [refreshKey, setRefreshKey] = useState(0);
  const [addOpen, setAddOpen] = useState(false);

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const result = await fetchWatchlistPage(page * PAGE_SIZE);
        if (cancelled) return;
        setItems(result.items);
        setHasMore(result.hasMore);
        setNames((prev) => ({ ...prev, ...result.names }));
      } catch (err) {
        if (!cancelled) {
          toast.error(errorMessage(err, "Failed to load your watchlist."));
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [page, refreshKey]);

  // Refetches the page currently on screen (e.g. after add/remove), without
  // changing which page that is.
  const refresh = useCallback(() => setRefreshKey((k) => k + 1), []);

  // Newly added/restored items always sort to the top (created_at DESC), so
  // jump back to page 0 to bring them into view.
  const goToFirstPage = useCallback(() => {
    setPage(0);
    setRefreshKey((k) => k + 1);
  }, []);

  const removeItem = useCallback(
    async (item: WatchlistItem) => {
      setItems((prev) => prev.filter((i) => i.id !== item.id));
      try {
        await api.deleteWatchlistItem(item.id);
        refresh();
        toast(`Removed ${item.symbol} from watchlist.`, {
          action: {
            label: "Undo",
            onClick: async () => {
              try {
                await api.addWatchlistItem(item.symbol);
                goToFirstPage();
              } catch (err) {
                toast.error(
                  errorMessage(err, `Could not restore ${item.symbol}.`),
                );
              }
            },
          },
        });
      } catch (err) {
        refresh();
        toast.error(errorMessage(err, `Could not remove ${item.symbol}.`));
      }
    },
    [refresh, goToFirstPage],
  );

  async function addSymbol(result: StockSearchResult) {
    try {
      await api.addWatchlistItem(result.symbol);
      setNames((prev) => ({ ...prev, [result.symbol]: result.name }));
      goToFirstPage();
    } catch (err) {
      if (err instanceof ApiError && err.status === 409) {
        toast.error(`${result.symbol} is already on your watchlist.`);
        return;
      }
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
            page={page}
            pageSize={PAGE_SIZE}
            hasMore={hasMore}
            onPrevPage={() => setPage((p) => Math.max(0, p - 1))}
            onNextPage={() => setPage((p) => p + 1)}
          />
        </section>
      </main>

      <CommandPalette
        open={addOpen}
        onClose={() => setAddOpen(false)}
        onSelect={addSymbol}
        livePrices={{}}
      />
    </>
  );
}
