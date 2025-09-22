import { useEffect, useState } from "preact/hooks";
import { BookInfo } from "../types/book.ts";
import { API_BASE_URL } from "../utils/api.ts";

export default function BookCard({ book, onDelete }: { book: BookInfo, onDelete?: (isbn: string) => void }) {
  const title = book.Title;
  const authors = book.Authors;
  const description = book.Description;
  const imageFileName = book.Image?.Path;
  const isbn = book.ISBN;
  const publishdate = book.Publishdate;

  const [imageUrl, setImageUrl] = useState<string>("");

  useEffect(() => {
    if (imageFileName) {
      setImageUrl(`${API_BASE_URL}/images/${imageFileName}`);
    } else {
      setImageUrl("");
    }
  }, [imageFileName]);

  let yearMonth = "";
  if (publishdate) {
    const [year, month] = publishdate.split("-");
    if (year && month) {
      yearMonth = `${year}-${month}`;
    }
  }

  const handleDelete = async () => {
    if (!isbn) return;
    const res = await fetch(`${API_BASE_URL}/book:${isbn}`, {
      method: "DELETE",
    });
    if (res.ok && onDelete) {
      onDelete(isbn);
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
      <div class="bookcard-horizontal-main">
        <div class="bookcard-horizontal-header">
          <h3 class="bookcard-horizontal-title">{title}</h3>
          <p class="bookcard-horizontal-authors">{Array.isArray(authors) ? authors.join(", ") : authors}</p>
        </div>
        <div class="bookcard-horizontal-desc">{description}</div>
        <div class="bookcard-horizontal-footer">
          <div class="bookcard-horizontal-meta">
            <span class="bookcard-horizontal-isbn">ISBN: {isbn}</span>
            <span class="bookcard-horizontal-date">{yearMonth}</span>
          </div>
          <button class="bookcard-delete-btn" onClick={handleDelete}>
            Delete
          </button>
        </div>
      </div>
    </div>
  );
}