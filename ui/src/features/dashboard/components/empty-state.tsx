import { Inbox } from 'lucide-react'

export function EmptyState({ title, description }: { title: string; description: string }) {
  return (
    <div className="grid min-h-36 place-items-center rounded-md border border-dashed border-border bg-muted/35 p-6 text-center">
      <div className="flex max-w-sm flex-col items-center gap-2">
        <div className="grid size-10 place-items-center rounded-md border border-border bg-background text-muted-foreground">
          <Inbox />
        </div>
        <p className="text-sm font-medium text-foreground">{title}</p>
        <p className="text-xs text-muted-foreground">{description}</p>
      </div>
    </div>
  )
}
