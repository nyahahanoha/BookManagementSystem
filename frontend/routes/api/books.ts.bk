// routes/api/books.ts
import { Handlers } from "$fresh/server.ts";
import { BooksResponse } from "../../types/book.ts";

const API_BASE_URL = Deno.env.get("BOOKS_API_BASE_URL") || "http://localhost:8080";
const TOKEN = Deno.env.get("BOOKS_API_TOKEN") || "";

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
    if (res.status === 200) {
      // 削除成功で 200 → body あり
      return new Response("Deleted", { status: 200 });
    } else if (res.status === 204) {
      // 削除成功で 204 → body なし
      return new Response(null, { status: 204 });
    } else {
      // 失敗
      return new Response("Failed", { status: res.status });
    }
  },
};
