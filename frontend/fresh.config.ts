// fresh.config.ts
import { defineConfig } from "$fresh/server.ts";
import tailwind from "$fresh/plugins/tailwind.ts";

const key = await Deno.readTextFile("./cert.key");
const cert = await Deno.readTextFile("./cert.pem");

export default defineConfig({
  key,
  cert,
  plugins: [tailwind()],
});
