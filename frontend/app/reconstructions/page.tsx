"use client";

import { LogoutOutlined } from "@ant-design/icons";
import { Alert, Button, Card, List, Space, Tag, Typography } from "antd";
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
    <main style={{ maxWidth: 900, margin: "0 auto", padding: 24 }}>
      <Space direction="vertical" size={16} style={{ width: "100%" }}>
        <Space style={{ width: "100%", justifyContent: "space-between" }} align="start">
          <div>
            <Typography.Title level={2} style={{ marginBottom: 0 }}>
              Reconstructions
            </Typography.Title>
            {user ? (
              <Typography.Text type="secondary">Signed in as {user.email}</Typography.Text>
            ) : null}
          </div>
          <Button
            icon={<LogoutOutlined />}
            onClick={() => {
              dispatch(logout());
              router.push("/login");
            }}
          >
            Logout
          </Button>
        </Space>

        {isError ? (
          <Alert type="error" showIcon message={`Failed to load reconstructions: ${JSON.stringify(error)}`} />
        ) : null}

        <Card>
          <List
            loading={isLoading}
            locale={{ emptyText: "No reconstructions yet." }}
            dataSource={data?.items ?? []}
            renderItem={(job) => (
              <List.Item>
                <List.Item.Meta
                  title={<Link href={`/reconstructions/${job.job_id}`}>{job.job_id}</Link>}
                  description={`Updated: ${new Date(job.updated_at).toLocaleString()}`}
                />
                <Tag>{job.status}</Tag>
              </List.Item>
            )}
          />
        </Card>
      </Space>
    </main>
  );
}
