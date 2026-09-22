import { fmtPct, fmtSigned, fmtUSD } from "@/lib/format";
import type { Candle } from "@/types";
import { Toggle } from "@/components/ui/toggle";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { StarIcon } from "../icons";
import { PriceChart } from "./PriceChart";

export type ChartRange = "1D" | "1W" | "1M" | "1Y";

const RANGES: { key: ChartRange; label: string }[] = [
  { key: "1D", label: "1D" },
  { key: "1W", label: "1W" },
  { key: "1M", label: "1M" },
  { key: "1Y", label: "1Y" },
];

const RANGE_DESCRIPTION: Record<ChartRange, string> = {
  "1D": "Today (5-minute bars)",
  "1W": "Past 5 trading days",
  "1M": "Past 30 days",
  "1Y": "Past 12 months",
};

interface SymbolPanelProps {
  symbol: string;
  name: string | null;
  price: number | null;
  prevClose: number | null;
  isWatched: boolean;
  watchBusy: boolean;
  onToggleWatch: () => void;
  range: ChartRange;
  onRangeChange: (range: ChartRange) => void;
  candles: Candle[];
  chartLoading: boolean;
}

export function SymbolPanel({
  symbol,
  name,
  price,
  prevClose,
  isWatched,
  watchBusy,
  onToggleWatch,
  range,
  onRangeChange,
  candles,
  chartLoading,
}: SymbolPanelProps) {
  const change =
    price !== null && prevClose !== null ? price - prevClose : null;
  const changePct =
    change !== null && prevClose ? (change / prevClose) * 100 : null;
  const up = (change ?? 0) >= 0;

  return (
    <div className="panel panel--symbol">
      <div className="symbol__head">
        <div>
          <div className="symbol__id">
            <h1>{symbol || "—"}</h1>
            <span className="symbol__name">{name ?? " "}</span>
          </div>
          <div className="symbol__price">
            <span className="price">
              {price === null ? "—" : fmtUSD(price)}
            </span>
            {change !== null && changePct !== null ? (
              <span className={`change ${up ? "change--up" : "change--down"}`}>
                {up ? "▲" : "▼"} {fmtSigned(change)} ({fmtPct(changePct)})
              </span>
            ) : (
              <span className="change">&nbsp;</span>
            )}
          </div>
        </div>
        <Toggle
          pressed={isWatched}
          onPressedChange={onToggleWatch}
          disabled={watchBusy || !symbol}
          aria-label={
            isWatched
              ? `Remove ${symbol} from watchlist`
              : `Add ${symbol} to watchlist`
          }
          variant="outline"
          className="px-0 data-[state=on]:border-(--color-accent) data-[state=on]:bg-accent/8 data-[state=on]:text-(--color-accent) [&_svg]:data-[state=on]:fill-(--color-accent)"
        >
          <StarIcon />
        </Toggle>
      </div>

      <div className="chart-card">
        <div className="chart-card__head">
          <Tabs
            value={range}
            onValueChange={(value) => onRangeChange(value as ChartRange)}
          >
            <TabsList className="h-auto gap-0.5 rounded-md border border-(--color-graphite-rule) bg-transparent p-0.75">
              {RANGES.map((r) => (
                <TabsTrigger
                  key={r.key}
                  value={r.key}
                  className="h-auto rounded px-3 py-1 font-mono text-xs text-(--color-graphite-muted) data-[state=active]:border-transparent data-[state=active]:bg-(--color-graphite-2) data-[state=active]:text-(--color-accent-2) data-[state=active]:shadow-none"
                >
                  {r.label}
                </TabsTrigger>
              ))}
            </TabsList>
          </Tabs>
          <span className="chart-card__range">{RANGE_DESCRIPTION[range]}</span>
        </div>
        <PriceChart candles={candles} loading={chartLoading} />
      </div>
    </div>
  );
}
