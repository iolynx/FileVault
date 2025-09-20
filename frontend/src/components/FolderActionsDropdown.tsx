import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { ContentItem } from "@/types/Content";
import { DownloadIcon, EllipsisVerticalIcon, FolderIcon, InfoIcon, PencilIcon, TrashIcon, UserRoundPlusIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { RenameDialogModal } from "./RenameDialogModal";
import { DeleteDialogModal } from "./DeleteDialogModal";
import { StopPropagationWrapper } from "./StopPropagationWrapper";
import api from "@/lib/axios";
import { toast } from "sonner";
import { useContentStore } from "@/stores/useContentStore";
import { InfoModal } from "./InfoModal";

interface ActionsDropDownProps {
	folder: ContentItem;
	onFolderChange: () => void;
}

export default function FolderActionsDropdown({ folder, onFolderChange }: ActionsDropDownProps) {
	const [isDeleteDialogOpen, setDeleteDialogOpen] = useState(false)
	const [isRenameDialogOpen, setRenameDialogOpen] = useState(false)
	const [isInfoModalOpen, setIsInfoModalOpen] = useState(false);

	const handleDelete = async () => {
		try {
			const res = await api.delete(`/folders/${folder.id}`, { withCredentials: true });
			if (res.status === 204) {
				toast.success("Deleted folder successfully");
				onFolderChange();
			} else {
				toast.error(res.data.error);
			}
		} catch (error: any) {
			console.log('Error while deleting file: ', error);
			toast.error(error.response.data.message);
		} finally {
			onFolderChange();
		}
	}

	const handleRename = async (newFolderName: string) => {
		try {
			if (newFolderName === "") {
				toast.error("Folder name cannot be empty")
				return
			}
			const res = await api.patch(`/folders/${folder.id}`,
				{ name: newFolderName },
				{ headers: { "Content-Type": "application/json" }, withCredentials: true }
			)
			toast.success(`Renamed folder to ${res.data.name}`);
		} catch (error: any) {
			toast.error(error.response.data.error)
		} finally {
			onFolderChange();
		}
	}

	return (
		<div className="">
			{/* Wrapper to stop event propagation to the parent TableRow (which has a listener to select the folder)*/}
			<StopPropagationWrapper>
				<DropdownMenu>
					<DropdownMenuTrigger asChild>
						<Button variant="ghost" onClick={(e) => e.stopPropagation()}><EllipsisVerticalIcon /></Button>
					</DropdownMenuTrigger>
					<DropdownMenuContent className="w-36" align="start">
						<DropdownMenuGroup>
							<DropdownMenuItem onSelect={() => setRenameDialogOpen(true)}>
								<PencilIcon />
								Rename
							</DropdownMenuItem >
							<DropdownMenuItem onSelect={() => setDeleteDialogOpen(true)}>
								<TrashIcon />
								Delete
							</DropdownMenuItem>
						</DropdownMenuGroup>

						<DropdownMenuSeparator />

						<DropdownMenuGroup>
							<DropdownMenuItem>
								<FolderIcon />
								Move
							</DropdownMenuItem>
						</DropdownMenuGroup>

						<DropdownMenuSeparator />

						<DropdownMenuGroup>
							<DropdownMenuItem onSelect={() => setIsInfoModalOpen(true)}>
								<InfoIcon />
								Info
							</DropdownMenuItem>
						</DropdownMenuGroup>
					</DropdownMenuContent>
				</DropdownMenu>

				<DeleteDialogModal
					isOpen={isDeleteDialogOpen}
					isOpenChange={setDeleteDialogOpen}
					onConfirm={() => handleDelete()}
					context="Folder"
				/>

				<RenameDialogModal
					isOpen={isRenameDialogOpen}
					context="Folder"
					isOpenChange={setRenameDialogOpen}
					originalFilename={folder.filename}
					onConfirm={(newFolderName) => handleRename(newFolderName)}
				/>

				<InfoModal
					isOpen={isInfoModalOpen}
					onOpenChange={setIsInfoModalOpen}
					item={folder}
				/>

			</StopPropagationWrapper>
		</div>
	)
}
