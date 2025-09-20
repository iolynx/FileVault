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
import { useFileUploader } from '@/hooks/useFileUploader';

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

	const { uploadFiles } = useFileUploader({
		onUploadComplete: onActionComplete,
	})

	const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		const files = event.target.files;
		if (files && files.length > 0) {
			uploadFiles(Array.from(files), currentFolderID);
		}
	};

	const handleCreateFolder = async () => {
		if (!newFolderName.trim()) {
			toast.error("Folder name cannot be empty.");
			return;
		}

		setIsCreating(true);
		try {
			await api.post('/folders', {
				name: newFolderName.trim(),
				parent_folder_id: currentFolderID,
			}, { withCredentials: true });

			toast.success(`Folder "${newFolderName.trim()}" created successfully!`);

			setFolderModalOpen(false);
			setNewFolderName('');

			// Refresh the file list in the parent component
			onActionComplete();

		} catch (error) {
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
