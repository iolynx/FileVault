/**
 * Represents a file or folder in the FilesTable
 * 
 * @interface ContentItem
 * 
 * @property {string} id - Unique identifier for the file or folder
 * @property {'file' | 'folder'} item_type - Indicates whether the item is a file or folder
 * @property {string} filename - Name of the file or folder
 * @property {bigint} [size] - Size of the file in bytes (undefined for folders)
 * @property {string} [content_type] - MIME type of the file (undefined for folders)
 * @property {string} uploaded_at - ISO string representing when the file/folder was uploaded
 * @property {boolean} user_owns_file - Indicates if the current user owns this item
 * @property {number | null} [download_count] - Number of times the file has been downloaded (undefined or null for folders)
 */
export interface ContentItem {
	id: string;
	item_type: 'file' | 'folder';
	filename: string;
	size?: bigint;
	content_type?: string;
	uploaded_at: string;
	user_owns_file: boolean;
	download_count?: number | null;
}
