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

/**
 * Props for FileUploadMenu
 * 
 * @typedef {Object} FileUploadMenuProps
 * @property {() => void} onActionComplete - Callback to trigger when an upload or folder creation completes
 * @property {string | null} currentFolderID - The ID of the current folder where files/folders will be uploaded
 */

/**
 * FileUploadMenu component
 * 
 * Provides UI and functionality for uploading files and creating new folders.
 * 
 * Features:
 * - Upload files to the current folder
 * - Create a new folder in the current folder
 * - Handles upload completion via callback
 * - Displays toast notifications for success/error
 * 
 * @param {FileUploadMenuProps} props - Component props
 * @returns {JSX.Element} JSX element rendering the file upload menu
 * 
 * @component
 */
export function FileUploadMenu({ onActionComplete, currentFolderID }: FileUploadMenuProps) {
	/** Reference to the hidden file input element */
	const inputRef = useRef<HTMLInputElement>(null);

	const [isFolderModalOpen, setFolderModalOpen] = useState(false);
	const [newFolderName, setNewFolderName] = useState('');
	const [isCreating, setIsCreating] = useState(false);


	/** 
	 * Uses uploadFiles from the useFileUploader hook to upload files
	 */
	const { uploadFiles } = useFileUploader({
		onUploadComplete: onActionComplete,
	})

	/**
	 * Handles the file input change event and uploads selected files.
	 * @param {React.ChangeEvent<HTMLInputElement>} event - Change event from the file input
	 */
	const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		const files = event.target.files;
		if (files && files.length > 0) {
			uploadFiles(Array.from(files), currentFolderID);
		}
	};

	/**
	 * Creates a new folder with the provided name in the current folder.
	 * Shows success/error toast notifications.
	 */
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
