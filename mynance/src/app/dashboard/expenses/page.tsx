"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card as UiCard, CardContent } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from "@/components/ui/select";
import { CalendarIcon, Pencil, Trash2 } from "lucide-react";
import { Popover, PopoverTrigger, PopoverContent } from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { format } from "date-fns";
import api from "@/lib/axios";
import { toast } from "sonner";
import axios from "axios";

interface Expense {
  id: number;
  description: string;
  amount: number;
  due_date: string;
  paid: boolean;
  category_id: number | null;
  category_name?: string;
}

interface Category {
  id: number;
  name: string;
}

export default function ExpensesPage() {
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [description, setDescription] = useState("");
  const [amount, setAmount] = useState("");
  const [dueDate, setDueDate] = useState("");
  const [categoryId, setCategoryId] = useState<string | undefined>(undefined);
  const [editId, setEditId] = useState<number | null>(null);


  const fetchExpenses = async () => {
    try {
      const res = await api.get("/expenses/all");

      const all = Array.isArray(res.data) ? res.data : [];

      setExpenses(all.filter((cat) => cat.active));
    } catch (err) {
      if (axios.isAxiosError(err) && err.response?.status !== 401) {
        toast.error("Erro ao carregar despesas, falha ao comunicar com o servidor");
      }

      setExpenses([]);
    }
  };


  const fetchCategories = async () => {
    try {
      const res = await api.get("/categories/all");

      const all = Array.isArray(res.data) ? res.data : [];

      setCategories(all.filter((cat) => cat.active));
    } catch (err) {
      if (axios.isAxiosError(err) && err.response?.status !== 401) {
        toast.error("Erro ao carregar categorias, falha ao comunicar com o servidor");
      }

      setCategories([]);
    }
  };

  useEffect(() => {
    fetchExpenses();
    fetchCategories();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();


    if (!description || !amount || !dueDate) {
      toast.error("Preencha todos os campos obrigatórios.");
      return;
    }
    // Validar todos campos e deixar vermelho  

    
    try {
      if (editId) {
        await api.put("/expenses", {
          id: editId,
          description,
          amount: parseFloat(amount),
          due_date: dueDate,
          category_id: categoryId && categoryId !== "none" ? parseInt(categoryId) : null,
        });
        toast.success("Despesa atualizada com sucesso");
      } else {
        await api.post("/expenses", {
          description,
          amount: parseFloat(amount),
          due_date: dueDate,
          category_id: categoryId && categoryId !== "none" ? parseInt(categoryId) : null,
        });
        toast.success("Despesa cadastrada com sucesso");
      }
      setDescription("");
      setAmount("");
      setDueDate("");
      setCategoryId(undefined);
      setEditId(null);
      fetchExpenses();
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Erro ao salvar despesa");
    }
  };

  const handleEdit = (expense: Expense) => {
    setEditId(expense.id);
    setDescription(expense.description);
    setAmount(String(expense.amount));
    setDueDate(expense.due_date);
    setCategoryId(expense.category_id ? String(expense.category_id) : undefined);
  };

  const handleDelete = async (id: number) => {
    try {
      await api.request({ method: "DELETE", url: "/expenses", data: { id } });
      toast.success("Despesa excluída com sucesso!");
      fetchExpenses();
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Erro ao excluir despesa.");
    }
  };


  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Despesas</h1>

      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="space-y-1">
            <Label>Descrição</Label>
            <Input value={description} onChange={(e) => setDescription(e.target.value)} />
          </div>
          <div className="space-y-1">
            <Label>Valor</Label>
            <Input type="number" value={amount} onChange={(e) => setAmount(e.target.value)} />
          </div>
          <div className="space-y-1">
            <Label>Data de vencimento</Label>
            <Popover>
              <PopoverTrigger asChild>
                <Button variant="outline" className="w-full justify-start text-left font-normal">
                  {dueDate ? format(new Date(dueDate), "dd/MM/yyyy") : <span className="text-muted-foreground">Selecione a data</span>}
                  <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0">
                <Calendar
                  mode="single"
                  selected={dueDate ? new Date(dueDate) : undefined}
                  onSelect={(date) => date && setDueDate(date.toISOString().split("T")[0])}
                  initialFocus
                />
              </PopoverContent>
            </Popover>
          </div>
          <div className="space-y-1">
            <Label>Categoria</Label>
            <Select value={categoryId} onValueChange={(val) => setCategoryId(val || undefined)}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Selecione uma categoria" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">Sem categoria</SelectItem>
                {categories.map((cat) => (
                  <SelectItem key={cat.id} value={String(cat.id)}>
                    {cat.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
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
                setDescription("");
                setAmount("");
                setDueDate("");
                setCategoryId(undefined);
              }}
            >
              Cancelar
            </Button>
          )}
        </div>
      </form>

      <div>
        <h2 className="text-lg font-semibold mt-4">Lista de Despesas</h2>
        <div className="grid grid-cols-1 sm:grid-cols-1 md:grid-cols-2 gap-3 mt-2">
          {expenses.map((expense) => (
            <UiCard key={expense.id} className="flex items-center px-4 py-2 gap-1">
              <div className="flex justify-between items-center w-full gap-4">
                <CardContent className="p-0">
                  <div className="font-medium">{expense.description}</div>
                  <div className="text-sm text-muted-foreground">
                    R$ {Number(expense.amount).toFixed(2)} - Vence em {new Date(expense.due_date).toLocaleDateString()}
                    {expense.category_name && ` - Categoria: ${expense.category_name}`}
                  </div>
                </CardContent>
                <div className="flex gap-2">
                  <Button size="icon" variant="ghost" onClick={() => handleEdit(expense)}>
                    <Pencil size={16} />
                  </Button>
                  <Button size="icon" variant="ghost" onClick={() => handleDelete(expense.id)}>
                    <Trash2 size={16} />
                  </Button>
                </div>
              </div>
            </UiCard>
          ))}
        </div>
      </div>
    </div>
  );
}