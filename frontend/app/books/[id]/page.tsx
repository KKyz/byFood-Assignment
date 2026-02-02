"use client";

import Link from "next/link";
import { use, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getBook, updateBook, deleteBook, type Book, type BookInput } from "@/lib/api";
import ErrorBanner from "@/components/ErrorBanner";
import BookFormModal from "@/components/BookFormModal";

export default function Page({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const idNum = Number(id);
  const router = useRouter();

  const [book, setBook] = useState<Book | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);

  async function load() {
    setActionError(null);
    setError(null);
    setBook(null);

    if (!Number.isFinite(idNum) || !Number.isInteger(idNum) || idNum <= 0) {
      setError("Invalid book id");
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      const b = await getBook(idNum);
      setBook(b);
    } catch (e: any) {
      setError(e.message ?? "Failed to load book");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [idNum]);

  async function handleDelete() {
    if (!book) return;
    setActionError(null);

    const ok = confirm(`Delete "${book.title}"?`);
    if (!ok) return;

    try {
      await deleteBook(book.id);
      router.push("/");
      router.refresh();
    } catch (e: any) {
      setActionError(e.message ?? "Failed to delete");
    }
  }

  return (
    <main className="mx-auto max-w-3xl p-6">
      <div className="mb-6 flex items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold">Book details</h1>
          <p className="text-sm text-gray-600">/books/{id}</p>
        </div>

        <div className="flex gap-2">
          <Link
            href="/"
            className="rounded-lg border px-4 py-2 hover:bg-gray-50"
          >
            ← Back
          </Link>
        </div>
      </div>

      {error && <ErrorBanner message={error} />}
      {actionError && <ErrorBanner message={actionError} />}

      {loading ? (
        <div className="rounded-lg border p-6 text-gray-700">Loading…</div>
      ) : !book ? (
        <div className="rounded-lg border p-6 text-gray-700">Book not found.</div>
      ) : (
        <div className="rounded-2xl border p-6">
          <div className="space-y-3">
            <div>
              <div className="text-sm text-gray-500">Title</div>
              <div className="text-lg font-semibold">{book.title}</div>
            </div>

            <div>
              <div className="text-sm text-gray-500">Author</div>
              <div className="text-lg">{book.author}</div>
            </div>

            <div>
              <div className="text-sm text-gray-500">Year</div>
              <div className="text-lg">{book.year}</div>
            </div>

            <div className="pt-2 text-sm text-gray-500">
              ID: <span className="font-mono">{book.id}</span>
            </div>

            <div className="flex gap-2 pt-4">
              <button
                className="rounded-lg border px-4 py-2 hover:bg-gray-50"
                onClick={() => setModalOpen(true)}
              >
                Edit
              </button>

              <button
                className="rounded-lg border border-red-300 px-4 py-2 text-red-700 hover:bg-red-50"
                onClick={handleDelete}
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}

      <BookFormModal
        open={modalOpen}
        mode="edit"
        initial={book}
        onClose={() => setModalOpen(false)}
        onSubmit={async (input: BookInput) => {
          if (!book) throw new Error("No book loaded");
          await updateBook(book.id, input);
          await load(); // refresh detail view
        }}
      />
    </main>
  );
}
