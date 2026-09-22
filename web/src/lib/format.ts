const usd = new Intl.NumberFormat("en-US", {
  style: "currency",
  currency: "USD",
  minimumFractionDigits: 2,
});

export function fmtUSD(n: number): string {
  return usd.format(n);
}

export function fmtSigned(n: number): string {
  const sign = n >= 0 ? "+" : "−";
  return `${sign}${usd.format(Math.abs(n))}`;
}

export function fmtPct(n: number): string {
  const sign = n >= 0 ? "+" : "−";
  return `${sign}${Math.abs(n).toFixed(2)}%`;
}

export function fmtQty(n: number): string {
  return new Intl.NumberFormat("en-US").format(n);
}

export function nowLabel(): string {
  return new Date().toLocaleTimeString("en-US", { hour12: false });
}
