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
import { toast } from "sonner";
import api from "@/lib/axios";
import { DeleteDialogModal } from "./DeleteDialogModal";
import { useEffect, useState } from "react";
import { RenameDialogModal } from "./RenameDialogModal";
import { mapUsersToOptions } from "@/lib/utils";
import { User } from "@/types/User";
import { MultiSelectOption } from "./multi-select";
import { ShareDialogModal } from "@/components/ShareDialogModal";
import { ContentItem } from "@/types/Content";
import { InfoModal } from "@/components/InfoModal";
import { useContentStore } from "@/stores/useContentStore";

interface ActionsDropDownProps {
	file: ContentItem;
	onFileChange: () => void;
	shareDialogOptions: MultiSelectOption[];
}

export default function FileActionsDropdown({ file, onFileChange, shareDialogOptions }: ActionsDropDownProps) {
	const [isDeleteDialogOpen, setDeleteDialogOpen] = useState(false)
	const [isRenameDialogOpen, setRenameDialogOpen] = useState(false)
	const [isShareDialogOpen, setShareDialogOpen] = useState(false)
	const [shareDialogDefautValue, setShareDialogDefaultValue] = useState<string[]>([]);
	const [shareDialogURL, setShareDialogURL] = useState<string>("");
	const [isInfoModalOpen, setIsInfoModalOpen] = useState(false);
	const { renameItem, deleteItem } = useContentStore();

	const fetchUsersWithAccessToFile = async () => {
		try {
			const res = await api.get(`/files/${file.id}/shares`,
				{ withCredentials: true },
			)
			const usersWithAccessToFile: User[] = await res.data
			const userIds = usersWithAccessToFile.map((user) => user.id);
			setShareDialogDefaultValue(userIds);
		} catch (error) {
			console.log("error while fetching users with access to file: ", error)
		}
	}

	const fetchFileURL = async () => {
		try {
			const res = await api.get(`/files/url/${file.id}`,
				{ withCredentials: true },
			)
			setShareDialogURL(res.data.url);
		} catch (error) {
			console.log("error while fetching link to file ", file.filename)
		}
	}
	useEffect(() => {
		if (file.user_owns_file && file.item_type === 'file' && isShareDialogOpen) {
			fetchUsersWithAccessToFile();
			fetchFileURL();
		}
	}, [isShareDialogOpen])

	const handleDelete = async () => {
		try {
			const res = await api.delete(`/files/${file.id}`, { withCredentials: true });
			if (res.status === 204) {
				toast.success("Deleted file successfully");
				deleteItem(file.id);
			} else {
				toast.error(res.data.error);
			}
		} catch (error: any) {
			console.log('Error while deleting file: ', error);
			toast.error(error.response.data.message);
		} finally {
			//onFileChange();
		}
	}

	const handleDownload = async () => {
		try {
			const res = await api.get(`/files/${file.id}`, {
				withCredentials: true,
				responseType: 'blob',
			});
			const url = window.URL.createObjectURL(new Blob([res.data]));

			const link = document.createElement('a');
			link.href = url;
			link.setAttribute('download', file.filename);
			document.body.appendChild(link);
			link.click();

			link.parentNode?.removeChild(link);
			window.URL.revokeObjectURL(url);
			toast.success("File Downloaded");
		} catch (error: any) {
			toast.error(error.response.data.message);
		}
	}

	const handleRename = async (newFilename: string) => {
		try {
			if (newFilename === "") {
				toast.error("Filename cannot be empty")
				return
			}
			const res = await api.patch(`/files/${file.id}`,
				{ name: newFilename },
				{ headers: { "Content-Type": "application/json" }, withCredentials: true }
			)
			console.log(res.data);
			renameItem(file.id, res.data);
			toast.success(`Renamed file to ${res.data.filename}`);
		} catch (error: any) {
			toast.error(error.response.data.message)
		} finally {
			//onFileChange();
		}
	}

	// sequential sharing for now, TODO: creat an endpoint that can accept users in bulk
	const handleShare = async (selectedUsers: string[]) => {
		// find the users that were added
		const added = selectedUsers.filter(id => !shareDialogDefautValue.includes(id));

		// find the users that were removed
		const removed = shareDialogDefautValue.filter(id => !selectedUsers.includes(id));
		try {
			for (const addedUser of added) {
				const res = await api.post(
					`/files/${file.id}/share`,
					{ target_user_id: parseInt(addedUser) },
					{
						headers: { "Content-Type": "application/json" },
						withCredentials: true,
					}
				);
				toast.success(res.data.message);
			}

			for (const removedUser of removed) {
				const res = await api.delete(`/files/${file.id}/share/${removedUser}`, { withCredentials: true },
				);
				toast.success(res.data.message);
			}
		} catch (error: any) {
			console.error(error);
			toast.error(error.response?.data?.error || "Failed to share file");
		} finally {
			// do nothing 
			// would've onFileChange() but
			// it exceeds rate limits :(
		}
	};
	return (
		<div className="">
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<Button variant="ghost"><EllipsisVerticalIcon /></Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent className="w-36" align="start">

					<DropdownMenuGroup>
						<DropdownMenuItem onSelect={() => handleDownload()}>
							<DownloadIcon />
							Download
						</DropdownMenuItem>
						<DropdownMenuItem onSelect={() => setRenameDialogOpen(true)} disabled={file.user_owns_file ? false : true}>
							<PencilIcon />
							Rename
						</DropdownMenuItem >
						<DropdownMenuItem onSelect={() => setDeleteDialogOpen(true)} disabled={file.user_owns_file ? false : true}>
							<TrashIcon />
							Delete
						</DropdownMenuItem>
					</DropdownMenuGroup>

					<DropdownMenuSeparator />

					<DropdownMenuGroup>
						<DropdownMenuItem onSelect={() => setShareDialogOpen(true)} disabled={file.user_owns_file ? false : true}>
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
				context="File"
			/>

			<RenameDialogModal
				isOpen={isRenameDialogOpen}
				isOpenChange={setRenameDialogOpen}
				originalFilename={file.filename}
				onConfirm={(newFilename) => handleRename(newFilename)}
				context="File"
			/>

			<ShareDialogModal
				isOpen={isShareDialogOpen}
				isOpenChange={setShareDialogOpen}
				userOptions={shareDialogOptions}
				onConfirm={(usersToShare: string[]) => handleShare(usersToShare)}
				defaultValue={shareDialogDefautValue}
				fileURL={shareDialogURL}
			/>

			<InfoModal
				isOpen={isInfoModalOpen}
				onOpenChange={setIsInfoModalOpen}
				item={file}
			/>

		</div>
	)
}
