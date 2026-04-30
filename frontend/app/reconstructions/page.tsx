"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { ThemeToggle } from "@/components/theme-toggle";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { useListJobsQuery } from "@/store/api/splatmakerApi";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { logout } from "@/store/slices/authSlice";

function statusVariant(status: string): "success" | "warning" | "danger" | "default" {
  if (status === "done") return "success";
  if (status === "in_progress" || status === "queued" || status === "new") return "warning";
  if (status === "failed" || status === "cancelled") return "danger";
  return "default";
}

export default function ReconstructionsPage() {
  const router = useRouter();
  const dispatch = useAppDispatch();
  const token = useAppSelector((state) => state.auth.token);
  const user = useAppSelector((state) => state.auth.user);
  const { data, isLoading, isError, error } = useListJobsQuery(undefined, { skip: !token });

  useEffect(() => {
    if (!token) router.replace("/login");
  }, [router, token]);

  return (
    <main className="mx-auto max-w-4xl p-6">
      <header className="mb-6 flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold">Reconstructions</h1>
          {user && <p className="text-sm text-muted-foreground">Signed in as {user.email}</p>}
        </div>
        <div className="flex items-center gap-2">
          <ThemeToggle />
          <Button
            onClick={() => {
              dispatch(logout());
              router.push("/login");
            }}
            variant="outline"
          >
            Logout
          </Button>
        </div>
      </header>

      {isLoading && <p className="text-muted-foreground">Loading reconstructions...</p>}
      {isError && (
        <Card className="border-rose-300 bg-rose-50 dark:border-rose-900 dark:bg-rose-950/40">
          <CardContent className="p-4 text-sm text-rose-700 dark:text-rose-300">
            Failed to load reconstructions: {JSON.stringify(error)}
          </CardContent>
        </Card>
      )}
      {!isLoading && !isError && data?.items.length === 0 && (
        <Card>
          <CardContent className="p-6 text-sm text-muted-foreground">No reconstructions yet.</CardContent>
        </Card>
      )}

      <section className="grid gap-3">
        {data?.items.map((job) => (
          <Link href={`/reconstructions/${job.job_id}`} key={job.job_id}>
            <Card className="transition hover:bg-accent/40">
              <CardHeader className="pb-3">
                <div className="flex items-start justify-between gap-3">
                  <CardTitle className="text-base">{job.job_id}</CardTitle>
                  <Badge variant={statusVariant(job.status)}>{job.status}</Badge>
                </div>
                <CardDescription>Updated: {new Date(job.updated_at).toLocaleString()}</CardDescription>
              </CardHeader>
            </Card>
          </Link>
        ))}
      </section>
    </main>
  );
}
