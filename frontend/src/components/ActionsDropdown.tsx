import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { DownloadIcon, EllipsisVerticalIcon, FolderIcon, InfoIcon, PencilIcon, TrashIcon, UserRoundPlusIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { File } from "@/types/File";
import { toast } from "sonner";
import api from "@/lib/axios";
import { DeleteDialogModal } from "./DeleteDialogModal";
import { useState } from "react";
import { RenameDialogModal } from "./RenameDialogModal";

interface ActionsDropDownProps {
	file: File;
	onFileChange: () => void;
}

export default function ActionsDropdown({ file, onFileChange }: ActionsDropDownProps) {
	const [isDeleteDialogOpen, setDeleteDialogOpen] = useState(false)
	const [isRenameDialogOpen, setRenameDialogOpen] = useState(false)

	const handleDelete = async (fileID: string | null) => {
		try {
			if (fileID === null) {
				throw new Error("null File ID");
			}
			const res = await api.delete(`/files/${fileID}`, { withCredentials: true });
			if (res.status === 200) {
				toast.success(res.data.message);
				onFileChange();
			} else {
				toast.error(res.data.error);
			}
		} catch (error) {
			console.log('Error while deleting file: ', error);
			toast.error('Failed to delete file.');
		} finally {
			onFileChange();
		}
	}

	// TODO: REMOVE THE USELESS PARAMS
	const handleDownload = async (fileID: string, filename: string) => {
		try {
			if (fileID === "") {
				throw new Error("FileID is empty");
			}

			const res = await api.get(`/files/${fileID}`, {
				withCredentials: true,
				responseType: 'blob',
			});
			const url = window.URL.createObjectURL(new Blob([res.data]));

			const link = document.createElement('a');
			link.href = url;
			link.setAttribute('download', filename);
			document.body.appendChild(link);
			link.click();

			link.parentNode?.removeChild(link);
			window.URL.revokeObjectURL(url);
			toast.success(res.data.message);
		} catch (error) {
			toast.error("Error while downloading");
		}
	}

	const handleRename = async (newFilename: string) => {
		try {
			if (newFilename == "") {
				throw new Error("New Filename cannot be Empty")
			}
			const res = await api.patch(`/files/${file.id}`,
				{ filename: newFilename },
				{ headers: { "Content-Type": "application/json" }, withCredentials: true }
			)
			if (res.status == 200) {
				console.log(res.data);
				toast.success(res.data.message);
			} else {
				toast.error(res.data.error);
			}
		} catch (error: any) {
			toast.error(error);
		} finally {
			onFileChange();
		}
	}

	return (
		<div className="">
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<Button variant="ghost"><EllipsisVerticalIcon /></Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent className="w-36" align="start">

					<DropdownMenuGroup>
						<DropdownMenuItem onSelect={() => handleDownload(file.id, file.filename)}>
							<DownloadIcon />
							Download
						</DropdownMenuItem>
						<DropdownMenuItem onSelect={() => setRenameDialogOpen(true)}>
							<PencilIcon />
							Rename
						</DropdownMenuItem>
						<DropdownMenuItem onSelect={() => setDeleteDialogOpen(true)}>
							<TrashIcon />
							Delete
						</DropdownMenuItem>
					</DropdownMenuGroup>

					<DropdownMenuSeparator />

					<DropdownMenuGroup>
						<DropdownMenuItem>
							<UserRoundPlusIcon />
							Share
						</DropdownMenuItem>
						<DropdownMenuItem>
							<FolderIcon />
							Move
						</DropdownMenuItem>
					</DropdownMenuGroup>

					<DropdownMenuSeparator />

					<DropdownMenuGroup>
						<DropdownMenuItem>
							<InfoIcon />
							Info
						</DropdownMenuItem>
					</DropdownMenuGroup>
				</DropdownMenuContent>
			</DropdownMenu>

			<DeleteDialogModal
				isOpen={isDeleteDialogOpen}
				isOpenChange={setDeleteDialogOpen}
				onConfirm={() => handleDelete(file.id)}
			/>

			<RenameDialogModal
				isOpen={isRenameDialogOpen}
				isOpenChange={setRenameDialogOpen}
				originalFilename={file.filename}
				onConfirm={(newFilename) => handleRename(newFilename)}
			/>
		</div>
	)
}
