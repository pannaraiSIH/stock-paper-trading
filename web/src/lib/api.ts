import { getStoredToken } from "@/stores/auth-storage";
import { useAuthStore } from "@/stores/authStore";
import { useMarketStore } from "@/stores/marketStore";
import { useTradingStore } from "@/stores/tradingStore";
import type {
  Account,
  Candle,
  CandleInterval,
  Order,
  OrderSide,
  Position,
  StockDetails,
  StockSearchResult,
  User,
  WatchlistItem,
} from "@/types";

export class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

interface Envelope<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export const PAGE_SIZE = 20;

export interface Page<T> {
  items: T[];
  hasMore: boolean;
}

async function request<T>(
  path: string,
  options: { method?: string; body?: unknown; auth?: boolean } = {},
): Promise<T> {
  const { method = "GET", body, auth = true } = options;

  const headers: Record<string, string> = {};
  if (body !== undefined) headers["Content-Type"] = "application/json";
  if (auth) {
    const token = getStoredToken();
    if (token) headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`/api${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  if (auth && res.status === 401) {
    useAuthStore.getState().logout();
    useTradingStore.getState().reset();
    useMarketStore.getState().setPrices({});
    useMarketStore.getState().setSymbols([]);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  let payload: Envelope<T> | undefined;
  try {
    payload = (await res.json()) as Envelope<T>;
  } catch {
    throw new ApiError(`request failed (${res.status})`, res.status);
  }

  if (!res.ok || !payload.success) {
    throw new ApiError(
      payload.error ?? `request failed (${res.status})`,
      res.status,
    );
  }

  return (payload.data ?? (null as T)) as T;
}

// Requests one row beyond `limit` so `hasMore` can be derived without a
// separate count query — the extra row is dropped before returning.
async function requestPage<T>(
  path: string,
  limit: number,
  offset: number,
): Promise<Page<T>> {
  const query = new URLSearchParams({
    limit: String(limit + 1),
    offset: String(offset),
  });
  const data = await request<T[] | null>(`${path}?${query}`);
  const items = data ?? [];
  return { items: items.slice(0, limit), hasMore: items.length > limit };
}

export const api = {
  register(email: string, password: string) {
    return request<User>("/auth/register", {
      method: "POST",
      body: { email, password },
      auth: false,
    });
  },

  login(email: string, password: string) {
    return request<{ accessToken: string }>("/auth/login", {
      method: "POST",
      body: { email, password },
      auth: false,
    });
  },

  getAccount() {
    return request<Account>("/account");
  },

  createAccount() {
    return request<Account>("/account", { method: "POST" });
  },

  async getPositions() {
    const data = await request<Position[] | null>(
      "/positions?limit=100&offset=0",
    );
    return data ?? [];
  },

  getPositionsPage(limit: number, offset: number) {
    return requestPage<Position>("/positions", limit, offset);
  },

  async getOrders() {
    const data = await request<Order[] | null>("/orders?limit=100&offset=0");
    return data ?? [];
  },

  getOrdersPage(limit: number, offset: number) {
    return requestPage<Order>("/orders", limit, offset);
  },

  createOrder(symbol: string, side: OrderSide, quantity: number) {
    return request<Order>("/orders", {
      method: "POST",
      body: { symbol, side, quantity },
    });
  },

  getWatchlist() {
    return request<{ id: number }>("/watchlist");
  },

  createWatchlist() {
    return request<{ id: number }>("/watchlist", { method: "POST" });
  },

  async getWatchlistItems() {
    const data = await request<WatchlistItem[] | null>(
      "/watchlist/items?limit=100&offset=0",
    );
    return data ?? [];
  },

  getWatchlistItemsPage(limit: number, offset: number) {
    return requestPage<WatchlistItem>("/watchlist/items", limit, offset);
  },

  addWatchlistItem(symbol: string) {
    return request<WatchlistItem>("/watchlist/items", {
      method: "POST",
      body: { symbol },
    });
  },

  deleteWatchlistItem(itemId: number) {
    return request<void>(`/watchlist/items/${itemId}`, { method: "DELETE" });
  },

  async searchStocks(query: string, outputSize = 8) {
    const params = new URLSearchParams({
      query,
      outputSize: String(outputSize),
    });
    const data = await request<StockSearchResult[] | null>(
      `/market/stocks?${params}`,
    );
    return data ?? [];
  },

  getStockDetails(symbol: string) {
    const params = new URLSearchParams({ symbol });
    return request<StockDetails>(`/market/stocks/details?${params}`);
  },

  async getCandles(
    symbol: string,
    interval: CandleInterval,
    outputSize: number,
  ) {
    // symbol is a query param, not a path segment — some symbols (crypto
    // pairs like BTC/USD) contain a literal "/" that breaks path routing.
    const params = new URLSearchParams({
      symbol,
      interval,
      outputSize: String(outputSize),
    });
    const data = await request<Candle[] | null>(
      `/market/stocks/candles?${params}`,
    );
    return data ?? [];
  },
};
