"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useEffect } from "react";
import { ThemeToggle } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { useAppSelector } from "@/store/hooks";

export default function ReconstructionDetailsPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const token = useAppSelector((state) => state.auth.token);

  useEffect(() => {
    if (!token) router.replace("/login");
  }, [router, token]);

  return (
    <main className="mx-auto max-w-4xl p-6">
      <div className="mb-4 flex items-center justify-between gap-3">
        <Button variant="outline" onClick={() => router.push("/reconstructions")}>← Back</Button>
        <ThemeToggle />
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Reconstruction {params.id}</CardTitle>
          <CardDescription>Detailed reconstruction view is the next step.</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            Placeholder page. You can return to the <Link className="underline" href="/reconstructions">reconstructions list</Link>.
          </p>
        </CardContent>
      </Card>
    </main>
  );
}
