import { useEffect, useState } from "preact/hooks";
import { BookInfo } from "../types/book.ts";

export default function BookCard({ book }: { book: BookInfo }) {
  const title = book.title || book.Title;
  const authors = book.authors || book.Authors;
  const description = book.description || book.Description;
  const imageFileName = book.image?.Path || book.Image?.Path;
  const isbn = book.isbn || book.ISBN;
  const publishdate = book.publishdate || book.Publishdate;

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
    const [year, month] = publishdate.split("-"); // "YYYY-MM-DD"
    if (year && month) {
      yearMonth = `${year}-${month}`;
    }
  }

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
      </div>
    </div>
  );
}