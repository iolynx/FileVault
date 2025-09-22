'use client';
import { useState, useEffect, useCallback } from 'react';

interface DragAndDropOptions {
	onDrop: (files: File[]) => void;
}

/**
 * Options for useDragAndDrop hook
 * 
 * @typedef {Object} DragAndDropOptions
 * @property {(files: File[]) => void} onDrop - Callback invoked when files are dropped
 */

/**
 * Custom React hook to handle drag-and-drop file uploads.
 * - Tracks whether files are currently being dragged over the window
 * - Calls the provided `onDrop` callback with dropped files
 * - Handles `dragenter`, `dragover`, `dragleave`, and `drop` events
 * - Uses preventDefault() and stopPropagation() to stop the event from propagating outwards.
 * 
 * @param {DragAndDropOptions} options - Options for the hook
 * @returns {{ isDragging: boolean }} - Object containing drag state
 */
export const useDragAndDrop = ({ onDrop }: DragAndDropOptions) => {
	const [isDragging, setIsDragging] = useState(false);

	const handleDragEnter = useCallback((e: DragEvent) => {
		e.preventDefault();
		e.stopPropagation();
		setIsDragging(true);
	}, []);

	const handleDragLeave = useCallback((e: DragEvent) => {
		e.preventDefault();
		e.stopPropagation();
		setIsDragging(false);
	}, []);

	const handleDragOver = useCallback((e: DragEvent) => {
		// fire the overlay event
		e.preventDefault();
		e.stopPropagation();
	}, []);

	const handleDrop = useCallback((e: DragEvent) => {
		e.preventDefault();
		e.stopPropagation();
		setIsDragging(false);

		if (e.dataTransfer?.files && e.dataTransfer.files.length > 0) {
			// Call the callback with the dropped files
			onDrop(Array.from(e.dataTransfer.files));
			e.dataTransfer.clearData();
		}
	}, [onDrop]);

	useEffect(() => {
		// Attach event listeners to the window
		window.addEventListener('dragenter', handleDragEnter);
		window.addEventListener('dragleave', handleDragLeave);
		window.addEventListener('dragover', handleDragOver);
		window.addEventListener('drop', handleDrop);

		// Clean up listeners on component unmount
		return () => {
			window.removeEventListener('dragenter', handleDragEnter);
			window.removeEventListener('dragleave', handleDragLeave);
			window.removeEventListener('dragover', handleDragOver);
			window.removeEventListener('drop', handleDrop);
		};
	}, [handleDragEnter, handleDragLeave, handleDragOver, handleDrop]);

	return { isDragging };
};
