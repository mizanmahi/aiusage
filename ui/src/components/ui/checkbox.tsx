import { type InputHTMLAttributes } from 'react'
import { cn } from '@/lib/utils'

export function Checkbox({ className, ...props }: InputHTMLAttributes<HTMLInputElement>) {
  return <input className={cn('size-4 cursor-pointer rounded border border-input accent-primary disabled:cursor-not-allowed disabled:opacity-50', className)} type="checkbox" {...props} />
}
