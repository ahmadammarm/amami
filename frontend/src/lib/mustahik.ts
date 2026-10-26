import type { HouseStatus } from '../types/jamaah';
import type { BadgeTone } from '../types/ui';

export type ScoreLevel = 'high' | 'medium' | 'low';

/** Buckets the 0-100 poverty score for display; higher means stronger claim to aid. */
export function scoreLevel(score: number): ScoreLevel {
  if (score >= 70) return 'high';
  if (score >= 40) return 'medium';
  return 'low';
}

const scoreTones: Record<ScoreLevel, BadgeTone> = {
  high: 'bata',
  medium: 'kuningan',
  low: 'netral',
};

export const scoreTone = (score: number): BadgeTone => scoreTones[scoreLevel(score)];

export const houseStatusLabels: Record<HouseStatus, string> = {
  OWNED: 'Milik sendiri',
  FAMILY: 'Menumpang keluarga',
  RENTED: 'Sewa / kontrak',
  NONE: 'Tidak punya tempat tinggal',
};
