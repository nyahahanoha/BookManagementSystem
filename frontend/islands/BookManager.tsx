// islands/BookManager.tsx
import { useEffect, useState } from "preact/hooks";
import { BookInfo, BooksResponse } from "../types/book.ts";
import { BookAPI } from "../utils/api.ts";
import BookCard from "../components/BookCard.tsx";
import SearchForm from "../components/SearchForm.tsx";
import LoadingSpinner from "../components/LoadingSpinner.tsx";

export default function BookManager() {
  const [books, setBooks] = useState<BookInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchResults, setSearchResults] = useState<BookInfo[] | null>(null);
  const [currentView, setCurrentView] = useState<'all' | 'search'>('all');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadAllBooks();
  }, []);

  const loadAllBooks = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await BookAPI.getAllBooks();
      if (response.Books) {
        setBooks(response.Books);
      } else {
        setBooks([]);
      }
    } catch (err) {
      setError("Failed to load books. Please check if the API server is running.");
      console.error("Error loading books:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async (query: string, type: 'isbn' | 'title') => {
    setSearchLoading(true);
    setError(null);
    
    try {
      let response: BooksResponse;
      
      if (type === 'isbn') {
        response = await BookAPI.getBookByISBN(query);
      } else {
        response = await BookAPI.searchBooksByTitle(query);
      }
      
      setSearchResults(response.Books || []);
      setCurrentView('search');
    } catch (err) {
      setError("Search failed. Please try again.");
      console.error("Search error:", err);
    } finally {
      setSearchLoading(false);
    }
  };

  const displayBooks = currentView === 'search' ? searchResults : books;
  const isEmpty = !displayBooks || displayBooks.length === 0;

  return (
    <>
      {/* Search Form */}
      <SearchForm onSearch={handleSearch} loading={searchLoading} />

      {/* View Toggle */}
      <div class="flex justify-center mb-8">
        <div class="bg-white rounded-lg p-1 shadow-md border border-gray-200">
          <button
            onClick={() => setCurrentView('all')}
            class={`px-6 py-2 rounded-md font-medium transition-all duration-200 ${
              currentView === 'all'
                ? 'bg-blue-500 text-white shadow-sm'
                : 'text-gray-600 hover:text-blue-500'
            }`}
          >
            All Books ({books.length})
          </button>
          <button
            onClick={() => setCurrentView('search')}
            class={`px-6 py-2 rounded-md font-medium transition-all duration-200 ${
              currentView === 'search'
                ? 'bg-blue-500 text-white shadow-sm'
                : 'text-gray-600 hover:text-blue-500'
            }`}
            disabled={!searchResults}
          >
            Search Results ({searchResults?.length || 0})
          </button>
        </div>
      </div>

      {/* Error Message */}
      {error && (
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg mb-6 text-center">
          <strong>Error:</strong> {error}
          <button
            onClick={loadAllBooks}
            class="ml-4 underline hover:no-underline"
          >
            Retry
          </button>
        </div>
      )}

      {/* Loading State */}
      {(loading || searchLoading) && <LoadingSpinner />}

      {/* Books Grid */}
      {!loading && !searchLoading && (
        <>
          {isEmpty ? (
            <div class="text-center py-16">
              <div class="text-6xl mb-4">📚</div>
              <h3 class="text-2xl font-semibold text-gray-700 mb-2">
                {currentView === 'search' ? 'No books found' : 'No books available'}
              </h3>
              <p class="text-gray-500 max-w-md mx-auto">
                {currentView === 'search' 
                  ? 'Try searching with a different keyword or ISBN.'
                  : 'There are no books in the system yet. Start by scanning some books!'
                }
              </p>
              {currentView === 'search' && (
                <button
                  onClick={() => setCurrentView('all')}
                  class="mt-4 px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                >
                  View All Books
                </button>
              )}
            </div>
          ) : (
            <div class="grid gap-6 md:gap-8">
              {displayBooks.map((book) => (
                <BookCard key={book.ISBN} book={book} />
              ))}
            </div>
          )}

          {/* Results Summary */}
          {!isEmpty && (
            <div class="text-center mt-10 py-6 border-t border-gray-200">
              <p class="text-gray-600">
                Showing <span class="font-semibold">{displayBooks.length}</span> 
                {currentView === 'search' ? ' search results' : ' books total'}
              </p>
            </div>
          )}
        </>
      )}
    </>
  );
}