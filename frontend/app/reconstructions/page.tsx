"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useListJobsQuery } from "@/store/api/splatmakerApi";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { logout } from "@/store/slices/authSlice";

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
    <main style={{ maxWidth: 860, margin: "0 auto", padding: 24 }}>
      <header style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <div>
          <h1 style={{ marginBottom: 4 }}>Reconstructions</h1>
          {user && <p style={{ marginTop: 0, color: "#667085" }}>Signed in as {user.email}</p>}
        </div>
        <button
          onClick={() => {
            dispatch(logout());
            router.push("/login");
          }}
        >
          Logout
        </button>
      </header>

      {isLoading && <p>Loading reconstructions...</p>}
      {isError && (
        <p style={{ color: "#b42318" }}>Failed to load reconstructions: {JSON.stringify(error)}</p>
      )}
      {!isLoading && !isError && data?.items.length === 0 && <p>No reconstructions yet.</p>}

      <ul style={{ listStyle: "none", padding: 0, display: "grid", gap: 12 }}>
        {data?.items.map((job) => (
          <li key={job.job_id} style={{ background: "white", padding: 16, borderRadius: 10 }}>
            <Link
              href={`/reconstructions/${job.job_id}`}
              style={{ textDecoration: "none", color: "inherit" }}
            >
              <div style={{ display: "flex", justifyContent: "space-between", gap: 12 }}>
                <strong>{job.job_id}</strong>
                <span>Status: {job.status}</span>
              </div>
              <small style={{ color: "#667085" }}>
                Updated: {new Date(job.updated_at).toLocaleString()}
              </small>
            </Link>
          </li>
        ))}
      </ul>
    </main>
  );
}
