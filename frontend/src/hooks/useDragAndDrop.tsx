'use client';
import { useState, useEffect, useCallback } from 'react';

interface DragAndDropOptions {
	onDrop: (files: File[]) => void;
}

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
