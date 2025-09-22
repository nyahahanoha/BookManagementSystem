// components/BookCard.tsx
import { BookInfo, LanguageMap } from "../types/book.ts";

interface BookCardProps {
  book: BookInfo;
}

export default function BookCard({ book }: BookCardProps) {
  const publishDate = new Date(book.Publishdate).toLocaleDateString('ja-JP');
  const language = LanguageMap[book.Language as keyof typeof LanguageMap] || "Unknown";
  
  // Image URL construction
  const imageUrl = book.Image?.Source ? 
    `${book.Image.Source.Scheme}://${book.Image.Source.Host}${book.Image.Source.Path}${book.Image.Source.RawQuery ? '?' + book.Image.Source.RawQuery : ''}` :
    '/placeholder-book.png';

  return (
    <div class="bg-white rounded-xl shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-100">
      <div class="md:flex">
        <div class="md:flex-shrink-0">
          <img
            class="h-48 w-full object-cover md:h-full md:w-48"
            src={imageUrl}
            alt={book.Title}
            onError={(e) => {
              (e.target as HTMLImageElement).src = '/placeholder-book.png';
            }}
          />
        </div>
        <div class="p-6 flex flex-col justify-between">
          <div class="flex-grow">
            <div class="flex items-center justify-between mb-2">
              <span class="inline-block bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full uppercase font-semibold tracking-wide">
                ISBN: {book.ISBN}
              </span>
              <span class="text-sm text-gray-500">{language}</span>
            </div>
            
            <h3 class="text-xl font-bold text-gray-900 mb-2 line-clamp-2">
              {book.Title}
            </h3>
            
            <div class="mb-3">
              <p class="text-sm font-medium text-gray-700 mb-1">Authors:</p>
              <p class="text-sm text-gray-600">
                {book.Authors && book.Authors.length > 0 ? book.Authors.join(", ") : "Unknown"}
              </p>
            </div>
            
            {book.Description && (
              <p class="text-gray-600 text-sm line-clamp-3 mb-4">
                {book.Description}
              </p>
            )}
          </div>
          
          <div class="flex items-center justify-between pt-4 border-t border-gray-100">
            <span class="text-sm text-gray-500">
              Published: {publishDate}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}