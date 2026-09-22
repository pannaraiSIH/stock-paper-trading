export interface User {
  id: number;
  email: string;
}

export interface Account {
  id: number;
  cashBalance: number;
  updatedAt: string;
}

export type OrderSide = "buy" | "sell";
export type OrderStatus = "pending" | "executed" | "rejected";

export interface Order {
  id: number;
  symbol: string;
  side: OrderSide;
  executionPrice: number | null;
  totalValue: number | null;
  status: OrderStatus;
}

export interface Position {
  id: number;
  symbol: string;
  quantity: number;
  averagePrice: number;
  updatedAt: string;
}

export interface WatchlistItem {
  id: number;
  symbol: string;
  price: number | null;
  createdAt: string;
}

export interface StockSearchResult {
  symbol: string;
  name: string;
  exchange: string;
  currency: string;
}

export interface StockDetails {
  symbol: string;
  name: string;
  exchange: string;
  micCode: string;
  sector: string | null;
  industry: string | null;
  website: string | null;
  description: string | null;
  type: string | null;
  ceo: string | null;
  address: string | null;
  city: string | null;
  state: string | null;
  country: string | null;
  phone: string | null;
}

export type CandleInterval = "1min" | "5min" | "1h" | "1day";

export interface Candle {
  datetime: string;
  open: string;
  high: string;
  low: string;
  close: string;
  volume: string;
}

export interface PriceEvent {
  event: string;
  symbol: string;
  exchange?: string;
  mic_code?: string;
  type?: string;
  currency?: string;
  timestamp: number;
  price: number;
  day_volume?: number;
  bid?: number;
  ask?: number;
}
