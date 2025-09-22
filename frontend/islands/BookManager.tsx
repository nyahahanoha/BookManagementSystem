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
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState<string>("");

  useEffect(() => {
    loadAllBooks();
  }, []);

  const loadAllBooks = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await BookAPI.getAllBooks();
      setBooks(response.books || []);
      setSearchResults(null);
      setSearchQuery("");
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
    setSearchQuery(query);

    if (query.trim() === "") {
      setSearchResults(null);
      setSearchLoading(false);
      return;
    }

    try {
      let response: BooksResponse;
      if (type === 'isbn') {
        response = await BookAPI.getBookByISBN(query);
      } else {
        response = await BookAPI.searchBooksByTitle(query);
      }
      setSearchResults(response.books || []);
    } catch (err) {
      setError("Search failed. Please try again.");
      console.error("Search error:", err);
    } finally {
      setSearchLoading(false);
    }
  };

  const displayBooks = searchResults === null ? books : searchResults;
  const isEmpty = !displayBooks || displayBooks.length === 0;

  return (
    <div class="bookmanager-bg">
      <div class="bookmanager-header">
        <h1 class="bookmanager-title">Book Management System</h1>
      </div>
      <div class="bookmanager-search">
        <SearchForm onSearch={handleSearch} loading={searchLoading} />
      </div>
      {error && (
        <div class="bookmanager-error">
          <span>{error}</span>
          <button onClick={loadAllBooks} class="bookmanager-error-btn">Retry</button>
        </div>
      )}
      {(loading || searchLoading) && (
        <div class="bookmanager-loading">
          <LoadingSpinner />
        </div>
      )}
      {!loading && !searchLoading && (
        <>
          {isEmpty ? (
            <div class="bookmanager-empty">
              <h3>No books found</h3>
              <p>
                {searchResults !== null && searchQuery.trim() !== ""
                  ? "Try searching with a different keyword or ISBN."
                  : "There are no books in the system yet."}
              </p>
              <button onClick={loadAllBooks} class="bookmanager-empty-btn">
                View All Books
              </button>
            </div>
          ) : (
            <div class="bookmanager-grid">
              {displayBooks.map((book) => (
                <BookCard key={book.isbn || book.ISBN} book={book} />
              ))}
            </div>
          )}
          {!isEmpty && (
            <div class="bookmanager-summary">
              <span>
                Showing <strong>{displayBooks.length}</strong> books
                {searchResults !== null && searchQuery.trim() !== "" && (
                  <> for "<span class="bookmanager-query">{searchQuery}</span>"</>
                )}
              </span>
            </div>
          )}
        </>
      )}
    </div>
  );
}