"use client";

import FilesTable from "@/components/FilesTable";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import api from "@/lib/axios";
import { APIError } from "@/types/APIError";
import { Filter } from "@/types/Filter";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { File } from "@/types/File"
import { FileUploadMenu } from "@/components/FileUploadMenu";
import { User } from "@/types/User";
import { Skeleton } from "@/components/ui/skeleton";
import { Router } from "next/router";
import { useRouter } from "next/navigation";

const DashboardPage = () => {
	const [loading, setLoading] = useState(true);
	const [files, setFiles] = useState<File[]>([]);
	const [search, setSearch] = useState("");
	const [filter, setFilter] = useState<Filter>();
	const [page, setPage] = useState(0);
	const router = useRouter()

	const fetchFiles = async () => {
		try {
			setLoading(true)
			const res = await api.get("/files", {
				params: { search: search, limit: 10, offset: 0 },
				headers: { "Content-Type": "application/json" },
				withCredentials: true,
			})

			// TODO: remove for prod
			console.log('Files RESULT IS: ', res.data)
			setFiles(res.data);
			router.refresh()
		} catch (error: any) {
			const err = error as APIError
			toast.error("Error: Failed to fetch files")
		} finally {
			setLoading(false)
		}
	};

	useEffect(() => {
		fetchFiles()
	}, [search])

	return (
		<div className="flex flex-col items-center">
			<div className="flex flex-col items-center my-10">
				<h1 className="text-3xl font-bold"> Dashboard </h1>
				<p> View and Manage your files here.</p>
			</div>
			<div>
				<div className="flex flex-row rounded-2xl border shadow-sm overflow-hidden w-full mt-2 p-4 gap-x-4">
					<div className="flex flex-row gap-x-2">
						<Input
							type="text"
							id="search"
							placeholder="Search by Filename"
							value={search}
							onChange={(e) => setSearch(e.target.value)}
						/>
					</div>
				</div>
			</div>
			<div className="m-4">
				<FileUploadMenu fetchFiles={fetchFiles} />
			</div>

			{loading ?
				<Skeleton className="h-[500px] w-[86%] rounded-2xl max-w-7xl mt-4" /> :
				<Card className="rounded-2xl border shadow-sm overflow-hidden w-full max-w-7xl mt-4 pt-1 pb-1">
					<FilesTable files={files} onFileChange={() => fetchFiles()} />

				</Card>
			}
		</div>

	)
}
export default DashboardPage;
