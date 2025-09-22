'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { formatBytes } from '@/lib/utils';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts';
import dynamic from 'next/dynamic';

interface StorageSavingsCardProps {
	totalLogicalBytes: number;
	totalPhysicalBytes: number;
}

const SavingsPieChart = dynamic(() => import('./SavingsPieChart'), {
	ssr: false,
	loading: () => <div className="h-[200px] w-full flex items-center justify-center">Loading chart...</div>,
});

const COLORS = ['hsl(var(--primary))', 'hsl(var(--primary) / 0.3)']; // Chart colors

export function StorageSavingsCard({
	totalLogicalBytes,
	totalPhysicalBytes,
}: StorageSavingsCardProps) {
	// --- Calculate the Savings Metrics ---
	const savingsBytes = Math.max(0, totalLogicalBytes - totalPhysicalBytes);
	const savingsPercentage = totalLogicalBytes > 0 ? (savingsBytes / totalLogicalBytes) * 100 : 0;

	// prepare data for the pie chart
	const chartData = [
		{ name: 'Physically Stored', value: totalPhysicalBytes },
		{ name: 'Saved via Deduplication', value: savingsBytes },
	];

	return (
		<Card>
			<CardHeader>
				<CardTitle>Global Storage Efficiency</CardTitle>
				<CardDescription>
					Overview of storage usage with file deduplication across all users.
				</CardDescription>
			</CardHeader>
			<CardContent>
				<div className="grid grid-cols-1 md:grid-cols-2 gap-8 items-center">

					<div className="h-[200px] w-full">
						<SavingsPieChart data={chartData} />
					</div>

					<div className="space-y-4">
						<div>
							<p className="text-sm text-muted-foreground">Total Savings</p>
							<p className="text-2xl font-bold text-green-500">
								{formatBytes(savingsBytes)}
							</p>
						</div>
						<div>
							<p className="text-sm text-muted-foreground">Efficiency</p>
							<p className="text-2xl font-bold">
								{savingsPercentage.toFixed(1)}%
							</p>
						</div>
						<div className="text-xs text-muted-foreground pt-2">
							<p>Original (Logical): {formatBytes(totalLogicalBytes)}</p>
							<p>Deduplicated (Physical): {formatBytes(totalPhysicalBytes)}</p>
						</div>
					</div>

				</div>		</CardContent>
		</Card>
	);
}
