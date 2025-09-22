'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Bar, BarChart, XAxis, YAxis } from 'recharts'; // Note: still uses recharts under the hood
import {
	ChartContainer,
	ChartTooltip,
	ChartTooltipContent,
	ChartConfig,
} from '@/components/ui/chart';
import api from '@/lib/axios';
import { toast } from 'sonner';

// The data shape from your API endpoint remains the same
interface ActivityStat {
	activity_day: string;
	event_count: number;
}

// Configuration for the chart's appearance and tooltips
const chartConfig = {
	events: {
		label: "Events",
		color: "hsl(var(--primary))",
	},
} satisfies ChartConfig;

export function ActivityChart() {
	const [data, setData] = useState<ActivityStat[]>([]);
	const [isLoading, setIsLoading] = useState(true);

	useEffect(() => {
		const fetchStats = async () => {
			setIsLoading(true);
			try {
				const response = await api.get('/admin/audit-logs/stats/activity-by-day');
				const formattedData = response.data.map((stat) => ({
					...stat,
					activity_day: new Date(stat.activity_day).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
				}));
				setData(formattedData);
			} catch (error) {
				toast.error('Failed to load activity stats.');
				console.error(error);
			} finally {
				setIsLoading(false);
			}
		};

		fetchStats();
	}, []);

	if (isLoading) {
		return (
			<Card>
				<CardHeader>
					<CardTitle>Activity Overview</CardTitle>
					<CardDescription>Loading chart data...</CardDescription>
				</CardHeader>
				<CardContent>
					<div className="h-[350px] w-full flex items-center justify-center bg-secondary rounded-md">
						<p className="text-muted-foreground">Loading...</p>
					</div>
				</CardContent>
			</Card>
		);
	}

	return (
		<Card>
			<CardHeader>
				<CardTitle>Activity Overview (Last 30 Days)</CardTitle>
				<CardDescription>Total events recorded per day.</CardDescription>
			</CardHeader>
			<CardContent>
				<ChartContainer config={chartConfig} className="h-[350px] w-full">
					<BarChart accessibilityLayer data={data}>
						<XAxis
							dataKey="activity_day"
							tickLine={false}
							tickMargin={10}
							axisLine={false}
							tickFormatter={(value) => value.slice(0, 6)}
						/>
						<YAxis
							tickLine={false}
							axisLine={false}
							tickMargin={10}
							allowDecimals={false}
						/>
						<ChartTooltip
							cursor={false}
							content={<ChartTooltipContent indicator="dot" />}
						/>
						<Bar dataKey="event_count" fill="var(--color-events)" radius={4} />
					</BarChart>
				</ChartContainer>
			</CardContent>
		</Card>
	);
}
