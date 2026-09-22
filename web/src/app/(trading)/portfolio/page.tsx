"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useTradingStore } from "@/stores/tradingStore";
import { useMarketStore } from "@/stores/marketStore";
import { api } from "@/lib/api";
import { errorMessage } from "@/lib/errors";
import { fmtPct, fmtSigned, fmtUSD } from "@/lib/format";
import { cn } from "@/lib/utils";
import type { Position } from "@/types";
import { PositionsTable } from "@/components/dashboard/PositionsTable";
import { PageLoading } from "@/components/PageLoading";

export default function PortfolioPage() {
  const { account } = useTradingStore();
  const { prices, setSymbols } = useMarketStore();
  const router = useRouter();

  const [loading, setLoading] = useState(true);
  const [positions, setPositions] = useState<Position[]>([]);

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

  const summary = useMemo(() => {
    let marketValue = 0;
    let costBasis = 0;
    for (const position of positions) {
      const current = livePriceMap[position.symbol] ?? position.averagePrice;
      marketValue += position.quantity * current;
      costBasis += position.quantity * position.averagePrice;
    }
    const pl = marketValue - costBasis;
    const plPct = costBasis ? (pl / costBasis) * 100 : 0;
    return { marketValue, pl, plPct };
  }, [positions, livePriceMap]);

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
    <>
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
          <div className="panel">
            <span className="panel__title">Cash available</span>
            <span className="price mono block">
              {fmtUSD(account?.cashBalance ?? 0)}
            </span>
          </div>
          <div className="panel">
            <span className="panel__title">Market value</span>
            <span className="price mono block">
              {fmtUSD(summary.marketValue)}
            </span>
          </div>
          <div className="panel">
            <span className="panel__title">Unrealized P&amp;L</span>
            <span
              className={cn(
                "price mono block",
                up ? "text-positive" : "text-negative",
              )}
            >
              {fmtSigned(summary.pl)} ({fmtPct(summary.plPct)})
            </span>
          </div>
        </section>

        <section
          className="row reveal"
          style={{ "--i": 2 } as React.CSSProperties}
        >
          <PositionsTable
            positions={positions}
            livePrices={livePriceMap}
            onSelectSymbol={handleSelectSymbol}
          />
        </section>
      </main>

      <footer className="statusbar">
        <span>Paperline &mdash; paper trading sandbox</span>
      </footer>
    </>
  );
}
