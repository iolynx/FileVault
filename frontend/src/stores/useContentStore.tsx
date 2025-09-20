import { create } from 'zustand';
import { getContent } from '@/lib/contentService';
import { toast } from 'sonner';

import { ContentItem } from '@/types/Content';
import api from '@/lib/axios';

interface PathSegment {
	id: string | null; // null for root (home) path
	name: string;
}

interface ContentState {
	path: PathSegment[];
	contents: ContentItem[];
	isLoading: boolean;
	fetchContents: (folderId: string | null, filters: any) => Promise<void>;
	navigateToFolder: (folder: { id: string; filename: string }) => void;
	navigateToPathIndex: (index: number) => void;
	renameItem: (itemId: string, newName: string, context: "File" | "Folder") => void;
	deleteItem: (itemId: string, itemType: 'file' | 'folder') => void;
}

export const useContentStore = create<ContentState>((set, get) => ({
	path: [{ id: null, name: 'Home' }],
	contents: [],
	isLoading: false,

	// --- ACTIONS ---

	// Fetches content for the current folder
	fetchContents: async (folderId, filters) => {
		set({ isLoading: true });
		try {
			const newContents = await getContent(folderId, filters);
			set({ contents: newContents, isLoading: false });
		} catch (error) {
			toast.error("Failed to fetch files and folders");
			console.error("Failed to fetch contents", error);
			set({ isLoading: false });
		}
	},

	// Takes the renamed file/folder and replaces the old one with it
	renameItem: async (itemId: string, newName: string, context: "File" | "Folder") => {
		try {
			let updatedItemResponse;
			if (context === "File") {
				updatedItemResponse = await api.patch(`/files/${itemId}`, { name: newName });
			} else {
				updatedItemResponse = await api.patch(`/folders/${itemId}`, { name: newName });
			}
			const updatedItem = updatedItemResponse.data;
			console.log('updated item: ', updatedItem);
			console.error("hiihihi");

			set((state) => ({
				contents: state.contents.map((item) =>
					item.id === itemId ? updatedItem : item
				),
			}));
			toast.success(`${context} renamed successfully!`);
		} catch (error) {
			toast.error("Failed to rename item.");
		}
	},

	// Deletes a file/folder from the store
	deleteItem: async (itemId: string, itemType: 'file' | 'folder') => {
		try {
			if (itemType === 'file') {
				await api.delete(`/files/${itemId}`);
			} else {
				await api.delete(`/folders/${itemId}`);
			}

			set((state) => ({
				contents: state.contents.filter((item) => item.id !== itemId),
			}));
			toast.success(`${itemType} deleted successfully!`);
		} catch (error) {
			toast.error(`Failed to delete ${itemType}.`);
		}
	},

	// Navigates into a subfolder
	navigateToFolder: (folder) => {
		const newPath = [...get().path, { id: folder.id, name: folder.filename }];
		set({ path: newPath });
	},

	// Navigates using the breadcrumbs
	navigateToPathIndex: (index) => {
		const newPath = get().path.slice(0, index + 1);
		set({ path: newPath });
	},
}));
