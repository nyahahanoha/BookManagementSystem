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
    
    if (query.trim()) {
      onSearch(query.trim(), searchType);
    }
  };

  return (
    <div class="bg-white rounded-xl shadow-lg p-6 mb-8 border border-gray-100">
      <h2 class="text-2xl font-bold text-gray-800 mb-6 text-center">
        📚 Book Search
      </h2>
      
      <form onSubmit={handleSubmit} class="space-y-4">
        <div class="flex flex-col sm:flex-row gap-4">
          <div class="flex-1">
            <input
              type="text"
              name="query"
              placeholder="Enter ISBN or book title..."
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
              required
            />
          </div>
          
          <div class="flex gap-2">
            <select
              name="searchType"
              class="px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
            >
              <option value="title">Title</option>
              <option value="isbn">ISBN</option>
            </select>
            
            <button
              type="submit"
              disabled={loading}
              class={`px-6 py-3 bg-blue-600 text-white rounded-lg font-medium transition-all duration-200 ${
                loading 
                  ? 'opacity-50 cursor-not-allowed' 
                  : 'hover:bg-blue-700 hover:shadow-lg transform hover:scale-105'
              }`}
            >
              {loading ? (
                <div class="flex items-center">
                  <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Searching...
                </div>
              ) : (
                '🔍 Search'
              )}
            </button>
          </div>
        </div>
      </form>
    </div>
  );
}
