import { memo } from "react";
import { fmtPct, fmtUSD } from "@/lib/format";
import { cn } from "@/lib/utils";
import type { WatchlistItem } from "@/types";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { TablePagination } from "@/components/dashboard/TablePagination";
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
  page: number;
  pageSize: number;
  hasMore: boolean;
  onPrevPage: () => void;
  onNextPage: () => void;
}

const headCls =
  "font-mono text-[0.68rem] tracking-wide uppercase text-muted-foreground";
const numHeadCls =
  "text-right font-mono text-[0.68rem] tracking-wide uppercase text-muted-foreground";
const numCellCls = "text-right font-mono tabular-nums";

export const WatchlistPanel = memo(function WatchlistPanel({
  rows,
  activeSymbol,
  onSelect,
  onRemove,
  page,
  pageSize,
  hasMore,
  onPrevPage,
  onNextPage,
}: WatchlistPanelProps) {
  return (
    <div className="panel panel--table min-w-0">
      <h2 className="panel__title">Watchlist</h2>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className={headCls}>Symbol</TableHead>
            <TableHead className={numHeadCls}>Price</TableHead>
            <TableHead className={numHeadCls}>Change</TableHead>
            <TableHead className="w-10" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.length === 0 ? (
            <TableRow className="hover:bg-transparent">
              <TableCell
                colSpan={4}
                className="py-4 text-center text-sm whitespace-normal text-muted-foreground"
              >
                No symbols yet. Search above to add one.
              </TableCell>
            </TableRow>
          ) : (
            rows.map(({ item, name, livePrice, prevClose }) => {
              const price = livePrice ?? item.price;
              const change =
                price !== null && prevClose !== null
                  ? price - prevClose
                  : null;
              const changePct =
                change !== null && prevClose
                  ? (change / prevClose) * 100
                  : null;
              const up = (change ?? 0) >= 0;
              const active = item.symbol === activeSymbol;

              return (
                <TableRow key={item.id}>
                  <TableCell className="font-mono">
                    <Button
                      type="button"
                      variant="link"
                      onClick={() => onSelect(item.symbol)}
                      className={cn(
                        "h-auto p-0 underline-offset-2",
                        active ? "text-primary" : "text-foreground",
                      )}
                    >
                      {item.symbol}
                    </Button>
                    {name && (
                      <span className="ml-2 font-sans text-muted-foreground">
                        {name}
                      </span>
                    )}
                  </TableCell>
                  <TableCell className={numCellCls}>
                    {price === null ? "—" : fmtUSD(price)}
                  </TableCell>
                  <TableCell
                    className={cn(
                      numCellCls,
                      changePct !== null && (up ? "text-positive" : "text-negative"),
                    )}
                  >
                    {changePct === null
                      ? "—"
                      : `${up ? "▲" : "▼"} ${fmtPct(changePct)}`}
                  </TableCell>
                  <TableCell className="text-right">
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      className="size-7 shrink-0 text-muted-foreground hover:text-negative"
                      aria-label={`Remove ${item.symbol} from watchlist`}
                      onClick={() => onRemove(item)}
                    >
                      <XIcon className="size-3.5" />
                    </Button>
                  </TableCell>
                </TableRow>
              );
            })
          )}
        </TableBody>
      </Table>
      <TablePagination
        page={page}
        pageSize={pageSize}
        itemCount={rows.length}
        hasMore={hasMore}
        onPrev={onPrevPage}
        onNext={onNextPage}
      />
    </div>
  );
});
