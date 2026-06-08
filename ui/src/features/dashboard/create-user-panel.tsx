import { useState } from 'react'
import type { FormEvent } from 'react'
import { Clipboard, UserPlus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import type { CreateUserInput, CreateUserResult } from '@/types'

type CreateUserPanelProps = {
  enabled: boolean
  isCreating: boolean
  error: Error | null
  onCreate: (input: CreateUserInput) => Promise<CreateUserResult>
}

export function CreateUserPanel({ enabled, isCreating, error, onCreate }: CreateUserPanelProps) {
  const [email, setEmail] = useState('')
  const [name, setName] = useState('')
  const [isAdmin, setIsAdmin] = useState(false)
  const [created, setCreated] = useState<CreateUserResult | null>(null)
  const canCreate = enabled && Boolean(email.trim()) && Boolean(name.trim()) && !isCreating

  async function submitUser(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const result = await onCreate({ email: email.trim(), name: name.trim(), is_admin: isAdmin })
    setCreated(result)
    setEmail('')
    setName('')
    setIsAdmin(false)
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Add User</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <form className="space-y-2" onSubmit={submitUser}>
          <Input type="email" value={email} onChange={(event) => setEmail(event.target.value)} placeholder="Email" required />
          <Input value={name} onChange={(event) => setName(event.target.value)} placeholder="Name" required />
          <label className="flex h-9 items-center gap-2 text-sm text-foreground">
            <input className="size-4 accent-primary" type="checkbox" checked={isAdmin} onChange={(event) => setIsAdmin(event.target.checked)} />
            Admin access
          </label>
          <Button className="w-full" type="submit" disabled={!canCreate}>
            <UserPlus className="size-4" />
            Create
          </Button>
        </form>

        {error && <p className="text-xs font-medium text-foreground">{error.message}</p>}
        {created && (
          <div className="space-y-2 rounded-md border border-border bg-muted p-3">
            <p className="text-xs font-medium text-foreground">{created.user.email}</p>
            <code className="block break-all text-xs text-muted-foreground">{created.api_key}</code>
            <Button className="w-full" type="button" variant="outline" onClick={() => void navigator.clipboard.writeText(created.api_key)}>
              <Clipboard className="size-4" />
              Copy key
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
