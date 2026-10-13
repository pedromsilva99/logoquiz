import { useEffect, useState, type FormEvent } from 'react'
import { useMutation } from '@tanstack/react-query'
import { guessLogo, logoImageUrl, type Logo } from './api'

type Props = {
  logo: Logo
  onSolved: (name: string) => void
  onClose: () => void
}

export function GuessPanel({ logo, onSolved, onClose }: Props) {
  const [answer, setAnswer] = useState('')
  // Incremented on every wrong answer; used as a key to replay the shake animation.
  const [wrongCount, setWrongCount] = useState(0)

  const guess = useMutation({
    mutationFn: (a: string) => guessLogo(logo.id, a),
    onSuccess: (res) => {
      if (res.correct && res.name) onSolved(res.name)
      else setWrongCount((n) => n + 1)
    },
  })

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && onClose()
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  function submit(e: FormEvent) {
    e.preventDefault()
    if (answer.trim()) guess.mutate(answer)
  }

  return (
    <div className="overlay" onClick={onClose}>
      <div className="panel" role="dialog" aria-modal="true" onClick={(e) => e.stopPropagation()}>
        <button className="close" onClick={onClose} aria-label="Fechar">
          ×
        </button>

        <img
          className={`panel-logo${logo.solved ? ' solved' : ''}`}
          src={logoImageUrl(logo)}
          alt={logo.solved ? logo.name : 'Logo por adivinhar'}
        />

        {logo.solved ? (
          <p className="answer-correct">
            ✓ <strong>{logo.name}</strong>
          </p>
        ) : (
          <form key={wrongCount} className={`guess${wrongCount > 0 ? ' wrong' : ''}`} onSubmit={submit}>
            <input
              autoFocus
              value={answer}
              onChange={(e) => setAnswer(e.target.value)}
              placeholder="Que marca é esta?"
              maxLength={100}
              aria-label="Resposta"
            />
            <button type="submit" disabled={guess.isPending || !answer.trim()}>
              Adivinhar
            </button>
          </form>
        )}

        {!logo.solved && wrongCount > 0 && <p className="answer-wrong">Não é essa… tenta outra vez!</p>}
        {guess.error && <p className="answer-wrong">{guess.error.message}</p>}
      </div>
    </div>
  )
}
