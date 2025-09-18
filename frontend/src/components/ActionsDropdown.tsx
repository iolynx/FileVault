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
import { useEffect, useState } from "react";
import { RenameDialogModal } from "./RenameDialogModal";
import { Axios, AxiosError } from "axios";
import { mapUsersToOptions } from "@/lib/utils";
import { User } from "@/types/User";
import { MultiSelectOption } from "./multi-select";
import { ShareDialogModal } from "./ShareDialogModal";

interface ActionsDropDownProps {
	file: File;
	onFileChange: () => void;
}

export default function ActionsDropdown({ file, onFileChange }: ActionsDropDownProps) {
	const [isDeleteDialogOpen, setDeleteDialogOpen] = useState(false)
	const [isRenameDialogOpen, setRenameDialogOpen] = useState(false)
	const [isShareDialogOpen, setShareDialogOpen] = useState(false)
	const [shareDialogOptions, setShareDialogOptions] = useState<MultiSelectOption[]>([]);
	const [shareDialogDefautValue, setShareDialogDefaultValue] = useState<string[]>([]);

	const fetchUsers = async () => {
		try {
			const res = await api.get("/users",
				{ withCredentials: true },
			)
			const users: User[] = await res.data
			setShareDialogOptions(mapUsersToOptions(users));
		} catch (error) {
			console.log("error while fetching users:", error)
		}
	}

	const fetchUsersWithAccessToFile = async () => {
		try {
			const res = await api.get(`/files/${file.id}/shares`,
				{ withCredentials: true },
			)
			console.log(res.data);
			const usersWithAccessToFile: User[] = await res.data
			const userIds = usersWithAccessToFile.map((user) => user.id);
			setShareDialogDefaultValue(userIds);
		} catch (error) {
			console.log("error while fetching users with access to file: ", error)
		}
	}
	useEffect(() => {
		fetchUsers()
		fetchUsersWithAccessToFile()
	}, [isShareDialogOpen])

	const handleDelete = async () => {
		try {
			const res = await api.delete(`/files/${file.id}`, { withCredentials: true });
			if (res.status === 200) {
				toast.success(res.data.message);
				onFileChange();
			} else {
				toast.error(res.data.error);
			}
		} catch (error: any) {
			console.log('Error while deleting file: ', error);
			toast.error(error.response.data.message);
		} finally {
			onFileChange();
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
			console.log(res.data)
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
				{ filename: newFilename },
				{ headers: { "Content-Type": "application/json" }, withCredentials: true }
			)
			toast.success(res.data.message);
		} catch (error: any) {
			toast.error(error.response.data.message)
		} finally {
			onFileChange();
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
					{ target_user_id: addedUser },
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
			onFileChange();
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
				onConfirm={() => handleDelete()}
			/>

			<RenameDialogModal
				isOpen={isRenameDialogOpen}
				isOpenChange={setRenameDialogOpen}
				originalFilename={file.filename}
				onConfirm={(newFilename) => handleRename(newFilename)}
			/>

			<ShareDialogModal
				isOpen={isShareDialogOpen}
				isOpenChange={setShareDialogOpen}
				userOptions={shareDialogOptions}
				onConfirm={(usersToShare) => handleShare(usersToShare)}
				defaultValue={shareDialogDefautValue}
			/>

		</div>
	)
}
