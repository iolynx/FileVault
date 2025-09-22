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
import { Button } from './ui/button';
import api from '@/lib/axios';
import { useRouter } from 'next/navigation';
import { useContentStore } from '@/stores/useContentStore';
import { User } from '@/types/User';
import { useAuthStore } from '@/stores/useAuthStore';
import Link from 'next/link';


export function UserAccountPopover() {
	// Return an empty div if user data hasn't loaded yet
	const user = useAuthStore((state) => state.user);
	const router = useRouter()

	if (!user) {
		return (
			<Link href="/login" passHref>
				<Button asChild>
					<a>Login</a>
				</Button>
			</Link>
		);
	}

	const usedStoragePercentage = (user.deduplicated_usage_bytes / user.storage_quota_bytes) * 100;
	const remainingBytes = Math.max(user.storage_quota_bytes - user.deduplicated_usage_bytes, 0);

	const logout = async () => {
		try {
			await api.post("/auth/logout");
		} catch (error) {
			console.log("Error trying to logout: ", error);
		} finally {
			useAuthStore.getState().setUser(null); // Clear the auth store
			useContentStore.getState().reset(); // Reset content
			router.push("/login");
		}
	};

	return (
		<Popover>
			<PopoverTrigger asChild className='m-3'>
				<Avatar className="cursor-pointer">
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
								<span className="font-semibold">{formatBytes(user.deduplicated_usage_bytes)}</span>
							</div>
							<div className="flex justify-between">
								<span className="text-muted-foreground">Original storage usage</span>
								<span className="font-semibold">{formatBytes(user.storage_used_bytes)}</span>
							</div>
							<div className="flex justify-between">
								<span className="text-muted-foreground">Storage savings</span>
								<span className="font-semibold text-green-600">
									{formatBytes(user.savings_bytes)} ({user.savings_percentage.toFixed(1)}%)
								</span>
							</div>
							<div className='flex justify-around mt-4'>
								<Button variant="outline" onClick={logout}>Logout</Button>
							</div>
						</div>
					</div>
				</div>
			</PopoverContent>
		</Popover>
	);
}
