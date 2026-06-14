import { Button } from '@/components/ui/button'

type TabValue = 'breakdown' | 'summary'

export function PanelTabs({ value, onChange }: { value: TabValue; onChange: (value: TabValue) => void }) {
  return (
    <div className="grid grid-cols-2 rounded-md border border-border bg-muted p-1">
      {(['breakdown', 'summary'] as TabValue[]).map((tab) => (
        <Button key={tab} type="button" variant={value === tab ? 'default' : 'ghost'} onClick={() => onChange(tab)}>
          {tab === 'breakdown' ? 'Breakdown' : 'Summary'}
        </Button>
      ))}
    </div>
  )
}
