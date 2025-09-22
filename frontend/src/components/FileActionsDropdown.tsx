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
import { useAuthStore } from "@/stores/useAuthStore";
import { MoveDialogModal } from "./MoveDialogModal";
import { parse } from "node:path/posix";

interface ActionsDropDownProps {
	file: ContentItem;
	onFileChange: () => void;
}

/**
 * Props for FileActionsDropdown
 * 
 * @typedef {Object} ActionsDropDownProps
 * @property {ContentItem} file - The file or folder item the dropdown actions will operate on
 * @property {() => void} onFileChange - Callback triggered when an action modifies the file (e.g., rename, delete)
 */

/**
 * FileActionsDropdown component
 * 
 * Renders a dropdown menu with actions that can be performed on a file, Including:
 * - Download
 * - Rename
 * - Delete
 * - Move
 * - Share
 * - Info
 * 
 * @param {ActionsDropDownProps} props - Component props
 * @returns {JSX.Element} JSX element rendering the file actions dropdown
 * 
 * @component
 */
export default function FileActionsDropdown({ file, onFileChange }: ActionsDropDownProps) {
	const [isDeleteDialogOpen, setDeleteDialogOpen] = useState(false)
	const [isRenameDialogOpen, setRenameDialogOpen] = useState(false)
	const [isShareDialogOpen, setShareDialogOpen] = useState(false)
	const [shareDialogDefautValue, setShareDialogDefaultValue] = useState<string[]>([]);
	const [shareDialogURL, setShareDialogURL] = useState<string>("");
	const [isInfoModalOpen, setIsInfoModalOpen] = useState(false);
	const { renameItem, deleteItem } = useContentStore();
	const [shareDialogOptions, setShareDialogOptions] = useState<MultiSelectOption[]>([]);
	const [isMoveDialogOpen, setMoveDialogOpen] = useState(false);
	const { fetchUser } = useAuthStore();
	const { path, fetchContents } = useContentStore();


	/**
	 * Fetches share information for the file.
	 * - Retrieves the list of users the file is already shared with
	 * - Sets default values for the share dialog
	 * - Retrieves the file's share URL
	 * - Retrieves the list of all possible users the file can be shared with
	 * 
	 * @async
	 * @function
	 */
	const fetchShareInfo = async () => {
		try {
			const res = await api.get(`/files/${file.id}/share-info`,
				{ withCredentials: true },
			)

			// Set the default valuel of the Share Dialog (the users the file is already shared with)
			const usersWithAccessToFile: User[] = await res.data.sharedWith;
			const userIds = usersWithAccessToFile.map((user) => user.id);
			setShareDialogDefaultValue(userIds);

			// Set the value of the Share URL
			setShareDialogURL(res.data.shareURL);

			// Set the value of the list of all possible users the file can be shared with
			setShareDialogOptions(mapUsersToOptions(res.data.allUsers));

		} catch (error) {
			console.log("error while fetching users with access to file: ", error)
		}
	}
	useEffect(() => {
		if (file.user_owns_file && file.item_type === 'file' && isShareDialogOpen) {
			fetchShareInfo()
		}
	}, [isShareDialogOpen])


	/**
	 * Deletes the current file.
	 * - Sends DELETE request to API
	 * - Shows success or error toast notifications
	 * - Updates the UI by removing the deleted file (without refetch)
	 * 
	 * @async
	 * @function
	 */
	const handleDelete = async () => {
		try {
			const res = await api.delete(`/files/${file.id}`, { withCredentials: true });
			if (res.status === 204) {
				toast.success("Deleted file successfully");
				fetchUser();
				deleteItem(file.id);
			} else {
				toast.error(res.data.error);
			}
		} catch (error: any) {
			console.log('Error while deleting file: ', error);
			toast.error(error.response.data.message);
		}
	}


	/**
	 * Downloads the current file.
	 * - Sends GET request to API for file blob
	 * - Creates a temporary link to trigger download
	 * - Shows success or error toast notifications
	 * 
	 * @async
	 * @function
	 */
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

	/**
	 * Renames the current file.
	 * - Sends PATCH request with new filename
	 * - Updates UI via `renameItem` callback
	 * - Shows success or error toast notifications
	 *
	 * @async
	 * @function
	 * @param {string} newFilename - The new filename for the file
	 */
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
		}
	}

	/**
	 * Shares the current file with selected users.
	 * - Sends PUT request to API to update file shares
	 * - Shows success or error toast notifications
	 * 
	 * @async
	 * @function
	 * @param {string[]} selectedUsers - Array of user IDs (as strings) to share the file with
	 */
	const handleShare = async (selectedUsers: string[]) => {
		try {
			const userIdsAsNumbers = selectedUsers.map(id => parseInt(id, 10))
			const res = await api.put(
				`/files/${file.id}/shares`,
				{ user_ids: userIdsAsNumbers },
				{
					headers: { "Content-Type": "application/json" },
					withCredentials: true,
				}
			);
			toast.success(res.data.message);

		} catch (error: any) {
			console.error(error);
			toast.error(error.response?.data?.error || "Failed to share file");
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
						<DropdownMenuItem onSelect={() => setMoveDialogOpen(true)} disabled={file.user_owns_file ? false : true}>
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

			<MoveDialogModal
				currentFolderId={path[path.length - 1].id}
				fileId={file.id}
				isOpen={isMoveDialogOpen}
				isOpenChange={setMoveDialogOpen}
				onConfirm={onFileChange}
				context={"file"}
			/>

		</div>
	)
}
