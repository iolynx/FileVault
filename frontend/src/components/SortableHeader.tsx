import { ArrowDown, ArrowUp } from "lucide-react";
import { TableHead } from "@/components/ui/table";
import { SortConfig } from "@/types/Sort";
import { cn } from "@/lib/utils";

type SortableHeaderProps = {
	children: React.ReactNode;
	columnKey: string;
	sortConfig: SortConfig;
	onSort: (key: string) => void;
	className?: string;
};

// Yhis is a helper component to enable sorter headers, 
// while keeping the actual table headers clean.
export const SortableHeader = ({ children, columnKey, sortConfig, onSort, className }: SortableHeaderProps) => {
	const isSorted = sortConfig?.key === columnKey;
	const Icon = sortConfig?.direction === 'asc' ? ArrowUp : ArrowDown;

	return (
		<TableHead onClick={() => onSort(columnKey)} className={cn("cursor-pointer hover:bg-muted/50" + className)}>
			<div className="flex items-center gap-2">
				{children}
				{isSorted && <Icon className="h-4 w-4" />}
			</div>
		</TableHead>
	);
};

