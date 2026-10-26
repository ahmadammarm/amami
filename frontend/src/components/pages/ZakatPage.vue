<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { ArrowDownToLine, ArrowUpFromLine, HandCoins } from 'lucide-vue-next';
import { zakatService } from '@/api/services/zakat';
import { financeService } from '@/api/services/finance';
import { queryKeys } from '@/api/queryKeys';
import { feedback } from '@/composables/useFeedback';
import { useAuthStore } from '@/stores/auth';
import { Permission } from '@/lib/permissions';
import { formatAmount, formatDateTime, formatRupiah } from '@/lib/format';
import { collectZakatSchema, distributeZakatSchema } from '@/schemas/zakat';
import type { ZakatType } from '@/types/zakat';
import PageHeader from '@/components/molecules/PageHeader.vue';
import DataPanel from '@/components/molecules/DataPanel.vue';
import PaginationBar from '@/components/molecules/PaginationBar.vue';
import EmptyState from '@/components/molecules/EmptyState.vue';
import StatusBadge from '@/components/molecules/StatusBadge.vue';
import FormDialog from '@/components/molecules/FormDialog.vue';
import FormField from '@/components/molecules/FormField.vue';
import AsyncJamaahSelect from '@/components/molecules/AsyncJamaahSelect.vue';
import { Button } from '@/components/atoms/button';

type Tab = 'donations' | 'distributions';

const ZAKAT_TYPE_LABELS: Record<ZakatType, string> = {
  FITRAH_BERAS: 'Zakat fitrah (beras)',
  MAAL: 'Zakat maal',
  INFAQ: 'Infaq / sedekah',
};

const authStore = useAuthStore();
const queryClient = useQueryClient();

const PAGE_SIZE = 10;
const activeTab = ref<Tab>('donations');
const donationPage = ref(1);
const distributionPage = ref(1);

const canManage = computed(() => authStore.can(Permission.ZakatManage));
const canReadFunds = computed(() => authStore.can(Permission.LedgerRead));

const { data: donations, isLoading: isLoadingDonations } = useQuery({
  queryKey: computed(() => queryKeys.zakat.donations(donationPage.value, PAGE_SIZE)),
  queryFn: () => zakatService.getDonations(donationPage.value, PAGE_SIZE),
  placeholderData: (previous) => previous,
});

const { data: distributions, isLoading: isLoadingDistributions } = useQuery({
  queryKey: computed(() => queryKeys.zakat.distributions(distributionPage.value, PAGE_SIZE)),
  queryFn: () => zakatService.getDistributions(distributionPage.value, PAGE_SIZE),
  placeholderData: (previous) => previous,
});

const { data: funds } = useQuery({
  queryKey: queryKeys.finance.funds,
  queryFn: financeService.getFunds,
  enabled: canReadFunds,
});

const zakatTypeLabel = (type: string) => ZAKAT_TYPE_LABELS[type as ZakatType] ?? type.replace(/_/g, ' ');

const refreshAfterZakat = () => {
  queryClient.invalidateQueries({ queryKey: queryKeys.zakat.all });
  queryClient.invalidateQueries({ queryKey: queryKeys.finance.all });
  queryClient.invalidateQueries({ queryKey: queryKeys.dashboard.metrics });
};

// Collect
const showCollectForm = ref(false);
const emptyCollect = { muzakki_id: '', zakat_type: 'MAAL', fund_id: null, amount: 0, unit: 'IDR', description: '' } as const;

const {
  handleSubmit: handleCollectSubmit,
  errors: collectErrors,
  defineField: defineCollectField,
  setFieldValue: setCollectField,
  resetForm: resetCollectForm,
} = useForm({
  validationSchema: toTypedSchema(collectZakatSchema),
  initialValues: emptyCollect,
});

