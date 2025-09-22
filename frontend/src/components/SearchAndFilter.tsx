import { Input } from "@/components/ui/input";
import { ActiveFilter, FilterOption } from "@/types/Filter";
import { useState } from "react";
import { format } from "date-fns"
import { FilterSelect } from "./FilterSelect";
import { DropdownMenuItem, DropdownMenuPortal, DropdownMenuSeparator, DropdownMenuSub, DropdownMenuSubContent, DropdownMenuTrigger } from "@radix-ui/react-dropdown-menu";
import { DropdownMenu, DropdownMenuContent, DropdownMenuSubTrigger } from "./ui/dropdown-menu";
import { Button } from "./ui/button";
import { DatePicker } from "./DatePicker";
import { ArrowRightIcon } from "lucide-react";
import { SizeFilter } from "./SizeFilter";
import { mimeTypeOptions } from "@/types/MimeTypes";

const locationOptions = [
	{ value: '0', label: 'All Locations' },
	{ value: '1', label: 'Your Vault' },
	{ value: '2', label: 'Shared with You' },
]

const PREDEFINED_RANGES = [
	{ label: 'Tiny (< 1 MB)', minBytes: 0, maxBytes: 1048576 },
	{ label: 'Small (1-10 MB)', minBytes: 1048576, maxBytes: 10485760 },
	{ label: 'Medium (10-100 MB)', minBytes: 10485760, maxBytes: 104857600 },
	{ label: 'Large (> 100 MB)', minBytes: 104857600, maxBytes: null }, // No upper limit
];


interface SearchAndFilterProps {
	activeFilters: ActiveFilter[];
	onFilterChange: (column: string, value: Date | string | undefined) => void;
	onSizeFilterChange: (filter: { min_size: number | null, max_size: number | null }) => void;
}


/**
 * Props for SearchAndFilterComponent
 * 
 * @typedef {Object} SearchAndFilterProps
 * @property {ActiveFilter[]} activeFilters - Array of currently active filters
 * @property {(column: string, value: Date | string | undefined) => void} onFilterChange - Callback when a filter changes
 * @property {(filter: { min_size: number | null, max_size: number | null }) => void} onSizeFilterChange - Callback when size filters change
 */

/**
 * SearchAndFilterComponent
 * 
 * Renders filter inputs for searching and filtering files/folders. 
 * Supports:
 * - Column-based filters (string or date)
 * - Size-based filters (min_size and max_size)
 * 
 * @param {SearchAndFilterProps} props - Component props
 * @returns {JSX.Element} JSX element rendering the search and filter inputs
 * 
 * @component
 */
export function SearchAndFilterComponent({ activeFilters, onFilterChange, onSizeFilterChange }: SearchAndFilterProps) {
	/**
	 * Returns the currently active value for the filter of any given column.
	 * @param {string} column - Column name
	 * @returns {string} The active filter value or empty string if none
	 */
	const getActiveValue = (column: string): string =>
		activeFilters.find((f) => f.column === column)?.value || '';


	const beforeDateString = getActiveValue('uploaded_before');
	const beforeDateObject = beforeDateString ? new Date(beforeDateString) : undefined;

	const afterDateString = getActiveValue('uploaded_after');
	const afterDateObject = afterDateString ? new Date(afterDateString) : undefined;

	return (
		<div className="flex flex-col rounded-xl border shadow-sm p-4 gap-4">
			<div className="flex flex-row gap-x-2 px-10 ">
				<Input
					type="text"
					className="h-9 rounded-xl"
					id="search"
					placeholder="Search by Filename"
					value={getActiveValue('search')}
					onChange={(e) => onFilterChange('search', e.target.value)}
				/>
			</div>

			<div className="flex flex-row gap-x-2 justify-center">
				<FilterSelect
					placeholder="Location"
					options={locationOptions}
					value={getActiveValue('user_owns_file')}
					onChange={(newValue) => onFilterChange('user_owns_file', newValue)}
					type="Location"
				/>
				<FilterSelect
					placeholder="File Type"
					options={mimeTypeOptions}
					value={getActiveValue('content_type')}
					onChange={(newValue) => onFilterChange('content_type', newValue)}
					type="Filetype"
				/>

				<SizeFilter
					ranges={PREDEFINED_RANGES}
					onApplyFilter={(filter) => { onSizeFilterChange(filter) }}
				/>

				<DropdownMenu>
					<DropdownMenuTrigger asChild>
						<Button variant="outline" className="justify-start font-normal">
							{afterDateObject ? format(afterDateObject, "PPP") : "Uploaded After..."}
						</Button>
					</DropdownMenuTrigger>
					<DropdownMenuContent>
						<DatePicker
							date={afterDateObject}
							onDateChange={(newDate) => onFilterChange('uploaded_after', newDate)}
						/>
					</DropdownMenuContent>
				</DropdownMenu>
				<ArrowRightIcon size={16} className="align-bottom mt-2" />
				<DropdownMenu>
					<DropdownMenuTrigger asChild>
						<Button variant="outline" className="justify-start font-normal">
							{beforeDateObject ? format(beforeDateObject, "PPP") : "Uploaded Before..."}
						</Button>
					</DropdownMenuTrigger>
					<DropdownMenuContent>
						<DatePicker
							date={beforeDateObject}
							onDateChange={(newDate) => onFilterChange('uploaded_before', newDate)}
						/>
					</DropdownMenuContent>
				</DropdownMenu>

			</div>
		</div>
	);
}
