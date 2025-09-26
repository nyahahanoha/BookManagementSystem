import { BookInfo, LanguageMap } from "../types/book.ts";
import { useState, useEffect } from "preact/hooks";

export default function BookCard({
  book,
  onRequestDelete, // サーバ経由で削除する関数を親から渡す
}: {
  book: BookInfo;
  onDelete?: (isbn: string) => void;
  onRequestDelete?: (isbn: string) => Promise<boolean>;
}) {
  const title = book.Title;
  const authors = book.Authors;
  const description = book.Description;
  const imageFileName = book.Image?.Path;
  const isbn = book.ISBN;
  const publishdate = book.Publishdate;
  const language = LanguageMap[book.Language] || "Unknown";
  const [deleted, setDeleted] = useState(false); // ← ボタン状態

  const [imageUrl, setImageUrl] = useState<string>(imageFileName || "");

  useEffect(() => {
  if (imageFileName) {
    setImageUrl(`/api/images?filename=${encodeURIComponent(imageFileName)}`);
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
    if (deleted) return; // ← 既に削除済みなら何もしない
    if (!isbn || !onRequestDelete) return;
    const success = await onRequestDelete(isbn); // サーバに削除依頼
    if (success) {
      setDeleted(true); // ← 削除成功なら状態更新
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
          <p class="bookcard-horizontal-authors">
            {Array.isArray(authors) ? authors.join(", ") : authors}
          </p>
        </div>
        <div class="bookcard-horizontal-desc">{description}</div>
        <div class="bookcard-horizontal-footer">
          <div class="bookcard-horizontal-meta">
            <span class="bookcard-horizontal-isbn">ISBN: {isbn}</span>
            <span class="bookcard-horizontal-date">{yearMonth}</span>
            <span class="bookcard-horizontal-date">{language}</span>
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
