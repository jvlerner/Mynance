import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function Home() {
  return (
    <div className="min-h-screen bg-background text-foreground">
      {/* Logo fixa no topo esquerdo */}
      <div className="px-10 py-4">
        <Image
          src="/logo.png"
          alt="Logo"
          width={200}
          height={80}
          className="rounded-full"
        />
      </div>

      {/* Hero section */}
      <main className="grid grid-cols-1 md:grid-cols-2 items-center px-10 py-20 gap-10">
        <div className="space-y-6">
          <h1 className="text-4xl font-bold leading-tight">
            Controle suas finanças de forma simples e inteligente
          </h1>
          <p className="text-muted-foreground">
            Com o MyNance, você acompanha seus gastos, receitas e cartões de crédito em um só lugar.
          </p>
          <div className="space-x-4">
            <Link href="/register">
              <Button>Comece agora</Button>
            </Link>
            <Link href="/login">
              <Button variant="outline">Já tenho conta</Button>
            </Link>
          </div>
        </div>
        <div className="flex justify-center">
          <Image
            src="/finance-illustration.png"
            alt="Ilustração finanças"
            width={400}
            height={400}
            className="rounded-lg"
          />
        </div>
      </main>
    </div>
  );
}