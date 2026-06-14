export function PagePatternBackground() {
  return (
    <div aria-hidden="true" className="pointer-events-none fixed inset-0 -z-10 hidden xl:block">
      <div className="page-pattern absolute inset-y-0 left-0 w-[calc((100vw-1180px)/2)] border-r border-border/60" />
      <div className="page-pattern absolute inset-y-0 right-0 w-[calc((100vw-1180px)/2)] border-l border-border/60" />
    </div>
  )
}
