'use client';

import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	DialogFooter,
} from '@/components/ui/dialog';
import {
	Table,
	TableBody,
	TableRow,
	TableCell,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { format } from 'date-fns';
import { formatBytes } from '@/lib/utils';

import { ContentItem } from '@/types/Content';

interface InfoModalProps {
	isOpen: boolean;
	onOpenChange: (open: boolean) => void;
	item: ContentItem | null;
}

export function InfoModal({ isOpen, onOpenChange, item }: InfoModalProps) {
	if (!item) {
		return null;
	}

	if (!item.item_type) {
		return null;
	}

	// Capitalize the item type for the title (File or Folder)
	const capitalizedType = item.item_type.charAt(0).toUpperCase() + item.item_type.slice(1);

	return (
		<Dialog open={isOpen} onOpenChange={onOpenChange}>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>{capitalizedType} Details</DialogTitle>
				</DialogHeader>

				{/* The borderless table for details */}
				<Table>
					<TableBody>
						<TableRow>
							<TableCell className="font-medium text-muted-foreground">Name</TableCell>
							<TableCell>{item.filename}</TableCell>
						</TableRow>
						<TableRow>
							<TableCell className="font-medium text-muted-foreground">Type</TableCell>
							<TableCell>{capitalizedType}</TableCell>
						</TableRow>

						{/* --- File-Specific Details --- */}
						{item.item_type === 'file' && item.size != null && (
							<TableRow>
								<TableCell className="font-medium text-muted-foreground">Size</TableCell>
								<TableCell>{formatBytes(Number(item.size))}</TableCell>
							</TableRow>
						)}

						{item.item_type === 'file' && item.content_type && (
							<TableRow>
								<TableCell className="font-medium text-muted-foreground">MIME Type</TableCell>
								<TableCell>{item.content_type}</TableCell>
							</TableRow>
						)}

						{/* --- Common Details --- */}
						<TableRow>
							<TableCell className="font-medium text-muted-foreground">Ownership</TableCell>
							<TableCell>{item.user_owns_file ? 'Owned by you' : 'Shared with you'}</TableCell>
						</TableRow>

						<TableRow>
							<TableCell className="font-medium text-muted-foreground">Created</TableCell>
							<TableCell>{format(new Date(item.uploaded_at), 'PPP p')}</TableCell>
						</TableRow>

						{/* --- Owner-Specific Details --- */}
						{item.user_owns_file && item.download_count != null && (
							<TableRow>
								<TableCell className="font-medium text-muted-foreground">Downloads</TableCell>
								<TableCell>{item.download_count}</TableCell>
							</TableRow>
						)}

					</TableBody>
				</Table>

				<DialogFooter>
					<Button variant="outline" onClick={() => onOpenChange(false)}>
						Close
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
