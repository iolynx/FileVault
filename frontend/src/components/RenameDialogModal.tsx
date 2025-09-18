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
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";

interface RenameDialogModalProps {
	isOpen: boolean;
	isOpenChange: (open: boolean) => void;
	originalFilename: string;
	onConfirm: (newFileName: string) => void;
}

export function RenameDialogModal({ isOpen, isOpenChange, originalFilename, onConfirm }: RenameDialogModalProps) {
	const [baseName, setBaseName] = useState('');
	const [extension, setExtension] = useState('');

	useEffect(() => {
		const lastDotIndex = originalFilename.lastIndexOf('.');

		// Handle files with no extension
		if (lastDotIndex === -1) {
			setBaseName(originalFilename);
			setExtension('');
		} else {
			setBaseName(originalFilename.substring(0, lastDotIndex));
			setExtension(originalFilename.substring(lastDotIndex + 1));
		}
	}, [originalFilename]);

	const handleSave = () => {
		if (!baseName.trim()) {
			alert("Filename cannot be empty.");
			return;
		}

		const newFilename = extension ? `${baseName.trim()}.${extension}` : baseName.trim();
		console.log(newFilename)
		onConfirm(newFilename);
	};

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
						<div className="flex flex-row align-middle text-center">
							<Input
								value={baseName}
								onChange={(e) => setBaseName(e.target.value)}
								className="flex-grow"
							/>
							{extension && <span className="text-gray-500 mt-2">.{extension}</span>}
						</div>
					</div>
				</div>
				<DialogFooter className="sm:justify-start">
					<DialogClose asChild>
						<Button type="button" variant="secondary" onClick={handleSave}>
							Done
						</Button>
					</DialogClose>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
};

