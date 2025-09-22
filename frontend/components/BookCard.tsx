import { useEffect, useState } from "preact/hooks";
import { BookInfo } from "../types/book.ts";

export default function BookCard({ book, onDelete }: { book: BookInfo, onDelete?: (isbn: string) => void }) {
  const title = book.Title;
  const authors = book.Authors;
  const description = book.Description;
  const imageFileName =  book.Image?.Path;
  const isbn = book.ISBN;
  const publishdate = book.Publishdate;

  const [imageUrl, setImageUrl] = useState<string>("");

  useEffect(() => {
    if (imageFileName) {
      setImageUrl(`http://localhost:8080/images/${imageFileName}`);
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
    const res = await fetch(`http://localhost:8080/book:${isbn}`, {
      method: "DELETE",
    });
    if (res.ok && onDelete) {
      onDelete(isbn);
    }
  };

  return (
    <div class="bookcard-dark">
      <div class="bookcard-img-wrap">
        {imageUrl ? (
          <img src={imageUrl} alt="" class="bookcard-img" />
        ) : (
          <div class="bookcard-img-placeholder">No Image</div>
        )}
      </div>
      <div class="bookcard-content">
        <h3 class="bookcard-title">{title}</h3>
        <p class="bookcard-authors">{Array.isArray(authors) ? authors.join(", ") : authors}</p>
        <p class="bookcard-desc">{description}</p>
        <div class="bookcard-meta">
          <span class="bookcard-isbn">ISBN: {isbn}</span>
          <span class="bookcard-date">{yearMonth}</span>
        </div>
        <button class="bookcard-delete-btn" onClick={handleDelete}>
          Delete
        </button>
      </div>
    </div>
  );
}