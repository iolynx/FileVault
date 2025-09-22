'use client';

import { useRef, useState } from 'react';
import api from '@/lib/axios';
import { toast } from 'sonner';
import { Progress } from '@/components/ui/progress';
import { useAuthStore } from '@/stores/useAuthStore';

interface UseFileUploaderProps {
	onUploadComplete: () => void;
}

/**
 * Props for useFileUploader hook
 * 
 * @typedef {Object} UseFileUploaderProps
 * @property {() => void} onUploadComplete - Callback invoked after a successful file upload
 */

/**
 * Custom React hook to handle file uploads.
 * - Manages uploading state (`isUploading`)
 * - Uses `onUploadComplete` callback when uploads finish
 * - Optionally integrates with auth/user store (`fetchUser`)
 * - Can display toast notifications for upload progress and errors
 * 
 * @param {UseFileUploaderProps} props - Options for the hook
 * @returns {{ isUploading: boolean, uploadFiles: (files: File[], folderId?: string | null) => Promise<void> }} 
 *   - `isUploading`: Indicates whether a file upload is currently in progress
 *   - `uploadFiles`: Function to upload an array of files to an optional folder ID
 */
export const useFileUploader = ({ onUploadComplete }: UseFileUploaderProps) => {
	const [isUploading, setIsUploading] = useState(false);
	const toastIdRef = useRef<string | number | null>(null);
	const { fetchUser } = useAuthStore();

	const uploadFiles = async (files: File[], currentFolderId: string | null) => {
		if (isUploading) return; // Prevent multiple simultaneous uploads

		setIsUploading(true);

		const formData = new FormData();
		files.forEach((file) => formData.append('files', file));
		if (currentFolderId) {
			formData.append('folder_id', currentFolderId);
		}

		try {
			toastIdRef.current = toast.custom(() => (
				<div className="bg-white border-2 rounded-lg shadow-lg p-4 w-64">
					<span>Uploading...</span>
					<Progress value={0} className="w-full mt-2" />
				</div>
			));

			await api.post('/files/upload', formData, {
				headers: { 'Content-Type': 'multipart/form-data' },
				withCredentials: true,
				onUploadProgress: (progressEvent) => {
					if (progressEvent.total) {
						const percentCompleted = Math.round(
							(progressEvent.loaded * 100) / progressEvent.total
						);
						if (toastIdRef.current) {
							toast.custom(() => (
								<div className="bg-black border-2 rounded-lg p-4 w-64">
									<span>Uploading... {percentCompleted}%</span>
									<Progress value={percentCompleted} className="w-full mt-2" />
								</div>
							), { id: toastIdRef.current });
						}
					}
				},
			});

			if (toastIdRef.current) {
				toast.dismiss(toastIdRef.current);
				toast.success('Files uploaded successfully!');
				setTimeout(fetchUser, 1200);
			}
		} catch (error: any) {
			console.log('Upload failed:', error);
			toast.error(`Upload Failed: ${error.response.data.error}`);
		} finally {
			setIsUploading(false);
			toastIdRef.current = null;
			onUploadComplete(); // Refresh the content list
		}
	};

	return { uploadFiles, isUploading };
};
