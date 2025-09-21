export interface User {
	id: string;
	name: string;
	email: string;
	role?: string;
	original_storage_bytes?: bigint;
	dedup_storage_bytes?: bigint;
	storage_quota?: 100000000;
}
