export function PageLoading({ label }: { label: string }) {
  return (
    <div className="authpage">
      <p className="watchlist__empty">Loading your {label}&hellip;</p>
    </div>
  );
}
