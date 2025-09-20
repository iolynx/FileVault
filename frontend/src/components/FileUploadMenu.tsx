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
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from './ui/dialog';

interface FileUploadMenuProps {
	onActionComplete: () => void;
	currentFolderID: string | null;
}
export function FileUploadMenu({ onActionComplete, currentFolderID }: FileUploadMenuProps) {
	const inputRef = useRef<HTMLInputElement>(null);

	const [uploadProgress, setUploadProgress] = useState<number | null>(null);
	const [isFolderModalOpen, setFolderModalOpen] = useState(false);
	const [newFolderName, setNewFolderName] = useState('');
	const [isCreating, setIsCreating] = useState(false);
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

		if (currentFolderID) {
			formData.append('folder_id', currentFolderID);
		}

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
			onActionComplete();
			if (toastIdRef.current) {
				toast.dismiss(toastIdRef.current);
				toast.success('Files uploaded successfully!');
			}
			setUploadProgress(null);
			toastIdRef.current = null;
		}
	};

	const handleCreateFolder = async () => {
		if (!newFolderName.trim()) {
			toast.error("Folder name cannot be empty.");
			return;
		}

		setIsCreating(true);
		try {
			// Make the API call to your backend
			await api.post('/folders', {
				name: newFolderName.trim(),
				parent_folder_id: currentFolderID, // Can be null for root
			}, { withCredentials: true });

			toast.success(`Folder "${newFolderName.trim()}" created successfully!`);

			// Reset and close the modal
			setFolderModalOpen(false);
			setNewFolderName('');

			// Refresh the file list in the parent component
			onActionComplete();

		} catch (error) {
			console.log('error:', error)
			console.error('Failed to create folder:', error);
			toast.error('Failed to create folder.');
		} finally {
			setIsCreating(false);
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
						onSelect={() => setFolderModalOpen(true)}
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
			<Dialog open={isFolderModalOpen} onOpenChange={setFolderModalOpen}>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>Create New Folder</DialogTitle>
						<DialogDescription>
							Enter a name for your new folder.
						</DialogDescription>
					</DialogHeader>
					<div className="py-4">
						<Input
							id="folderName"
							placeholder="e.g., 'Projects' or 'Reports'"
							value={newFolderName}
							onChange={(e) => setNewFolderName(e.target.value)}
							onKeyDown={(e) => e.key === 'Enter' && handleCreateFolder()}
						/>
					</div>
					<DialogFooter>
						<Button variant="ghost" onClick={() => setFolderModalOpen(false)}>Cancel</Button>
						<Button onClick={handleCreateFolder} disabled={isCreating}>
							{isCreating ? 'Creating...' : 'Create'}
						</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>
		</div>
	);
}
