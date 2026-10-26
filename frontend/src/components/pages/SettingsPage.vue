<script setup lang="ts">
import { computed, ref } from 'vue';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { Pencil, RotateCw } from 'lucide-vue-next';
import { useMosqueProfileQuery, useUpdateMosqueProfile } from '../../composables/useSettingsQuery';
import { feedback } from '../../composables/useFeedback';
import { formatRupiah } from '../../lib/format';
import { mosqueProfileSchema, type MosqueProfileFormValues } from '../../schemas/settings';
import type { MosqueProfile } from '../../types/settings';
import PageHeader from '../molecules/PageHeader.vue';
import DataPanel from '../molecules/DataPanel.vue';
import EmptyState from '../molecules/EmptyState.vue';
import FormDialog from '../molecules/FormDialog.vue';
import FormField from '../molecules/FormField.vue';
import { Button } from '../atoms/button';

interface Row {
  label: string;
  value: string;
  mono?: boolean;
}

const { data: profile, isLoading, isError, refetch } = useMosqueProfileQuery();
const { mutate: updateProfile, isPending: isSaving } = useUpdateMosqueProfile();

const orDash = (value: string | undefined) => value?.trim() || '–';

const identityRows = computed<Row[]>(() => [
  { label: 'Nama masjid', value: orDash(profile.value?.name) },
  { label: 'Alamat', value: orDash(profile.value?.address) },
  { label: 'Telepon', value: orDash(profile.value?.phone) },
  { label: 'Nomor legal yayasan', value: orDash(profile.value?.legal_yayasan_id), mono: true },
]);

const bankRows = computed<Row[]>(() => [
  { label: 'Bank', value: orDash(profile.value?.bank_name) },
  { label: 'Nomor rekening', value: orDash(profile.value?.bank_account_number), mono: true },
  { label: 'Atas nama', value: orDash(profile.value?.bank_account_name) },
]);

const operationRows = computed<Row[]>(() => {
  const fitrah = Number(profile.value?.zakat_fitrah_amount);
  return [
    { label: 'Zakat fitrah per jiwa', value: Number.isFinite(fitrah) && fitrah > 0 ? formatRupiah(fitrah) : '–' },
    { label: 'Tahun hijriah aktif', value: profile.value?.active_hijri_year ? `${profile.value.active_hijri_year} H` : '–' },
  ];
});

// Edit
const showForm = ref(false);

const toFormValues = (p: MosqueProfile | undefined): MosqueProfileFormValues => ({
  name: p?.name ?? '',
  address: p?.address ?? '',
  phone: p?.phone ?? '',
  legal_yayasan_id: p?.legal_yayasan_id ?? '',
  bank_name: p?.bank_name ?? '',
  bank_account_name: p?.bank_account_name ?? '',
  bank_account_number: p?.bank_account_number ?? '',
  zakat_fitrah_amount: p?.zakat_fitrah_amount ?? '',
  active_hijri_year: p?.active_hijri_year ?? '',
});

const { handleSubmit, errors, defineField, resetForm, meta } = useForm({
  validationSchema: toTypedSchema(mosqueProfileSchema),
  initialValues: toFormValues(undefined),
});

const [name] = defineField('name');
const [address] = defineField('address');
const [phone] = defineField('phone');
const [legalId] = defineField('legal_yayasan_id');
const [bankName] = defineField('bank_name');
const [bankAccountName] = defineField('bank_account_name');
const [bankAccountNumber] = defineField('bank_account_number');
const [zakatFitrah] = defineField('zakat_fitrah_amount');
const [hijriYear] = defineField('active_hijri_year');

const openForm = () => {
  resetForm({ values: toFormValues(profile.value) });
  showForm.value = true;
};

/** Closing with unsaved edits asks first, so a stray click does not discard them. */
const onFormOpenChange = async (open: boolean) => {
  if (open) return;
  if (meta.value.dirty) {
    const discard = await feedback.confirm({
      title: 'Buang perubahan?',
      message: 'Perubahan profil yang belum disimpan akan hilang.',
      confirmLabel: 'Buang perubahan',
      cancelLabel: 'Lanjut mengubah',
      danger: true,
    });
    if (!discard) return;
  }
  showForm.value = false;
};

const submitForm = handleSubmit((values) => {
  updateProfile({ ...values, logo: profile.value?.logo }, {
    onSuccess: () => {
      showForm.value = false;
      feedback.success('Profil masjid tersimpan', 'Perubahan langsung tampil di halaman depan masjid.');
    },
  });
});
</script>

