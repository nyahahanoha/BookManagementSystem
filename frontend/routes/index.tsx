import { Handlers, PageProps } from "$fresh/server.ts";
import BookManager from "../islands/BookManager.tsx";
import { decode } from "https://deno.land/x/djwt@v3.0.1/mod.ts";

const ADMIN_EMAILS = (Deno.env.get("ADMIN_EMAILS") || "example@gmail.com")
  .split(",")
  .map(email => email.trim())


interface Data {
  claims: Record<string, unknown> | null;
}

export const handler: Handlers<Data> = {
  GET(req, ctx) {
    const jwt = req.headers.get("X-Pomerium-Jwt-Assertion");

    let claims: Record<string, unknown> | null = null;
    if (jwt) {
      try {
        const [, payload] = decode<{ [name: string]: string }>(jwt);
        claims = payload;
        console.log("Decoded JWT claims:", claims);
      } catch (e) {
        console.error("JWT decode error:", e);
      }
    }
    return ctx.render({ claims });
  },
};



export default function Index({ data }: PageProps<Data>) {
  const email = data.claims?.email as string | undefined ?? null;
  const canEdit = email !== null && ADMIN_EMAILS.includes(email);

  return (
    <main class="index-dark-bg">
      <BookManager canEdit={canEdit} />
    </main>
  );
}
