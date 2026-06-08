import { useEffect, useMemo, useState } from 'react'
import type { FormEvent } from 'react'
import { getDailySummary, getUserProjects, getUsers } from './api'
import './App.css'
import type { DailyPoint, ProjectSummary, UserSummary } from './types'
import { formatCost, formatDate, formatTokens } from './format'

const storedAPIKey = 'aiusage.admin.apiKey'

function App() {
  const [apiKey, setAPIKey] = useState(() => sessionStorage.getItem(storedAPIKey) ?? '')
  const [users, setUsers] = useState<UserSummary[]>([])
  const [projects, setProjects] = useState<ProjectSummary[]>([])
  const [daily, setDaily] = useState<DailyPoint[]>([])
  const [selectedUserID, setSelectedUserID] = useState('')
  const [from, setFrom] = useState('2026-01-01')
  const [to, setTo] = useState('2026-12-31')
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState('Enter admin API key')

  const selectedUser = users.find((user) => user.id === selectedUserID) ?? users[0]
  const totals = useMemo(() => summarize(users), [users])
  const maxProjectTokens = Math.max(1, ...projects.map((project) => project.total_tokens))
  const maxDailyTokens = Math.max(1, ...daily.map((point) => point.total_tokens))

  useEffect(() => {
    if (apiKey) {
      void refresh(apiKey, selectedUserID, from, to)
    }
  }, [])

  async function refresh(key = apiKey, userID = selectedUser?.id ?? selectedUserID, start = from, end = to) {
    if (!key) {
      setMessage('Enter admin API key')
      return
    }

    setLoading(true)
    setMessage('Loading')
    try {
      sessionStorage.setItem(storedAPIKey, key)
      const nextUsers = await getUsers({ apiKey: key })
      const nextUserID = userID || nextUsers[0]?.id || ''
      const [nextProjects, nextDaily] = await Promise.all([
        nextUserID ? getUserProjects(nextUserID, { apiKey: key }) : Promise.resolve([]),
        getDailySummary(start, end, { apiKey: key }),
      ])

      setUsers(nextUsers)
      setSelectedUserID(nextUserID)
      setProjects(nextProjects)
      setDaily(nextDaily)
      setMessage(`Updated ${new Date().toLocaleTimeString()}`)
    } catch (error) {
      setMessage(error instanceof Error ? error.message : 'Request failed')
    } finally {
      setLoading(false)
    }
  }

  function submitKey(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    void refresh(apiKey, selectedUserID, from, to)
  }

  function selectUser(userID: string) {
    setSelectedUserID(userID)
    void refresh(apiKey, userID, from, to)
  }

  return (
    <main className="dashboard">
      <header className="topbar">
        <div>
          <p className="eyebrow">aiusage</p>
          <h1>Admin Dashboard</h1>
        </div>
        <form className="auth-form" onSubmit={submitKey}>
          <label htmlFor="api-key">API key</label>
          <input
            id="api-key"
            type="password"
            value={apiKey}
            onChange={(event) => setAPIKey(event.target.value)}
            autoComplete="current-password"
          />
          <button type="submit" disabled={loading}>
            Refresh
          </button>
        </form>
      </header>

      <section className="metrics" aria-label="Usage totals">
        <Metric label="Users" value={users.length.toString()} />
        <Metric label="Tokens" value={formatTokens(totals.tokens)} />
        <Metric label="Cost" value={formatCost(totals.cost)} />
        <Metric label="Last active" value={formatDate(totals.lastActive)} />
      </section>

      <section className="workarea">
        <aside className="panel users-panel">
          <div className="panel-header">
            <h2>Users</h2>
            <span>{message}</span>
          </div>
          <div className="user-list">
            {users.map((user) => (
              <button
                key={user.id}
                type="button"
                className={user.id === selectedUser?.id ? 'user-row active' : 'user-row'}
                onClick={() => selectUser(user.id)}
              >
                <span>
                  <strong>{user.name || user.email}</strong>
                  <small>{user.email}</small>
                </span>
                <span className="row-stat">{formatTokens(user.total_tokens)}</span>
              </button>
            ))}
            {!users.length && <EmptyState label="No users loaded" />}
          </div>
        </aside>

        <section className="panel detail-panel">
          <div className="panel-header">
            <h2>{selectedUser ? selectedUser.name || selectedUser.email : 'Projects'}</h2>
            <span>{projects.length} projects</span>
          </div>
          <div className="project-list">
            {projects.map((project) => (
              <ProjectRow key={`${project.user_id}-${project.tool}-${project.project}`} project={project} max={maxProjectTokens} />
            ))}
            {!projects.length && <EmptyState label="No projects loaded" />}
          </div>
        </section>
      </section>

      <section className="panel trend-panel">
        <div className="panel-header trend-controls">
          <h2>Daily Summary</h2>
          <div className="date-controls">
            <input type="date" value={from} onChange={(event) => setFrom(event.target.value)} />
            <input type="date" value={to} onChange={(event) => setTo(event.target.value)} />
            <button type="button" onClick={() => refresh(apiKey, selectedUserID, from, to)} disabled={loading}>
              Apply
            </button>
          </div>
        </div>
        <div className="bars" aria-label="Daily token totals">
          {daily.map((point) => (
            <div className="bar-row" key={point.date}>
              <span>{point.date}</span>
              <div className="bar-track">
                <div className="bar-fill" style={{ width: `${(point.total_tokens / maxDailyTokens) * 100}%` }} />
              </div>
              <strong>{formatTokens(point.total_tokens)}</strong>
            </div>
          ))}
          {!daily.length && <EmptyState label="No daily data loaded" />}
        </div>
      </section>
    </main>
  )
}

function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div className="metric">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  )
}

function ProjectRow({ project, max }: { project: ProjectSummary; max: number }) {
  return (
    <div className="project-row">
      <div className="project-main">
        <strong>{project.project}</strong>
        <span>
          {project.tool} · {formatDate(project.last_active)}
        </span>
      </div>
      <div className="project-meter">
        <div style={{ width: `${(project.total_tokens / max) * 100}%` }} />
      </div>
      <div className="project-cost">
        <strong>{formatTokens(project.total_tokens)}</strong>
        <span>{formatCost(project.total_cost_usd)}</span>
      </div>
    </div>
  )
}

function EmptyState({ label }: { label: string }) {
  return <div className="empty-state">{label}</div>
}

function summarize(users: UserSummary[]) {
  return users.reduce(
    (total, user) => ({
      tokens: total.tokens + user.total_tokens,
      cost: total.cost + user.total_cost_usd,
      lastActive: maxDate(total.lastActive, user.last_active),
    }),
    { tokens: 0, cost: 0, lastActive: '' },
  )
}

function maxDate(left: string, right: string) {
  if (!left) return right
  if (!right) return left
  return left > right ? left : right
}

export default App
