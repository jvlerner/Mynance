"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Menu,
  X,
  Home,
  Layers,
  CreditCard,
  List,
  BarChart,
  LogOut,
} from "lucide-react";
import Image from "next/image";
import clsx from "clsx";
import api, { setUnauthorizedHandler } from "@/lib/axios";
import { toast } from "sonner";

const navItems = [
  { label: "Início", href: "/dashboard", icon: Home },
  { label: "Categorias", href: "/dashboard/categories", icon: Layers },
  { label: "Cartões", href: "/dashboard/cards", icon: CreditCard },
  { label: "Despesas", href: "/dashboard/expenses", icon: List },
  { label: "Despesas Crédito", href: "/dashboard/expenses", icon: List },
  { label: "Relatórios", href: "/dashboard/reports", icon: BarChart },
];

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const [collapsed, setCollapsed] = useState(false);
  const router = useRouter();

  useEffect(() => {
    setUnauthorizedHandler(() => {
      toast.error("Sessão expirada. Faça login novamente.");
      router.push("/login");
    });
  }, [router]);

  const handleLogout = async () => {
    try {
      await api.post("/auth/logout");
      router.push("/login");
    } catch (error) {
      console.error("Erro ao fazer logout:", error);
    }
  };

  return (
    <div className="flex min-h-screen">
      {/* Sidebar */}
      <aside
        className={clsx(
          "shadow-lg transition-all duration-300 ease-in-out flex flex-col",
          "bg-[color:var(--sidebar)] text-[color:var(--sidebar-foreground)]",
          collapsed ? "w-20" : "w-64"
        )}
      >
        {/* Header da Sidebar */}
        <div className="flex items-center justify-center p-4 border-b">
          {!collapsed && (
            <Image
              src="/logo-m.png"
              alt="Logo"
              width={120}
              height={40}
              className="transition-all duration-300"
            />
          )}
          <Button
            size="icon"
            variant="ghost"
            onClick={() => setCollapsed(!collapsed)}
            className={collapsed ? undefined : "ml-auto"}
          >
            {collapsed ? <Menu size={20} /> : <X size={20} />}
          </Button>
        </div>

        {/* Navegação */}
        <nav className="flex flex-col gap-1 p-2 flex-1">
          {navItems.map(({ label, href, icon: Icon }) => (
            <Link
              key={href}
              href={href}
              className={clsx(
                "flex items-center rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-muted hover:text-foreground transition",
                collapsed ? "justify-center" : "gap-3"
              )}
            >
              <Icon size={20} />
              {!collapsed && label}
            </Link>
          ))}
        </nav>

        {/* Botão de Logout na parte inferior */}
        <div className="p-2 border-t">
          <Button
            variant="ghost"
            onClick={handleLogout}
            className={clsx(
              "w-full flex items-center px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-muted hover:text-foreground transition",
              collapsed ? "justify-center" : "gap-3"
            )}
          >
            <LogOut size={20} />
            {!collapsed && "Sair"}
          </Button>
        </div>
      </aside>

      {/* Conteúdo principal */}
      <main className="flex-1 h-screen overflow-y-auto p-10">
        {children}
      </main>
    </div>
  );
}
