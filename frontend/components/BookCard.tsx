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

  // バックエンドの画像URLを直接参照するように変更
  const imageUrl = book.imageurl
    ? `${apiBaseUrl}/images/${book.imageurl}` : null;

  return (
    <div class="book-card">
      {imageUrl ? (
        <img src={imageUrl} alt={`Cover of ${book.title}`} class="book-cover" />
      ) : (
        <div class="book-cover book-cover-placeholder">
          <span>No Image</span>
        </div>
      )}
      <div class="book-info">
        {isEditing ? (
          <div class="book-title-edit">
            <input
              type="text"
              value={newTitle}
              onInput={(e) => setNewTitle((e.target as HTMLInputElement).value)}
              onKeyDown={(e) => e.key === 'Enter' && handleRenameRequest()}
              class="book-title-input"
            />
            <button onClick={handleRenameRequest} class="book-title-save">Save</button>
            <button
              onClick={() => {
                setIsEditing(false);
                setNewTitle(book.title);
              }}
              class="book-title-cancel"
            >
              Cancel
            </button>
          </div>
        ) : (
          <h3 class="book-title" onClick={() => onRequestRename && setIsEditing(true)}>
            {book.title}
          </h3>
        )}
        <p class="book-author">{book.authors?.join(", ")}</p>
        <p class="book-isbn">ISBN: {book.isbn}</p>
        {onRequestDelete && (
          <button onClick={handleDeleteRequest} class="delete-button">Delete</button>
        )}
      </div>
      {isModalOpen && (
        <DeleteConfirmationModal onConfirm={handleConfirmDelete} onCancel={() => setIsModalOpen(false)} />
      )}
    </div>
  );
}