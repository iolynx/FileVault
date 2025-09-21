'use client';

import { useAuthStore } from '@/stores/useAuthStore';
import { Progress } from '@/components/ui/progress';
import { Button } from '@/components/ui/button';
import { XIcon, HardDrive } from 'lucide-react';
import { formatBytes } from '@/lib/utils';
import { User } from '@/types/User';

interface StorageQuotaBannerProps {
	user: User | null | undefined;
}

export function StorageQuotaBanner({ user }: StorageQuotaBannerProps) {
	const { isBannerDismissed, dismissBanner } = useAuthStore();
	const storage_quota = user?.storage_quota !== undefined ? user.storage_quota : 10000000;

	if (!user || isBannerDismissed || !user.dedup_storage_bytes) {
		console.log("User info: ", user, isBannerDismissed, user?.storage_quota, user?.dedup_storage_bytes);
		return null;
	}

	const percentageUsed = (Number(user.dedup_storage_bytes) / storage_quota) * 100;

	// Only show the banner if usage is above a certain threshold (e.g., 80%)
	console.warn('percent used: ', percentageUsed);
	if (percentageUsed < 80) {
		return null;
	}

	return (
		<div className="fixed top-0 left-0 right-0 z-50 bg-secondary p-2 shadow-md flex items-center justify-center gap-4">
			<HardDrive className="h-5 w-5 text-muted-foreground" />
			<div className="flex-grow">
				<p className="text-sm font-medium">
					You are using {formatBytes(Number(user.dedup_storage_bytes))} of {formatBytes(storage_quota)}.
				</p>
				<Progress value={percentageUsed} className="h-2 mt-1" />
			</div>
			<Button variant="ghost" size="icon" onClick={dismissBanner}>
				<XIcon className="h-4 w-4" />
			</Button>
		</div>
	);
}
