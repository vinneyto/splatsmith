import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

export type JobStatus = 'queued' | 'running' | 'succeeded' | 'failed' | 'canceled';

export interface JobSummary {
  job_id: string;
  status: JobStatus;
  progress_percent?: number;
  current_step?: string;
  idempotency_key?: string;
  created_at: string;
  updated_at: string;
}

export interface OutputFileRef {
  key: string;
  file_name: string;
  size_bytes?: number;
}

export interface JobDetails {
  summary: JobSummary;
  attempt: number;
  source_ref: string;
  simulate_failure?: boolean;
  error_message?: string;
  started_at?: string;
  finished_at?: string;
  last_heartbeat_at?: string;
  output_files: OutputFileRef[];
}

export interface JobResultURL {
  key: string;
  file_name: string;
  url: string;
  expires_at: string;
}

export const splatmakerApi = createApi({
  reducerPath: 'splatmakerApi',
  baseQuery: fetchBaseQuery({
    baseUrl: process.env.NEXT_PUBLIC_API_BASE_URL || '',
  }),
  tagTypes: ['Jobs'],
  endpoints: (builder) => ({
    listJobs: builder.query<{ items: JobSummary[] }, { status?: JobStatus; limit?: number; offset?: number } | undefined>({
      query: (params) => ({
        url: '/v1/jobs',
        params: params ?? undefined,
      }),
      providesTags: ['Jobs'],
    }),
    getJob: builder.query<JobDetails, string>({
      query: (jobId) => `/v1/jobs/${jobId}`,
      providesTags: (_r, _e, id) => [{ type: 'Jobs', id }],
    }),
    getJobResultUrls: builder.query<{ items: JobResultURL[] }, { jobId: string; ttlSeconds?: number }>({
      query: ({ jobId, ttlSeconds }) => ({
        url: `/v1/jobs/${jobId}/result-urls`,
        params: ttlSeconds ? { ttl_seconds: ttlSeconds } : undefined,
      }),
    }),
  }),
});

export const { useListJobsQuery, useGetJobQuery, useGetJobResultUrlsQuery } = splatmakerApi;
