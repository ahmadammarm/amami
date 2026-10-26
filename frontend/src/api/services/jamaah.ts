import { isNotFoundError, request } from '../client';
import type { Page } from '../../types/api';
import type { Jamaah, JamaahListParams, MustahikSurvey } from '../../types/jamaah';
import type { JamaahPayload, MustahikSurveyPayload } from '../../schemas/jamaah';

export const jamaahService = {
  getJamaahs(params: JamaahListParams): Promise<Page<Jamaah>> {
    return request<Page<Jamaah>>({ method: 'GET', url: '/jamaah', params });
  },

  getJamaahByID(id: string): Promise<Jamaah> {
    return request<Jamaah>({ method: 'GET', url: `/jamaah/${id}` });
  },

  createJamaah(payload: JamaahPayload): Promise<Jamaah> {
    return request<Jamaah>({ method: 'POST', url: '/jamaah', data: payload });
  },

  updateJamaah(id: string, payload: JamaahPayload): Promise<Jamaah> {
    return request<Jamaah>({ method: 'PUT', url: `/jamaah/${id}`, data: payload });
  },

  /** Returns null when the jamaah has not been surveyed yet. */
  async getSurvey(id: string): Promise<MustahikSurvey | null> {
    try {
      return await request<MustahikSurvey>({ method: 'GET', url: `/jamaah/${id}/mustahik-survey` });
    } catch (error) {
      if (isNotFoundError(error)) return null;
      throw error;
    }
  },

  submitSurvey(id: string, payload: MustahikSurveyPayload): Promise<MustahikSurvey> {
    return request<MustahikSurvey>({ method: 'PUT', url: `/jamaah/${id}/mustahik-survey`, data: payload });
  },
};
