"use client";

import { LockOutlined, UserOutlined } from "@ant-design/icons";
import { Alert, Button, Card, Form, Input, Space, Typography } from "antd";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useLoginMutation } from "@/store/api/splatmakerApi";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { setAuth } from "@/store/slices/authSlice";

type LoginFormValues = {
  username: string;
  password: string;
};

export default function LoginPage() {
  const router = useRouter();
  const dispatch = useAppDispatch();
  const token = useAppSelector((state) => state.auth.token);
  const [login, { isLoading }] = useLoginMutation();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (token) {
      router.replace("/reconstructions");
    }
  }, [router, token]);

  async function handleSubmit(values: LoginFormValues) {
    setError(null);

    try {
      const result = await login(values).unwrap();
      dispatch(
        setAuth({
          token: result.access_token,
          user: { userId: result.user.user_id, email: result.user.email },
        })
      );
      router.push("/reconstructions");
    } catch {
      setError("Login failed: invalid username/password or backend unavailable");
    }
  }

  return (
    <main style={{ minHeight: "100vh", display: "grid", placeItems: "center", padding: 16 }}>
      <Card style={{ width: 380 }}>
        <Space direction="vertical" size={16} style={{ width: "100%" }}>
          <Typography.Title level={3} style={{ margin: 0 }}>
            Splatmaker Login
          </Typography.Title>

          {error ? <Alert type="error" message={error} showIcon /> : null}

          <Form<LoginFormValues>
            layout="vertical"
            initialValues={{ username: "dev", password: "devpass" }}
            onFinish={handleSubmit}
          >
            <Form.Item label="Username" name="username" rules={[{ required: true }]}>
              <Input prefix={<UserOutlined />} autoComplete="username" />
            </Form.Item>

            <Form.Item label="Password" name="password" rules={[{ required: true }]}>
              <Input.Password prefix={<LockOutlined />} autoComplete="current-password" />
            </Form.Item>

            <Form.Item style={{ marginBottom: 0 }}>
              <Button type="primary" htmlType="submit" loading={isLoading} block>
                Sign in
              </Button>
            </Form.Item>
          </Form>
        </Space>
      </Card>
    </main>
  );
}
