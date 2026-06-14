import { Monitor, Moon, Sun } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useTheme } from '@/components/theme'

export function ThemeToggle() {
  const { theme, setTheme } = useTheme()

  const options = [
    { value: 'light' as const, label: 'Light', icon: Sun },
    { value: 'dark' as const, label: 'Dark', icon: Moon },
    { value: 'system' as const, label: 'System', icon: Monitor },
  ]

  return (
    <div className="inline-flex rounded-md border border-border bg-card p-1" aria-label="Theme">
      {options.map((option) => {
        const Icon = option.icon
        return (
          <Button
            key={option.value}
            type="button"
            variant={theme === option.value ? 'default' : 'ghost'}
            size="icon"
            className="size-8"
            onClick={() => setTheme(option.value)}
            title={option.label}
            aria-label={option.label}
          >
            <Icon data-icon="inline-start" />
          </Button>
        )
      })}
    </div>
  )
}
