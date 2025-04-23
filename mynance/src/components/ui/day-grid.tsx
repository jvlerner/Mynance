"use client"

import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

interface DayGridProps {
  selected?: number
  onSelect: (day: number) => void
  className?: string
}

export function DayGrid({ selected, onSelect, className }: DayGridProps) {
  return (
    <div className={cn("grid grid-cols-7 gap-2 p-3", className)}>
      {Array.from({ length: 31 }, (_, i) => {
        const day = i + 1
        const isSelected = selected === day
        return (
          <Button
            key={day}
            variant={isSelected ? "default" : "ghost"}
            size="icon"
            className={cn("h-9 w-9", isSelected && "bg-primary text-primary-foreground")}
            onClick={() => onSelect(day)}
          >
            {day}
          </Button>
        )
      })}
    </div>
  )
}
