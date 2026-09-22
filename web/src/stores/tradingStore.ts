import { api } from "@/lib/api";
import { bootstrapAccount, bootstrapWatchlist } from "@/lib/bootstrap";
import { Account } from "@/types";
import { create } from "zustand";

interface TradingStore {
  account: Account | null;
  status: "idle" | "loading" | "ready" | "error";
  error: unknown;
  initialize: () => Promise<void>;
  refreshAccount: () => Promise<void>;
  reset: () => void;
}

let sessionVersion = 0;

export const useTradingStore = create<TradingStore>((set, get) => ({
  account: null,
  status: "idle",
  error: null,

  initialize: async () => {
    const { status } = get();
    if (status === "loading" || status === "ready") return;
    set({ status: "loading", error: null });

    const requestVersion = sessionVersion;

    try {
      const [account] = await Promise.all([
        bootstrapAccount(),
        bootstrapWatchlist(),
      ]);
      if (requestVersion !== sessionVersion) return;
      set({ account, status: "ready" });
    } catch (error) {
      if (requestVersion !== sessionVersion) return;
      set({ error, status: "error" });
    }
  },

  refreshAccount: async () => {
    const requestVersion = sessionVersion;

    try {
      const account = await api.getAccount();
      if (requestVersion !== sessionVersion) return;
      set({ account, status: "ready" });
    } catch (error) {
      if (requestVersion !== sessionVersion) return;
      set({ error, status: "error" });
    }
  },

  reset: () => {
    sessionVersion += 1;
    set({ account: null, status: "idle", error: null });
  },
}));
