import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { fmtUSD } from "@/lib/format";
import type { StockSearchResult } from "@/types";
import {
  CommandDialog,
  Command,
  CommandInput,
  CommandList,
  CommandEmpty,
  CommandGroup,
  CommandItem,
  CommandShortcut,
} from "@/components/ui/command";

interface CommandPaletteProps {
  open: boolean;
  onClose: () => void;
  onSelect: (result: StockSearchResult) => void;
  livePrices: Record<string, number>;
}

export function CommandPalette({ open, onClose, onSelect, livePrices }: CommandPaletteProps) {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<StockSearchResult[]>([]);
  const [searching, setSearching] = useState(false);

  useEffect(() => {
    // Resets the search on every open so a stale query never lingers.
    if (open) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setQuery("");
      setResults([]);
    }
  }, [open]);

  useEffect(() => {
    if (!open) return;
    const trimmed = query.trim();
    if (trimmed.length < 1) {
      // eslint-disable-next-line react-hooks/set-state-in-effect -- clears stale results when the query empties
      setResults([]);
      return;
    }
    setSearching(true);
    const handle = window.setTimeout(async () => {
      try {
        const data = await api.searchStocks(trimmed, 8);
        setResults(data);
      } catch {
        setResults([]);
      } finally {
        setSearching(false);
      }
    }, 250);
    return () => window.clearTimeout(handle);
  }, [query, open]);

  return (
    <CommandDialog
      open={open}
      onOpenChange={(next) => {
        if (!next) onClose();
      }}
      title="Search symbols"
      description="Search for a stock symbol or company name"
    >
      <Command shouldFilter={false}>
        <CommandInput
          value={query}
          onValueChange={setQuery}
          placeholder={"Search symbol or company…"}
        />
        <CommandList>
          <CommandEmpty>
            {searching ? "Searching…" : query.trim() ? "No symbols match your search." : "Type to search a symbol or company."}
          </CommandEmpty>
          <CommandGroup>
            {results.map((result) => (
              <CommandItem
                key={`${result.symbol}-${result.exchange}`}
                value={`${result.symbol}-${result.exchange}`}
                onSelect={() => {
                  onSelect(result);
                  onClose();
                }}
              >
                <span className="font-mono font-semibold text-foreground">{result.symbol}</span>
                <span className="flex-1 truncate text-muted-foreground">
                  {result.name} &middot; {result.exchange}
                </span>
                {livePrices[result.symbol] !== undefined && (
                  <CommandShortcut className="font-mono text-xs tabular-nums text-foreground">
                    {fmtUSD(livePrices[result.symbol])}
                  </CommandShortcut>
                )}
              </CommandItem>
            ))}
          </CommandGroup>
        </CommandList>
      </Command>
    </CommandDialog>
  );
}
