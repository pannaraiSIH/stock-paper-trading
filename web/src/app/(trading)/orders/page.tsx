"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { api } from "@/lib/api";
import { errorMessage } from "@/lib/errors";
import type { Order } from "@/types";
import { OrdersTable } from "@/components/dashboard/OrdersTable";
import { PageLoading } from "@/components/PageLoading";

export default function OrdersPage() {
  const [loading, setLoading] = useState(true);
  const [orders, setOrders] = useState<Order[]>([]);

  useEffect(() => {
    (async () => {
      setLoading(true);
      try {
        const orderList = await api.getOrders();
        setOrders(orderList);
      } catch (err) {
        toast.error(errorMessage(err, "Failed to load your orders."));
      } finally {
        setLoading(false);
      }
    })();
  }, []);

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
          <OrdersTable orders={orders} limit={100} />
        </section>
      </main>

      <footer className="statusbar">
        <span>Paperline &mdash; paper trading sandbox</span>
      </footer>
    </>
  );
}
