import { memo } from "react";
import { fmtUSD } from "@/lib/format";
import { cn } from "@/lib/utils";
import type { Order, OrderSide, OrderStatus } from "@/types";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";

interface OrdersTableProps {
  orders: Order[];
  limit?: number;
}

const numHeadCls =
  "text-right font-mono text-[0.68rem] tracking-wide uppercase text-muted-foreground";
const headCls =
  "font-mono text-[0.68rem] tracking-wide uppercase text-muted-foreground";

const SIDE_STYLES: Record<OrderSide, string> = {
  buy: "border-transparent bg-[var(--color-positive)]/12 text-[var(--color-positive)]",
  sell: "border-transparent bg-[var(--color-negative)]/12 text-[var(--color-negative)]",
};

const STATUS_STYLES: Record<OrderStatus, string> = {
  executed: "border-[var(--color-positive)]/40 text-[var(--color-positive)]",
  rejected: "border-[var(--color-negative)]/40 text-[var(--color-negative)]",
  pending: "border-[var(--color-warning)]/50 text-[var(--color-warning)]",
};

export const OrdersTable = memo(function OrdersTable({
  orders,
  limit = 12,
}: OrdersTableProps) {
  return (
    <div className="panel panel--table">
      <h2 className="panel__title">Order history</h2>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className={headCls}>Order</TableHead>
            <TableHead className={headCls}>Symbol</TableHead>
            <TableHead className={headCls}>Side</TableHead>
            <TableHead className={numHeadCls}>Price</TableHead>
            <TableHead className={numHeadCls}>Total</TableHead>
            <TableHead className={headCls}>Status</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {orders.length === 0 ? (
            <TableRow className="hover:bg-transparent">
              <TableCell
                colSpan={6}
                className="py-4 text-center text-sm whitespace-normal text-muted-foreground"
              >
                No orders yet.
              </TableCell>
            </TableRow>
          ) : (
            orders.slice(0, limit).map((order) => (
              <TableRow key={order.id}>
                <TableCell className="font-mono">ORD-{order.id}</TableCell>
                <TableCell className="font-mono">{order.symbol}</TableCell>
                <TableCell>
                  <Badge
                    variant="outline"
                    className={cn(
                      "font-mono font-semibold",
                      SIDE_STYLES[order.side],
                    )}
                  >
                    {order.side.toUpperCase()}
                  </Badge>
                </TableCell>
                <TableCell className="text-right font-mono tabular-nums">
                  {order.executionPrice === null
                    ? "—"
                    : fmtUSD(order.executionPrice)}
                </TableCell>
                <TableCell className="text-right font-mono tabular-nums">
                  {order.totalValue === null ? "—" : fmtUSD(order.totalValue)}
                </TableCell>
                <TableCell>
                  <Badge
                    variant="outline"
                    className={cn(
                      "font-mono uppercase",
                      STATUS_STYLES[order.status],
                    )}
                  >
                    {order.status}
                  </Badge>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
});
