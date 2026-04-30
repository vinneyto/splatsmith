"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useListJobsQuery } from "@/store/api/splatmakerApi";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { logout } from "@/store/slices/authSlice";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function ReconstructionsPage() {
  const router = useRouter();
  const dispatch = useAppDispatch();
  const token = useAppSelector((state) => state.auth.token);
  const user = useAppSelector((state) => state.auth.user);
  const { data, isLoading, isError, error } = useListJobsQuery(undefined, {
    skip: !token,
  });

  useEffect(() => {
    if (!token) {
      router.replace("/login");
    }
  }, [router, token]);

  return (
    <main className="mx-auto max-w-3xl p-6">
      <header className="mb-6 flex items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">Reconstructions</h1>
          {user ? <p className="text-sm text-muted-foreground">Signed in as {user.email}</p> : null}
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={() => {
            dispatch(logout());
            router.push("/login");
          }}
        >
          Logout
        </Button>
      </header>

      {isLoading ? <p className="text-sm text-muted-foreground">Loading reconstructions...</p> : null}
      {isError ? (
        <p className="text-sm text-red-600">Failed to load reconstructions: {JSON.stringify(error)}</p>
      ) : null}
      {!isLoading && !isError && data?.items.length === 0 ? (
        <p className="text-sm text-muted-foreground">No reconstructions yet.</p>
      ) : null}

      <div className="grid gap-3">
        {data?.items.map((job) => (
          <Link key={job.job_id} href={`/reconstructions/${job.job_id}`} className="group block">
            <Card className="transition-colors group-hover:bg-accent/40">
              <CardHeader className="pb-3">
                <CardTitle className="text-base">{job.job_id}</CardTitle>
              </CardHeader>
              <CardContent className="flex items-center justify-between pt-0">
                <p className="text-xs text-muted-foreground">
                  Updated: {new Date(job.updated_at).toLocaleString()}
                </p>
                <Badge variant="secondary">{job.status}</Badge>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>
    </main>
  );
}
