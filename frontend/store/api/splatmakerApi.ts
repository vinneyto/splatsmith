import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import type { RootState } from "@/store/store";

const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export type LoginRequest = {
  username: string;
  password: string;
};

export type LoginResponse = {
  access_token: string;
  token_type: "Bearer";
  user: {
    user_id: string;
    email: string;
  };
};

export type JobStatus = "new" | "queued" | "in_progress" | "done" | "failed" | "cancelled";

export type JobItem = {
  job_id: string;
  status: JobStatus;
  progress_percent: number;
  current_step?: string | null;
  idempotency_key?: string | null;
  created_at: string;
  updated_at: string;
};

export type ListJobsResponse = {
  items: JobItem[];
};

export type OutputFileRef = {
  key: string;
  file_name: string;
  size_bytes?: number | null;
};

export type JobDetails = {
  summary: JobItem;
  attempt: number;
  source_ref?: string | null;
  simulate_failure: boolean;
  error_message?: string | null;
  started_at?: string | null;
  finished_at?: string | null;
  last_heartbeat_at?: string | null;
  output_files: OutputFileRef[];
};

export type ResultFileURL = {
  key: string;
  file_name: string;
  url: string;
  expires_at: string;
};

export type JobResultUrlsResponse = {
  items: ResultFileURL[];
};

export const splatmakerApi = createApi({
  reducerPath: "splatmakerApi",
  baseQuery: fetchBaseQuery({
    baseUrl,
    prepareHeaders: (headers, { getState }) => {
      const state = getState() as RootState;
      const token = state.auth.token;
      if (token) {
        headers.set("Authorization", `Bearer ${token}`);
      }
      headers.set("Content-Type", "application/json");
      return headers;
    },
  }),
  tagTypes: ["Jobs"],
  endpoints: (builder) => ({
    login: builder.mutation<LoginResponse, LoginRequest>({
      query: (body) => ({
        url: "/v1/auth/login",
        method: "POST",
        body,
      }),
    }),
    listJobs: builder.query<ListJobsResponse, void>({
      query: () => ({
        url: "/v1/jobs",
        method: "GET",
      }),
      providesTags: ["Jobs"],
    }),
    getJob: builder.query<JobDetails, string>({
      query: (jobId) => ({
        url: `/v1/jobs/${jobId}`,
        method: "GET",
      }),
      providesTags: (_result, _error, jobId) => [{ type: "Jobs", id: jobId }],
    }),
    getJobResultUrls: builder.query<JobResultUrlsResponse, { jobId: string; ttlSeconds?: number }>({
      query: ({ jobId, ttlSeconds }) => ({
        url: `/v1/jobs/${jobId}/result-urls`,
        method: "GET",
        params: ttlSeconds ? { ttl_seconds: ttlSeconds } : undefined,
      }),
      providesTags: (_result, _error, arg) => [{ type: "Jobs", id: arg.jobId }],
    }),
  }),
});

export const { useLoginMutation, useListJobsQuery, useGetJobQuery, useGetJobResultUrlsQuery } = splatmakerApi;
