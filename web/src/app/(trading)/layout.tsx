import { TradingShell } from "@/components/TradingShell";

export default function TradingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <TradingShell>{children}</TradingShell>;
}
