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

export type JobItem = {
  job_id: string;
  status: "new" | "in_progress" | "done" | "failed";
  created_at: string;
  updated_at: string;
};

export type ListJobsResponse = {
  items: JobItem[];
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
  }),
});

export const { useLoginMutation, useListJobsQuery } = splatmakerApi;
