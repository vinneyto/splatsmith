"use client";

import { ArrowLeftOutlined } from "@ant-design/icons";
import { Button, Card, Space, Typography } from "antd";
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
    <main style={{ maxWidth: 900, margin: "0 auto", padding: 24 }}>
      <Space direction="vertical" size={16} style={{ width: "100%" }}>
        <Button icon={<ArrowLeftOutlined />} href="/reconstructions">
          Back to list
        </Button>

        <Card>
          <Typography.Title level={3}>Reconstruction {params.id}</Typography.Title>
          <Typography.Paragraph type="secondary" style={{ marginBottom: 0 }}>
            Placeholder: detailed reconstruction view will be implemented next.
          </Typography.Paragraph>
        </Card>
      </Space>
    </main>
  );
}
