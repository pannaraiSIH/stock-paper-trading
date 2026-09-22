import { useMemo, useState, type FormEvent } from "react";
import { fmtQty, fmtUSD } from "@/lib/format";
import { cn } from "@/lib/utils";
import type { OrderSide } from "@/types";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";

interface OrderTicketProps {
  symbol: string;
  price: number | null;
  cash: number;
  ownedQty: number;
  submitting: boolean;
  onSubmit: (side: OrderSide, quantity: number) => Promise<void>;
}

export function OrderTicket({
  symbol,
  price,
  cash,
  ownedQty,
  submitting,
  onSubmit,
}: OrderTicketProps) {
  const [side, setSide] = useState<OrderSide>("buy");
  const [qtyInput, setQtyInput] = useState("");

  const qty = Number(qtyInput);
  const total = price !== null && qty > 0 ? qty * price : 0;

  const error = useMemo(() => {
    if (qtyInput === "") return "";
    if (!Number.isInteger(qty) || qty <= 0)
      return "Enter a whole number of shares greater than zero.";
    if (!symbol) return "Pick a symbol first.";
    if (price === null) return "Waiting for a live price for this symbol.";
    if (side === "buy" && total > cash) {
      const maxQty = Math.floor(cash / price);
      return `You have ${fmtUSD(cash)} available — that buys at most ${maxQty} share${maxQty === 1 ? "" : "s"} of ${symbol}.`;
    }
    if (side === "sell" && qty > ownedQty) {
      return `You own ${fmtQty(ownedQty)} share${ownedQty === 1 ? "" : "s"} of ${symbol}. Reduce the quantity to sell.`;
    }
    return "";
  }, [qtyInput, qty, symbol, price, side, total, cash, ownedQty]);

  const valid = qty > 0 && Number.isInteger(qty) && !error && price !== null;

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!valid) return;
    await onSubmit(side, qty);
    setQtyInput("");
  }

  return (
    <div className="panel panel--ticket">
      <Tabs value={side} onValueChange={(value) => setSide(value as OrderSide)}>
        <TabsList className="mb-4 h-auto w-full rounded-md bg-muted p-0.75">
          <TabsTrigger
            value="buy"
            className="h-9 font-medium data-[state=active]:bg-positive/12 data-[state=active]:text-positive data-[state=active]:shadow-none"
          >
            Buy
          </TabsTrigger>
          <TabsTrigger
            value="sell"
            className="h-9 font-medium data-[state=active]:bg-negative/12 data-[state=active]:text-negative data-[state=active]:shadow-none"
          >
            Sell
          </TabsTrigger>
        </TabsList>
      </Tabs>

      <form className="flex flex-col gap-4" onSubmit={handleSubmit} noValidate>
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="qtyInput">Quantity</Label>
          <Input
            id="qtyInput"
            name="qty"
            type="number"
            min={1}
            step={1}
            inputMode="numeric"
            placeholder="0"
            aria-describedby="qtyHelper"
            aria-invalid={Boolean(error)}
            value={qtyInput}
            onChange={(e) => setQtyInput(e.target.value)}
            disabled={!symbol}
            className="font-mono text-base tabular-nums"
          />
          <p
            id="qtyHelper"
            className={cn(
              "min-h-lh text-xs text-muted-foreground",
              error && "text-destructive",
            )}
          >
            {error || " "}
          </p>
        </div>

        <dl className="ticket__summary">
          <div>
            <dt>Order type</dt>
            <dd>Market</dd>
          </div>
          <div>
            <dt>Est. price</dt>
            <dd className="mono">{price === null ? "—" : fmtUSD(price)}</dd>
          </div>
          <div>
            <dt>Est. total</dt>
            <dd className="mono">{fmtUSD(total)}</dd>
          </div>
          <div>
            <dt>{side === "buy" ? "Cash available" : "Shares owned"}</dt>
            <dd className="mono">
              {side === "buy" ? fmtUSD(cash) : fmtQty(ownedQty)}
            </dd>
          </div>
        </dl>

        <Button
          type="submit"
          disabled={!valid || submitting}
          className={cn(
            "h-11 font-semibold",
            side === "buy"
              ? "bg-positive text-(--color-positive-ink) hover:bg-positive/90"
              : "bg-negative text-(--color-negative-ink) hover:bg-negative/90",
          )}
        >
          {side === "buy" ? `Buy ${symbol || ""}` : `Sell ${symbol || ""}`}
        </Button>
      </form>
    </div>
  );
}