const [muzakkiId] = defineCollectField('muzakki_id');
const [zakatType] = defineCollectField('zakat_type');
const [collectFundId] = defineCollectField('fund_id');
const [collectAmount] = defineCollectField('amount');
const [collectUnit] = defineCollectField('unit');
const [collectDescription] = defineCollectField('description');

// Rice is counted in KG and never posted to a fund; cash is always IDR.
watch(zakatType, (type) => {
  const isRice = type === 'FITRAH_BERAS';
  setCollectField('unit', isRice ? 'KG' : 'IDR');
  if (isRice) setCollectField('fund_id', null);
});

const collectMutation = useMutation({
  mutationFn: zakatService.collectZakat,
  onSuccess: (donation) => {
    refreshAfterZakat();
    showCollectForm.value = false;
    feedback.success('Penerimaan zakat tercatat', `${formatAmount(donation.amount, donation.unit)} dari ${donation.muzakki_name}.`);
  },
  onError: (error: unknown) => feedback.error('Penerimaan zakat belum tercatat', error),
});

const openCollectForm = () => {
  resetCollectForm({ values: emptyCollect });
  showCollectForm.value = true;
};

const submitCollect = handleCollectSubmit((values) => collectMutation.mutate(values));

// Distribute
const showDistributeForm = ref(false);
const emptyDistribute = { mustahik_id: '', fund_id: null, amount: 0, unit: 'IDR', description: '' } as const;

const {
  handleSubmit: handleDistributeSubmit,
  errors: distributeErrors,
  defineField: defineDistributeField,
  setFieldValue: setDistributeField,
  resetForm: resetDistributeForm,
} = useForm({
  validationSchema: toTypedSchema(distributeZakatSchema),
  initialValues: emptyDistribute,
});

const [mustahikId] = defineDistributeField('mustahik_id');
const [distributeFundId] = defineDistributeField('fund_id');
const [distributeAmount] = defineDistributeField('amount');
const [distributeUnit] = defineDistributeField('unit');
const [distributeDescription] = defineDistributeField('description');

watch(distributeUnit, (unit) => {
  if (unit !== 'IDR') setDistributeField('fund_id', null);
});

const distributeMutation = useMutation({
  mutationFn: zakatService.distributeZakat,
  onSuccess: (distribution) => {
    refreshAfterZakat();
    showDistributeForm.value = false;
    feedback.success('Penyaluran zakat tercatat', `${formatAmount(distribution.amount, distribution.unit)} untuk ${distribution.mustahik_name}.`);
  },
  onError: (error: unknown) => feedback.error('Penyaluran zakat belum tercatat', error),
});

const openDistributeForm = () => {
  resetDistributeForm({ values: emptyDistribute });
  showDistributeForm.value = true;
};

const submitDistribute = handleDistributeSubmit((values) => distributeMutation.mutate(values));

const tabs: { id: Tab; label: string }[] = [
  { id: 'donations', label: 'Penerimaan' },
  { id: 'distributions', label: 'Penyaluran' },
];
</script>

