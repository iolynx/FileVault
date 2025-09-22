import { create } from 'zustand';
import { getContent } from '@/lib/contentService';
import { toast } from 'sonner';

import { ContentItem } from '@/types/Content';
import { retry } from '@/lib/retry';

interface PathSegment {
	id: string | null; // null for root (home) path
	name: string;
}

interface ContentState {
	path: PathSegment[];
	contents: ContentItem[];
	totalCount: number;
	isLoading: boolean;
	setLoading: (loading: boolean) => void;
	fetchContents: (folderId: string | null, filters: any, pagination: { pageIndex: number, pageSize: number }) => Promise<void>;
	navigateToFolder: (folder: { id: string; filename: string }) => void;
	navigateToPathIndex: (index: number) => void;
	renameItem: (itemId: string, updatedItem: ContentItem) => void;
	deleteItem: (itemId: string) => void;
	reset: () => void;
}

const initialState = {
	path: [{ id: null, name: 'Home' }],
	contents: [],
	isLoading: false,
};

export const useContentStore = create<ContentState>((set, get) => ({
	path: [{ id: null, name: 'Home' }],
	contents: [],
	isLoading: false,
	totalCount: 0,

	// --- ACTIONS ---

	// Fetches content for the current folder
	// and assigns totalCount 
	fetchContents: async (folderId, filters, pagination) => {
		set({ isLoading: true });
		try {
			const fetcher = () => getContent(folderId, filters, pagination)

			const response = await retry(fetcher, {
				retries: 3,
				initialDelay: 1500,
				shouldRetry: (error: any) => error.response?.status === 429,
			})

			set({
				contents: response.data,
				totalCount: response.totalCount,
				isLoading: false
			});
		} catch (error: any) {
			toast.error("Failed to fetch files and folders");
			console.error("Failed to fetch contents after retries", error);
			set({ isLoading: false });
		}
	},

	// Takes the renamed file/folder and replaces the old one with it
	renameItem: async (itemId: string, updatedItem: ContentItem) => {
		set((state) => ({
			contents: state.contents.map((item) =>
				item.id === itemId ? updatedItem : item
			),
		}));
	},

	// Deletes a file/folder from the store
	deleteItem: async (itemId: string) => {
		set((state) => ({
			contents: state.contents.filter((item) => item.id !== itemId),
		}));
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

	// sets loading to true
	setLoading: (loading: boolean) => {
		set((state) => ({ isLoading: loading }));
	},

	// Resets to initial state
	reset: () => {
		set(initialState);
	}
}));
