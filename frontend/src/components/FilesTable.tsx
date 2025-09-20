import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { File } from "@/types/File"
import { formatBytes } from "@/lib/utils";
import { getFileIcon } from "./FileIcon";
import { ContentItem } from "@/types/Content";
import { FolderIcon } from "lucide-react";
import FileActionsDropdown from "./FileActionsDropdown";
import { useContentStore } from "@/stores/useContentStore";
import { useEffect } from "react";
import FolderActionsDropdown from "./FolderActionsDropdown";

interface FilesTableProps {
	contents: ContentItem[];
	onDataChange: () => void;
}

export default function FilesTable({ contents, onDataChange }: FilesTableProps) {
	const { navigateToFolder } = useContentStore();

	const handleRowClick = (item: ContentItem) => {
		if (item.item_type === 'folder') {
			navigateToFolder({ id: item.id, filename: item.filename });
		} else {
			// do nothing
		}
	};

	return (
		<div>
			<Table>
				<TableHeader>
					<TableRow>
						<TableHead className="px-4 w-[100%]">Name</TableHead>
						<TableHead>Location</TableHead>
						<TableHead>Size</TableHead>
						<TableHead>Uploaded On</TableHead>
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
							className="hover:bg-gray-900 group"
						>
							<TableCell className="flex flex-row gap-x-2 px-4 pt-4">
								{contentItem.item_type === "file"
									? (getFileIcon(contentItem.content_type))
									: (<FolderIcon />)
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
