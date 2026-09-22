import { create } from "zustand";
import { api } from "@/lib/api";
import { clearStoredToken, getStoredToken, setStoredToken } from "./auth-storage";

const EMAIL_KEY = "paperline.email";

function getStoredEmail(): string | null {
  try {
    return window.localStorage.getItem(EMAIL_KEY);
  } catch {
    return null;
  }
}

interface AuthState {
  token: string | null;
  email: string | null;
  isHydrated: boolean;
  hydrate: () => void;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  token: null,
  email: null,
  isHydrated: false,

  // One-time hydration from localStorage, which isn't available during SSR.
  // Called imperatively (not via the hook) from Providers on mount.
  hydrate: () => {
    if (get().isHydrated) return;
    set({ token: getStoredToken(), email: getStoredEmail(), isHydrated: true });
  },

  login: async (loginEmail, password) => {
    const { accessToken } = await api.login(loginEmail, password);
    setStoredToken(accessToken);
    try {
      window.localStorage.setItem(EMAIL_KEY, loginEmail);
    } catch {
      // ignore
    }
    set({ token: accessToken, email: loginEmail });
  },

  register: async (registerEmail, password) => {
    await api.register(registerEmail, password);
    await get().login(registerEmail, password);
  },

  logout: () => {
    clearStoredToken();
    try {
      window.localStorage.removeItem(EMAIL_KEY);
    } catch {
      // ignore
    }
    set({ token: null, email: null });
  },
}));
