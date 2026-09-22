"use client";

import { useEffect, type ReactNode } from "react";
import { useAuthStore } from "@/stores/authStore";
import { Toaster } from "@/components/ui/sonner";

export function Providers({ children }: { children: ReactNode }) {
  useEffect(() => {
    // Imperative store access (not the hook) — a one-time hydration side
    // effect, not something a component needs to re-render on.
    useAuthStore.getState().hydrate();
  }, []);

  return (
    <>
      {children}
      <Toaster position="bottom-right" />
    </>
  );
}
