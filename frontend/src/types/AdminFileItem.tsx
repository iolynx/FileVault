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
