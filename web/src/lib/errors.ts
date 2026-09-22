import { ApiError } from "./api";

// Never surface raw backend error text to the user by accident — callers
// pass a sane fallback for anything that isn't a recognized ApiError.
export function errorMessage(err: unknown, fallback: string): string {
  return err instanceof ApiError ? err.message : fallback;
}
