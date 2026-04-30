"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAppSelector } from "@/store/hooks";
import { buttonVariants } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

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
    <main className="mx-auto max-w-3xl p-6">
      <Link href="/reconstructions" className={buttonVariants({ variant: "outline", size: "sm" }) + " mb-4 inline-flex"}>
        Back to list
      </Link>

      <Card>
        <CardHeader>
          <CardTitle>Reconstruction {params.id}</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            Placeholder: detailed reconstruction view will be implemented next.
          </p>
        </CardContent>
      </Card>
    </main>
  );
}
