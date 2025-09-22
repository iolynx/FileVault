import { getCurrentUser } from "@/lib/auth";
import { redirect } from "next/navigation";
import { ActivityChart } from "@/components/admin/ActivityChart";
import { AuditLogTable } from "@/components/admin/AuditLogTable";

// This remains a Server Component for the initial auth check
export default async function AuditLogPage() {
  const user = await getCurrentUser();

  if (user?.role !== 'admin') {
    redirect('/dashboard');
  }

  return (
    <div className="container mx-auto py-10 space-y-8 flex flex-col">
      <div className="flex flex-col items-center mb-10 pb-10">
        <h1 className="text-3xl font-bold">Audit Logs</h1>
        <p className="text-muted-foreground">Monitor application activity and user actions.</p>
      </div>

      <div className="space-y-4 ">
        <ActivityChart />
        <AuditLogTable />
      </div>
    </div>
  );
}
