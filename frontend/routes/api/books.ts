// routes/api/books.ts
import { Handlers } from "$fresh/server.ts";
import { BooksResponse } from "../../types/book.ts";

const API_BASE_URL = Deno.env.get("API_BASE_URL") || "https://api.example.com";
const TOKEN = Deno.env.get("TOKEN") || "";

export const handler: Handlers = {
  GET: async (req) => {
    const url = new URL(req.url);
    const title = url.searchParams.get("title");
    const isbn = url.searchParams.get("isbn");

    let apiUrl = `${API_BASE_URL}/books`;

    if (isbn) {
      // ISBN検索
      apiUrl = `${API_BASE_URL}/book:${encodeURIComponent(isbn)}`;
    } else if (title) {
      // タイトル検索
      apiUrl = `${API_BASE_URL}/books/search:${encodeURIComponent(title)}`;
    }

    const res = await fetch(apiUrl, {
      headers: { "Authorization": TOKEN },
    });

    const data: BooksResponse = await res.json();
    return new Response(JSON.stringify(data), { headers: { "Content-Type": "application/json" } });
  },

  DELETE: async (req) => {
    const body = await req.json();
    const res = await fetch(`${API_BASE_URL}/book:${body.isbn}`, {
      method: "DELETE",
      headers: { "Authorization": TOKEN },
    });
    return new Response(res.status === 200 ? "Deleted" : "Failed", { status: res.status });
  },
};
