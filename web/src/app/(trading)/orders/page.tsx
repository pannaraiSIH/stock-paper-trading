"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { api, PAGE_SIZE } from "@/lib/api";
import { errorMessage } from "@/lib/errors";
import type { Order } from "@/types";
import { OrdersTable } from "@/components/dashboard/OrdersTable";
import { PageLoading } from "@/components/PageLoading";

export default function OrdersPage() {
  const [loading, setLoading] = useState(true);
  const [orders, setOrders] = useState<Order[]>([]);
  const [page, setPage] = useState(0);
  const [hasMore, setHasMore] = useState(false);

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const result = await api.getOrdersPage(PAGE_SIZE, page * PAGE_SIZE);
        if (cancelled) return;
        setOrders(result.items);
        setHasMore(result.hasMore);
      } catch (err) {
        if (!cancelled) {
          toast.error(errorMessage(err, "Failed to load your orders."));
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [page]);

  if (loading) {
    return <PageLoading label="orders" />;
  }

  return (
    <>
      <main id="main" className="main">
        <section
          className="row reveal"
          style={{ "--i": 0 } as React.CSSProperties}
        >
          <h1 className="text-lg font-display font-semibold text-foreground">
            Order history
          </h1>
        </section>

        <section
          className="row reveal"
          style={{ "--i": 1 } as React.CSSProperties}
        >
          <OrdersTable
            orders={orders}
            page={page}
            pageSize={PAGE_SIZE}
            hasMore={hasMore}
            onPrevPage={() => setPage((p) => Math.max(0, p - 1))}
            onNextPage={() => setPage((p) => p + 1)}
          />
        </section>
      </main>

      <footer className="statusbar">
        <span>Paperline &mdash; paper trading sandbox</span>
      </footer>
    </>
  );
}
