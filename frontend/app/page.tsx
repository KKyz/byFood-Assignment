"use client";

import Link from "next/link";
import { useMemo, useState } from "react";
import BookFormModal from "@/components/BookFormModal";
import ErrorBanner from "@/components/ErrorBanner";
import { useBooks } from "@/context/BooksContext";
import type { Book } from "@/lib/api";

export default function Page() {
  const { books, loading, error, refresh, addBook, editBook, removeBook } = useBooks();

  const [modalOpen, setModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState<"create" | "edit">("create");
  const [editing, setEditing] = useState<Book | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);

  const sortedBooks = useMemo(() => {
    return [...books].sort((a, b) => a.id - b.id);
  }, [books]);

  function openCreate() {
    setActionError(null);
    setModalMode("create");
    setEditing(null);
    setModalOpen(true);
  }

  function openEdit(b: Book) {
    setActionError(null);
    setModalMode("edit");
    setEditing(b);
    setModalOpen(true);
  }

  async function handleDelete(b: Book) {
    setActionError(null);
    const ok = confirm(`Delete "${b.title}"?`);
    if (!ok) return;

    try {
      await removeBook(b.id);
    } catch (e: any) {
      setActionError(e.message ?? "Failed to delete");
    }
  }

  return (
    <main className="mx-auto max-w-4xl p-6">
      <div className="mb-6 flex items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold">Books</h1>
          <p className="text-sm text-gray-600">
            Manage your library (CRUD). Backend: Go + SQLite.
          </p>
        </div>

        <div className="flex gap-2">
          <button
            className="rounded-lg border px-4 py-2 hover:bg-gray-50"
            onClick={() => refresh()}
            disabled={loading}
          >
            Refresh
          </button>
          <button
            className="rounded-lg bg-black px-4 py-2 text-white"
            onClick={openCreate}
          >
            + Add book
          </button>
        </div>
      </div>

      {error && <ErrorBanner message={error} />}
      {actionError && <ErrorBanner message={actionError} />}

      {loading ? (
        <div className="rounded-lg border p-6 text-gray-700">Loading…</div>
      ) : sortedBooks.length === 0 ? (
        <div className="rounded-lg border p-6 text-gray-700">
          No books yet. Click <span className="font-semibold">Add book</span>.
        </div>
      ) : (
        <div className="overflow-hidden rounded-lg border">
          <table className="w-full text-left text-sm">
            <thead className="bg-gray-50 text-gray-700">
              <tr>
                <th className="px-4 py-3">Title</th>
                <th className="px-4 py-3">Author</th>
                <th className="px-4 py-3">Year</th>
                <th className="px-4 py-3"></th>
              </tr>
            </thead>
            <tbody>
              {sortedBooks.map((b) => (
                <tr key={b.id} className="border-t">
                  <td className="px-4 py-3 font-medium">{b.title}</td>
                  <td className="px-4 py-3">{b.author}</td>
                  <td className="px-4 py-3">{b.year}</td>
                  <td className="px-4 py-3">
                    <div className="flex justify-end gap-2">
                      <Link
                        className="rounded-lg border px-3 py-1 hover:bg-gray-50"
                        href={`/books/${b.id}`}
                      >
                        View
                      </Link>
                      <button
                        className="rounded-lg border px-3 py-1 hover:bg-gray-50"
                        onClick={() => openEdit(b)}
                      >
                        Edit
                      </button>
                      <button
                        className="rounded-lg border border-red-300 px-3 py-1 text-red-700 hover:bg-red-50"
                        onClick={() => handleDelete(b)}
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <BookFormModal
        open={modalOpen}
        mode={modalMode}
        initial={editing}
        onClose={() => setModalOpen(false)}
        onSubmit={async (input) => {
          if (modalMode === "create") {
            await addBook(input);
          } else {
            if (!editing) throw new Error("No book selected");
            await editBook(editing.id, input);
          }
        }}
      />
    </main>
  );
}
