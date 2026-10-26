<script setup lang="ts">
import { computed, ref } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { Plus, Undo2, ReceiptText } from 'lucide-vue-next';
import { financeService } from '@/api/services/finance';
import { queryKeys } from '@/api/queryKeys';
import { feedback } from '@/composables/useFeedback';
import { useAuthStore } from '@/stores/auth';
import { Permission } from '@/lib/permissions';
import { formatDateTime, formatRupiah } from '@/lib/format';
import { reversalSchema, transactionSchema } from '@/schemas/finance';
import type { Transaction } from '@/types/finance';
import PageHeader from '@/components/molecules/PageHeader.vue';
import DataPanel from '@/components/molecules/DataPanel.vue';
import PaginationBar from '@/components/molecules/PaginationBar.vue';
import EmptyState from '@/components/molecules/EmptyState.vue';
import StatusBadge from '@/components/molecules/StatusBadge.vue';
import FormDialog from '@/components/molecules/FormDialog.vue';
import FormField from '@/components/molecules/FormField.vue';
import { Button } from '@/components/atoms/button';
import { Dialog, DialogContent, DialogDescription, DialogTitle } from '@/components/atoms/dialog';

const authStore = useAuthStore();
const queryClient = useQueryClient();

const PAGE_SIZE = 10;
const currentPage = ref(1);

const { data: funds, isLoading: isLoadingFunds } = useQuery({
  queryKey: queryKeys.finance.funds,
  queryFn: financeService.getFunds,
});

const { data: transactionPage, isLoading: isLoadingTransactions } = useQuery({
  queryKey: computed(() => queryKeys.finance.transactions(currentPage.value, PAGE_SIZE)),
  queryFn: () => financeService.getTransactions(currentPage.value, PAGE_SIZE),
  placeholderData: (previous) => previous,
});

const transactions = computed(() => transactionPage.value?.items ?? []);
const canWrite = computed(() => authStore.can(Permission.LedgerWrite));

const fundName = (id: number) => funds.value?.find(f => f.id === id)?.name ?? `Kas #${id}`;

const refreshLedger = () => {
  queryClient.invalidateQueries({ queryKey: queryKeys.finance.all });
  queryClient.invalidateQueries({ queryKey: queryKeys.dashboard.metrics });
};

// New transaction
const showTransactionForm = ref(false);
const emptyTransaction = { fund_id: 0, type: 'CREDIT', amount: 0, category: '', description: '' } as const;

const {
  handleSubmit: handleTransactionSubmit,
  errors: transactionErrors,
  defineField: defineTransactionField,
  resetForm: resetTransactionForm,
} = useForm({
  validationSchema: toTypedSchema(transactionSchema),
  initialValues: emptyTransaction,
});

const [fundId] = defineTransactionField('fund_id');
const [txType] = defineTransactionField('type');
const [amount] = defineTransactionField('amount');
const [category] = defineTransactionField('category');
const [description] = defineTransactionField('description');

const createTxMutation = useMutation({
  mutationFn: financeService.createTransaction,
  onSuccess: (tx) => {
    refreshLedger();
    showTransactionForm.value = false;
    feedback.success('Transaksi tersimpan', `${formatRupiah(tx.amount)} dicatat di ${fundName(tx.fund_id)}.`);
  },
  onError: (error: unknown) => feedback.error('Transaksi belum tersimpan', error),
});

const openTransactionForm = () => {
  resetTransactionForm({ values: emptyTransaction });
  showTransactionForm.value = true;
};

const submitTransaction = handleTransactionSubmit((values) => createTxMutation.mutate(values));

// Detail and reversal
const selected = ref<Transaction | null>(null);
const showReversalForm = ref(false);

const {
  handleSubmit: handleReversalSubmit,
  errors: reversalErrors,
  defineField: defineReversalField,
  resetForm: resetReversalForm,
} = useForm({
  validationSchema: toTypedSchema(reversalSchema),
  initialValues: { reason: '' },
});

const [reason] = defineReversalField('reason');

/** Each transaction can be reversed once; reversal entries cannot be reversed. */
const canReverse = (tx: Transaction) => canWrite.value && tx.reversal_of_id === null && tx.reversed_by_id === null;

const reverseTxMutation = useMutation({
  mutationFn: ({ id, reason }: { id: string; reason: string }) => financeService.reverseTransaction(id, { reason }),
  onSuccess: () => {
    refreshLedger();
    showReversalForm.value = false;
    selected.value = null;
    feedback.success('Transaksi dibatalkan', 'Entri pembalik sudah dicatat dan saldo kas disesuaikan.');
  },
  onError: (error: unknown) => feedback.error('Transaksi belum dibatalkan', error),
});

const openReversal = () => {
  resetReversalForm({ values: { reason: '' } });
  showReversalForm.value = true;
};

