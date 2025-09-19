export interface File {
	id: string;
	filename: string;
	size: bigint;
	content_type: string;
	uploaded_at: string;
	user_owns_file: boolean;
	download_count?: number | null;
}
