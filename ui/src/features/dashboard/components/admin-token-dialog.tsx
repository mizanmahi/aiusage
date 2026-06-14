import { useState } from 'react'
import type { FormEvent } from 'react'
import { createPortal } from 'react-dom'
import { KeyRound, X, Zap } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

type AdminTokenDialogProps = {
  token: string
  isLoading: boolean
  isPowered: boolean
  onLoad: (token: string) => void
}

export function AdminTokenDialog({ token, isLoading, isPowered, onLoad }: AdminTokenDialogProps) {
  const [open, setOpen] = useState(false)
  const [draftToken, setDraftToken] = useState(token)

  function submitToken(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const nextToken = draftToken.trim()
    if (!nextToken) return
    onLoad(nextToken)
    setOpen(false)
  }

  function openDialog() {
    setDraftToken(token)
    setOpen(true)
  }

  return (
    <>
      <Button type="button" variant={isPowered ? 'outline' : 'default'} onClick={openDialog}>
        {isPowered ? <KeyRound data-icon="inline-start" /> : <Zap data-icon="inline-start" />}
        {isPowered ? 'Admin loaded' : 'Load admin token'}
      </Button>

      {open &&
        createPortal(
          <div className="power-dialog fixed inset-0 z-50 grid place-items-center bg-background/25 p-4 backdrop-blur-[2px]" onMouseDown={() => setOpen(false)}>
          <section
            aria-modal="true"
            className="w-full max-w-md rounded-lg border border-primary/20 bg-card/95 p-4 text-card-foreground shadow-lg"
            role="dialog"
            onMouseDown={(event) => event.stopPropagation()}
          >
            <div className="flex items-start justify-between gap-4">
              <div>
                <p className="text-base font-semibold text-foreground">Load admin token</p>
                <p className="mt-1 text-sm text-muted-foreground">Unlock live usage data for this dashboard session.</p>
              </div>
              <Button type="button" variant="ghost" size="icon" onClick={() => setOpen(false)} aria-label="Close">
                <X data-icon="inline-start" />
              </Button>
            </div>

            <form className="mt-4 flex flex-col gap-3" onSubmit={submitToken}>
              <Input
                type="password"
                value={draftToken}
                onChange={(event) => setDraftToken(event.target.value)}
                placeholder="Admin API token"
                autoComplete="current-password"
                autoFocus
              />
              <Button type="submit" disabled={!draftToken.trim() || isLoading}>
                <Zap data-icon="inline-start" />
                Power dashboard
              </Button>
            </form>
          </section>
          </div>,
          document.body,
        )}
    </>
  )
}
