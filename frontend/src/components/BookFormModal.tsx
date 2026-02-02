"use client";

import { useEffect, useMemo, useState } from "react";
import type { Book, BookInput } from "@/lib/api";

type Props = {
  open: boolean;
  mode: "create" | "edit";
  initial: Book | null;
  onClose: () => void;
  onSubmit: (input: BookInput) => Promise<void>;
};

export default function BookFormModal({
  open,
  mode,
  initial,
  onClose,
  onSubmit,
}: Props) {
  const titleText = mode === "create" ? "Add book" : "Edit book";
  const submitText = mode === "create" ? "Create" : "Save";

  const [title, setTitle] = useState("");
  const [author, setAuthor] = useState("");
  const [year, setYear] = useState<string>("");

  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!open) return;

    setError(null);

    if (mode === "edit" && initial) {
      setTitle(initial.title);
      setAuthor(initial.author);
      setYear(String(initial.year));
    } else {
      setTitle("");
      setAuthor("");
      setYear("");
    }
  }, [open, mode, initial]);

  // Escape to close + lock background scroll while modal open
  useEffect(() => {
    if (!open) return;

    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };

    window.addEventListener("keydown", onKeyDown);

    const prevOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";

    return () => {
      window.removeEventListener("keydown", onKeyDown);
      document.body.style.overflow = prevOverflow;
    };
  }, [open, onClose]);

  const yearNum = useMemo(() => Number(year), [year]);

  function validate(): string | null {
    const t = title.trim();
    const a = author.trim();

    if (!t) return "Title is required.";
    if (!a) return "Author is required.";
    if (!year.trim()) return "Year is required.";
    if (!Number.isFinite(yearNum) || !Number.isInteger(yearNum))
      return "Year must be an integer.";
    if (yearNum < 0) return "Year must be 0 or greater.";
    if (yearNum > 3000) return "Year looks too large.";
    return null;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (submitting) return;

    const msg = validate();
    if (msg) {
      setError(msg);
      return;
    }

    setError(null);
    setSubmitting(true);

    try {
      await onSubmit({
        title: title.trim(),
        author: author.trim(),
        year: yearNum,
      });
      onClose();
    } catch (e: any) {
      setError(e?.message ?? "Request failed");
    } finally {
      setSubmitting(false);
    }
  }

  // Don’t render anything if closed
  if (!open) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4"
      aria-modal="true"
      role="dialog"
      aria-label={titleText}
    >
      {/* Backdrop */}
      <button
        className="absolute inset-0 bg-black/30"
        onClick={onClose}
        aria-label="Close modal"
        type="button"
      />

      {/* Modal */}
      <div className="relative z-10 w-full max-w-lg rounded-2xl bg-white p-6 shadow-lg">
        <div className="mb-4 flex items-start justify-between gap-4">
          <div>
            <h2 className="text-xl font-bold">{titleText}</h2>
            <p className="text-sm text-gray-600">
              {mode === "create"
                ? "Create a new book entry."
                : "Update the selected book."}
            </p>
          </div>

          <button
            className="rounded-lg border px-3 py-1 hover:bg-gray-50"
            onClick={onClose}
            type="button"
            disabled={submitting}
          >
            ✕
          </button>
        </div>

        {error && (
          <div className="mb-4 rounded-lg border border-red-300 bg-red-50 px-4 py-3 text-sm text-red-800">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="mb-1 block text-sm font-medium">Title</label>
            <input
              className="w-full rounded-lg border px-3 py-2"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Dune"
              disabled={submitting}
              autoFocus
            />
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium">Author</label>
            <input
              className="w-full rounded-lg border px-3 py-2"
              value={author}
              onChange={(e) => setAuthor(e.target.value)}
              placeholder="Frank Herbert"
              disabled={submitting}
            />
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium">Year</label>
            <input
              className="w-full rounded-lg border px-3 py-2"
              value={year}
              onChange={(e) => setYear(e.target.value)}
              placeholder="1965"
              inputMode="numeric"
              disabled={submitting}
            />
          </div>

          <div className="flex items-center justify-end gap-2 pt-2">
            <button
              className="rounded-lg border px-4 py-2 hover:bg-gray-50"
              onClick={onClose}
              type="button"
              disabled={submitting}
            >
              Cancel
            </button>

            <button
              className="rounded-lg bg-black px-4 py-2 text-white disabled:opacity-60"
              type="submit"
              disabled={submitting}
            >
              {submitting ? "Saving…" : submitText}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
