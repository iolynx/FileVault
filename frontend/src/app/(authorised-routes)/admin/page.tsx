import AdminFilesTable from "@/components/AdminFilesTable";
import { AdminPageClient } from "@/components/AdminPageClient";
import { Card } from "@/components/ui/card";
import { getCurrentUser } from "@/lib/auth";
import { redirect } from "next/navigation";

export default async function AdminPage() {
  const user = await getCurrentUser();

  if (user?.role !== 'admin') {
    redirect('/dashboard');
  }

  return (
    <div className="flex flex-col items-center">
      <div className="flex flex-col items-center my-10">
        <h1 className="text-3xl font-bold"> Admin Panel </h1>
        <p className="text-muted-foreground">Welcome, admin. Here you can view the files of all the users.</p>
        <Card className="rounded-2xl border shadow-sm overflow-hidden w-full max-w-7xl min-w-6xl mt-4 pt-1 pb-1">
          <AdminPageClient />
        </Card>
      </div>
    </div>

  );
}
