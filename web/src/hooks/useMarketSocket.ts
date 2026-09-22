"use client";

import { useEffect, useRef, useState } from "react";
import type { PriceEvent } from "@/types";
import { nowLabel } from "@/lib/format";

const WS_URL =
  process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8089/api/market/ws";
const RECONNECT_DELAY_MS = 2500;

interface MarketSocketState {
  prices: Record<string, PriceEvent>;
  connected: boolean;
  syncedTime: string;
}

/**
 * Maintains a single WebSocket connection to the market hub and keeps it
 * subscribed to exactly the given symbol set, resubscribing after
 * reconnects (the hub only knows about subscriptions on the live socket).
 */
export function useMarketSocket(
  symbols: string[],
  token: string | null,
): MarketSocketState {
  const [prices, setPrices] = useState<Record<string, PriceEvent>>({});
  const [connected, setConnected] = useState(false);
  const [syncedTime, setSyncedTime] = useState("");

  const symbolsRef = useRef<string[]>(symbols);
  useEffect(() => {
    symbolsRef.current = symbols;
  }, [symbols]);

  const socketRef = useRef<WebSocket | null>(null);
  const subscribedRef = useRef<Set<string>>(new Set());

  useEffect(() => {
    if (!token) return;

    let cancelled = false;
    let reconnectTimer: number | undefined;

    const connect = () => {
      if (cancelled) return;

      const socket = new WebSocket(
        `${WS_URL}?token=${encodeURIComponent(token)}`,
      );
      socketRef.current = socket;

      socket.addEventListener("open", () => {
        if (cancelled) return;
        setConnected(true);
        subscribedRef.current = new Set();
        for (const symbol of symbolsRef.current) {
          socket.send(JSON.stringify({ action: "subscribe", symbol }));
          subscribedRef.current.add(symbol);
        }
      });

      socket.addEventListener("message", (event) => {
        try {
          const data = JSON.parse(event.data) as PriceEvent;
          if (data.event === "price" && data.symbol) {
            setPrices((prev) => ({ ...prev, [data.symbol]: data }));
            setSyncedTime(nowLabel());
          }
        } catch {
          // ignore malformed frames
        }
      });

      const scheduleReconnect = () => {
        if (cancelled) return;
        setConnected(false);
        reconnectTimer = window.setTimeout(connect, RECONNECT_DELAY_MS);
      };

      socket.addEventListener("close", scheduleReconnect);
      socket.addEventListener("error", () => socket.close());
    };

    connect();

    return () => {
      cancelled = true;
      if (reconnectTimer) window.clearTimeout(reconnectTimer);
      socketRef.current?.close();
      socketRef.current = null;
    };
  }, [token]);

  useEffect(() => {
    const socket = socketRef.current;
    if (!socket || socket.readyState !== WebSocket.OPEN) return;

    const desired = new Set(symbols);
    const current = subscribedRef.current;

    for (const symbol of desired) {
      if (!current.has(symbol)) {
        socket.send(JSON.stringify({ action: "subscribe", symbol }));
        current.add(symbol);
      }
    }
    for (const symbol of current) {
      if (!desired.has(symbol)) {
        socket.send(JSON.stringify({ action: "unsubscribe", symbol }));
        current.delete(symbol);
      }
    }
    // symbols is compared by identity via its joined key below.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [symbols.join(",")]);

  return { prices, connected, syncedTime };
}
