"use client";

import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react";
import {
  type Book,
  type BookInput,
  listBooks,
  createBook,
  updateBook,
  deleteBook,
} from "@/lib/api";

type BooksContextType = {
  books: Book[];
  loading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
  addBook: (input: BookInput) => Promise<void>;
  editBook: (id: number, input: BookInput) => Promise<void>;
  removeBook: (id: number) => Promise<void>;
};

const BooksContext = createContext<BooksContextType | null>(null);

export function BooksProvider({ children }: { children: ReactNode }) {
  const [books, setBooks] = useState<Book[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  async function refresh() {
    try {
      setLoading(true);
      setError(null);
      setBooks(await listBooks());
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  async function addBook(input: BookInput) {
    await createBook(input);
    await refresh();
  }

  async function editBook(id: number, input: BookInput) {
    await updateBook(id, input);
    await refresh();
  }

  async function removeBook(id: number) {
    await deleteBook(id);
    await refresh();
  }

  useEffect(() => {
    refresh();
  }, []);

  return (
    <BooksContext.Provider
      value={{ books, loading, error, refresh, addBook, editBook, removeBook }}
    >
      {children}
    </BooksContext.Provider>
  );
}

export function useBooks() {
  const ctx = useContext(BooksContext);
  if (!ctx) {
    throw new Error("useBooks must be used within BooksProvider");
  }
  return ctx;
}
