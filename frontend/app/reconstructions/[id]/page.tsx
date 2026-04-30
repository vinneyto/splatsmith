"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useEffect } from "react";
import { useGetJobQuery, useGetJobResultUrlsQuery } from "@/store/api/splatmakerApi";
import { useAppSelector } from "@/store/hooks";
import { buttonVariants } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function ReconstructionDetailsPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const token = useAppSelector((state) => state.auth.token);
  const jobId = params.id;

  const { data, isLoading, isError, error } = useGetJobQuery(jobId, {
    skip: !token || !jobId,
  });

  const shouldLoadResultUrls = Boolean(token && jobId && data?.summary.status === "done");
  const {
    data: resultUrls,
    isLoading: resultLoading,
    isError: resultError,
    error: resultErrorData,
  } = useGetJobResultUrlsQuery(
    { jobId, ttlSeconds: 900 },
    {
      skip: !shouldLoadResultUrls,
    }
  );

  useEffect(() => {
    if (!token) {
      router.replace("/login");
    }
  }, [router, token]);

  return (
    <main className="mx-auto max-w-3xl p-6">
      <Link
        href="/reconstructions"
        className={buttonVariants({ variant: "outline", size: "sm" }) + " mb-4 inline-flex"}
      >
        Back to list
      </Link>

      <Card>
        <CardHeader>
          <CardTitle>Reconstruction {jobId}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {isLoading ? <p className="text-sm text-muted-foreground">Loading job details...</p> : null}
          {isError ? (
            <p className="text-sm text-red-600">Failed to load job: {JSON.stringify(error)}</p>
          ) : null}

          {data ? (
            <>
              <p className="text-sm">
                Status: <strong>{data.summary.status}</strong>
              </p>
              <p className="text-xs text-muted-foreground">
                Updated: {new Date(data.summary.updated_at).toLocaleString()}
              </p>

              {data.summary.status === "done" ? (
                <div className="pt-2">
                  <h3 className="mb-2 text-sm font-semibold">Result files</h3>
                  {resultLoading ? <p className="text-sm text-muted-foreground">Loading files...</p> : null}
                  {resultError ? (
                    <p className="text-sm text-red-600">
                      Failed to load files: {JSON.stringify(resultErrorData)}
                    </p>
                  ) : null}

                  {!resultLoading && !resultError && (resultUrls?.items.length ?? 0) === 0 ? (
                    <p className="text-sm text-muted-foreground">No files available.</p>
                  ) : null}

                  <ul className="space-y-2 text-sm">
                    {resultUrls?.items.map((file) => (
                      <li key={file.key} className="rounded-md border p-2">
                        <a href={file.url} target="_blank" rel="noreferrer" className="font-medium hover:underline">
                          {file.file_name}
                        </a>
                        <p className="text-xs text-muted-foreground">{file.key}</p>
                      </li>
                    ))}
                  </ul>
                </div>
              ) : (
                <p className="text-sm text-muted-foreground">
                  Files will be available when the job status becomes done.
                </p>
              )}
            </>
          ) : null}
        </CardContent>
      </Card>
    </main>
  );
}
