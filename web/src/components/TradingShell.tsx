"use client";

import { useAuthStore } from "@/stores/authStore";
import { Sidebar } from "./dashboard/Sidebar";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { useTradingStore } from "@/stores/tradingStore";
import { Topbar } from "./dashboard/Topbar";
import { useMarketSocket } from "@/hooks/useMarketSocket";
import { useMarketStore } from "@/stores/marketStore";
import { errorMessage } from "@/lib/errors";
import { PageLoading } from "./PageLoading";

export function TradingShell({ children }: { children: React.ReactNode }) {
  const { email, token, isHydrated, logout } = useAuthStore();
  const { status, error, initialize, reset } = useTradingStore();
  const { symbols, setPrices, setSymbols } = useMarketStore();
  const router = useRouter();

  const [mobileNavOpen, setMobileNavOpen] = useState(false);

  const { prices, connected, syncedTime } = useMarketSocket(symbols, token);

  const handleOpenMobileNav = useCallback(() => setMobileNavOpen(true), []);

  const handleSignOut = useCallback(() => {
    logout();
    reset();
    setPrices({});
    setSymbols([]);
    router.replace("/login");
  }, [logout, reset, router, setPrices, setSymbols]);

  useEffect(() => {
    if (isHydrated && !token) router.replace("/login");
  }, [isHydrated, token, router]);

  useEffect(() => {
    if (token) void initialize();
  }, [token, initialize]);

  useEffect(() => {
    setPrices(prices);
  }, [prices, setPrices]);

  if (!isHydrated || !token) return null;

  return (
    <div className="app">
      <Sidebar
        email={email}
        onSignOut={handleSignOut}
        mobileOpen={mobileNavOpen}
        onMobileOpenChange={setMobileNavOpen}
      />

      <div className="workspace">
        <Topbar
          connected={connected}
          syncedTime={syncedTime}
          onOpenMobileNav={handleOpenMobileNav}
          avatarInitials={email ? email.slice(0, 2).toUpperCase() : "?"}
        />

        {status === "error" ? (
          <div>
            <p>{errorMessage(error, "Failed to load your account.")}</p>
            <button onClick={() => void initialize()}>Retry</button>
          </div>
        ) : status !== "ready" ? (
          <PageLoading label="account" />
        ) : (
          children
        )}
      </div>
    </div>
  );
}
