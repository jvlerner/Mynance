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

export default function RegisterPage() {
  const [form, setForm] = useState({ name: "", email: "", password: "", confirmPassword: "" });
  const [errors, setErrors] = useState<{ name?: boolean; email?: boolean; password?: boolean; confirmPassword?: boolean }>({});
  const router = useRouter();

  function validateEmail(email: string) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  function validatePassword(password: string) {
    const strongRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*(),.?":{}|<>]).{8,}$/;
    return strongRegex.test(password);
  }

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    const newErrors: typeof errors = {};
    if (form.name.length < 4) {
      newErrors.name = true;
      toast.error("Nome deve ter pelo menos 4 caracteres.");
      setErrors(newErrors);
      return;
    }

    if (!validateEmail(form.email)) {
      newErrors.email = true;
      toast.error("Digite um e-mail válido.");
      setErrors(newErrors);
      return;
    }

    if (!validatePassword(form.password)) {
      newErrors.password = true;
      toast.error("Digite uma senha válida com pelo menos 8 caracteres, uma letra maiúscula, uma minúscula, um número e um caractere especial.");
      setErrors(newErrors);
      return;
    }

    if (form.password !== form.confirmPassword) {
      newErrors.confirmPassword = true;
      toast.error("As senhas devem coincidir.");
      setErrors(newErrors);
      return;
    }

    try {
      const res = await api.post("/auth/register", {
        name: form.name,
        email: form.email,
        password: form.password,
      });
      if (res.status === 200) {
        toast.success("Cadastro realizado com sucesso! Realize o login para continuar.");
        router.push("/login");
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
          <h1 className="text-2xl font-bold text-center">Registrar</h1>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Nome</Label>
              <Input
                id="name"
                type="text"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder="Seu nome"
                className={errors.name ? "border-red-500" : ""}
              />
            </div>
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
            <div className="space-y-2">
              <Label htmlFor="confirmPassword">Confirmar senha</Label>
              <Input
                id="confirmPassword"
                type="password"
                value={form.confirmPassword}
                onChange={(e) => setForm({ ...form, confirmPassword: e.target.value })}
                placeholder="••••••••"
                className={errors.confirmPassword ? "border-red-500" : ""}
              />
            </div>
            <Button type="submit" className="w-full mt-4">
              Criar conta
            </Button>
          </form>
          <p className="text-sm text-center mt-4">
            Já tem uma conta?{' '}
            <Link href="/login" className="text-blue-600 hover:underline">
              Entrar
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}