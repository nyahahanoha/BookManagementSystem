interface Props {
  onConfirm: () => void;
  onCancel: () => void;
}

export default function DeleteConfirmationModal({ onConfirm, onCancel }: Props) {
  return (
    <div class="modal-overlay">
      <div class="modal-content">
        <p class="modal-text">Do you really want to delete this book?</p>
        <div class="modal-actions">
          <button onClick={onCancel} class="modal-button cancel">
            No
          </button>
          <button onClick={onConfirm} class="modal-button confirm">
            Yes
          </button>
        </div>
      </div>
    </div>
  );
}