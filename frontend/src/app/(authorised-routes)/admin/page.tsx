import { getCurrentUser } from "@/lib/auth";
import { redirect } from "next/navigation";

export default async function AdminPage() {
  const user = await getCurrentUser();

  if (user?.role !== 'admin') {
    redirect('/dashboard');
  }

  return (
    <div className="flex flex-col justify-center">
      <h1 className="text-2xl font-bold">Admin Panel</h1>
      <p className="text-muted-foreground">Welcome, admin. Here you can manage all files and users.</p>

    </div>
  );
}
