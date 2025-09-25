import { Handlers } from "$fresh/server.ts";

const API_BASE_URL = Deno.env.get("BOOKS_API_BASE_URL") || "http://localhost:8080";
const TOKEN = Deno.env.get("BOOKS_API_TOKEN") || "";

export const handler: Handlers = {
  GET: async (req) => {
    const url = new URL(req.url);
    const filename = url.searchParams.get("filename");

    if (!filename) {
      return new Response("Filename required", { status: 400 });
    }

    try {
      // バックエンド API から画像を取得
      const res = await fetch(`${API_BASE_URL}/images/${filename}`, {
        headers: { "Authorization": TOKEN },
      });

      if (!res.ok) {
        return new Response("Image not found", { status: res.status });
      }

      const arrayBuffer = await res.arrayBuffer();

      return new Response(arrayBuffer, {
        headers: {
          "Content-Type": res.headers.get("Content-Type") || "application/octet-stream",
        },
      });
    } catch (err) {
      console.error("Error fetching image:", err);
      return new Response("Internal Server Error", { status: 500 });
    }
  },
};
