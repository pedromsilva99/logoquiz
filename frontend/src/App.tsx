import { useState } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { fetchLevel, logoImageUrl, type Level } from './api'
import { GuessPanel } from './GuessPanel'

const LEVEL = 1

export default function App() {
  const queryClient = useQueryClient()
  const { data: level, isPending, error } = useQuery({
    queryKey: ['level', LEVEL],
    queryFn: () => fetchLevel(LEVEL),
  })
  const [selectedId, setSelectedId] = useState<number | null>(null)

  if (isPending) return <main className="app status">A carregar…</main>
  if (error) return <main className="app status">Erro: {error.message}</main>

  const solvedCount = level.logos.filter((l) => l.solved).length
  const selected = level.logos.find((l) => l.id === selectedId)

  function markSolved(id: number, name: string) {
    queryClient.setQueryData<Level>(['level', LEVEL], (old) =>
      old && {
        ...old,
        logos: old.logos.map((l) => (l.id === id ? { ...l, solved: true, name } : l)),
      },
    )
  }

  return (
    <main className="app">
      <header className="header">
        <h1>
          Logo Quiz <span>Portugal</span>
        </h1>
        <p className="level-title">
          Nível {level.number} · {level.title}
        </p>
        <div className="progress" aria-label={`${solvedCount} de ${level.logos.length} acertados`}>
          <div className="progress-bar" style={{ width: `${(solvedCount / level.logos.length) * 100}%` }} />
        </div>
        <p className="progress-label">
          {solvedCount} / {level.logos.length}
        </p>
      </header>

      <ul className="grid">
        {level.logos.map((logo) => (
          <li key={logo.id}>
            <button
              className={`tile${logo.solved ? ' solved' : ''}`}
              onClick={() => setSelectedId(logo.id)}
            >
              <img src={logoImageUrl(logo)} alt={logo.solved ? logo.name : `Logo ${logo.position}`} />
            </button>
          </li>
        ))}
      </ul>

      {selected && (
        <GuessPanel
          key={selected.id}
          logo={selected}
          onSolved={(name) => markSolved(selected.id, name)}
          onClose={() => setSelectedId(null)}
        />
      )}
    </main>
  )
}
