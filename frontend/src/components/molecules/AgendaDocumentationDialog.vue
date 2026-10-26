<script setup lang="ts">
import { computed, ref } from 'vue';
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import { toast } from 'vue-sonner';
import { FileText, Trash2, UploadCloud, XIcon } from 'lucide-vue-next';
import { agendaService } from '../../api/services/agenda';
import { assetUrl, getErrorMessage } from '../../api/client';
import { queryKeys } from '../../api/queryKeys';
import {
  DOCUMENTATION_ACCEPT,
  DOCUMENTATION_MAX_BYTES,
  DOCUMENTATION_MAX_FILES,
} from '../../schemas/agenda';
import type { Agenda, AgendaDocumentation } from '../../types/agenda';

const props = defineProps<{
  agenda: Agenda;
  canManage: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

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
    toast.success(`${uploaded.length} file berhasil diunggah`);
    queryClient.invalidateQueries({ queryKey: docsKey.value });
  },
  onError: (error: unknown) => toast.error(getErrorMessage(error, 'Gagal mengunggah dokumentasi')),
});

const deleteMutation = useMutation({
  mutationFn: (doc: AgendaDocumentation) => agendaService.deleteDocumentation(props.agenda.id, doc.id),
  onSuccess: () => {
    toast.success('Dokumentasi dihapus');
    queryClient.invalidateQueries({ queryKey: docsKey.value });
  },
  onError: (error: unknown) => toast.error(getErrorMessage(error, 'Gagal menghapus dokumentasi')),
});

/** Client-side pre-check mirroring the server rules, for faster feedback. */
function validateFiles(files: File[]): string | null {
  if (files.length > DOCUMENTATION_MAX_FILES) {
    return `Maksimal ${DOCUMENTATION_MAX_FILES} file sekali unggah`;
  }
  const invalidType = files.find((file) => !allowedTypes.includes(file.type));
  if (invalidType) return `${invalidType.name}: hanya JPG, PNG, WEBP, atau PDF`;
  const tooLarge = files.find((file) => file.size > DOCUMENTATION_MAX_BYTES);
  if (tooLarge) return `${tooLarge.name} melebihi ${DOCUMENTATION_MAX_BYTES / 1024 / 1024} MB`;
  return null;
}

const onFilesSelected = (event: Event) => {
  const input = event.target as HTMLInputElement;
  const files = Array.from(input.files ?? []);
  input.value = '';
  if (files.length === 0) return;

  const problem = validateFiles(files);
  if (problem) {
    toast.error(problem);
    return;
  }
  uploadMutation.mutate(files);
};

const confirmDelete = (doc: AgendaDocumentation) => {
  if (window.confirm(`Hapus ${doc.file_name}?`)) {
    deleteMutation.mutate(doc);
  }
};

const isImage = (doc: AgendaDocumentation) => doc.file_type.startsWith('image/');
const formatSize = (bytes: number) => `${(bytes / 1024 / 1024).toFixed(2)} MB`;
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-2xl overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
        <div>
          <h3 class="font-bold text-gray-900">Dokumentasi Kegiatan</h3>
          <p class="text-xs text-gray-500">{{ agenda.title }}</p>
        </div>
        <button @click="emit('close')" class="text-gray-400 hover:text-gray-600">
          <XIcon class="w-5 h-5" />
        </button>
      </div>

      <div class="p-6 space-y-4 max-h-[70vh] overflow-y-auto">
        <div v-if="canManage">
          <input
            ref="fileInput"
            type="file"
            multiple
            :accept="DOCUMENTATION_ACCEPT"
            class="hidden"
            @change="onFilesSelected"
          />
          <button
            type="button"
            :disabled="uploadMutation.isPending.value"
            @click="fileInput?.click()"
            class="w-full border-2 border-dashed border-gray-200 rounded-xl py-6 flex flex-col items-center gap-2 text-sm text-gray-500 hover:border-emerald-400 hover:text-emerald-600 transition-colors disabled:opacity-50"
          >
            <UploadCloud class="w-6 h-6" />
            <span v-if="uploadMutation.isPending.value">Mengunggah...</span>
            <span v-else>Unggah foto atau PDF (maks. {{ DOCUMENTATION_MAX_FILES }} file, {{ DOCUMENTATION_MAX_BYTES / 1024 / 1024 }} MB per file)</span>
          </button>
        </div>

        <div v-if="isLoading" class="text-center text-sm text-gray-500 py-6">Memuat dokumentasi...</div>
        <div v-else-if="!documents || documents.length === 0" class="text-center text-sm text-gray-500 py-6">
          Belum ada dokumentasi untuk kegiatan ini.
        </div>
        <div v-else class="grid grid-cols-2 sm:grid-cols-3 gap-4">
          <div v-for="doc in documents" :key="doc.id" class="group relative rounded-xl border border-gray-100 overflow-hidden bg-gray-50">
            <a :href="assetUrl(doc.file_url)" target="_blank" rel="noopener" class="block aspect-square">
              <img v-if="isImage(doc)" :src="assetUrl(doc.file_url)" :alt="doc.file_name" class="w-full h-full object-cover" loading="lazy" />
              <div v-else class="w-full h-full flex flex-col items-center justify-center gap-2 text-gray-500">
                <FileText class="w-10 h-10" />
                <span class="text-xs">PDF</span>
              </div>
            </a>
            <div class="px-2 py-1.5 bg-white border-t border-gray-100">
              <p class="text-xs font-medium text-gray-700 truncate">{{ doc.file_name }}</p>
              <p class="text-[10px] text-gray-400">{{ formatSize(doc.file_size) }}</p>
            </div>
            <button
              v-if="canManage"
              @click="confirmDelete(doc)"
              :disabled="deleteMutation.isPending.value"
              class="absolute top-2 right-2 p-1.5 rounded-lg bg-white/90 text-red-600 opacity-0 group-hover:opacity-100 transition-opacity"
              title="Hapus"
            >
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
