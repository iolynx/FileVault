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
}

export function DeleteDialogModal({ isOpen, isOpenChange, onConfirm }: DeleteDialogModalProps) {
	return (
		<AlertDialog open={isOpen} onOpenChange={isOpenChange}>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>
						Are you sure you want to delete this File?
					</AlertDialogTitle>
					<AlertDialogDescription>
						This action cannot be undone. This will permanently delete the File.
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

