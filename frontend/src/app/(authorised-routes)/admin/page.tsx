import { AdminPageClient } from "@/components/AdminPageClient";
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
        <p className="text-muted-foreground mb-12">Welcome, admin. Here you can view the files of all the users.</p>
        <div >
          <AdminPageClient />
        </div>
      </div>
    </div>

  );
}
