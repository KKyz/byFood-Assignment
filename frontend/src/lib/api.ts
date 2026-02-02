// frontend/src/lib/api.ts

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:8080";

async function apiFetch<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    cache: "no-store",
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  if (!res.ok) {
  let message = `${res.status} ${res.statusText}`;

  // try JSON error: { error: "..." }
  try {
    const data = await res.json();
    if (data?.error) {
      message = `${message}: ${data.error}`;
    }
  } catch {
    // fallback to plain text body (sometimes servers return text/html)
    try {
      const text = await res.text();
      if (text) message = `${message}: ${text}`;
    } catch {
      // ignore
    }
  }

  throw new Error(message);
}


  // 204 No Content
  if (res.status === 204) {
    return null as T;
  }

  return res.json();
}

//Types

export type Book = {
  id: number;
  title: string;
  author: string;
  year: number;
};

export type BookInput = {
  title: string;
  author: string;
  year: number;
};

//API functions

export function listBooks(): Promise<Book[]> {
  return apiFetch<Book[]>("/books");
}

export function getBook(id: number): Promise<Book> {
  return apiFetch<Book>(`/books/${id}`);
}

export function createBook(input: BookInput): Promise<Book> {
  return apiFetch<Book>("/books", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function updateBook(
  id: number,
  input: BookInput
): Promise<Book> {
  return apiFetch<Book>(`/books/${id}`, {
    method: "PUT",
    body: JSON.stringify(input),
  });
}

export function deleteBook(id: number): Promise<void> {
  return apiFetch<void>(`/books/${id}`, {
    method: "DELETE",
  });
}
