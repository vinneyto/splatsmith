"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { PipelineSettingsForm } from "@/components/pipeline-settings-form";
import { Button, buttonVariants } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  type PipelineSettings,
  useGetStandardPipelineSettingsQuery,
  usePutStandardPipelineSettingsMutation,
} from "@/store/api/splatmakerApi";
import { useAppSelector } from "@/store/hooks";

export default function UserDefaultPipelineSettingsPage() {
  const router = useRouter();
  const token = useAppSelector((state) => state.auth.token);
  const { data, isLoading, isError, error } = useGetStandardPipelineSettingsQuery(undefined, { skip: !token });
  const [save, { isLoading: isSaving }] = usePutStandardPipelineSettingsMutation();
  const [form, setForm] = useState<PipelineSettings | null>(null);
  const [status, setStatus] = useState<string | null>(null);

  useEffect(() => {
    if (!token) {
      router.replace("/login");
    }
  }, [router, token]);

  useEffect(() => {
    if (data?.settings) {
      setForm(data.settings);
    }
  }, [data]);

  const canSave = useMemo(() => Boolean(form && token && !isSaving), [form, token, isSaving]);

  async function onSave() {
    if (!form) return;
    setStatus(null);
    try {
      await save(form).unwrap();
      setStatus("Settings saved");
    } catch {
      setStatus("Failed to save settings");
    }
  }

  return (
    <main className="mx-auto max-w-5xl p-6">
      <div className="mb-4 flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Default Pipeline Settings</h1>
        <Link href="/reconstructions" className={buttonVariants({ variant: "outline", size: "sm" })}>
          Back to reconstructions
        </Link>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>User defaults</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {isLoading ? <p className="text-sm text-muted-foreground">Loading settings...</p> : null}
          {isError ? <p className="text-sm text-red-600">Failed to load settings: {JSON.stringify(error)}</p> : null}

          {form ? <PipelineSettingsForm value={form} onChange={setForm} /> : null}

          <div className="flex items-center gap-3">
            <Button onClick={onSave} disabled={!canSave}>
              {isSaving ? "Saving..." : "Save settings"}
            </Button>
            {status ? <p className="text-sm text-muted-foreground">{status}</p> : null}
          </div>
        </CardContent>
      </Card>
    </main>
  );
}
