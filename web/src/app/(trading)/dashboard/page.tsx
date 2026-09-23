"use client";

import {
  Suspense,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
  type CSSProperties,
} from "react";
import { useSearchParams } from "next/navigation";
import { toast } from "sonner";
import { api } from "@/lib/api";
import { errorMessage } from "@/lib/errors";
import { fmtUSD } from "@/lib/format";
import type {
  Candle,
  CandleInterval,
  OrderSide,
  Position,
  StockSearchResult,
  WatchlistItem,
} from "@/types";
import {
  SymbolPanel,
  type ChartRange,
} from "@/components/dashboard/SymbolPanel";
import { OrderTicket } from "@/components/dashboard/OrderTicket";
import { CommandPalette } from "@/components/dashboard/CommandPalette";
import { useTradingStore } from "@/stores/tradingStore";
import { Button } from "@/components/ui/button";
import { SearchIcon } from "@/components/icons";
import { useMarketStore } from "@/stores/marketStore";
import { PageLoading } from "@/components/PageLoading";

const RANGE_CONFIG: Record<
  ChartRange,
  { interval: CandleInterval; outputSize: number }
> = {
  "1D": { interval: "5min", outputSize: 80 },
  "1W": { interval: "1h", outputSize: 50 },
  "1M": { interval: "1day", outputSize: 30 },
  "1Y": { interval: "1day", outputSize: 260 },
};

export default function DashboardPage() {
  return (
    <Suspense fallback={null}>
      <DashboardPageInner />
    </Suspense>
  );
}