<template>
  <div class="space-y-8">
    <PageHeader title="Zakat" description="Penerimaan dari muzakki dan penyaluran kepada mustahik. Zakat uang otomatis tercatat di kas.">
      <template v-if="canManage" #actions>
        <Button variant="outline" @click="openDistributeForm"><ArrowUpFromLine /> Salurkan</Button>
        <Button @click="openCollectForm"><ArrowDownToLine /> Terima zakat</Button>
      </template>
    </PageHeader>

    <div class="inline-flex rounded-lg border border-batu bg-white p-1" role="tablist" aria-label="Jenis catatan zakat">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        type="button"
        role="tab"
        :aria-selected="activeTab === tab.id"
        class="rounded-md px-4 py-1.5 text-sm font-medium transition-colors"
        :class="activeTab === tab.id ? 'bg-sajadah text-white' : 'text-jati-muda hover:text-jati'"
        @click="activeTab = tab.id"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- Donations -->
    <DataPanel v-if="activeTab === 'donations'">
      <div v-if="isLoadingDonations" class="p-8 text-center text-sm text-jati-muda">Memuat penerimaan…</div>
      <EmptyState
        v-else-if="!donations || donations.items.length === 0"
        :icon="HandCoins"
        title="Belum ada penerimaan zakat"
        description="Zakat dan infaq yang diterima pengurus akan tercatat di sini."
      />
      <div v-else class="overflow-x-auto">
        <table class="tabel min-w-[720px]">
          <thead>
            <tr>
              <th>Waktu</th>
              <th>Muzakki</th>
              <th>Jenis</th>
              <th>Keterangan</th>
              <th class="text-right">Jumlah</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="d in donations.items" :key="d.id">
              <td class="whitespace-nowrap text-jati-muda">{{ formatDateTime(d.created_at) }}</td>
              <td class="font-medium">{{ d.muzakki_name }}</td>
              <td><StatusBadge tone="sajadah">{{ zakatTypeLabel(d.zakat_type) }}</StatusBadge></td>
              <td class="max-w-xs truncate text-jati-muda">{{ d.description || '–' }}</td>
              <td class="tabular whitespace-nowrap text-right font-semibold text-sajadah">+ {{ formatAmount(d.amount, d.unit) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <template v-if="donations && donations.total > 0" #footer>
        <PaginationBar v-model:page="donationPage" :page-size="PAGE_SIZE" :total="donations.total" :total-pages="donations.total_pages" noun="penerimaan" />
      </template>
    </DataPanel>

    <!-- Distributions -->
    <DataPanel v-else>
      <div v-if="isLoadingDistributions" class="p-8 text-center text-sm text-jati-muda">Memuat penyaluran…</div>
      <EmptyState
        v-else-if="!distributions || distributions.items.length === 0"
        :icon="HandCoins"
        title="Belum ada penyaluran zakat"
        description="Bantuan yang disalurkan kepada mustahik akan tercatat di sini."
      />
      <div v-else class="overflow-x-auto">
        <table class="tabel min-w-[640px]">
          <thead>
            <tr>
              <th>Waktu</th>
              <th>Mustahik</th>
              <th>Keterangan</th>
              <th class="text-right">Jumlah</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="d in distributions.items" :key="d.id">
              <td class="whitespace-nowrap text-jati-muda">{{ formatDateTime(d.distributed_at) }}</td>
              <td class="font-medium">{{ d.mustahik_name }}</td>
              <td class="max-w-xs truncate text-jati-muda">{{ d.description || '–' }}</td>
              <td class="tabular whitespace-nowrap text-right font-semibold text-bata">− {{ formatAmount(d.amount, d.unit) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <template v-if="distributions && distributions.total > 0" #footer>
        <PaginationBar v-model:page="distributionPage" :page-size="PAGE_SIZE" :total="distributions.total" :total-pages="distributions.total_pages" noun="penyaluran" />
      </template>
    </DataPanel>

    <!-- Collect -->
    <FormDialog
      v-model:open="showCollectForm"
      title="Terima zakat"
      description="Zakat uang dicatat sebagai pemasukan kas yang dipilih."
      submit-label="Simpan penerimaan"
      :submitting="collectMutation.isPending.value"
      @submit="submitCollect"
    >
      <FormField label="Muzakki" :error="collectErrors.muzakki_id" required>
        <AsyncJamaahSelect v-model="muzakkiId" muzakki-only placeholder="Cari nama muzakki…" />
      </FormField>
      <FormField label="Jenis zakat" for="collect-type" required>
        <select id="collect-type" v-model="zakatType" class="field-control">
          <option v-for="(label, value) in ZAKAT_TYPE_LABELS" :key="value" :value="value">{{ label }}</option>
        </select>
      </FormField>
      <FormField v-if="collectUnit === 'IDR'" label="Masuk ke kas" for="collect-fund" :error="collectErrors.fund_id" required>
        <select id="collect-fund" v-model="collectFundId" class="field-control" :aria-invalid="Boolean(collectErrors.fund_id)">
          <option :value="null" disabled>Pilih kas</option>
          <option v-for="fund in funds" :key="fund.id" :value="fund.id">{{ fund.name }} ({{ formatRupiah(fund.current_balance) }})</option>
        </select>
      </FormField>
      <FormField
        :label="collectUnit === 'IDR' ? 'Nominal (Rp)' : 'Jumlah beras (kg)'"
        for="collect-amount"
        :error="collectErrors.amount || collectErrors.unit"
        required
      >
        <input
          id="collect-amount"
          v-model.number="collectAmount"
          type="number"
          min="0"
          :step="collectUnit === 'IDR' ? 1 : 0.1"
          class="field-control tabular"
          :aria-invalid="Boolean(collectErrors.amount)"
        />
      </FormField>
      <FormField label="Keterangan" for="collect-description" :error="collectErrors.description">
        <textarea id="collect-description" v-model="collectDescription" rows="2" class="field-control" />
      </FormField>
    </FormDialog>

    <!-- Distribute -->
    <FormDialog
      v-model:open="showDistributeForm"
      title="Salurkan zakat"
      description="Bantuan uang dicatat sebagai pengeluaran dari kas yang dipilih."
      submit-label="Simpan penyaluran"
      :submitting="distributeMutation.isPending.value"
      @submit="submitDistribute"
    >
      <FormField
        label="Mustahik"
        :error="distributeErrors.mustahik_id"
        hint="Hanya mustahik terdaftar, diurutkan dari skor kebutuhan tertinggi."
        required
      >
        <AsyncJamaahSelect v-model="mustahikId" mustahik-only placeholder="Cari nama mustahik…" />
      </FormField>
      <fieldset>
        <legend class="mb-1.5 text-sm font-medium text-jati">Bentuk bantuan</legend>
        <div class="grid grid-cols-2 gap-2">
          <label
            v-for="option in [{ value: 'IDR', label: 'Uang' }, { value: 'KG', label: 'Beras (kg)' }]"
            :key="option.value"
            class="flex cursor-pointer items-center gap-2 rounded-lg border px-3 py-2.5 text-sm font-medium transition-colors"
            :class="distributeUnit === option.value ? 'border-sajadah bg-sajadah-muda text-sajadah-tua' : 'border-batu-tua text-jati'"
          >
            <input v-model="distributeUnit" type="radio" :value="option.value" class="accent-sajadah" />
            {{ option.label }}
          </label>
        </div>
      </fieldset>
      <FormField v-if="distributeUnit === 'IDR'" label="Diambil dari kas" for="distribute-fund" :error="distributeErrors.fund_id" required>
        <select id="distribute-fund" v-model="distributeFundId" class="field-control" :aria-invalid="Boolean(distributeErrors.fund_id)">
          <option :value="null" disabled>Pilih kas</option>
          <option v-for="fund in funds" :key="fund.id" :value="fund.id">{{ fund.name }} (tersedia {{ formatRupiah(fund.current_balance) }})</option>
        </select>
      </FormField>
      <FormField
        :label="distributeUnit === 'IDR' ? 'Nominal (Rp)' : 'Jumlah beras (kg)'"
        for="distribute-amount"
        :error="distributeErrors.amount || distributeErrors.unit"
        required
      >
        <input
          id="distribute-amount"
          v-model.number="distributeAmount"
          type="number"
          min="0"
          :step="distributeUnit === 'IDR' ? 1 : 0.1"
          class="field-control tabular"
          :aria-invalid="Boolean(distributeErrors.amount)"
        />
      </FormField>
      <FormField label="Keterangan" for="distribute-description" :error="distributeErrors.description">
        <textarea id="distribute-description" v-model="distributeDescription" rows="2" class="field-control" />
      </FormField>
    </FormDialog>
  </div>
</template>
