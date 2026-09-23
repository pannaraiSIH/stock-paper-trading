import type { Position } from "@/types";

export interface PositionPL {
  current: number;
  marketValue: number;
  costBasis: number;
  pl: number;
  plPct: number;
}

export function calculatePositionPL(
  position: Position,
  livePrice: number | undefined,
): PositionPL {
  const current = livePrice ?? position.averagePrice;
  const marketValue = position.quantity * current;
  const costBasis = position.quantity * position.averagePrice;
  const pl = marketValue - costBasis;
  const plPct = costBasis ? (pl / costBasis) * 100 : 0;
  return { current, marketValue, costBasis, pl, plPct };
}

export interface PortfolioSummary {
  marketValue: number;
  pl: number;
  plPct: number;
}

export function calculatePortfolioSummary(
  positions: Position[],
  livePrices: Record<string, number>,
): PortfolioSummary {
  let marketValue = 0;
  let costBasis = 0;
  for (const position of positions) {
    const pos = calculatePositionPL(position, livePrices[position.symbol]);
    marketValue += pos.marketValue;
    costBasis += pos.costBasis;
  }
  const pl = marketValue - costBasis;
  const plPct = costBasis ? (pl / costBasis) * 100 : 0;
  return { marketValue, pl, plPct };
}
