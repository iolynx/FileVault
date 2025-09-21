import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from "@/components/ui/alert-dialog"

interface DeleteDialogModalProps {
	isOpen: boolean;
	isOpenChange: (open: boolean) => void;
	onConfirm: () => void;
	context: "Folder" | "File"
}

export function DeleteDialogModal({ isOpen, isOpenChange, onConfirm, context }: DeleteDialogModalProps) {
	return (
		<AlertDialog open={isOpen} onOpenChange={isOpenChange}>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>
						Are you sure you want to delete this {context}?
					</AlertDialogTitle>
					<AlertDialogDescription>
						This will permanently delete the {context}, {context === "Folder" && (<p>and recursively delete all files and subfolders within this folder.</p>)}
						This action cannot be undone.
					</AlertDialogDescription>
				</AlertDialogHeader>
				<AlertDialogFooter>
					<AlertDialogCancel className="text-red hover:bg-red-100">Cancel</AlertDialogCancel>
					<AlertDialogAction onClick={onConfirm}>Yes, Delete it</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	);
};

