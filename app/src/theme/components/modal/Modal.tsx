import React, { ReactNode } from 'react';
import './Modal.css'; // Optional for styling

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  children: ReactNode; // To allow custom content inside the modal
}

const Modal: React.FC<ModalProps> = ({ isOpen, onClose, children }) => {
  if (!isOpen) return null; // Don't render the modal if it's closed

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close-button" onClick={onClose}>
          &times;
        </button>
        {children}
      </div>
    </div>
  );
};

export default Modal;
