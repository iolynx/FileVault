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
import { Progress } from '@/components/ui/progress';
import { Input } from './ui/input';

interface FileUploadMenuProps {
	fetchFiles: () => void;
}
export function FileUploadMenu({ fetchFiles }: FileUploadMenuProps) {
	const inputRef = useRef<HTMLInputElement>(null);

	const [uploadProgress, setUploadProgress] = useState<number | null>(null);
	const toastIdRef = useRef<string | number>(null);

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
			toastIdRef.current = toast.custom((t) => (
				<div className="bg-white border-2 rounded-lg shadow-lg p-4 w-64">
					<span>Uploading...</span>
					<Progress value={0} className="w-full mt-2" />
				</div>
			));
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
						if (toastIdRef.current) {
							toast.custom((t) => (
								<div className="bg-black border-2 rounded-lg p-4 w-64">
									<span>Uploading... {percentCompleted}%</span>
									<Progress value={percentCompleted} className="w-full mt-2" />
								</div>
							), { id: toastIdRef.current });
						}
					}
				},
			});

			console.log('Upload successful!', response.data);
		} catch (error) {
			console.error('Upload failed:', error);
			toast.error("Upload Failed!");
		} finally {
			fetchFiles();
			if (toastIdRef.current) {
				toast.dismiss(toastIdRef.current);
				toast.success('Files uploaded successfully!');
			}
			setUploadProgress(null);
			toastIdRef.current = null;
			//setTimeout(() => setUploadProgress(null), 2000);
		}
	};

	return (
		<div className="justify-items-center justify-center">
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

			<Input
				type="file"
				multiple
				ref={inputRef}
				onChange={handleFileChange}
				style={{ display: 'none' }}
			/>
		</div>
	);
}
