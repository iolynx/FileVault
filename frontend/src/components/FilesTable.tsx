import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { File } from "@/types/File"
import { formatBytes } from "@/lib/utils";
import { getFileIcon } from "./FileIcon";
import ActionsDropdown from "./ActionsDropdown";

interface FilesTableProps {
	files: File[];
	onFileChange: () => void;
}

export default function FilesTable({ files, onFileChange }: FilesTableProps) {
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
					{files.length == 0 && (
						<div className="justify-center text-center my-4">
							No files here (yet)
						</div>
					)}
					{files.map((file, index) => (
						<TableRow key={index}>
							<TableCell className="flex flex-row gap-x-2 px-4 pt-4">{getFileIcon(file.content_type)} {file.filename}</TableCell>
							<TableCell>{file.user_owns_file ? "Your Vault" : "Shared with You"}</TableCell>
							<TableCell>{formatBytes(Number(file.size))}</TableCell>
							<TableCell>{new Date(file.uploaded_at).toLocaleDateString()}</TableCell>
							<TableCell>
								<ActionsDropdown
									file={file}
									onFileChange={onFileChange}
								/>
							</TableCell>
						</TableRow>
					))}
				</TableBody>
			</Table>

		</div>
	)
};
