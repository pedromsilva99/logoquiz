export type Logo = {
  id: number
  position: number
  solved: boolean
  name?: string // only present once solved
}

export type Level = {
  number: number
  title: string
  logos: Logo[]
}

export type GuessResult = {
  correct: boolean
  name?: string
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, init)
  if (!res.ok) {
    throw new Error(`${init?.method ?? 'GET'} ${path} falhou (${res.status})`)
  }
  return res.json() as Promise<T>
}

export function fetchLevel(number: number) {
  return request<Level>(`/api/levels/${number}`)
}

export function guessLogo(id: number, answer: string) {
  return request<GuessResult>(`/api/logos/${id}/guess`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ answer }),
  })
}

// The server decides which variant to send; the query param only busts the
// browser cache when a logo goes from hidden to solved.
export function logoImageUrl(logo: Logo) {
  return `/api/logos/${logo.id}/image?v=${logo.solved ? 'full' : 'hidden'}`
}
