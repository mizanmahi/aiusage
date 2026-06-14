import { type InputHTMLAttributes } from 'react'
import { cn } from '@/lib/utils'

export function Checkbox({ className, ...props }: InputHTMLAttributes<HTMLInputElement>) {
  return <input className={cn('size-4 rounded border border-input accent-primary', className)} type="checkbox" {...props} />
}
