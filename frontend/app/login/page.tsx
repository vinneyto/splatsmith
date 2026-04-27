"use client";

import { FormEvent, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
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
    if (token) {
      router.replace("/reconstructions");
    }
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
        }),
      );
      router.push("/reconstructions");
    } catch {
      setError("Login failed: invalid username/password or backend unavailable");
    }
  }

  return (
    <main style={{ minHeight: "100vh", display: "grid", placeItems: "center" }}>
      <form onSubmit={handleSubmit} style={{ width: 360, background: "white", padding: 24, borderRadius: 12 }}>
        <h1 style={{ marginTop: 0 }}>Splatmaker Login</h1>
        <label style={{ display: "block", marginBottom: 12 }}>
          Username
          <input
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            style={{ width: "100%", marginTop: 6, padding: 8 }}
          />
        </label>
        <label style={{ display: "block", marginBottom: 16 }}>
          Password
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            style={{ width: "100%", marginTop: 6, padding: 8 }}
          />
        </label>
        {error && <p style={{ color: "#b42318" }}>{error}</p>}
        <button type="submit" disabled={isLoading} style={{ width: "100%", padding: 10 }}>
          {isLoading ? "Signing in..." : "Sign in"}
        </button>
      </form>
    </main>
  );
}
