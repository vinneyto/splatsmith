'use client';

import { useParams } from 'next/navigation';
import { useGetJobQuery, useGetJobResultUrlsQuery } from '@/store/api/splatmakerApi';

export default function ReconstructionDetailsPage() {
  const params = useParams<{ id: string }>();
  const jobId = params.id;

  const { data: job, isLoading, error } = useGetJobQuery(jobId);
  const { data: resultUrls } = useGetJobResultUrlsQuery({ jobId, ttlSeconds: 3600 }, { skip: !jobId });

  if (isLoading) return <main className="p-6">Loading…</main>;
  if (error || !job) return <main className="p-6 text-red-600">Failed to load job details.</main>;

  return (
    <main className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Job {job.summary.job_id}</h1>
      <div className="text-sm text-gray-700">status: {job.summary.status}</div>
      <div className="text-sm text-gray-700">source: {job.source_ref}</div>
      {job.error_message ? <div className="text-sm text-red-600">error: {job.error_message}</div> : null}

      <section className="space-y-2">
        <h2 className="text-lg font-medium">Files</h2>
        {(resultUrls?.items ?? []).length === 0 ? <p className="text-sm text-gray-500">No files</p> : null}
        <ul className="space-y-2">
          {(resultUrls?.items ?? []).map((f) => (
            <li key={f.key} className="border rounded p-3">
              <div className="text-sm">{f.file_name}</div>
              <a className="text-blue-600 underline text-sm break-all" href={f.url} target="_blank" rel="noreferrer">
                {f.url}
              </a>
            </li>
          ))}
        </ul>
      </section>
    </main>
  );
}