<template>
  <div class="space-y-8">
    <PageHeader title="Pengaturan" description="Identitas masjid dan rekening yang ditampilkan kepada jamaah di halaman depan.">
      <template v-if="profile" #actions>
        <Button @click="openForm"><Pencil /> Ubah profil</Button>
      </template>
    </PageHeader>

    <div v-if="isLoading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="h-40 animate-pulse rounded-xl bg-batu/60" />
    </div>

    <div v-else-if="isError" class="rounded-xl border border-batu bg-white">
      <EmptyState title="Profil masjid tidak dapat dimuat" description="Server tidak merespons. Periksa koneksi, lalu coba lagi.">
        <Button variant="outline" @click="refetch()"><RotateCw /> Coba lagi</Button>
      </EmptyState>
    </div>

    <template v-else>
      <DataPanel title="Identitas masjid">
        <dl class="shaf px-5">
          <div v-for="row in identityRows" :key="row.label" class="grid gap-1 py-3.5 sm:grid-cols-[14rem_1fr]">
            <dt class="text-sm text-jati-muda">{{ row.label }}</dt>
            <dd class="text-sm text-jati" :class="{ 'font-mono': row.mono }">{{ row.value }}</dd>
          </div>
        </dl>
      </DataPanel>

      <div class="grid gap-6 lg:grid-cols-2">
        <DataPanel title="Rekening masjid" description="Tampil di halaman depan untuk transfer zakat dan infaq.">
          <dl class="shaf px-5">
            <div v-for="row in bankRows" :key="row.label" class="grid gap-1 py-3.5 sm:grid-cols-[10rem_1fr]">
              <dt class="text-sm text-jati-muda">{{ row.label }}</dt>
              <dd class="text-sm text-jati" :class="{ 'font-mono tracking-wide': row.mono }">{{ row.value }}</dd>
            </div>
          </dl>
        </DataPanel>
        <DataPanel title="Zakat dan kalender" description="Acuan untuk penerimaan zakat fitrah dan paket qurban.">
          <dl class="shaf px-5">
            <div v-for="row in operationRows" :key="row.label" class="grid gap-1 py-3.5 sm:grid-cols-[12rem_1fr]">
              <dt class="text-sm text-jati-muda">{{ row.label }}</dt>
              <dd class="tabular text-sm text-jati">{{ row.value }}</dd>
            </div>
          </dl>
        </DataPanel>
      </div>
    </template>

    <FormDialog
      :open="showForm"
      title="Ubah profil masjid"
      submit-label="Simpan profil"
      size="lg"
      :submitting="isSaving"
      @update:open="onFormOpenChange"
      @submit="submitForm"
    >
      <p class="text-xs font-semibold uppercase tracking-[0.14em] text-jati-muda">Identitas</p>
      <FormField label="Nama masjid" for="profile-name" :error="errors.name" required>
        <input id="profile-name" v-model="name" type="text" class="field-control" :aria-invalid="Boolean(errors.name)" />
      </FormField>
      <FormField label="Alamat" for="profile-address" :error="errors.address" required>
        <textarea id="profile-address" v-model="address" rows="2" class="field-control" :aria-invalid="Boolean(errors.address)" />
      </FormField>
      <div class="grid gap-4 sm:grid-cols-2">
        <FormField label="Telepon" for="profile-phone" :error="errors.phone" required>
          <input id="profile-phone" v-model="phone" type="tel" class="field-control" :aria-invalid="Boolean(errors.phone)" />
        </FormField>
        <FormField label="Nomor legal yayasan" for="profile-legal" :error="errors.legal_yayasan_id" required>
          <input id="profile-legal" v-model="legalId" type="text" class="field-control font-mono" :aria-invalid="Boolean(errors.legal_yayasan_id)" />
        </FormField>
      </div>

      <p class="pt-2 text-xs font-semibold uppercase tracking-[0.14em] text-jati-muda">Rekening</p>
      <div class="grid gap-4 sm:grid-cols-2">
        <FormField label="Bank" for="profile-bank" :error="errors.bank_name">
          <input id="profile-bank" v-model="bankName" type="text" class="field-control" placeholder="Contoh: Bank Syariah Indonesia" />
        </FormField>
        <FormField label="Nomor rekening" for="profile-account-number" :error="errors.bank_account_number">
          <input id="profile-account-number" v-model="bankAccountNumber" type="text" inputmode="numeric" class="field-control font-mono" :aria-invalid="Boolean(errors.bank_account_number)" />
        </FormField>
      </div>
      <FormField label="Atas nama" for="profile-account-name" :error="errors.bank_account_name">
        <input id="profile-account-name" v-model="bankAccountName" type="text" class="field-control" />
      </FormField>

      <p class="pt-2 text-xs font-semibold uppercase tracking-[0.14em] text-jati-muda">Zakat dan kalender</p>
      <div class="grid gap-4 sm:grid-cols-2">
        <FormField label="Zakat fitrah per jiwa (Rp)" for="profile-fitrah" :error="errors.zakat_fitrah_amount" hint="Setara 2,5 kg beras.">
          <input id="profile-fitrah" v-model="zakatFitrah" type="text" inputmode="numeric" class="field-control tabular" :aria-invalid="Boolean(errors.zakat_fitrah_amount)" />
        </FormField>
        <FormField label="Tahun hijriah aktif" for="profile-hijri" :error="errors.active_hijri_year">
          <input id="profile-hijri" v-model="hijriYear" type="text" inputmode="numeric" class="field-control tabular" :aria-invalid="Boolean(errors.active_hijri_year)" />
        </FormField>
      </div>
    </FormDialog>
  </div>
</template>
