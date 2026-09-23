import { ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";

interface TablePaginationProps {
  page: number;
  pageSize: number;
  itemCount: number;
  hasMore: boolean;
  onPrev: () => void;
  onNext: () => void;
}

export function TablePagination({
  page,
  pageSize,
  itemCount,
  hasMore,
  onPrev,
  onNext,
}: TablePaginationProps) {
  const start = page * pageSize + 1;
  const end = page * pageSize + itemCount;

  return (
    <div className="mt-4 flex items-center justify-between border-t border-border pt-3">
      <span className="font-mono text-[0.68rem] tracking-wide text-muted-foreground uppercase">
        {itemCount === 0 ? "No rows" : `Showing ${start}–${end}`}
      </span>
      <div className="flex items-center gap-1.5">
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={onPrev}
          disabled={page === 0}
        >
          <ChevronLeft className="size-4" />
          Previous
        </Button>
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={onNext}
          disabled={!hasMore}
        >
          Next
          <ChevronRight className="size-4" />
        </Button>
      </div>
    </div>
  );
}
