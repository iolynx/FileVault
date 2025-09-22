export interface User {
	id: string;
	name: string;
	email: string;
	role: string;
	storage_used_bytes: number;
	deduplicated_usage_bytes: number;
	storage_quota_bytes: number;
	savings_bytes: number;
	savings_percentage: number;
}
