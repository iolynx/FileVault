"use client";

import FilesTable from "@/components/FilesTable";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import api from "@/lib/axios";
import { APIError } from "@/types/APIError";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { File } from "@/types/File"
import { FileUploadMenu } from "@/components/FileUploadMenu";
import { Skeleton } from "@/components/ui/skeleton";
import { useRouter } from "next/navigation";
import { ActiveFilter, FilterOption } from "@/types/Filter";
import { SearchAndFilterComponent } from "@/components/SearchAndFilter";

const filterOptions: FilterOption[] = [
	{ value: 'content_type', label: 'MIME Type' },
	{ value: 'user_owns_file', label: 'Location' },
	{ value: 'uploaded_before', label: 'Before:' },
	{ value: 'uploaded_after', label: 'After:' }
];

const DashboardPage = () => {
	const [loading, setLoading] = useState(true);
	const [files, setFiles] = useState<File[]>([]);
	const [search, setSearch] = useState("");
	const [activeFilters, setActiveFilters] = useState<ActiveFilter[]>([]);
	const [page, setPage] = useState(0);
	const router = useRouter();

	const handleFilterChange = (column: string, value: string | Date | undefined) => {
		setActiveFilters(prevFilters => {
			const existingFilterIndex = prevFilters.findIndex(f => f.column === column);

			let stringValue: string;
			if (value instanceof Date) {
				stringValue = value.toISOString();
			} else {
				stringValue = value || '';
			}

			// If the value is empty, remove the filter from the array
			if (stringValue === '') {
				return prevFilters.filter(f => f.column !== column);
			}

			// If the filter already exists, update its value
			if (existingFilterIndex > -1) {
				const newFilters = [...prevFilters];
				newFilters[existingFilterIndex].value = stringValue;
				return newFilters;
			}

			// Otherwise, add the new filter to the array
			return [...prevFilters, { column, value: stringValue }];
		});
	};

	const fetchFiles = async (filters: ActiveFilter[]) => {
		try {
			setLoading(true);

			// Reduce filters to set of params 
			const params = filters.reduce((acc, filter) => {
				acc[filter.column] = filter.value;
				return acc;
			}, {} as Record<string, string>);

			// Add pagination params
			params.limit = '10';
			params.offset = '0';

			const res = await api.get("/files", {
				params: params,
				headers: { "Content-Type": "application/json" },
				withCredentials: true,
			});

			setFiles(res.data || []);
			router.refresh();
		} catch (error: any) {
			toast.error("Error: Failed to fetch files");
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		// Set up a timer to delay the API call
		const debounceTimer = setTimeout(() => {
			fetchFiles(activeFilters);
		}, 500); // debouce of 500ms for user to stop typing

		// Clear the timer if the user types again
		return () => clearTimeout(debounceTimer);
	}, [activeFilters]);

	return (
		<div className="flex flex-col items-center">
			<div className="flex flex-col items-center my-10">
				<h1 className="text-3xl font-bold"> Dashboard </h1>
				<p> View and Manage your files here.</p>
			</div>
			<div>
				<SearchAndFilterComponent
					activeFilters={activeFilters}
					onFilterChange={handleFilterChange}
				/>
			</div>
			<div className="m-4">
				<FileUploadMenu fetchFiles={() => fetchFiles(activeFilters)} />
			</div>

			{loading ?
				<Skeleton className="h-[500px] w-[86%] rounded-2xl max-w-7xl mt-4" /> :
				<Card className="rounded-2xl border shadow-sm overflow-hidden w-full max-w-7xl mt-4 pt-1 pb-1">
					<FilesTable files={files} onFileChange={() => fetchFiles(activeFilters)} />

				</Card>
			}
		</div>

	)
}
export default DashboardPage;
