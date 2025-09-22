import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { File } from "@/types/File"
import { cn, formatBytes, mapUsersToOptions } from "@/lib/utils";
import { getFileIcon } from "./FileIcon";
import { ContentItem } from "@/types/Content";
import { ArrowDown, ArrowUp, FolderIcon, FolderOpen } from "lucide-react";
import FileActionsDropdown from "./FileActionsDropdown";
import { useContentStore } from "@/stores/useContentStore";
import FolderActionsDropdown from "./FolderActionsDropdown";
import { SortConfig } from "@/types/Sort";
import { SortableHeader } from "@/components/SortableHeader";


interface FilesTableProps {
	contents: ContentItem[];
	onDataChange: () => void;
	sortConfig: SortConfig;
	onSort: (key: string) => void;
}

/**
 * Props for FilesTable
 * 
 * @typedef {Object} FilesTableProps
 * @property {ContentItem[]} contents - Array of file/folder items to display in the table
 * @property {() => void} onDataChange - Callback to refresh data after changes (e.g., file upload, rename)
 * @property {SortConfig} sortConfig - Current sorting configuration for the table
 * @property {(key: string) => void} onSort - Callback triggered when the user sorts by a column
 */

/**
 * FilesTable component
 * 
 * Renders a table of files and folders.
 * This component supports;
 * - Clicking folders to navigate into them
 * - Sorting by a few columns using the provided sortConfig and onSort callback
 * - Refreshing data via the onDataChange callback
 * 
 * Uses the `useContentStore` hook to navigate into folders.
 * 
 * @param {FilesTableProps} props - Component props
 * @returns {JSX.Element} JSX element rendering the files/folders table
 * 
 * @component
 */

export default function FilesTable({ contents, onDataChange, sortConfig, onSort }: FilesTableProps) {
	const { navigateToFolder } = useContentStore();

	/**
	 * Handles a row click in the table.
	 * If the clicked item is a folder, navigates into that folder (using useContentStore).
	 * 
	 * @param {ContentItem} item - The clicked content item (file or folder)
	 */
	const handleRowClick = (item: ContentItem) => {
		if (item.item_type === 'folder') {
			navigateToFolder({ id: item.id, filename: item.filename });
		}
	};

	return (
		<div>
			<Table>
				<TableHeader>
					<TableRow>
						<SortableHeader columnKey="filename" sortConfig={sortConfig} onSort={onSort} className="px-4 w-[100%]">
							Name
						</SortableHeader>

						<TableHead>Location</TableHead>
						<SortableHeader columnKey="size" sortConfig={sortConfig} onSort={onSort}>
							Size
						</SortableHeader>
						<SortableHeader columnKey="uploaded_at" sortConfig={sortConfig} onSort={onSort}>
							Uploaded On
						</SortableHeader>
						<TableHead className="w-0.5 pl-0 ml-0"></TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{contents.length == 0 && (
						<div className="justify-center text-center my-4">
							No files here (yet)
						</div>
					)}
					{contents.map((contentItem, index) => (
						<TableRow
							key={index}
							onClick={() => handleRowClick(contentItem)}
							className={
								cn({
									"group cursor-pointer": contentItem.item_type === "folder"
								})}
						>
							<TableCell className="flex flex-row gap-x-2 px-4 pt-4">
								{contentItem.item_type === "file"
									? (getFileIcon(contentItem.content_type || ""))
									: (<FolderIcon fill="grey" strokeOpacity={0.3} />)
								}
								{contentItem.filename}
							</TableCell>
							<TableCell>
								{contentItem.item_type === "file"
									? (contentItem.user_owns_file ? "Your Vault" : "Shared with You")
									: ("-")
								}
							</TableCell>
							<TableCell>
								{contentItem.item_type === "file"
									? (formatBytes(Number(contentItem.size)))
									: ("- ")
								}
							</TableCell>
							<TableCell>{new Date(contentItem.uploaded_at).toLocaleDateString()}</TableCell>
							<TableCell>
								{contentItem.item_type === "file"
									? (
										<FileActionsDropdown
											file={contentItem}
											onFileChange={onDataChange}
										/>
									) : (
										<FolderActionsDropdown
											folder={contentItem}
											onFolderChange={onDataChange}

										/>
									)
								}
							</TableCell>
						</TableRow>
					))}
				</TableBody>
			</Table>

		</div>
	)
};
