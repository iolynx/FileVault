/**
 * Represents a file in the admin view, including ownership and metadata details.
 * 
 * @interface AdminFileItem
 * 
 * @property {string} id - Unique identifier for the file
 * @property {string} filename - Name of the file
 * @property {number} size - Size of the file in bytes
 * @property {string} uploaded_at - ISO string representing when the file was uploaded
 * @property {number} owner_id - ID of the user who owns the file
 * @property {string} owner_email - Email of the user who owns the file
 * @property {number} download_count - Number of times the file has been downloaded
 * @property {string | null} declared_mime - Declared MIME type of the file, or null if not set
 */
export interface AdminFileItem {
	id: string;
	filename: string;
	size: number;
	uploaded_at: string;
	owner_id: number;
	owner_email: string;
	download_count: number;
	declared_mime: string | null;
}
