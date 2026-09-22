import { PriceEvent } from "@/types";
import { create } from "zustand";

interface MarketStore {
  prices: Record<string, PriceEvent> | null;
  symbols: string[];
  setPrices: (prices: Record<string, PriceEvent>) => void;
  setSymbols: (symbols: string[]) => void;
}

export const useMarketStore = create<MarketStore>((set) => ({
  prices: null,
  symbols: [],

  setPrices: (prices) => {
    set({ prices });
  },

  setSymbols: (symbols) => {
    set({ symbols });
  },
}));
