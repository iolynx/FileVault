'use client';

import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { SortConfig } from '@/types/Sort';
import { SortableHeader } from '@/components/SortableHeader'; // Assuming you moved SortableHeader to its own file
import { AdminFileItem } from '@/types/AdminFileItem';
import { formatBytes } from '@/lib/utils';
import { MoreHorizontal } from 'lucide-react';
import { getFileIcon } from '@/components/FileIcon';

interface AdminFilesTableProps {
	files: AdminFileItem[];
	sortConfig: SortConfig;
	onSort: (key: string) => void;
	onDataChange: () => void; // For potential future actions like manual refresh
}

export default function AdminFilesTable({ files, sortConfig, onSort, onDataChange }: AdminFilesTableProps) {
	return (
		<div>
			<Table className='w-full'>
				<TableHeader>
					<TableRow>
						<SortableHeader columnKey="filename" sortConfig={sortConfig} onSort={onSort} className='w-[100%] ml-6'>
							Name
						</SortableHeader>
						<SortableHeader columnKey="owner_email" sortConfig={sortConfig} onSort={onSort}>
							Owner
						</SortableHeader>
						<SortableHeader columnKey="size" sortConfig={sortConfig} onSort={onSort}>
							Size
						</SortableHeader>
						<SortableHeader columnKey="uploaded_at" sortConfig={sortConfig} onSort={onSort}>
							Uploaded On
						</SortableHeader>
						<SortableHeader columnKey="download_count" sortConfig={sortConfig} onSort={onSort}>
							Downloads
						</SortableHeader>
						<TableHead>Actions</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{files.length === 0 && (
						<TableRow>
							<TableCell colSpan={6} className="h-24 text-center">
								No files found.
							</TableCell>
						</TableRow>
					)}
					{files.map((file) => (
						<TableRow key={file.id} className="group">
							<TableCell className="font-medium flex items-center gap-2">
								{getFileIcon(file.declared_mime || "")}
								{file.filename}
							</TableCell>
							<TableCell>{file.owner_email}</TableCell>
							<TableCell>{formatBytes(file.size)}</TableCell>
							<TableCell>{new Date(file.uploaded_at).toLocaleDateString()}</TableCell>
							<TableCell>{file.download_count}</TableCell>
							<TableCell>
								{/* Admin actions would be different, e.g., delete, view owner, etc. */}
								{/* For now, a placeholder: */}
								<MoreHorizontal className="h-4 w-4" />
							</TableCell>
						</TableRow>
					))}
				</TableBody>
			</Table>
		</div>
	);
}
