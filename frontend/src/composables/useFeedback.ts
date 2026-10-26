import { computed, shallowRef } from 'vue';
import { getErrorMessage } from '../api/client';

export type FeedbackKind = 'success' | 'error' | 'confirm';

export interface FeedbackRequest {
  kind: FeedbackKind;
  title: string;
  message?: string;
  confirmLabel: string;
  cancelLabel?: string;
  /** Destructive confirmations get the bata (red) styling. */
  danger: boolean;
  resolve: (confirmed: boolean) => void;
}

export interface ConfirmOptions {
  title: string;
  message?: string;
  confirmLabel: string;
  cancelLabel?: string;
  danger?: boolean;
}

// Module-level queue so one <FeedbackDialog> in App.vue serves the whole app,
// and requests fired back-to-back are shown one after another.
const queue = shallowRef<FeedbackRequest[]>([]);
const current = computed(() => queue.value[0] ?? null);

function enqueue(request: Omit<FeedbackRequest, 'resolve'>): Promise<boolean> {
  return new Promise((resolve) => {
    queue.value = [...queue.value, { ...request, resolve }];
  });
}

function settle(confirmed: boolean): void {
  const [head, ...rest] = queue.value;
  if (!head) return;
  queue.value = rest;
  head.resolve(confirmed);
}

export const feedback = {
  /** Reports a completed action. Resolves when the dialog is closed. */
  async success(title: string, message?: string): Promise<void> {
    await enqueue({ kind: 'success', title, message, confirmLabel: 'Oke', danger: false });
  },

  /** Reports a failure; `cause` may be a message or any error thrown by the API client. */
  async error(title: string, cause?: unknown): Promise<void> {
    const message = cause === undefined ? undefined : typeof cause === 'string' ? cause : getErrorMessage(cause);
    await enqueue({ kind: 'error', title, message, confirmLabel: 'Tutup', danger: false });
  },

  /** Asks before an action; resolves true only when the user confirms. */
  confirm(options: ConfirmOptions): Promise<boolean> {
    return enqueue({
      kind: 'confirm',
      title: options.title,
      message: options.message,
      confirmLabel: options.confirmLabel,
      cancelLabel: options.cancelLabel ?? 'Batal',
      danger: options.danger ?? false,
    });
  },
};

export function useFeedback() {
  return feedback;
}

/** Internal: the state the FeedbackDialog renders. */
export function useFeedbackQueue() {
  return { current, settle };
}
