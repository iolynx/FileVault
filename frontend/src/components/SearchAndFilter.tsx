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

const locationOptions = [
	{ value: '0', label: 'All Locations' },
	{ value: '1', label: 'Your Vault' },
	{ value: '2', label: 'Shared with You' },
]
const mimeTypeOptions = [
	{ value: 'image/jpeg', label: 'JPEG Image' },
	{ value: 'application/pdf', label: 'PDF Document' },
	{ value: 'text/plain', label: 'Text File' },
];

interface SearchAndFilterProps {
	activeFilters: ActiveFilter[];
	onFilterChange: (column: string, value: Date | string | undefined) => void;
}

export function SearchAndFilterComponent({ activeFilters, onFilterChange }: SearchAndFilterProps) {
	// Helper function to get the current value for any filter
	const getActiveValue = (column: string) =>
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
					value={getActiveValue('filename')}
					onChange={(e) => onFilterChange('filename', e.target.value)}
				/>
			</div>

			<div className="flex flex-row gap-x-2 justify-center">
				<FilterSelect
					placeholder="Location"
					options={locationOptions}
					value={getActiveValue('user_owns_file')}
					onChange={(newValue) => onFilterChange('user_owns_file', newValue)}
				/>
				<FilterSelect
					placeholder="File Type"
					options={mimeTypeOptions}
					value={getActiveValue('content_type')}
					onChange={(newValue) => onFilterChange('content_type', newValue)}
				/>

				<DropdownMenu>
					<DropdownMenuTrigger asChild>
						<Button variant="outline" className="w-[180px] justify-start font-normal">
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
						<Button variant="outline" className="w-[180px] justify-start font-normal">
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
