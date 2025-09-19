'use client';

import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Progress } from '@/components/ui/progress';
import { Separator } from '@/components/ui/separator';
import { formatBytes } from '@/lib/utils';
import { User } from '@/types/User'
import { Button } from './ui/button';
import api from '@/lib/axios';
import { useRouter } from 'next/navigation';

interface UserAccountPopoverProps {
	user: User | null;
}

export function UserAccountPopover({ user }: UserAccountPopoverProps) {
	// Return an empty div if user data hasn't loaded yet
	if (user === null) {
		return (
			<div>
			</div>
		)
	}

	const router = useRouter()
	const dedupBytes = Number(user.dedup_storage_bytes ?? 0);
	const originalBytes = Number(user.original_storage_bytes ?? 0);

	const totalAllowedStorage = 10 * 1024 * 1024; // 10 MB in bytes

	const usedStoragePercentage = (dedupBytes / totalAllowedStorage) * 100;
	const remainingBytes = Math.max(totalAllowedStorage - dedupBytes, 0);
	const savedBytes = Math.max(originalBytes - dedupBytes, 0);
	const savingsPercentage = originalBytes > 0 ? (savedBytes / originalBytes) * 100 : 0;

	const Logout = async () => {
		try {
			const res = await api.post("/auth/logout")
			router.push("/login")
		} catch (error) {
			console.log("Error trying to logout: ", error)
		}
	}

	return (
		<Popover>
			<PopoverTrigger asChild className='m-3'>
				<Avatar className="cursor-pointer ">
					<AvatarFallback>
						{user.name.charAt(0).toUpperCase()}
					</AvatarFallback>
				</Avatar>
			</PopoverTrigger>

			<PopoverContent className="w-80" align="end">
				<div className="flex flex-col items-center gap-4">
					<p className="text-sm text-muted-foreground">{user.email}</p>
					<Avatar className="h-20 w-20">
						<AvatarFallback className="text-3xl">
							{user.name.charAt(0).toUpperCase()}
						</AvatarFallback>
					</Avatar>
					<p className="text-lg font-medium">Hi, {user.name}!</p>
					<Separator />

					<div className="w-full space-y-4">
						<div className="space-y-2">
							<div className="flex justify-between text-sm">
								<span className="font-medium">Storage Usage</span>
								<span className="text-muted-foreground">
									{formatBytes(remainingBytes)} remaining
								</span>
							</div>
							<Progress value={usedStoragePercentage} className="h-2" />
						</div>
						<div className="space-y-2 text-sm">
							<div className="flex justify-between">
								<span className="text-muted-foreground">Total used (deduplicated)</span>
								<span className="font-semibold">{formatBytes(dedupBytes)}</span>
							</div>
							<div className="flex justify-between">
								<span className="text-muted-foreground">Original storage usage</span>
								<span className="font-semibold">{formatBytes(originalBytes)}</span>
							</div>
							<div className="flex justify-between">
								<span className="text-muted-foreground">Storage savings</span>
								<span className="font-semibold text-green-600">
									{formatBytes(savedBytes)} ({savingsPercentage.toFixed(1)}%)
								</span>
							</div>
							<div className='flex justify-around mt-4'>
								<Button variant="outline" onClick={() => Logout()}>Logout</Button>
							</div>
						</div>
					</div>
				</div>
			</PopoverContent>
		</Popover >
	);
}
