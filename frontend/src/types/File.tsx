export interface File {
	id: string;
	filename: string;
	size: bigint;
	content_type: string;
	uploaded_at: string;
	userOwnsFile: boolean;
}
