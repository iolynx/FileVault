import { Skeleton } from "@/components/ui/skeleton";
import { Card } from "@/components/ui/card";


/**
 * FileSkeleton: Displays a skeleton component when the data is being fetched
 */
export default function FilesSkeleton() {
	return (
		<div className="flex flex-col p-0 m-0">
			{/* Breadcrumb skeleton */}
			<div className="flex items-center gap-2 mt-2">
				<Skeleton className="h-5 w-16" />
			</div>

			{/* Card skeleton */}
			<Card className="rounded-2xl justify-around overflow-hidden max-w-7xl min-w-6xl mt-4 pt-1 pb-1 p-4">
				{/* Table header skeleton */}
				<div className="flex justify-between mb-4">
					<Skeleton className="h-6 w-24" />
					<Skeleton className="h-6 w-20" />
				</div>


				{/* Pagination skeleton */}
				<div className="flex justify-between items-center mt-6">
					<Skeleton className="h-5 w-24" />
					<div className="flex gap-2">
						<Skeleton className="h-8 w-8 rounded-md" />
						<Skeleton className="h-8 w-8 rounded-md" />
						<Skeleton className="h-8 w-8 rounded-md" />
					</div>
				</div>
			</Card>
		</div>
	);
}