function DashboardPageInner() {
  const { account, refreshAccount } = useTradingStore();
  const { prices, setSymbols } = useMarketStore();
  const searchParams = useSearchParams();

  const [loading, setLoading] = useState(true);
  const [positions, setPositions] = useState<Position[]>([]);
  const [watchlistItems, setWatchlistItems] = useState<WatchlistItem[]>([]);

  const [activeSymbol, setActiveSymbol] = useState("");
  const [range, setRange] = useState<ChartRange>("1D");
  const [candles, setCandles] = useState<Candle[]>([]);
  const [chartLoading, setChartLoading] = useState(false);
  const [watchBusy, setWatchBusy] = useState(false);
  const [orderSubmitting, setOrderSubmitting] = useState(false);
  const [paletteOpen, setPaletteOpen] = useState(false);

  const [symbolMeta, setSymbolMeta] = useState<Record<string, string>>({});
  const [prevCloseMap, setPrevCloseMap] = useState<Record<string, number>>({});
  const metaFetchedRef = useRef<Set<string>>(new Set());
  const prevCloseFetchedRef = useRef<Set<string>>(new Set());

  const symbolFromUrl = searchParams.get("symbol");

  useEffect(() => {
    (async () => {
      setLoading(true);
      try {
        const [items, pos] = await Promise.all([
          api.getWatchlistItems(),
          api.getPositions(),
        ]);
        setWatchlistItems(items);
        setPositions(pos);
        setActiveSymbol(
          (current) =>
            current ||
            symbolFromUrl ||
            items[0]?.symbol ||
            pos[0]?.symbol ||
            "AAPL",
        );
      } catch (err) {
        toast.error(errorMessage(err, "Failed to load your account."));
      } finally {
        setLoading(false);
      }
    })();
  }, [symbolFromUrl]);

  const subscribedSymbols = useMemo(() => {
    const set = new Set<string>();
    if (activeSymbol) set.add(activeSymbol);
    for (const item of watchlistItems) set.add(item.symbol);
    for (const position of positions) set.add(position.symbol);
    return Array.from(set).sort();
  }, [activeSymbol, watchlistItems, positions]);

  const ensureMeta = useCallback((symbol: string, knownName?: string) => {
    if (!symbol) return;
    if (knownName) {
      metaFetchedRef.current.add(symbol);
      setSymbolMeta((prev) =>
        prev[symbol] === knownName ? prev : { ...prev, [symbol]: knownName },
      );
      return;
    }
    if (metaFetchedRef.current.has(symbol)) return;
    metaFetchedRef.current.add(symbol);
    api
      .getStockDetails(symbol)
      .then((details) =>
        setSymbolMeta((prev) => ({ ...prev, [symbol]: details.name })),
      )
      .catch(() => metaFetchedRef.current.delete(symbol));
  }, []);

  const ensurePrevClose = useCallback((symbol: string) => {
    if (!symbol || prevCloseFetchedRef.current.has(symbol)) return;
    prevCloseFetchedRef.current.add(symbol);
    api
      .getCandles(symbol, "1day", 2)
      .then((data) => {
        if (data.length === 0) return;
        const reference = data[1] ?? data[0];
        const value = Number(reference.close);
        if (Number.isFinite(value)) {
          setPrevCloseMap((prev) => ({ ...prev, [symbol]: value }));
        }
      })
      .catch(() => prevCloseFetchedRef.current.delete(symbol));
  }, []);

  useEffect(() => {
    for (const symbol of subscribedSymbols) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      ensureMeta(symbol);
      ensurePrevClose(symbol);
    }
  }, [subscribedSymbols, ensureMeta, ensurePrevClose]);

  useEffect(() => {
    if (!activeSymbol) return;
    let cancelled = false;
    const cfg = RANGE_CONFIG[range];
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setChartLoading(true);
    api
      .getCandles(activeSymbol, cfg.interval, cfg.outputSize)
      .then((data) => {
        if (!cancelled) setCandles(data);
      })
      .catch(() => {
        if (!cancelled) setCandles([]);
      })
      .finally(() => {
        if (!cancelled) setChartLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [activeSymbol, range]);

  useEffect(() => {
    setSymbols(subscribedSymbols);
    return () => setSymbols([]);
  }, [subscribedSymbols, setSymbols]);

  async function removeWatchlistItem(item: WatchlistItem) {
    setWatchlistItems((prev) => prev.filter((i) => i.id !== item.id));
    try {
      await api.deleteWatchlistItem(item.id);
      toast(`Removed ${item.symbol} from watchlist.`, {
        action: {
          label: "Undo",
          onClick: async () => {
            try {
              const restored = await api.addWatchlistItem(item.symbol);
              setWatchlistItems((prev) => [...prev, restored]);
            } catch (err) {
              toast.error(
                errorMessage(err, `Could not restore ${item.symbol}.`),
              );
            }
          },
        },
      });
    } catch (err) {
      setWatchlistItems((prev) => [...prev, item]);
      toast.error(errorMessage(err, `Could not remove ${item.symbol}.`));
    }
  }

  async function toggleWatch() {
    if (!activeSymbol) return;
    const existing = watchlistItems.find((i) => i.symbol === activeSymbol);
    setWatchBusy(true);
    try {
      if (existing) {
        await removeWatchlistItem(existing);
      } else {
        const item = await api.addWatchlistItem(activeSymbol);
        setWatchlistItems((prev) => [...prev, item]);
      }
    } catch (err) {
      toast.error(errorMessage(err, "Could not update watchlist."));
    } finally {
      setWatchBusy(false);
    }
  }

  async function handleOrderSubmit(side: OrderSide, quantity: number) {
    if (!activeSymbol) return;
    setOrderSubmitting(true);
    try {
      const order = await api.createOrder(activeSymbol, side, quantity);
      const [pos] = await Promise.all([api.getPositions(), refreshAccount()]);
      setPositions(pos);
      if (order.status === "rejected") {
        toast.error(`Order for ${quantity} ${activeSymbol} was rejected.`);
      } else {
        const verb = side === "buy" ? "Bought" : "Sold";
        const priceLabel =
          order.executionPrice !== null
            ? ` at ${fmtUSD(order.executionPrice)}`
            : "";
        toast(`${verb} ${quantity} ${activeSymbol}${priceLabel}.`);
      }
    } catch (err) {
      toast.error(errorMessage(err, "Order failed."));
    } finally {
      setOrderSubmitting(false);
    }
  }

  function handleSearchSelect(result: StockSearchResult) {
    ensureMeta(result.symbol, result.name);
    setActiveSymbol(result.symbol);
  }

  function handleOpenSearch() {
    setPaletteOpen(true);
  }

  const livePriceMap = useMemo(() => {
    const map: Record<string, number> = {};
    if (!prices) return map;
    for (const [symbol, evt] of Object.entries(prices)) map[symbol] = evt.price;
    return map;
  }, [prices]);

  const latestCandleClose =
    candles.length > 0 ? Number(candles[0].close) : null;
  const activePrice = activeSymbol
    ? (livePriceMap[activeSymbol] ??
      watchlistItems.find((i) => i.symbol === activeSymbol)?.price ??
      (latestCandleClose !== null && Number.isFinite(latestCandleClose)
        ? latestCandleClose
        : null))
    : null;
  const activeOwnedQty =
    positions.find((p) => p.symbol === activeSymbol)?.quantity ?? 0;
  const cash = account?.cashBalance ?? 0;

  if (loading) {
    return <PageLoading label="dashboard" />;
  }

  return (
    <>
      <main id="main" className="main">
        <section
          className="row reveal"
          style={{ "--i": 0 } as unknown as CSSProperties}
        >
          <div className="flex items-center">
            <Button
              variant="outline"
              onClick={handleOpenSearch}
              aria-haspopup="dialog"
              className="h-11 min-w-60 max-w-88 justify-start gap-2 px-3 font-normal text-muted-foreground max-sm:min-w-0 max-sm:flex-1"
            >
              <SearchIcon className="size-4 shrink-0" />
              <span className="flex-1 text-left">Search symbols&hellip;</span>
              <kbd className="rounded border border-border px-1.5 py-0.5 font-mono text-[0.66rem] text-muted-foreground">
                &#8984;K
              </kbd>
            </Button>
          </div>
        </section>

        <section
          className="row row--primary reveal"
          style={{ "--i": 1 } as unknown as CSSProperties}
        >
          <SymbolPanel
            symbol={activeSymbol}
            name={activeSymbol ? (symbolMeta[activeSymbol] ?? null) : null}
            price={activePrice}
            prevClose={
              activeSymbol ? (prevCloseMap[activeSymbol] ?? null) : null
            }
            isWatched={watchlistItems.some((i) => i.symbol === activeSymbol)}
            watchBusy={watchBusy}
            onToggleWatch={toggleWatch}
            range={range}
            onRangeChange={setRange}
            candles={candles}
            chartLoading={chartLoading}
          />
          <OrderTicket
            symbol={activeSymbol}
            price={activePrice}
            cash={cash}
            ownedQty={activeOwnedQty}
            submitting={orderSubmitting}
            onSubmit={handleOrderSubmit}
          />
        </section>
      </main>

      <CommandPalette
        open={paletteOpen}
        onClose={() => setPaletteOpen(false)}
        onSelect={handleSearchSelect}
        livePrices={livePriceMap}
      />
    </>
  );
}
