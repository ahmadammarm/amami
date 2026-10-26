import type { HouseStatus } from '../types/jamaah';

export type ScoreLevel = 'high' | 'medium' | 'low';

/** Buckets the 0-100 poverty score for display; higher means stronger claim to aid. */
export function scoreLevel(score: number): ScoreLevel {
  if (score >= 70) return 'high';
  if (score >= 40) return 'medium';
  return 'low';
}

export const scoreBadgeClass: Record<ScoreLevel, string> = {
  high: 'bg-red-100 text-red-700',
  medium: 'bg-amber-100 text-amber-700',
  low: 'bg-gray-100 text-gray-600',
};

export const houseStatusLabels: Record<HouseStatus, string> = {
  OWNED: 'Milik sendiri',
  FAMILY: 'Menumpang keluarga',
  RENTED: 'Sewa / kontrak',
  NONE: 'Tidak punya tempat tinggal',
};
