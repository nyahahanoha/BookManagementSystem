import { Head } from "$fresh/runtime.ts";
import BookManager from "../islands/BookManager.tsx";

export default function Home() {
  return (
    <>
      <Head>
        <title>Book Management System</title>
        <meta name="description" content="Modern book management and search system" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      
      <div class="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50">
        <div class="container mx-auto px-4 py-8">
          {/* Header */}
          <header class="text-center mb-10">
            <h1 class="text-5xl font-bold text-gray-800 mb-4 bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
              📖 Book Management System
            </h1>
            <p class="text-xl text-gray-600 max-w-2xl mx-auto">
              Discover, search, and manage your book collection with ease
            </p>
          </header>

          <BookManager />
        </div>
      </div>
    </>
  );
}