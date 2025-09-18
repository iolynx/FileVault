'use client';

import { useRef, useState } from 'react';
import axios from 'axios';

import { Button } from '@/components/ui/button';
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { FileUpIcon, FolderUpIcon, PlusIcon, Upload } from 'lucide-react';
import api from '@/lib/axios';
import { toast } from 'sonner';

interface FileUploadMenuProps {
	fetchFiles: () => void;
}
export function FileUploadMenu({ fetchFiles }: FileUploadMenuProps) {
	const inputRef = useRef<HTMLInputElement>(null);

	const [uploadProgress, setUploadProgress] = useState<number | null>(null);

	const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		const files = event.target.files;
		if (files && files.length > 0) {
			handleUpload(Array.from(files));
		}
	};

	const handleUpload = async (files: File[]) => {
		const formData = new FormData();
		files.forEach((file) => {
			console.log('adding:', file);
			formData.append('files', file);
		});

		setUploadProgress(0); // Start progress tracking

		try {
			const response = await api.post('/files/upload', formData, {
				headers: {
					'Content-Type': 'multipart/form-data',
				},
				withCredentials: true,
				onUploadProgress: (progressEvent) => {
					if (progressEvent.total) {
						const percentCompleted = Math.round(
							(progressEvent.loaded * 100) / progressEvent.total
						);
						setUploadProgress(percentCompleted);
					}
				},
			});

			console.log('Upload successful!', response.data);
		} catch (error) {
			console.error('Upload failed:', error);
			toast.error("Upload Failed!");
		} finally {
			fetchFiles();
			setTimeout(() => setUploadProgress(null), 2000);
		}
	};

	return (
		<div>
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<Button variant="outline"><PlusIcon /> New</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent className='w-56' align="start">
					<DropdownMenuItem
						onSelect={() => inputRef.current?.click()}
					>
						<FileUpIcon />
						<span>Upload Files</span>
					</DropdownMenuItem>
					<DropdownMenuItem
					>
						<FolderUpIcon />
						<span>Create Folder</span>
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>

			<input
				type="file"
				multiple
				ref={inputRef}
				onChange={handleFileChange}
				style={{ display: 'none' }}
			/>

			{uploadProgress !== null && (
				<div className="mt-4 w-full bg-gray-200 rounded-full h-2.5">
					<div
						className="bg-blue-600 h-2.5 rounded-full"
						style={{ width: `${uploadProgress}%` }}
					></div>
					<p className="text-sm text-center mt-1">{uploadProgress}%</p>
				</div>
			)}
		</div>
	);
}
