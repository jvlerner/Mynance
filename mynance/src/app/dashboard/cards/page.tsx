"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card as UiCard, CardContent } from "@/components/ui/card";
import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { Pencil, Trash2, RotateCcw } from "lucide-react";
import api from "@/lib/axios";
import { toast } from "sonner";
import axios from "axios";
import clsx from "clsx";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { DayGrid } from "@/components/ui/day-grid";
import Image from "next/image";

type CreditCard = {
    id: number;
    name: string;
    bank: string;
    limitAmount: number;
    dueDay: number;
    active: boolean;
};

export default function CardsPage() {
    const [cards, setCards] = useState<CreditCard[]>([]);
    const [inactiveCards, setInactiveCards] = useState<CreditCard[]>([]);
    const [banks, setBanks] = useState<string[]>([]);
    const [form, setForm] = useState({
        name: "",
        bank: "",
        limitAmount: "",
        dueDay: "",
    });
    const [errors, setErrors] = useState<{
        name?: boolean;
        bank?: boolean;
        limitAmount?: boolean;
        dueDay?: boolean;
    }>({});
    const [editId, setEditId] = useState<number | null>(null);

    const fetchCards = async () => {
        try {
            const res = await api.get("/credit-cards/all");
            const all = Array.isArray(res.data) ? res.data : [];

            setCards(all.filter(card => card.active));
            setInactiveCards(all.filter(card => !card.active));
        } catch (err) {
            if (axios.isAxiosError(err) && err.response?.status !== 401) {
                toast.error("Erro ao carregar cartões de crédito.");
            }
            setCards([]);
            setInactiveCards([]);
        }
    };

    const fetchBanks = async () => {
        try {
            const res = await api.get("/banks");
            const all = Array.isArray(res.data) ? res.data : [];
            setBanks(all);
        } catch (err) {
            if (axios.isAxiosError(err) && err.response?.status !== 401) {
                toast.error("Erro ao carregar bancos.");
            }
            setBanks([]);
        }
    }

    useEffect(() => {
        fetchCards();
        fetchBanks();
    }, []);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setForm(prev => ({
            ...prev,
            [e.target.name]: e.target.value,
        }));

        // Remove erro ao digitar
        setErrors(prev => ({
            ...prev,
            [e.target.name]: false,
        }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        const newErrors = {
            name: !form.name.trim(),
            bank: !form.bank.trim(),
            limitAmount: !form.limitAmount,
            dueDay: !form.dueDay,
        };
        setErrors(newErrors);

        if (newErrors.name) {
            toast.error("O nome do cartão é obrigatório.");
            return;
        }

        if (newErrors.bank) {
            toast.error("Selecione um banco.");
            return;
        }

        if (newErrors.limitAmount) {
            toast.error("Informe um limite válido maior que 0.");
            return;
        }

        if (newErrors.dueDay) {
            toast.error("Informe um dia de vencimento entre 1 e 31.");
            return;
        }

        try {
            if (editId) {
                await api.put("/credit-cards", {
                    id: editId,
                    name: form.name,
                    bank: form.bank,
                    limitAmount: parseFloat(form.limitAmount),
                    dueDay: parseInt(form.dueDay),
                });
                toast.success("Cartão atualizado com sucesso!");
            } else {
                await api.post("/credit-cards", {
                    name: form.name,
                    bank: form.bank,
                    limitAmount: parseFloat(form.limitAmount),
                    dueDay: parseInt(form.dueDay),
                });
                toast.success("Cartão cadastrado com sucesso!");
            }

            setForm({ name: "", bank: "", limitAmount: "", dueDay: "" });
            setEditId(null);
            setErrors({});
            await fetchCards();
        } catch (err: any) {
            toast.error(err.response?.data?.error || "Erro ao salvar cartão.");
        }
    };

    const handleEdit = (card: CreditCard) => {
        setEditId(card.id);
        setForm({
            name: card.name,
            bank: card.bank,
            limitAmount: String(card.limitAmount),
            dueDay: String(card.dueDay),
        });
        setErrors({});
    };

    const handleDelete = async (id: number) => {
        try {
            await api.request({ method: "DELETE", url: "/credit-cards", data: { id } });
            toast.success("Cartão desativado com sucesso!");
            await fetchCards();
        } catch (err: any) {
            toast.error(err.response?.data?.error || "Erro ao desativar cartão.");
        }
    };

    const handleActivate = async (id: number) => {
        try {
            await api.post("/credit-cards/activate", { id });
            toast.success("Cartão ativado com sucesso!");
            await fetchCards();
        } catch (err: any) {
            toast.error(err.response?.data?.error || "Erro ao ativar cartão.");
        }
    };

    return (
        <div className="space-y-6">
            <h1 className="text-2xl font-bold">Cartões de Crédito</h1>

            <form onSubmit={handleSubmit} className="flex flex-col gap-2">
                <div className="flex gap-2">
                    {/* Nome */}
                    <div className="space-y-1">
                        <Label htmlFor="name">Nome do cartão</Label>
                        <Input
                            id="name"
                            name="name"
                            value={form.name}
                            onChange={handleChange}
                            className={clsx("w-full", errors.name && "border-red-500")}
                        />
                    </div>

                    {/* Banco */}
                    <div className="space-y-1">
                        <Label htmlFor="bank">Banco</Label>
                        <Select
                            onValueChange={(value: string) => {
                                setForm((prev) => ({ ...prev, bank: value }));
                                setErrors((prev) => ({ ...prev, bank: false }));
                            }}
                            value={form.bank}
                        >
                            <SelectTrigger className={clsx(errors.bank && "border-red-500")}>
                                <SelectValue placeholder="Selecione um banco" />
                            </SelectTrigger>
                            <SelectContent>
                                {banks.map((bank) => (
                                    <SelectItem key={bank} value={bank}>
                                        {bank}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    {/* Limite */}
                    <div className="space-y-1">
                        <Label htmlFor="limitAmount">Limite</Label>
                        <Input
                            id="limitAmount"
                            name="limitAmount"
                            type="number"
                            step="0.01"
                            value={form.limitAmount}
                            onChange={handleChange}
                            className={clsx("w-full", errors.limitAmount && "border-red-500")}
                        />
                    </div>

                    {/* Dia de vencimento */}
                    <div className="space-y-1">
                        <Label htmlFor="dueDay">Dia de vencimento</Label>
                        <Popover>
                            <PopoverTrigger asChild>
                                <Button
                                    variant="outline"
                                    className={clsx(
                                        "w-full justify-start text-left font-normal",
                                        !form.dueDay && "text-muted-foreground",
                                        errors.dueDay && "border-red-500"
                                    )}
                                >
                                    {form.dueDay ? `Todo dia ${form.dueDay}` : "Selecione o dia"}
                                </Button>
                            </PopoverTrigger>
                            <PopoverContent className="w-auto p-0">
                                <DayGrid
                                    selected={Number(form.dueDay)}
                                    onSelect={(day) => {
                                        setForm((prev) => ({ ...prev, dueDay: String(day) }))
                                        setErrors((prev) => ({ ...prev, dueDay: false }))
                                    }}
                                />
                            </PopoverContent>
                        </Popover>
                    </div>
                </div>


                {/* Botões */}
                <div className="flex gap-2 items-end">
                    <Button type="submit">{editId ? "Atualizar" : "Cadastrar"}</Button>
                    {editId && (
                        <Button
                            type="button"
                            variant="outline"
                            onClick={() => {
                                setEditId(null);
                                setForm({ name: "", bank: "", limitAmount: "", dueDay: "" });
                                setErrors({});
                            }}
                        >
                            Cancelar
                        </Button>
                    )}
                </div>

            </form>


            {/* Cartões Ativos */}
            <div>
                <h2 className="text-lg font-semibold">Ativos</h2>
                <div className="grid grid-cols-1 sm:grid-cols-1 md:grid-cols-2 gap-3 mt-2">
                    {cards.map((card) => (
                        <UiCard key={card.id} className="flex items-center px-4 py-2 gap-1">
                            <div className="flex justify-between items-center w-full gap-4">
                                <CardContent className="p-0 flex gap-2">
                                    {/* Imagem do banco à esquerda */}
                                    <Image
                                        src={`/bancos/${card.bank.toLowerCase()}.png`}
                                        alt={card.bank}
                                        width={80}
                                        height={80}
                                        className="rounded-md object-contain"
                                    />
                                    {/* Informaçoes do cartao */}
                                    <div className="flex flex-col">
                                        <span className="font-medium">{card.name}</span>
                                        <span className="text-sm text-muted-foreground">{card.bank}</span>
                                        <span className="text-sm text-muted-foreground">
                                            Limite de R$ {card.limitAmount.toFixed(2)}
                                        </span>
                                        <span className="text-sm text-muted-foreground">
                                            Vencimento todo dia {card.dueDay}
                                        </span>
                                    </div>
                                </CardContent>
                                {/* Ações */}
                                <div className="flex flex-col gap-2">
                                    <Button size="icon" variant="ghost" onClick={() => handleEdit(card)}>
                                        <Pencil size={16} />
                                    </Button>
                                    <Button size="icon" variant="ghost" onClick={() => handleDelete(card.id)}>
                                        <Trash2 size={16} />
                                    </Button>
                                </div>
                            </div>


                        </UiCard>


                    ))}
                </div>
            </div>

            {/* Cartões Inativos */}
            {inactiveCards.length > 0 && (
                <div>
                    <h2 className="text-lg font-semibold mt-8">Inativos</h2>
                    <div className="grid grid-cols-1 sm:grid-cols-1 md:grid-cols-2 gap-3 mt-2">
                        {inactiveCards.map((card) => (
                            <UiCard
                                key={card.id}
                                className="flex items-center px-4 py-2 gap-1">
                                <div className="flex justify-between items-center w-full gap-4">
                                    <CardContent className="p-0 flex gap-2">
                                        <Image
                                            src={`/bancos/${card.bank.toLowerCase()}.png`}
                                            alt={card.bank}
                                            width={80}
                                            height={80}
                                            className="rounded-md object-contain"
                                        />
                                        <div className="flex flex-col">
                                            <span className="line-through text-muted-foreground">{card.name}</span>
                                            <span className="text-sm text-muted-foreground">{card.bank}</span>
                                            <span className="text-sm text-muted-foreground">
                                                Limite de R$ {card.limitAmount.toFixed(2)}
                                            </span>
                                            <span className="text-sm text-muted-foreground">
                                                Vencimento todo dia {card.dueDay}
                                            </span>
                                        </div>
                                    </CardContent>
                                    <div className="flex flex-col gap-2">
                                        <Button size="icon" variant="ghost" onClick={() => handleActivate(card.id)}>
                                            <RotateCcw size={16} />
                                        </Button>
                                    </div>
                                </div >
                            </UiCard>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
