import { useState } from "preact/hooks";
import { Book } from "../types/book.ts";
import DeleteConfirmationModal from "./DeleteConfirmationModal.tsx";

interface Props {
  book: Book;
  apiBaseUrl: string;
  onDelete?: (isbn: string) => void;
  onRequestDelete?: (isbn: string) => Promise<boolean>;
  onRequestRename?: (isbn: string, newTitle: string) => Promise<boolean>;
}

export default function BookCard({ book, apiBaseUrl, onDelete, onRequestDelete, onRequestRename }: Props) {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [newTitle, setNewTitle] = useState(book.title);

  const handleDeleteRequest = () => {
    if (onRequestDelete) {
      setIsModalOpen(true);
    }
  };

  const handleConfirmDelete = async () => {
    if (onRequestDelete) {
      const success = await onRequestDelete(book.isbn);
      if (success) {
        onDelete?.(book.isbn);
      }
    }
    setIsModalOpen(false);
  };

  const handleRenameRequest = async () => {
    if (!onRequestRename || newTitle.trim() === "" || newTitle === book.title) {
      setIsEditing(false);
      setNewTitle(book.title); // Reset title on cancel or invalid input
      return;
    }
    const success = await onRequestRename(book.isbn, newTitle);
    if (success) {
      setIsEditing(false);
    }
    // If it fails, the BookManager will show an error, and we keep the editing state
    // for the user to retry or cancel.
  };

  // YYYY-MM-DD or YYYY-MM or YYYY to YYYY-MM
  const formatPublishDate = (dateStr: string) => {
    if (!dateStr) return "";
    const parts = dateStr.split('-');
    if (parts.length >= 2) {
      return `${parts[0]}-${parts[1]}`;
    }
    return parts[0];
  };

  // バックエンドの画像URLを直接参照するように変更
  const imageUrl = book.imageurl
    ? `${apiBaseUrl}/images/${book.imageurl}` : null;

  return (
    <div class="bookcard-horizontal">
      <div class="bookcard-horizontal-img">
        {imageUrl ? (
          <img src={imageUrl} alt={`Cover of ${book.title}`} class="bookcard-horizontal-img-el" />
        ) : (
          <div class="bookcard-img-placeholder">
            <span>No Image</span>
          </div>
        )}
      </div>
      <div class="bookcard-horizontal-main">
        <div class="bookcard-horizontal-header">
          {isEditing ? (
            <div class="bookcard-title-edit">
              <input
                type="text"
                value={newTitle}
                onInput={(e) => setNewTitle((e.target as HTMLInputElement).value)}
                onKeyDown={(e) => e.key === 'Enter' && handleRenameRequest()}
                class="searchform-blue-input"
              />
              <button onClick={handleRenameRequest} class="bookcard-save-btn">Save</button>
              <button
                onClick={() => {
                  setIsEditing(false);
                  setNewTitle(book.title);
                }}
                class="bookcard-cancel-btn"
              >
                Cancel
              </button>
            </div>
          ) : (
            <h3 class="bookcard-horizontal-title" onClick={() => onRequestRename && setIsEditing(true)}>
              {book.title}
            </h3>
          )}
          <p class="bookcard-horizontal-authors">{book.authors?.join(", ")}</p>
        </div>
        <div class="bookcard-horizontal-desc">{book.description}</div>
        <div class="bookcard-horizontal-footer">
          <div class="bookcard-horizontal-meta">
            <span class="bookcard-horizontal-isbn">ISBN: {book.isbn}</span>
            <span class="bookcard-horizontal-date">{formatPublishDate(book.publishdate)}</span>
          </div>
          {onRequestDelete && (
            <button onClick={handleDeleteRequest} class="bookcard-delete-btn">Delete</button>
          )}
        </div>
      </div>
      {isModalOpen && (
        <DeleteConfirmationModal onConfirm={handleConfirmDelete} onCancel={() => setIsModalOpen(false)} />
      )}
    </div>
  );
}