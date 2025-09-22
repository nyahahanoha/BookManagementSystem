// components/SearchForm.tsx
import { JSX } from "preact";

interface SearchFormProps {
  onSearch: (query: string, type: 'isbn' | 'title') => void;
  loading: boolean;
}

export default function SearchForm({ onSearch, loading }: SearchFormProps) {
  const handleSubmit = (e: JSX.TargetedEvent<HTMLFormElement, Event>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const query = formData.get("query") as string;
    const searchType = formData.get("searchType") as 'isbn' | 'title';
    onSearch(query.trim(), searchType);
  };

  return (
    <div class="searchform-blue-bg">
      <h2 class="searchform-blue-title">
        Book Search
      </h2>
      <form onSubmit={handleSubmit} class="searchform-blue-form">
        <input
          type="text"
          name="query"
          placeholder="Enter ISBN or book title..."
          class="searchform-blue-input"
        />
        <select
          name="searchType"
          class="searchform-blue-select"
        >
          <option value="title">Title</option>
          <option value="isbn">ISBN</option>
        </select>
        <button
          type="submit"
          disabled={loading}
          class={`searchform-blue-btn${loading ? " searchform-blue-btn-loading" : ""}`}
        >
          {loading ? (
            <span class="searchform-blue-spinner"></span>
          ) : (
            'Search'
          )}
        </button>
      </form>
    </div>
  );
}