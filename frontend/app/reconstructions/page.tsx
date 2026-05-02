'use client';

import Link from 'next/link';
import { useState } from 'react';
import { useListJobsQuery, type JobStatus } from '@/store/api/splatmakerApi';

const statuses: Array<JobStatus | ''> = ['', 'queued', 'running', 'succeeded', 'failed', 'canceled'];

export default function ReconstructionsPage() {
  const [status, setStatus] = useState<JobStatus | ''>('');
  const { data, isLoading, isFetching, error } = useListJobsQuery(status ? { status, limit: 100 } : { limit: 100 });

  return (
    <main className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Jobs Viewer</h1>

      <div>
        <label className="mr-2 text-sm">Status:</label>
        <select className="border rounded px-2 py-1" value={status} onChange={(e) => setStatus(e.target.value as JobStatus | '')}>
          {statuses.map((item) => (
            <option key={item || 'all'} value={item}>
              {item || 'all'}
            </option>
          ))}
        </select>
      </div>

      {isLoading || isFetching ? <p>Loading…</p> : null}
      {error ? <p className="text-red-600">Failed to load jobs.</p> : null}

      <ul className="space-y-2">
        {(data?.items ?? []).map((job) => (
          <li key={job.job_id} className="border rounded p-3">
            <div className="font-medium">{job.job_id}</div>
            <div className="text-sm text-gray-600">status: {job.status}</div>
            <div className="text-sm text-gray-600">updated: {new Date(job.updated_at).toLocaleString()}</div>
            <Link className="text-blue-600 underline text-sm" href={`/reconstructions/${job.job_id}`}>
              Open details
            </Link>
          </li>
        ))}
      </ul>
    </main>
  );
}
