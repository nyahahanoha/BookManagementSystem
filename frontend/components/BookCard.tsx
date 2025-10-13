import { Book } from "../types/book.ts";
import { useState, useEffect } from "preact/hooks";

export const LanguageMap = {
  0: "Unknown",
  1: "Japanese",
  2: "English",
} as const;

export default function BookCard({ book, onRequestDelete }: {
  book: Book;
  onDelete?: (isbn: string) => void;
  onRequestDelete?: (isbn: string) => Promise<boolean>;
}) {
  const { title, authors, description, imageurl, isbn, publishdate, language } = book;
  const [deleted, setDeleted] = useState(false);
  const [imageUrl, setImageUrl] = useState<string>("");

  useEffect(() => {
    if (imageurl) {
      setImageUrl(`/api/images?filename=${encodeURIComponent(imageurl)}`);
    } else {
      setImageUrl("");
    }
  }, [imageurl]);

  let yearMonth = "";
  if (publishdate) {
    const [year, month] = publishdate.split("-");
    if (year && month) {
      yearMonth = `${year}-${month}`;
    }
  }

  const handleDelete = async () => {
    if (deleted) return;
    if (!isbn || !onRequestDelete) return;
    const success = await onRequestDelete(isbn);
    if (success) {
      setDeleted(true);
    }
  };

  return (
    <div class="bookcard-horizontal">
      <div class="bookcard-horizontal-img">
        {imageUrl ? (
          <img src={imageUrl} alt="" class="bookcard-horizontal-img-el" />
        ) : (
          <div class="bookcard-img-placeholder">No Image</div>
        )}
      </div>
      <div class_="bookcard-horizontal-main">
        <div class="bookcard-horizontal-header">
          <h3 class="bookcard-horizontal-title">{title}</h3>
          <p class="bookcard-horizontal-authors">
            {Array.isArray(authors) ? authors.join(", ") : authors}
          </p>
        </div>
        <div class="bookcard-horizontal-desc">{description}</div>
        <div class="bookcard-horizontal-footer">
          <div class="bookcard-horizontal-meta">
            <span class="bookcard-horizontal-isbn">ISBN: {isbn}</span>
            <span class="bookcard-horizontal-date">{yearMonth}</span>
            <span class="bookcard-horizontal-date">{LanguageMap[language]}</span>
          </div>
          {onRequestDelete && (
            <button
              class={`bookcard-delete-btn ${deleted ? "deleted" : ""}`}
              onClick={handleDelete}
              disabled={deleted}
            >
              {deleted ? "Deleted" : "Delete"}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}