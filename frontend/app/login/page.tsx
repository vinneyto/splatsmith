"use client";

import { FormEvent, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { ThemeToggle } from "@/components/theme-toggle";
import { useLoginMutation } from "@/store/api/splatmakerApi";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { setAuth } from "@/store/slices/authSlice";

export default function LoginPage() {
  const router = useRouter();
  const dispatch = useAppDispatch();
  const token = useAppSelector((state) => state.auth.token);
  const [login, { isLoading }] = useLoginMutation();

  const [username, setUsername] = useState("dev");
  const [password, setPassword] = useState("devpass");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (token) router.replace("/reconstructions");
  }, [router, token]);

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setError(null);

    try {
      const result = await login({ username, password }).unwrap();
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
    <main className="relative grid min-h-screen place-items-center p-6">
      <div className="absolute right-6 top-6">
        <ThemeToggle />
      </div>
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Splatmaker Login</CardTitle>
          <CardDescription>Sign in to manage reconstructions.</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-4" onSubmit={handleSubmit}>
            <div className="space-y-2">
              <label className="text-sm font-medium" htmlFor="username">
                Username
              </label>
              <Input
                id="username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                autoComplete="username"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium" htmlFor="password">
                Password
              </label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                autoComplete="current-password"
              />
            </div>
            {error && <p className="text-sm text-rose-600 dark:text-rose-300">{error}</p>}
            <Button className="w-full" disabled={isLoading} type="submit">
              {isLoading ? "Signing in..." : "Sign in"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </main>
  );
}