const submitReversal = handleReversalSubmit((values) => {
  if (selected.value) reverseTxMutation.mutate({ id: selected.value.id, reason: values.reason });
});

const onDetailOpenChange = (open: boolean) => {
  if (!open) selected.value = null;
};
</script>

<template>
  <div class="space-y-8">
    <PageHeader title="Keuangan" description="Buku kas masjid. Setiap entri bersifat permanen; kesalahan dikoreksi dengan entri pembalik.">
      <template v-if="canWrite" #actions>
        <Button @click="openTransactionForm"><Plus /> Catat transaksi</Button>
      </template>
    </PageHeader>

    <!-- Funds -->
    <div v-if="isLoadingFunds" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="i in 3" :key="i" class="h-28 animate-pulse rounded-xl bg-batu/60" />
    </div>
    <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="fund in funds" :key="fund.id" class="rounded-xl border border-batu bg-white p-5 shadow-sm">
        <div class="flex items-center justify-between gap-3">
          <p class="font-medium text-jati">{{ fund.name }}</p>
          <span class="font-mono text-xs text-jati-muda">{{ fund.code }}</span>
        </div>
        <p class="tabular mt-3 text-2xl font-semibold text-jati">{{ formatRupiah(fund.current_balance) }}</p>
      </div>
    </div>

    <!-- Ledger -->
    <DataPanel title="Riwayat transaksi" description="Urut dari yang terbaru.">
      <div v-if="isLoadingTransactions" class="p-8 text-center text-sm text-jati-muda">Memuat transaksi…</div>
      <EmptyState
        v-else-if="transactions.length === 0"
        :icon="ReceiptText"
        title="Belum ada transaksi"
        :description="canWrite ? 'Catat pemasukan atau pengeluaran pertama untuk mulai mengisi buku kas.' : 'Transaksi yang dicatat bendahara akan muncul di sini.'"
      />
      <div v-else class="scrollbar-tipis overflow-x-auto overscroll-x-contain">
        <table class="tabel min-w-[760px]">
          <thead>
            <tr>
              <th>Tanggal</th>
              <th>Kas</th>
              <th>Kategori</th>
              <th class="text-right">Nominal</th>
              <th class="w-24"><span class="sr-only">Aksi</span></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="tx in transactions" :key="tx.id">
              <td class="whitespace-nowrap text-jati-muda">{{ formatDateTime(tx.created_at) }}</td>
              <td>{{ fundName(tx.fund_id) }}</td>
              <td>
                <span class="font-medium">{{ tx.category }}</span>
                <StatusBadge v-if="tx.reversal_of_id" tone="kuningan" class="ml-2">Pembalik</StatusBadge>
                <StatusBadge v-else-if="tx.reversed_by_id" class="ml-2">Dibatalkan</StatusBadge>
              </td>
              <td
                class="tabular whitespace-nowrap text-right font-semibold"
                :class="[tx.type === 'CREDIT' ? 'text-sajadah' : 'text-bata', { 'line-through opacity-60': tx.reversed_by_id }]"
              >
                {{ tx.type === 'CREDIT' ? '+' : '−' }} {{ formatRupiah(tx.amount) }}
              </td>
              <td class="text-right">
                <button type="button" class="text-sm font-medium text-sajadah hover:text-sajadah-tua" @click="selected = tx">
                  Detail
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <template v-if="transactionPage && transactionPage.total > 0" #footer>
        <PaginationBar
          v-model:page="currentPage"
          :page-size="PAGE_SIZE"
          :total="transactionPage.total"
          :total-pages="transactionPage.total_pages"
          noun="transaksi"
        />
      </template>
    </DataPanel>

    <!-- New transaction -->
    <FormDialog
      v-model:open="showTransactionForm"
      title="Catat transaksi"
      description="Transaksi tidak bisa diubah setelah disimpan. Periksa nominal dan kas tujuannya."
      submit-label="Simpan transaksi"
      :submitting="createTxMutation.isPending.value"
      @submit="submitTransaction"
    >
      <FormField label="Kas" for="tx-fund" :error="transactionErrors.fund_id" required>
        <select id="tx-fund" v-model.number="fundId" class="field-control" :aria-invalid="Boolean(transactionErrors.fund_id)">
          <option :value="0" disabled>Pilih kas</option>
          <option v-for="fund in funds" :key="fund.id" :value="fund.id">{{ fund.name }} ({{ formatRupiah(fund.current_balance) }})</option>
        </select>
      </FormField>
      <fieldset>
        <legend class="mb-1.5 text-sm font-medium text-jati">Jenis</legend>
        <div class="grid grid-cols-2 gap-2">
          <label
            v-for="option in [{ value: 'CREDIT', label: 'Pemasukan' }, { value: 'DEBIT', label: 'Pengeluaran' }]"
            :key="option.value"
            class="flex cursor-pointer items-center gap-2 rounded-lg border px-3 py-2.5 text-sm font-medium transition-colors"
            :class="txType === option.value ? 'border-sajadah bg-sajadah-muda text-sajadah-tua' : 'border-batu-tua text-jati'"
          >
            <input v-model="txType" type="radio" :value="option.value" class="accent-sajadah" />
            {{ option.label }}
          </label>
        </div>
      </fieldset>
      <FormField label="Nominal (Rp)" for="tx-amount" :error="transactionErrors.amount" required>
        <input id="tx-amount" v-model.number="amount" type="number" min="1" step="1" class="field-control tabular" :aria-invalid="Boolean(transactionErrors.amount)" />
      </FormField>
      <FormField label="Kategori" for="tx-category" :error="transactionErrors.category" hint="Contoh: Infaq Jumat, Listrik, Kegiatan." required>
        <input id="tx-category" v-model="category" type="text" class="field-control" :aria-invalid="Boolean(transactionErrors.category)" />
      </FormField>
      <FormField label="Keterangan" for="tx-description" :error="transactionErrors.description">
        <input id="tx-description" v-model="description" type="text" class="field-control" />
      </FormField>
    </FormDialog>

    <!-- Transaction detail -->
    <Dialog :open="selected !== null && !showReversalForm" @update:open="onDetailOpenChange">
      <DialogContent v-if="selected" class="gap-0 rounded-2xl border-batu p-0 sm:max-w-md">
        <div class="border-b border-batu px-6 pt-6 pb-4 pr-12">
          <DialogTitle class="font-display text-xl font-semibold tracking-wide">Detail transaksi</DialogTitle>
          <DialogDescription class="mt-1 font-mono text-xs text-jati-muda">{{ selected.id }}</DialogDescription>
        </div>
        <dl class="shaf px-6 text-sm">
          <div class="flex justify-between gap-4 py-3">
            <dt class="text-jati-muda">Nominal</dt>
            <dd class="tabular font-semibold" :class="selected.type === 'CREDIT' ? 'text-sajadah' : 'text-bata'">
              {{ selected.type === 'CREDIT' ? 'Pemasukan' : 'Pengeluaran' }} {{ formatRupiah(selected.amount) }}
            </dd>
          </div>
          <div class="flex justify-between gap-4 py-3">
            <dt class="text-jati-muda">Kas</dt>
            <dd>{{ fundName(selected.fund_id) }}</dd>
          </div>
          <div class="flex justify-between gap-4 py-3">
            <dt class="text-jati-muda">Tanggal</dt>
            <dd>{{ formatDateTime(selected.created_at) }}</dd>
          </div>
          <div class="flex justify-between gap-4 py-3">
            <dt class="text-jati-muda">Kategori</dt>
            <dd>{{ selected.category }}</dd>
          </div>
          <div v-if="selected.description" class="flex justify-between gap-4 py-3">
            <dt class="text-jati-muda">Keterangan</dt>
            <dd class="text-right">{{ selected.description }}</dd>
          </div>
          <div v-if="selected.reversal_of_id" class="py-3">
            <dt class="text-jati-muda">Membatalkan transaksi</dt>
            <dd class="mt-1 font-mono text-xs">{{ selected.reversal_of_id }}</dd>
          </div>
          <div v-if="selected.reversed_by_id" class="py-3">
            <dt class="text-jati-muda">Dibatalkan oleh entri</dt>
            <dd class="mt-1 font-mono text-xs">{{ selected.reversed_by_id }}</dd>
          </div>
        </dl>
        <div class="flex justify-end gap-2 border-t border-batu bg-kapur/60 px-6 py-4">
          <Button v-if="canReverse(selected)" variant="outline" class="text-bata hover:text-bata" @click="openReversal">
            <Undo2 /> Batalkan transaksi
          </Button>
          <Button @click="selected = null">Tutup</Button>
        </div>
      </DialogContent>
    </Dialog>

    <!-- Reversal -->
    <FormDialog
      v-model:open="showReversalForm"
      title="Batalkan transaksi"
      :description="selected ? `Entri ${selected.type === 'CREDIT' ? 'pengeluaran' : 'pemasukan'} sebesar ${formatRupiah(selected.amount)} akan dicatat di ${fundName(selected.fund_id)} untuk menyeimbangkan transaksi ini.` : undefined"
      submit-label="Catat pembatalan"
      tone="danger"
      :submitting="reverseTxMutation.isPending.value"
      @submit="submitReversal"
    >
      <FormField label="Alasan" for="reversal-reason" :error="reversalErrors.reason" hint="Akan tersimpan sebagai keterangan entri pembalik." required>
        <textarea id="reversal-reason" v-model="reason" rows="3" class="field-control" placeholder="Contoh: Tercatat dua kali" :aria-invalid="Boolean(reversalErrors.reason)" />
      </FormField>
    </FormDialog>
  </div>
</template>
