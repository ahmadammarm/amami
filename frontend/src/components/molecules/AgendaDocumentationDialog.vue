<script setup lang="ts">
import { computed, ref } from 'vue';
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import { FileText, Trash2, UploadCloud, Images } from 'lucide-vue-next';
import { feedback } from '../../composables/useFeedback';
import { agendaService } from '../../api/services/agenda';
import { assetUrl } from '../../api/client';
import { queryKeys } from '../../api/queryKeys';
import { formatDate } from '../../lib/format';
import {
  DOCUMENTATION_ACCEPT,
  DOCUMENTATION_MAX_BYTES,
  DOCUMENTATION_MAX_FILES,
} from '../../schemas/agenda';
import type { Agenda, AgendaDocumentation } from '../../types/agenda';
import { Dialog, DialogContent, DialogDescription, DialogTitle } from '../atoms/dialog';
import EmptyState from './EmptyState.vue';

const props = defineProps<{
  agenda: Agenda;
  canManage: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const MAX_MB = DOCUMENTATION_MAX_BYTES / 1024 / 1024;

const queryClient = useQueryClient();
const fileInput = ref<HTMLInputElement | null>(null);
const docsKey = computed(() => queryKeys.agenda.documentation(props.agenda.id));
const allowedTypes = DOCUMENTATION_ACCEPT.split(',');

const { data: documents, isLoading } = useQuery({
  queryKey: docsKey,
  queryFn: () => agendaService.getDocumentation(props.agenda.id),
});

const uploadMutation = useMutation({
  mutationFn: (files: File[]) => agendaService.uploadDocumentation(props.agenda.id, files),
  onSuccess: (uploaded) => {
    queryClient.invalidateQueries({ queryKey: docsKey.value });
    feedback.success('Dokumentasi diunggah', `${uploaded.length} file ditambahkan ke ${props.agenda.title}.`);
  },
  onError: (error: unknown) => feedback.error('Dokumentasi belum diunggah', error),
});

const deleteMutation = useMutation({
  mutationFn: (doc: AgendaDocumentation) => agendaService.deleteDocumentation(props.agenda.id, doc.id),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: docsKey.value });
    feedback.success('Dokumentasi dihapus');
  },
  onError: (error: unknown) => feedback.error('Dokumentasi belum dihapus', error),
});

/** Client-side pre-check mirroring the server rules, for faster feedback. */
function validateFiles(files: File[]): string | null {
  if (files.length > DOCUMENTATION_MAX_FILES) {
    return `Pilih paling banyak ${DOCUMENTATION_MAX_FILES} file sekaligus.`;
  }
  const invalidType = files.find((file) => !allowedTypes.includes(file.type));
  if (invalidType) return `${invalidType.name} bukan JPG, PNG, WEBP, atau PDF.`;
  const tooLarge = files.find((file) => file.size > DOCUMENTATION_MAX_BYTES);
  if (tooLarge) return `${tooLarge.name} lebih dari ${MAX_MB} MB.`;
  return null;
}

const onFilesSelected = (event: Event) => {
  const input = event.target as HTMLInputElement;
  const files = Array.from(input.files ?? []);
  input.value = '';
  if (files.length === 0) return;

  const problem = validateFiles(files);
  if (problem) {
    feedback.error('File tidak dapat diunggah', problem);
    return;
  }
  uploadMutation.mutate(files);
};

const confirmDelete = async (doc: AgendaDocumentation) => {
  const confirmed = await feedback.confirm({
    title: 'Hapus dokumentasi?',
    message: `${doc.file_name} akan dihapus permanen dari arsip kegiatan.`,
    confirmLabel: 'Hapus file',
    danger: true,
  });
  if (confirmed) deleteMutation.mutate(doc);
};

const isImage = (doc: AgendaDocumentation) => doc.file_type.startsWith('image/');
const formatSize = (bytes: number) => `${(bytes / 1024 / 1024).toFixed(1)} MB`;

