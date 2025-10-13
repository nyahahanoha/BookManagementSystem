import { Handlers, PageProps } from "$fresh/server.ts";
import BookManager from "../islands/BookManager.tsx";
import { createRemoteJWKSet, jwtVerify } from "jose";

const ADMIN_EMAILS = (Deno.env.get("ADMIN_EMAILS") || "example@gmail.com")
  .split(",")
  .map(email => email.trim());
const JWKS_URL = (Deno.env.get("JWKS_URL") || "https://auth.example.com/.well-known/pomerium/jwks.json");
const jwks = createRemoteJWKSet(new URL(JWKS_URL));

interface Data {
  claims: Record<string, unknown> | null;
}

export const handler: Handlers<Data> = {
  async GET(req, ctx) {
    const jwt = req.headers.get("X-Pomerium-Jwt-Assertion");

    let claims: Record<string, unknown> | null = null;
    if (jwt) {
      try {
        const { payload } = await jwtVerify(jwt, jwks);
        claims = payload as Record<string, unknown>;
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
