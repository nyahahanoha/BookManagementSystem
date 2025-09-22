import { BookInfo } from "../types/book.ts";

export default function BookCard({ book }: { book: BookInfo }) {
  const title = book.title || book.Title;
  const authors = book.authors || book.Authors;
  const description = book.description || book.Description;
  const imagePath = book.image?.Path || book.Image?.Path;
  const isbn = book.isbn || book.ISBN;
  const publishdate = book.publishdate || book.Publishdate;

  return (
    <div class="bookcard-dark">
      <div class="bookcard-img-wrap">
        {imagePath && imagePath !== "" ? (
          <img src={imagePath} alt="" class="bookcard-img" />
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
          <span class="bookcard-date">{publishdate?.slice(0, 10)}</span>
        </div>
      </div>
    </div>
  );
}