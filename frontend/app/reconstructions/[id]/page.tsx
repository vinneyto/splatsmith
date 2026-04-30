"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAppSelector } from "@/store/hooks";

export default function ReconstructionDetailsPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const token = useAppSelector((state) => state.auth.token);

  useEffect(() => {
    if (!token) {
      router.replace("/login");
    }
  }, [router, token]);

  return (
    <main style={{ maxWidth: 860, margin: "0 auto", padding: 24 }}>
      <p>
        <Link href="/reconstructions">← Back to list</Link>
      </p>
      <h1>Reconstruction {params.id}</h1>
      <p>Placeholder: detailed reconstruction view will be implemented next.</p>
    </main>
  );
}
