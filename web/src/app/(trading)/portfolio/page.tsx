"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useTradingStore } from "@/stores/tradingStore";
import { useMarketStore } from "@/stores/marketStore";
import { api, PAGE_SIZE } from "@/lib/api";
import { errorMessage } from "@/lib/errors";
import { fmtPct, fmtSigned, fmtUSD } from "@/lib/format";
import { calculatePortfolioSummary } from "@/lib/portfolio";
import { cn } from "@/lib/utils";
import type { Position } from "@/types";
import { PositionsTable } from "@/components/dashboard/PositionsTable";
import { PageLoading } from "@/components/PageLoading";
import { Card, CardContent } from "@/components/ui/card";

export default function PortfolioPage() {
  const { account } = useTradingStore();
  const { prices, setSymbols } = useMarketStore();
  const router = useRouter();

  const [loading, setLoading] = useState(true);
  // Full set, used for the summary cards — independent of the table's page.
  const [positions, setPositions] = useState<Position[]>([]);
  // Table-only slice.
  const [pagePositions, setPagePositions] = useState<Position[]>([]);
  const [page, setPage] = useState(0);
  const [hasMore, setHasMore] = useState(false);

  useEffect(() => {
    (async () => {
      setLoading(true);
      try {
        const pos = await api.getPositions();
        setPositions(pos);
      } catch (err) {
        toast.error(errorMessage(err, "Failed to load your portfolio."));
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const result = await api.getPositionsPage(PAGE_SIZE, page * PAGE_SIZE);
        if (cancelled) return;
        setPagePositions(result.items);
        setHasMore(result.hasMore);
      } catch (err) {
        if (!cancelled) {
          toast.error(errorMessage(err, "Failed to load your positions."));
        }
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [page]);

  const subscribedSymbols = useMemo(
    () => Array.from(new Set(positions.map((p) => p.symbol))).sort(),
    [positions],
  );

  useEffect(() => {
    setSymbols(subscribedSymbols);
    return () => setSymbols([]);
  }, [subscribedSymbols, setSymbols]);

  const livePriceMap = useMemo(() => {
    const map: Record<string, number> = {};
    if (!prices) return map;
    for (const [symbol, evt] of Object.entries(prices)) map[symbol] = evt.price;
    return map;
  }, [prices]);

  const summary = useMemo(
    () => calculatePortfolioSummary(positions, livePriceMap),
    [positions, livePriceMap],
  );

  const handleSelectSymbol = useCallback(
    (symbol: string) =>
      router.push(`/dashboard?symbol=${encodeURIComponent(symbol)}`),
    [router],
  );

  if (loading) {
    return <PageLoading label="portfolio" />;
  }

  const up = summary.pl >= 0;

  return (
    <main id="main" className="main">
      <section
        className="row reveal"
        style={{ "--i": 0 } as React.CSSProperties}
      >
        <h1 className="text-lg font-display font-semibold text-foreground">
          Portfolio
        </h1>
      </section>

      <section
        className="row reveal grid-cols-3 max-[40rem]:grid-cols-1"
        style={{ "--i": 1 } as React.CSSProperties}
      >
        <Card className="min-w-0">
          <CardContent className="gap-1.5">
            <span className="panel__title" style={{ marginBottom: 0 }}>
              Cash available
            </span>
            <span className="block wrap-break-word font-mono text-base font-semibold tabular-nums text-foreground max-[40rem]:text-xl xl:text-lg">
              {fmtUSD(account?.cashBalance ?? 0)}
            </span>
          </CardContent>
        </Card>
        <Card className="min-w-0">
          <CardContent className="gap-1.5">
            <span className="panel__title" style={{ marginBottom: 0 }}>
              Market value
            </span>
            <span className="block wrap-break-word font-mono text-base font-semibold tabular-nums text-foreground max-[40rem]:text-xl xl:text-lg">
              {fmtUSD(summary.marketValue)}
            </span>
          </CardContent>
        </Card>
        <Card className="min-w-0">
          <CardContent className="gap-1.5">
            <span className="panel__title" style={{ marginBottom: 0 }}>
              Unrealized P&amp;L
            </span>
            <span
              className={cn(
                "block wrap-break-word font-mono text-base font-semibold tabular-nums max-[40rem]:text-xl xl:text-lg",
                up ? "text-positive" : "text-negative",
              )}
            >
              {fmtSigned(summary.pl)} ({fmtPct(summary.plPct)})
            </span>
          </CardContent>
        </Card>
      </section>

      <section
        className="row reveal"
        style={{ "--i": 2 } as React.CSSProperties}
      >
        <PositionsTable
          positions={pagePositions}
          livePrices={livePriceMap}
          onSelectSymbol={handleSelectSymbol}
          page={page}
          pageSize={PAGE_SIZE}
          hasMore={hasMore}
          onPrevPage={() => setPage((p) => Math.max(0, p - 1))}
          onNextPage={() => setPage((p) => p + 1)}
        />
      </section>
    </main>
  );
}
