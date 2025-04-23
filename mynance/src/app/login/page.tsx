"use client";

import { useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import api from "@/lib/axios";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

export default function LoginPage() {
  const [form, setForm] = useState({ email: "", password: "" });
  const [errors, setErrors] = useState<{ email?: boolean; password?: boolean }>({});
  const router = useRouter();

  async function handleLogin(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setErrors({});

    if (!form.email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
      setErrors({ email: true });
      toast.error("Digite um e-mail válido.");
      return;
    }

    if (!form.password) {
      setErrors({ password: true });
      toast.error("Digite sua senha.");
      return;
    }

    try {
      const res = await api.post("/auth/login", form);
      if (res.status === 200) {
        toast.success("Login realizado com sucesso!");
        router.push("/dashboard");
      }
    } catch (err: any) {
      if (err.response?.data?.error) {
        toast.error(err.response.data.error);
      } else {
        toast.error("Erro ao conectar com o servidor.");
      }
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <Card className="w-full max-w-md p-6 shadow-lg">
        <CardContent className="space-y-4">
          <div className="flex justify-center mb-2">
            <Image
              src="/logo.png"
              alt="Logo"
              width={250}
              height={90}
              className="rounded-full"
            />
          </div>
          <h1 className="text-2xl font-bold text-center">Login</h1>
          <form onSubmit={handleLogin} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                value={form.email}
                onChange={(e) => setForm({ ...form, email: e.target.value })}
                placeholder="you@example.com"
                className={errors.email ? "border-red-500" : ""}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Senha</Label>
              <Input
                id="password"
                type="password"
                value={form.password}
                onChange={(e) => setForm({ ...form, password: e.target.value })}
                placeholder="••••••••"
                className={errors.password ? "border-red-500" : ""}
              />
            </div>
            <Button type="submit" className="w-full mt-4">
              Entrar
            </Button>
          </form>
          <p className="text-sm text-center mt-4">
            Não tem uma conta?{' '}
            <Link href="/register" className="text-blue-600 hover:underline">
              Registrar
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}