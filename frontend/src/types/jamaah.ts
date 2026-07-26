export interface Jamaah {
  id: string;
  user_id?: string;
  full_name: string;
  phone: string;
  address: string;
  is_mustahik: boolean;
  mustahik_score: number;
  created_at: string;
  updated_at: string;
}

export interface CreateJamaahPayload {
  full_name: string;
  phone: string;
  address: string;
  is_mustahik: boolean;
}

export interface UpdateJamaahPayload {
  full_name: string;
  phone: string;
  address: string;
  is_mustahik: boolean;
}

export interface PaginatedJamaahResponse {
  data: Jamaah[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}
