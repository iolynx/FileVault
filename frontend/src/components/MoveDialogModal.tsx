import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select"
import { useEffect, useState } from "react"
import api from "@/lib/axios"
import { toast } from "sonner"
import { toSentenceCase } from "@/lib/utils"

interface MoveDialogModalProps {
	isOpen: boolean
	isOpenChange: (open: boolean) => void
	fileId: string
	currentFolderId: string | null
	onConfirm: () => void
	context: "file" | "folder"
}

interface Folder {
	id: string
	name: string
}

export function MoveDialogModal({
	isOpen,
	isOpenChange,
	fileId,
	currentFolderId,
	onConfirm,
	context,
}: MoveDialogModalProps) {
	const [folders, setFolders] = useState<Folder[]>([])
	const [selectedFolderId, setSelectedFolderId] = useState<string>("")

	useEffect(() => {
		if (isOpen) {
			fetchFolders()
		}
	}, [isOpen])

	const fetchFolders = async () => {
		try {
			let res;
			if (currentFolderId === null) {
				res = await api.get(`/folders/`, {
					withCredentials: true,
				})
			} else {
				res = await api.get(`/folders/${fileId}`, {
					withCredentials: true,
				})
			}
			let fetchedFolders = res.data; // all other folders
			// if we're inside a folder, add a "Home" option
			if (currentFolderId !== null) {
				fetchedFolders = [
					{
						id: null,
						name: "Home",		//for now, calling "/" as "Home"
						createdAt: null,
						parentFolderId: null
					},
					...fetchedFolders
				]
			}

			setFolders(fetchedFolders.filter((folder: Folder) => (folder.id !== fileId)))

		} catch (error) {
			console.error("Error fetching folders:", error)
			toast.error("Failed to load folders")
		}
	}

	const handleMove = async () => {
		if (!selectedFolderId && selectedFolderId != null) {
			toast.error("Please select a folder")
			return
		}
		try {
			await api.patch(
				`/${context}s/${fileId}/move`,
				{ target_folder_id: selectedFolderId },
				{ headers: { "Content-Type": "application/json" }, withCredentials: true }
			)
			toast.success(`${toSentenceCase(context)} moved successfully`)
			isOpenChange(false)
			onConfirm()
		} catch (error) {
			console.error(error)
			toast.error("Failed to move file")
		}
	}

	return (
		<Dialog open={isOpen} onOpenChange={isOpenChange}>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Move {toSentenceCase(context)}</DialogTitle>
					<DialogDescription>
						Select a folder where this {context} should be moved
					</DialogDescription>
				</DialogHeader>

				<Select onValueChange={setSelectedFolderId}>
					<SelectTrigger className="w-full">
						<SelectValue placeholder="Choose folder" />
					</SelectTrigger>
					<SelectContent>
						{folders.map((folder) => (
							<SelectItem key={folder.id} value={folder.id}>
								{folder.name}
							</SelectItem>
						))}
					</SelectContent>
				</Select>

				<DialogFooter>
					<Button variant="outline" onClick={() => isOpenChange(false)}>
						Cancel
					</Button>
					<Button onClick={handleMove}>Move</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	)
}