const onOpenChange = (open: boolean) => {
  if (!open) emit('close');
};
</script>

<template>
  <Dialog :open="true" @update:open="onOpenChange">
    <DialogContent class="max-h-[92dvh] gap-0 overflow-hidden rounded-2xl border-batu p-0 sm:max-w-3xl">
      <div class="border-b border-batu px-5 pt-5 pb-4 pr-12 sm:px-6 sm:pt-6">
        <DialogTitle class="font-display text-xl font-semibold tracking-wide">Dokumentasi kegiatan</DialogTitle>
        <DialogDescription class="mt-1 text-sm text-jati-muda">{{ agenda.title }} · {{ agenda.day }}, {{ agenda.time }}</DialogDescription>
      </div>

      <div class="scrollbar-tipis max-h-[75dvh] space-y-5 overflow-y-auto overscroll-contain px-5 py-5 sm:px-6">
        <div v-if="canManage">
          <input ref="fileInput" type="file" multiple :accept="DOCUMENTATION_ACCEPT" class="sr-only" @change="onFilesSelected" />
          <button
            type="button"
            :disabled="uploadMutation.isPending.value"
            class="flex w-full flex-col items-center gap-2 rounded-xl border-2 border-dashed border-batu-tua bg-kapur/50 px-4 py-6 text-sm text-jati-muda transition-colors hover:border-sajadah hover:text-sajadah disabled:opacity-50"
            @click="fileInput?.click()"
          >
            <UploadCloud class="size-6" />
            <span class="font-medium">{{ uploadMutation.isPending.value ? 'Mengunggah…' : 'Pilih foto atau PDF' }}</span>
            <span class="text-xs">Maksimal {{ DOCUMENTATION_MAX_FILES }} file, {{ MAX_MB }} MB per file.</span>
          </button>
        </div>

        <p v-if="isLoading" class="py-6 text-center text-sm text-jati-muda">Memuat dokumentasi…</p>
        <EmptyState
          v-else-if="!documents || documents.length === 0"
          :icon="Images"
          title="Belum ada dokumentasi"
          :description="canManage ? 'Unggah foto atau laporan kegiatan untuk arsip pertanggungjawaban.' : 'Foto dan laporan kegiatan akan tampil di sini.'"
        />
        <ul v-else class="grid grid-cols-2 gap-4 sm:grid-cols-3">
          <li v-for="doc in documents" :key="doc.id" class="group relative overflow-hidden rounded-xl border border-batu bg-white">
            <a :href="assetUrl(doc.file_url)" target="_blank" rel="noopener" class="block aspect-[4/3] bg-kapur">
              <img v-if="isImage(doc)" :src="assetUrl(doc.file_url)" :alt="doc.file_name" class="h-full w-full object-cover" loading="lazy" />
              <span v-else class="flex h-full w-full flex-col items-center justify-center gap-1 text-jati-muda">
                <FileText class="size-8" />
                <span class="text-xs font-semibold">PDF</span>
              </span>
            </a>
            <div class="border-t border-batu px-3 py-2">
              <p class="truncate text-xs font-medium text-jati">{{ doc.file_name }}</p>
              <p class="text-[11px] text-jati-muda">{{ formatSize(doc.file_size) }} · {{ formatDate(doc.uploaded_at) }}</p>
            </div>
            <button
              v-if="canManage"
              type="button"
              :disabled="deleteMutation.isPending.value"
              class="absolute top-2 right-2 rounded-lg bg-white/95 p-2 text-bata shadow-sm transition-opacity focus-visible:opacity-100 [@media(hover:hover)]:opacity-0 [@media(hover:hover)]:group-hover:opacity-100"
              :aria-label="`Hapus ${doc.file_name}`"
              @click="confirmDelete(doc)"
            >
              <Trash2 class="size-4" />
            </button>
          </li>
        </ul>
      </div>
    </DialogContent>
  </Dialog>
</template>
