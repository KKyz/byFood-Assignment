"use client";

import { useState } from "react";

export default function ErrorBanner({ message }: { message: string }) {
  const [open, setOpen] = useState(true);
  if (!open) return null;

  return (
    <div className="mb-4 rounded-lg border border-red-300 bg-red-50 px-4 py-3 text-red-800">
      <div className="flex items-start justify-between gap-4">
        <p className="text-sm">{message}</p>
        <button
          className="text-sm font-semibold underline"
          onClick={() => setOpen(false)}
          type="button"
        >
          Dismiss
        </button>
      </div>
    </div>
  );
}
