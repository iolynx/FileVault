import { Dispatch, SetStateAction, useState } from "react";
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";

interface RenameDialogModalProps {
	isOpen: boolean;
	isOpenChange: (open: boolean) => void;
	originalFilename: string;
	onConfirm: (newFileName: string) => void;
}

export function RenameDialogModal({ isOpen, isOpenChange, originalFilename, onConfirm }: RenameDialogModalProps) {
	const [newFilename, setNewFilename] = useState("")
	return (
		<Dialog open={isOpen} onOpenChange={isOpenChange}>
			<DialogContent className="sm:max-w-md">
				<DialogHeader>
					<DialogTitle>Rename File</DialogTitle>
				</DialogHeader>
				<div className="flex items-center gap-2">
					<div className="grid flex-1 gap-2">
						<Label htmlFor="link" className="sr-only">
							New Filename
						</Label>
						<Input
							id="filename"
							defaultValue={originalFilename}
							onChange={(e) => setNewFilename(e.target.value)}
						/>
					</div>
				</div>
				<DialogFooter className="sm:justify-start">
					<DialogClose asChild>
						<Button type="button" variant="secondary" onClick={(e) => onConfirm(newFilename)}>
							Done
						</Button>
					</DialogClose>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
};

