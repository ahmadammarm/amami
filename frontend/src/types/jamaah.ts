import type { ISODateString } from './api';

export interface Jamaah {
  id: string;
  user_id: string | null;
  full_name: string;
  phone: string;
  address: string;
  is_mustahik: boolean;
  mustahik_score: number;
  created_at: ISODateString;
  updated_at: ISODateString;
}

export type HouseStatus = 'OWNED' | 'FAMILY' | 'RENTED' | 'NONE';

export type JamaahSort = 'newest' | 'score_desc' | 'name_asc';

export interface JamaahListParams {
  page: number;
  limit: number;
  query?: string;
  mustahik?: boolean;
  min_score?: number;
  sort?: JamaahSort;
}

export interface MustahikSurvey {
  jamaah_id: string;
  total_income: number;
  dependents_count: number;
  house_status: HouseStatus;
  has_disability: boolean;
  criteria_notes: string;
  mustahik_score: number;
  surveyed_by: string | null;
  updated_at: ISODateString;
}
