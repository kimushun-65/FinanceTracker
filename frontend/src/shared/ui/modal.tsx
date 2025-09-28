'use client';

import { createContext, useContext, useState, ReactNode } from 'react';
import { cn } from '@/shared/utils';

interface ModalContextType {
  isOpen: boolean;
  openModal: () => void;
  closeModal: () => void;
}

const ModalContext = createContext<ModalContextType | undefined>(undefined);

interface ModalProps {
  children: ReactNode;
  defaultOpen?: boolean;
  onClose?: () => void;
}

export function Modal({ children, defaultOpen = false, onClose }: ModalProps) {
  const [isOpen, setIsOpen] = useState(defaultOpen);

  const openModal = () => setIsOpen(true);
  const closeModal = () => {
    setIsOpen(false);
    onClose?.();
  };

  return (
    <ModalContext.Provider value={{ isOpen, openModal, closeModal }}>
      {children}
    </ModalContext.Provider>
  );
}

interface ModalTriggerProps {
  children: ReactNode;
  asChild?: boolean;
}

export function ModalTrigger({ children, asChild }: ModalTriggerProps) {
  const context = useContext(ModalContext);
  if (!context) {
    throw new Error('ModalTrigger must be used within a Modal');
  }

  const { openModal } = context;

  if (asChild) {
    return (
      <div onClick={openModal} role="button" tabIndex={0}>
        {children}
      </div>
    );
  }

  return (
    <button onClick={openModal} type="button">
      {children}
    </button>
  );
}

interface ModalContentProps {
  children: ReactNode;
  className?: string;
  size?: 'sm' | 'md' | 'lg' | 'xl';
}

export function ModalContent({
  children,
  className,
  size = 'md',
}: ModalContentProps) {
  const context = useContext(ModalContext);
  if (!context) {
    throw new Error('ModalContent must be used within a Modal');
  }

  const { isOpen, closeModal } = context;

  if (!isOpen) return null;

  const sizeClasses = {
    sm: 'max-w-md',
    md: 'max-w-lg',
    lg: 'max-w-2xl',
    xl: 'max-w-4xl',
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black bg-opacity-50"
        onClick={closeModal}
        aria-hidden="true"
      />
      
      {/* Modal */}
      <div
        className={cn(
          'relative w-full mx-4 bg-white rounded-lg shadow-xl',
          sizeClasses[size],
          className
        )}
        role="dialog"
        aria-modal="true"
      >
        {children}
      </div>
    </div>
  );
}

interface ModalHeaderProps {
  children: ReactNode;
  className?: string;
}

export function ModalHeader({ children, className }: ModalHeaderProps) {
  const context = useContext(ModalContext);
  if (!context) {
    throw new Error('ModalHeader must be used within a Modal');
  }

  const { closeModal } = context;

  return (
    <div className={cn('flex items-center justify-between p-6 border-b', className)}>
      <div className="text-lg font-semibold">{children}</div>
      <button
        onClick={closeModal}
        className="text-gray-400 hover:text-gray-600 text-xl font-bold"
        aria-label="Close modal"
      >
        ×
      </button>
    </div>
  );
}

interface ModalBodyProps {
  children: ReactNode;
  className?: string;
}

export function ModalBody({ children, className }: ModalBodyProps) {
  return (
    <div className={cn('p-6', className)}>
      {children}
    </div>
  );
}

interface ModalFooterProps {
  children: ReactNode;
  className?: string;
}

export function ModalFooter({ children, className }: ModalFooterProps) {
  return (
    <div className={cn('flex justify-end gap-3 p-6 border-t bg-gray-50', className)}>
      {children}
    </div>
  );
}

export { useModalContext };

function useModalContext() {
  const context = useContext(ModalContext);
  if (!context) {
    throw new Error('useModalContext must be used within a Modal');
  }
  return context;
}