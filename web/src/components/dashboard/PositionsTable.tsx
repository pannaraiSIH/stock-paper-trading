import { memo } from "react";
import { fmtPct, fmtQty, fmtSigned, fmtUSD } from "@/lib/format";
import { calculatePositionPL } from "@/lib/portfolio";
import { cn } from "@/lib/utils";
import type { Position } from "@/types";
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

interface PositionsTableProps {
  positions: Position[];
  livePrices: Record<string, number>;
  onSelectSymbol?: (symbol: string) => void;
  page: number;
  pageSize: number;
  hasMore: boolean;
  onPrevPage: () => void;
  onNextPage: () => void;
}

const numHeadCls =
  "text-right font-mono text-[0.68rem] tracking-wide uppercase text-muted-foreground";
const numCellCls = "text-right font-mono tabular-nums";

export const PositionsTable = memo(function PositionsTable({
  positions,
  livePrices,
  onSelectSymbol,
  page,
  pageSize,
  hasMore,
  onPrevPage,
  onNextPage,
}: PositionsTableProps) {
  return (
    <div className="panel panel--table min-w-0">
      <h2 className="panel__title">Open positions</h2>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="font-mono text-[0.68rem] tracking-wide uppercase text-muted-foreground">
              Symbol
            </TableHead>
            <TableHead className={numHeadCls}>Qty</TableHead>
            <TableHead className={numHeadCls}>Avg price</TableHead>
            <TableHead className={numHeadCls}>Current</TableHead>
            <TableHead className={numHeadCls}>Market value</TableHead>
            <TableHead className={numHeadCls}>P&amp;L</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {positions.length === 0 ? (
            <TableRow className="hover:bg-transparent">
              <TableCell
                colSpan={6}
                className="py-4 text-center text-sm whitespace-normal text-muted-foreground"
              >
                No open positions. Place a trade to get started.
              </TableCell>
            </TableRow>
          ) : (
            positions.map((position) => {
              const { current, marketValue, pl, plPct } = calculatePositionPL(
                position,
                livePrices[position.symbol],
              );
              const up = pl >= 0;

              return (
                <TableRow key={position.id}>
                  <TableCell className="font-mono">
                    {onSelectSymbol ? (
                      <Button
                        type="button"
                        variant="link"
                        onClick={() => onSelectSymbol(position.symbol)}
                        className="h-auto p-0 text-foreground underline-offset-2"
                      >
                        {position.symbol}
                      </Button>
                    ) : (
                      position.symbol
                    )}
                  </TableCell>
                  <TableCell className={numCellCls}>
                    {fmtQty(position.quantity)}
                  </TableCell>
                  <TableCell className={numCellCls}>
                    {fmtUSD(position.averagePrice)}
                  </TableCell>
                  <TableCell className={numCellCls}>
                    {fmtUSD(current)}
                  </TableCell>
                  <TableCell className={numCellCls}>
                    {fmtUSD(marketValue)}
                  </TableCell>
                  <TableCell
                    className={cn(
                      numCellCls,
                      up ? "text-positive" : "text-negative",
                    )}
                  >
                    {fmtSigned(pl)} ({fmtPct(plPct)})
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
        itemCount={positions.length}
        hasMore={hasMore}
        onPrev={onPrevPage}
        onNext={onNextPage}
      />
    </div>
  );
});
