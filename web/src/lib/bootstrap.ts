import { api } from "./api";
import type { Account } from "@/types";

// GET-then-create: the backend has no "get or create" endpoint, so every
// page that needs the user's account/watchlist to exist does this same
// idempotent dance (a 500 on the GET just means "doesn't exist yet").
export async function bootstrapAccount(): Promise<Account> {
  try {
    return await api.getAccount();
  } catch {
    return await api.createAccount();
  }
}

export async function bootstrapWatchlist(): Promise<void> {
  try {
    await api.getWatchlist();
  } catch {
    await api.createWatchlist();
  }
}
