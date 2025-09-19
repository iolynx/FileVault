import { Dispatch, SetStateAction, useEffect, useState } from "react";
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
import { MultiSelect, MultiSelectOption } from "./multi-select";
import { toast } from "sonner";

interface ShareDialogModalProps {
	isOpen: boolean;
	isOpenChange: (open: boolean) => void;
	userOptions: MultiSelectOption[];
	onConfirm: (usersToShare: string[]) => void;
	defaultValue: string[]
	fileURL: string;
}

export function ShareDialogModal({ isOpen, isOpenChange, userOptions, onConfirm, defaultValue, fileURL }: ShareDialogModalProps) {
	const handleCopy = async () => {
		try {
			await navigator.clipboard.writeText(fileURL);
			toast.success("Copied file URL to clipboard")
		} catch (err) {
			console.log("Failed to copy file URL", err)
			toast.error("Could not copy file URL to clipboard")
		}
	}
	const [usersToShare, setUsersToShare] = useState<string[]>([]);
	return (
		<Dialog open={isOpen} onOpenChange={isOpenChange}>
			<DialogContent className="sm:max-w-md">
				<DialogHeader>
					<DialogTitle>Share To:</DialogTitle>
				</DialogHeader>
				<div className="flex items-center gap-2">
					<div className="grid flex-1 gap-2">
						<MultiSelect
							options={userOptions}
							onValueChange={setUsersToShare}
							defaultValue={defaultValue}
							placeholder="Search For Users"
							animationConfig={{
								badgeAnimation: "none",
							}}
						/>
					</div>
				</div>
				<DialogFooter className="sm:justify-between">
					<DialogClose asChild>
						<Button type="button" variant="outline" onClick={() => onConfirm(usersToShare)}>
							Done
						</Button>
					</DialogClose>
					<div>
						<Button variant="outline" onClick={handleCopy}>
							Copy Link
						</Button>
					</div>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
};

