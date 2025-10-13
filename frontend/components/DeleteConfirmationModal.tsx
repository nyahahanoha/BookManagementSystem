interface Props {
  onConfirm: () => void;
  onCancel: () => void;
}

export default function DeleteConfirmationModal({ onConfirm, onCancel }: Props) {
  return (
    <div class="modal-overlay">
      <div class="modal-content">
        <h2>Confirm Deletion</h2>
        <p>Are you sure you want to delete this book? This action cannot be undone.</p>
        <div class="modal-actions">
          <button onClick={onCancel} class="modal-button cancel">
            Cancel
          </button>
          <button onClick={onConfirm} class="modal-button confirm">
            Delete
          </button>
        </div>
      </div>
    </div>
  );
}