"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Pencil, Trash2, RotateCcw } from "lucide-react";
import api from "@/lib/axios";
import { toast } from "sonner";
import axios from "axios";
import { RgbaColorPicker } from "react-colorful";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";


type Category = {
    id: number;
    name: string;
    color: string;
    active: boolean;
};

export default function CategoriasPage() {
    const [categories, setCategories] = useState<Category[]>([]);
    const [inactiveCategories, setInactiveCategories] = useState<Category[]>([]);
    const [name, setName] = useState("");
    const [hex, setHex] = useState("#ff5733");
    const [alpha, setAlpha] = useState(1);
    const [editId, setEditId] = useState<number | null>(null);
    const [popoverOpen, setPopoverOpen] = useState(false);


    const color = `${hex}${Math.round(alpha * 255)
        .toString(16)
        .padStart(2, "0")}`; // => #RRGGBBAA


    // Converte RGBA para #RRGGBBAA
    const rgbaToHex = ({ r, g, b, a }: { r: number; g: number; b: number; a: number }) => {
        const toHex = (v: number) => v.toString(16).padStart(2, "0");
        return `#${toHex(r)}${toHex(g)}${toHex(b)}${toHex(Math.round(a * 255))}`;
    };

    // Converte #RRGGBBAA para objeto RGBA
    const hexToRgba = (hex: string) => {
        const r = parseInt(hex.slice(1, 3), 16);
        const g = parseInt(hex.slice(3, 5), 16);
        const b = parseInt(hex.slice(5, 7), 16);
        const a = parseInt(hex.slice(7, 9) || "ff", 16) / 255;
        return { r, g, b, a };
    };

    const fetchCategories = async () => {
        try {
            const res = await api.get("/categories/all");

            const all = Array.isArray(res.data) ? res.data : [];

            setCategories(all.filter((cat) => cat.active));
            setInactiveCategories(all.filter((cat) => !cat.active));
        } catch (err) {
            if (axios.isAxiosError(err) && err.response?.status !== 401) {
                toast.error("Erro ao carregar categorias, falha ao comunicar com o servidor");
            }

            setCategories([]);
            setInactiveCategories([]);
        }
    };

    useEffect(() => {
        fetchCategories();
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!name.trim()) {
            toast.error("O nome da categoria é obrigatório.");
            return;
        }

        try {
            if (editId) {
                await api.put("/categories", { id: editId, name, color });
                toast.success("Categoria atualizada com sucesso!");
            } else {
                await api.post("/categories", { name, color });
                toast.success("Categoria cadastrada com sucesso!");
            }

            setName("");
            setHex("#ff5733");
            setAlpha(1);
            setEditId(null);
            await fetchCategories();
        } catch (err: any) {
            toast.error(err.response?.data?.error || "Erro ao salvar categoria.");
        }
    };

    const handleEdit = (category: Category) => {
        setName(category.name);
        setHex(category.color.slice(0, 7));
        const alphaHex = category.color.slice(7, 9);
        const alphaDecimal = parseInt(alphaHex || "ff", 16) / 255;
        setAlpha(alphaDecimal);
        setEditId(category.id);
    };

    const handleDelete = async (id: number) => {
        try {
            await api.request({ method: "DELETE", url: "/categories", data: { id } });
            toast.success("Categoria deletada com sucesso!");
            await fetchCategories();
        } catch (err: any) {
            toast.error(err.response?.data?.error || "Erro ao deletar categoria.");
        }
    };

    const handleActivate = async (id: number) => {
        try {
            await api.post("/categories/activate", { id });
            toast.success("Categoria ativada com sucesso!");
            await fetchCategories();
        } catch (err: any) {
            toast.error(err.response?.data?.error || "Erro ao ativar categoria.");
        }
    };

    return (
        <div className="space-y-6">
            <h1 className="text-2xl font-bold">Categorias</h1>

            <form onSubmit={handleSubmit} className="flex flex-col gap-2">
                <div className="flex flex-wrap gap-2">
                    <div className="space-y-1">
                        <Label htmlFor="name">Nome da categoria</Label>
                        <Input
                            id="name"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="Ex: Alimentação"
                            className="w-full"
                        />
                    </div>

                    <div className="space-y-1">
                        <Label>Cor da categoria</Label>

                        <Popover  open={popoverOpen} onOpenChange={setPopoverOpen}>
                            <PopoverTrigger className="w-full" asChild>
                                <Button 
                                    type="button"
                                    variant="outline"
                                    className="w-full justify-start"
                                >
                                    <div
                                        className="w-4 h-4 rounded-full mr-2 border"
                                        style={{ backgroundColor: color }}
                                    />
                                    {color}
                                </Button>
                            </PopoverTrigger>

                            <PopoverContent className="w-auto p-4">
                                <RgbaColorPicker
                                    color={hexToRgba(color)}
                                    onChange={(newColor) => {
                                        setHex(rgbaToHex(newColor).slice(0, 7));
                                        setAlpha(newColor.a);
                                    }}
                                />
                            </PopoverContent>
                        </Popover>
                    </div>
                </div>

                <div className="flex gap-2">
                    <Button type="submit">{editId ? "Atualizar" : "Cadastrar"}</Button>
                    {editId && (
                        <Button
                            type="button"
                            variant="outline"
                            onClick={() => {
                                setEditId(null);
                                setName("");
                                setHex("#ff5733");
                                setAlpha(1);
                            }}
                        >
                            Cancelar
                        </Button>
                    )}
                </div>
            </form>

            <div>
                <h2 className="text-lg font-semibold">Ativas</h2>
                <div className="grid grid-cols-1 sm:grid-cols-1 md:grid-cols-2 gap-3 mt-2">
                    {categories.map((cat) => (
                        <Card key={cat.id} className="flex items-center px-4 py-2 gap-1">
                            <div className="flex justify-between items-center w-full gap-4">
                                <CardContent className="p-0 flex gap-2 items-center">
                                    <span
                                        className="w-4 h-4 rounded-full border"
                                        style={{ backgroundColor: cat.color }}
                                    />
                                    <span>{cat.name}</span>
                                </CardContent>
                                <div className="flex gap-2">
                                    <Button size="icon" variant="ghost" onClick={() => handleEdit(cat)}>
                                        <Pencil size={16} />
                                    </Button>
                                    <Button size="icon" variant="ghost" onClick={() => handleDelete(cat.id)}>
                                        <Trash2 size={16} />
                                    </Button>
                                </div>
                            </div>
                        </Card>
                    ))}
                </div>
            </div>

            {inactiveCategories.length > 0 && (
                <div>
                    <h2 className="text-lg font-semibold mt-8">Inativas</h2>
                    <div className="grid grid-cols-1 sm:grid-cols-1 md:grid-cols-2 gap-3 mt-2">
                        {inactiveCategories.map((cat) => (
                            <Card key={cat.id} className="flex items-center px-4 py-2 gap-1">
                                <div className="flex justify-between items-center w-full gap-4">
                                    <CardContent className="p-0 flex gap-2 items-center">
                                        <span
                                            className="w-4 h-4 rounded-full border"
                                            style={{ backgroundColor: cat.color }}
                                        />
                                        <span className="line-through text-muted-foreground">{cat.name}</span>
                                    </CardContent>
                                    <div className="flex flex-col gap-2">
                                        <Button size="icon" variant="ghost" onClick={() => handleActivate(cat.id)}>
                                            <RotateCcw size={16} />
                                        </Button>
                                    </div>
                                </div>
                            </Card>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
