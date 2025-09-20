'use client';
import { UploadCloud } from 'lucide-react';

interface DropzoneOverlayProps {
	isDragging: boolean;
}

export function DropzoneOverlay({ isDragging }: DropzoneOverlayProps) {
	if (!isDragging) {
		return null;
	}

	return (
		<div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
			<div className="flex flex-col items-center justify-center p-8 border-2 border-dashed border-white rounded-lg text-white">
				<UploadCloud className="w-16 h-16 mb-4" />
				<p className="text-xl font-semibold">Drop files here to upload</p>
			</div>
		</div>
	);
}
