'use client';

import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts';
import { formatBytes } from '@/lib/utils';

interface ChartData {
	name: string;
	value: number;
	[key: string]: unknown;
}

interface SavingsPieChartProps {
	data: ChartData[];
}

const COLORS = ['hsl(var(--primary))', 'hsl(var(--primary) / 0.3)'];

export default function SavingsPieChart({ data }: SavingsPieChartProps) {
	return (
		<ResponsiveContainer width="100%" height="100%">
			<PieChart>
				<Tooltip
					cursor={{ fill: 'transparent' }}
					content={({ active, payload }) => {
						if (active && payload && payload.length) {
							return (
								<div className="rounded-lg border bg-background p-2 shadow-sm">
									<p className="text-sm font-medium">{`${payload[0].name}: ${formatBytes(payload[0].value as number)}`}</p>
								</div>
							);
						}
						return null;
					}}
				/>
				<Pie
					data={data}
					cx="50%"
					cy="50%"
					innerRadius={60}
					outerRadius={80}
					paddingAngle={5}
					dataKey="value"
				>
					{data.map((entry, index) => (
						<Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
					))}
				</Pie>
			</PieChart>
		</ResponsiveContainer>
	);
}
