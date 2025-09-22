'use client';

import { useState, useEffect, useCallback } from 'react';
import { Table, TableBody, TableCell, TableHeader, TableHead, TableRow } from '@/components/ui/table';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { DataTablePagination } from '@/components/DataTablePagination';
import api from '@/lib/axios';
import { toast } from 'sonner';

// A helper type for nullable numbers from the Go Backend
/**
 * AuditLog
 * - The shape of the AuditLog data type
 *   returned from the endpoint /admin/audit-logs.
 * - Represents a record of an audited action performed by a user.
 */
interface AuditLog {
	id: number;
	user_id: { Int64: number; Valid: boolean; };
	action: string;
	target_id: string;
	details: Record<string, unknown> | null; // It's now a regular object or null
	created_at: string;
}

export function AuditLogTable() {
	const [logs, setLogs] = useState<AuditLog[]>([]);
	const [isLoading, setIsLoading] = useState(true);
	const [totalCount, setTotalCount] = useState(0);
	const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 });

	const fetchLogs = useCallback(async () => {
		setIsLoading(true);
		try {
			const params = new URLSearchParams({
				page: String(pagination.pageIndex + 1),
				limit: String(pagination.pageSize),
			});
			const response = await api.get(`/admin/audit-logs?${params.toString()}`);
			setLogs(response.data);
			// For simplicity, we'll assume total count comes from a header or another call
			// In a real app, you'd get this from your API response
			setTotalCount(100); // Placeholder, update with your actual total count logic
		} catch (error) {
			toast.error('Failed to fetch audit logs.');
			console.error(error);
		} finally {
			setIsLoading(false);
		}
	}, [pagination]);

	useEffect(() => {
		fetchLogs();
	}, [fetchLogs]);

	return (
		<Card>
			<CardHeader>
				<CardTitle>Recent Activity</CardTitle>
			</CardHeader>
			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>User ID</TableHead>
							<TableHead>Action</TableHead>
							<TableHead>Target ID</TableHead>
							<TableHead>Details</TableHead>
							<TableHead>Timestamp</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						{isLoading ? (
							<TableRow><TableCell colSpan={5} className="h-24 text-center">Loading...</TableCell></TableRow>
						) : (
							logs.map((log) => (
								<TableRow key={log.id}>
									<TableCell>
										{log.user_id.Int64}
									</TableCell>
									<TableCell><span className="font-mono bg-muted px-2 py-1 rounded-md">{log.action}</span></TableCell>
									<TableCell className="font-mono text-xs">{log.target_id}</TableCell>
									<TableCell className="font-mono text-xs">
										{log.details ? JSON.stringify(log.details) : 'N/A'}
									</TableCell>
									<TableCell>{new Date(log.created_at).toLocaleString()}</TableCell>
								</TableRow>
							))
						)}
					</TableBody>
				</Table>
				<DataTablePagination
					pageIndex={pagination.pageIndex}
					pageSize={pagination.pageSize}
					totalCount={totalCount}
					setPageSize={(size) => setPagination({ pageIndex: 0, pageSize: size })}
					setPageIndex={(page) => setPagination(p => ({ ...p, pageIndex: page }))}
				/>
			</CardContent>
		</Card>
	);
}
