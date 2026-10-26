<script setup lang="ts">
import { computed, ref } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { Plus, Pencil, Trash2, Beef, UserPlus } from 'lucide-vue-next';
import { qurbanService } from '../../api/services/qurban';
import { queryKeys } from '../../api/queryKeys';
import { feedback } from '../../composables/useFeedback';
import { useAuthStore } from '../../stores/auth';
import { Permission } from '../../lib/permissions';
import { formatDate, formatNumber, formatRupiah } from '../../lib/format';
import {
  animalSchema,
  bookingSchema,
  packageSchema,
  type AnimalPayload,
  type BookingPayload,
  type PackagePayload,
} from '../../schemas/qurban';
import type { AnimalStatus, AnimalType, PackageType, PaymentStatus, QurbanAnimal, QurbanPackage } from '../../types/qurban';
import type { BadgeTone } from '../../types/ui';
import PageHeader from '../molecules/PageHeader.vue';
import DataPanel from '../molecules/DataPanel.vue';
import EmptyState from '../molecules/EmptyState.vue';
import StatCard from '../molecules/StatCard.vue';
import StatusBadge from '../molecules/StatusBadge.vue';
import FormDialog from '../molecules/FormDialog.vue';
import FormField from '../molecules/FormField.vue';
import AsyncJamaahSelect from '../molecules/AsyncJamaahSelect.vue';
import { Button } from '../atoms/button';

type TabId = 'packages' | 'bookings' | 'animals';

const PACKAGE_TYPE_LABELS: Record<PackageType, string> = {
  SAPI: 'Sapi',
  SAPI_KOLEKTIF: 'Sapi kolektif (1/7)',
  KAMBING: 'Kambing',
  DOMBA: 'Domba',
};

const ANIMAL_TYPE_LABELS: Record<AnimalType, string> = {
  SAPI: 'Sapi',
  KAMBING: 'Kambing',
  DOMBA: 'Domba',
};

const ANIMAL_STATUS: Record<AnimalStatus, { label: string; tone: BadgeTone }> = {
  AVAILABLE: { label: 'Tersedia', tone: 'sajadah' },
  ASSIGNED: { label: 'Dialokasikan', tone: 'kuningan' },
  SLAUGHTERED: { label: 'Disembelih', tone: 'netral' },
};

const PAYMENT_STATUS: Record<PaymentStatus, { label: string; tone: BadgeTone }> = {
  PENDING: { label: 'Belum lunas', tone: 'kuningan' },
  PAID: { label: 'Lunas', tone: 'sajadah' },
};

const authStore = useAuthStore();
const queryClient = useQueryClient();

const canManage = computed(() => authStore.can(Permission.QurbanManage));

const tabs: { id: TabId; label: string }[] = [
  { id: 'packages', label: 'Paket' },
  { id: 'bookings', label: 'Pendaftaran' },
  { id: 'animals', label: 'Hewan' },
];
const activeTab = ref<TabId>('packages');

const { data: packages, isLoading: isLoadingPackages } = useQuery({
  queryKey: queryKeys.qurban.packages,
  queryFn: () => qurbanService.getPackages(1, 100),
});

const { data: bookings, isLoading: isLoadingBookings } = useQuery({
  queryKey: queryKeys.qurban.bookings,
  queryFn: qurbanService.getBookings,
});

const { data: animals, isLoading: isLoadingAnimals } = useQuery({
  queryKey: queryKeys.qurban.animals,
  queryFn: qurbanService.getAnimals,
});

const packageItems = computed(() => packages.value?.items ?? []);
const availablePackages = computed(() => packageItems.value.filter(p => p.stock_remaining > 0));
const remainingStock = computed(() => packageItems.value.reduce((sum, p) => sum + p.stock_remaining, 0));
const paidBookings = computed(() => (bookings.value ?? []).filter(b => b.payment_status === 'PAID').length);

const refreshQurban = () => {
  queryClient.invalidateQueries({ queryKey: queryKeys.qurban.all });
  queryClient.invalidateQueries({ queryKey: queryKeys.dashboard.metrics });
};

// Package form
const showPackageForm = ref(false);
const editingPackage = ref<QurbanPackage | null>(null);
const emptyPackage: PackagePayload = { name: '', type: 'SAPI', price: 0, year_hijri: 1447, stock_total: 1 };

const {
  handleSubmit: handlePackageSubmit,
  errors: packageErrors,
  defineField: definePackageField,
  resetForm: resetPackageForm,
} = useForm({ validationSchema: toTypedSchema(packageSchema), initialValues: emptyPackage });

const [packageName] = definePackageField('name');
const [packageType] = definePackageField('type');
const [packagePrice] = definePackageField('price');
const [packageYear] = definePackageField('year_hijri');
const [packageStock] = definePackageField('stock_total');

const savePackageMutation = useMutation({
  mutationFn: (payload: PackagePayload) =>
    editingPackage.value ? qurbanService.updatePackage(editingPackage.value.id, payload) : qurbanService.createPackage(payload),
  onSuccess: (pkg) => {
    const wasEditing = editingPackage.value !== null;
    refreshQurban();
    showPackageForm.value = false;
    feedback.success(wasEditing ? 'Paket diperbarui' : 'Paket ditambahkan', `${pkg.name}, ${formatRupiah(pkg.price)}.`);
  },
  onError: (error: unknown) => feedback.error('Paket belum disimpan', error),
});

const openPackageForm = (pkg: QurbanPackage | null) => {
  editingPackage.value = pkg;
  resetPackageForm({
    values: pkg
      ? { name: pkg.name, type: pkg.type, price: pkg.price, year_hijri: pkg.year_hijri, stock_total: pkg.stock_total }
      : emptyPackage,
  });
  showPackageForm.value = true;
};

const submitPackage = handlePackageSubmit((values) => savePackageMutation.mutate(values));

const deletePackageMutation = useMutation({
  mutationFn: qurbanService.deletePackage,
  onSuccess: () => {
    refreshQurban();
    feedback.success('Paket dihapus');
  },
  onError: (error: unknown) => feedback.error('Paket belum dihapus', error),
});

const confirmDeletePackage = async (pkg: QurbanPackage) => {
  const confirmed = await feedback.confirm({
    title: 'Hapus paket?',
    message: `${pkg.name} tidak bisa dipilih lagi untuk pendaftaran baru. Pendaftaran yang sudah ada tetap tersimpan.`,
    confirmLabel: 'Hapus paket',
    danger: true,
  });
  if (confirmed) deletePackageMutation.mutate(pkg.id);
};

// Animal form
const showAnimalForm = ref(false);
const editingAnimal = ref<QurbanAnimal | null>(null);
const emptyAnimal: AnimalPayload = { tag_number: '', type: 'SAPI', weight: 0, status: 'AVAILABLE', vendor_info: '' };

const {
  handleSubmit: handleAnimalSubmit,
  errors: animalErrors,
  defineField: defineAnimalField,
  resetForm: resetAnimalForm,
} = useForm({ validationSchema: toTypedSchema(animalSchema), initialValues: emptyAnimal });

const [animalTag] = defineAnimalField('tag_number');
const [animalType] = defineAnimalField('type');
const [animalWeight] = defineAnimalField('weight');
const [animalStatus] = defineAnimalField('status');
const [animalVendor] = defineAnimalField('vendor_info');

const saveAnimalMutation = useMutation({
  mutationFn: (payload: AnimalPayload) =>
    editingAnimal.value ? qurbanService.updateAnimal(editingAnimal.value.id, payload) : qurbanService.createAnimal(payload),
  onSuccess: (animal) => {
    const wasEditing = editingAnimal.value !== null;
    refreshQurban();
    showAnimalForm.value = false;
    feedback.success(wasEditing ? 'Data hewan diperbarui' : 'Hewan ditambahkan', `Tag ${animal.tag_number}.`);
  },
  onError: (error: unknown) => feedback.error('Data hewan belum disimpan', error),
});

const openAnimalForm = (animal: QurbanAnimal | null) => {
  editingAnimal.value = animal;
  resetAnimalForm({
    values: animal
      ? { tag_number: animal.tag_number, type: animal.type, weight: animal.weight, status: animal.status, vendor_info: animal.vendor_info }
      : emptyAnimal,
  });
  showAnimalForm.value = true;
};

const submitAnimal = handleAnimalSubmit((values) => saveAnimalMutation.mutate(values));

const deleteAnimalMutation = useMutation({
  mutationFn: qurbanService.deleteAnimal,
  onSuccess: () => {
    refreshQurban();
    feedback.success('Data hewan dihapus');
  },
  onError: (error: unknown) => feedback.error('Data hewan belum dihapus', error),
});

const confirmDeleteAnimal = async (animal: QurbanAnimal) => {
  const confirmed = await feedback.confirm({
    title: 'Hapus data hewan?',
    message: `Hewan dengan tag ${animal.tag_number} akan dihapus dari daftar.`,
    confirmLabel: 'Hapus hewan',
    danger: true,
  });
  if (confirmed) deleteAnimalMutation.mutate(animal.id);
};

// Booking form
const showBookingForm = ref(false);
const emptyBooking: BookingPayload = { shohibul_id: '', package_id: 0, payment_status: 'PENDING' };

const {
  handleSubmit: handleBookingSubmit,
  errors: bookingErrors,
  defineField: defineBookingField,
  resetForm: resetBookingForm,
} = useForm({ validationSchema: toTypedSchema(bookingSchema), initialValues: emptyBooking });

const [shohibulId] = defineBookingField('shohibul_id');
const [bookingPackageId] = defineBookingField('package_id');
const [paymentStatus] = defineBookingField('payment_status');

const selectedPackagePrice = computed(
  () => packageItems.value.find(p => p.id === bookingPackageId.value)?.price ?? 0,
);

const createBookingMutation = useMutation({
  mutationFn: qurbanService.createBooking,
  onSuccess: (booking) => {
    refreshQurban();
    showBookingForm.value = false;
    feedback.success('Shohibul terdaftar', `${booking.shohibul_name} terdaftar untuk ${booking.package_name}.`);
  },
  onError: (error: unknown) => feedback.error('Shohibul belum terdaftar', error),
});

const openBookingForm = () => {
  resetBookingForm({ values: emptyBooking });
  showBookingForm.value = true;
};

const submitBooking = handleBookingSubmit((values) => createBookingMutation.mutate(values));
</script>

<template>
  <div class="space-y-8">
    <PageHeader title="Qurban" description="Paket qurban, pendaftaran shohibul, dan data hewan sampai hari penyembelihan.">
      <template v-if="canManage" #actions>
        <Button variant="outline" @click="openAnimalForm(null)"><Plus /> Tambah hewan</Button>
        <Button variant="outline" @click="openPackageForm(null)"><Plus /> Tambah paket</Button>
        <Button @click="openBookingForm"><UserPlus /> Daftarkan shohibul</Button>
      </template>
    </PageHeader>

    <div class="grid gap-4 sm:grid-cols-3">
      <StatCard label="Sisa kuota" :value="formatNumber(remainingStock)" :hint="`dari ${formatNumber(packageItems.length)} paket`" />
      <StatCard label="Shohibul" :value="formatNumber(bookings?.length ?? 0)" :hint="`${formatNumber(paidBookings)} sudah lunas`" />
      <StatCard label="Hewan" :value="formatNumber(animals?.length ?? 0)" hint="tercatat dengan nomor tag" />
    </div>

    <div class="inline-flex rounded-lg border border-batu bg-white p-1" role="tablist" aria-label="Data qurban">
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

    <!-- Packages -->
    <DataPanel v-if="activeTab === 'packages'">
      <div v-if="isLoadingPackages" class="p-8 text-center text-sm text-jati-muda">Memuat paket…</div>
      <EmptyState v-else-if="packageItems.length === 0" :icon="Beef" title="Belum ada paket qurban" description="Tambahkan paket beserta harga dan kuotanya sebelum membuka pendaftaran." />
      <div v-else class="scrollbar-tipis overflow-x-auto overscroll-x-contain">
        <table class="tabel min-w-[720px]">
          <thead>
            <tr>
              <th>Paket</th>
              <th>Jenis</th>
              <th class="text-right">Harga</th>
              <th>Tahun</th>
              <th class="text-right">Sisa kuota</th>
              <th v-if="canManage" class="w-24"><span class="sr-only">Aksi</span></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="pkg in packageItems" :key="pkg.id">
              <td class="font-medium">{{ pkg.name }}</td>
              <td class="text-jati-muda">{{ PACKAGE_TYPE_LABELS[pkg.type] }}</td>
              <td class="tabular text-right">{{ formatRupiah(pkg.price) }}</td>
              <td class="tabular text-jati-muda">{{ pkg.year_hijri }} H</td>
              <td class="text-right">
                <StatusBadge :tone="pkg.stock_remaining > 0 ? 'sajadah' : 'bata'" class="tabular">
                  {{ pkg.stock_remaining }} / {{ pkg.stock_total }}
                </StatusBadge>
              </td>
              <td v-if="canManage">
                <div class="flex justify-end gap-1">
                  <button type="button" class="tombol-ikon" :aria-label="`Ubah ${pkg.name}`" title="Ubah" @click="openPackageForm(pkg)">
                    <Pencil class="size-4" />
                  </button>
                  <button type="button" class="tombol-ikon tombol-ikon-bahaya" :aria-label="`Hapus ${pkg.name}`" title="Hapus" @click="confirmDeletePackage(pkg)">
                    <Trash2 class="size-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </DataPanel>

    <!-- Bookings -->
    <DataPanel v-else-if="activeTab === 'bookings'">
      <div v-if="isLoadingBookings" class="p-8 text-center text-sm text-jati-muda">Memuat pendaftaran…</div>
      <EmptyState v-else-if="!bookings || bookings.length === 0" :icon="UserPlus" title="Belum ada shohibul" description="Pendaftaran shohibul qurban akan tampil di sini." />
      <div v-else class="scrollbar-tipis overflow-x-auto overscroll-x-contain">
        <table class="tabel min-w-[680px]">
          <thead>
            <tr>
              <th>Tanggal</th>
              <th>Shohibul</th>
              <th>Paket</th>
              <th class="text-right">Total</th>
              <th>Pembayaran</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="b in bookings" :key="b.id">
              <td class="whitespace-nowrap text-jati-muda">{{ formatDate(b.booking_date) }}</td>
              <td class="font-medium">{{ b.shohibul_name }}</td>
              <td class="text-jati-muda">{{ b.package_name }}</td>
              <td class="tabular text-right">{{ formatRupiah(b.total_amount) }}</td>
              <td><StatusBadge :tone="PAYMENT_STATUS[b.payment_status].tone">{{ PAYMENT_STATUS[b.payment_status].label }}</StatusBadge></td>
            </tr>
          </tbody>
        </table>
      </div>
    </DataPanel>

    <!-- Animals -->
    <DataPanel v-else>
      <div v-if="isLoadingAnimals" class="p-8 text-center text-sm text-jati-muda">Memuat data hewan…</div>
      <EmptyState v-else-if="!animals || animals.length === 0" :icon="Beef" title="Belum ada hewan tercatat" description="Catat hewan dengan nomor tag saat hewan tiba dari penjual." />
      <div v-else class="scrollbar-tipis overflow-x-auto overscroll-x-contain">
        <table class="tabel min-w-[700px]">
          <thead>
            <tr>
              <th>No. tag</th>
              <th>Jenis</th>
              <th class="text-right">Berat</th>
              <th>Penjual</th>
              <th>Status</th>
              <th v-if="canManage" class="w-24"><span class="sr-only">Aksi</span></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="animal in animals" :key="animal.id">
              <td class="font-mono text-sm">{{ animal.tag_number }}</td>
              <td>{{ ANIMAL_TYPE_LABELS[animal.type] }}</td>
              <td class="tabular text-right">{{ formatNumber(animal.weight) }} kg</td>
              <td class="max-w-xs truncate text-jati-muda">{{ animal.vendor_info || '–' }}</td>
              <td><StatusBadge :tone="ANIMAL_STATUS[animal.status].tone">{{ ANIMAL_STATUS[animal.status].label }}</StatusBadge></td>
              <td v-if="canManage">
                <div class="flex justify-end gap-1">
                  <button type="button" class="tombol-ikon" :aria-label="`Ubah hewan ${animal.tag_number}`" title="Ubah" @click="openAnimalForm(animal)">
                    <Pencil class="size-4" />
                  </button>
                  <button type="button" class="tombol-ikon tombol-ikon-bahaya" :aria-label="`Hapus hewan ${animal.tag_number}`" title="Hapus" @click="confirmDeleteAnimal(animal)">
                    <Trash2 class="size-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </DataPanel>

    <!-- Package form -->
    <FormDialog
      v-model:open="showPackageForm"
      :title="editingPackage ? 'Ubah paket' : 'Tambah paket'"
      :submit-label="editingPackage ? 'Simpan perubahan' : 'Tambah paket'"
      :submitting="savePackageMutation.isPending.value"
      @submit="submitPackage"
    >
      <FormField label="Nama paket" for="package-name" :error="packageErrors.name" required>
        <input id="package-name" v-model="packageName" type="text" class="field-control" placeholder="Contoh: Kambing premium" :aria-invalid="Boolean(packageErrors.name)" />
      </FormField>
      <FormField label="Jenis" for="package-type" required>
        <select id="package-type" v-model="packageType" class="field-control">
          <option v-for="(label, value) in PACKAGE_TYPE_LABELS" :key="value" :value="value">{{ label }}</option>
        </select>
      </FormField>
      <div class="grid gap-4 sm:grid-cols-2">
        <FormField label="Harga (Rp)" for="package-price" :error="packageErrors.price" required>
          <input id="package-price" v-model.number="packagePrice" type="number" min="1" step="1" class="field-control tabular" :aria-invalid="Boolean(packageErrors.price)" />
        </FormField>
        <FormField label="Tahun hijriah" for="package-year" :error="packageErrors.year_hijri" required>
          <input id="package-year" v-model.number="packageYear" type="number" class="field-control tabular" :aria-invalid="Boolean(packageErrors.year_hijri)" />
        </FormField>
      </div>
      <FormField
        label="Kuota"
        for="package-stock"
        :error="packageErrors.stock_total"
        :hint="editingPackage ? 'Tidak boleh lebih kecil dari jumlah yang sudah terdaftar.' : 'Jumlah shohibul yang bisa mendaftar.'"
        required
      >
        <input id="package-stock" v-model.number="packageStock" type="number" min="1" step="1" class="field-control tabular" :aria-invalid="Boolean(packageErrors.stock_total)" />
      </FormField>
    </FormDialog>

    <!-- Animal form -->
    <FormDialog
      v-model:open="showAnimalForm"
      :title="editingAnimal ? 'Ubah data hewan' : 'Tambah hewan'"
      :submit-label="editingAnimal ? 'Simpan perubahan' : 'Tambah hewan'"
      :submitting="saveAnimalMutation.isPending.value"
      @submit="submitAnimal"
    >
      <FormField label="Nomor tag" for="animal-tag" :error="animalErrors.tag_number" required>
        <input id="animal-tag" v-model="animalTag" type="text" class="field-control font-mono" :aria-invalid="Boolean(animalErrors.tag_number)" />
      </FormField>
      <div class="grid gap-4 sm:grid-cols-2">
        <FormField label="Jenis" for="animal-type" required>
          <select id="animal-type" v-model="animalType" class="field-control">
            <option v-for="(label, value) in ANIMAL_TYPE_LABELS" :key="value" :value="value">{{ label }}</option>
          </select>
        </FormField>
        <FormField label="Berat (kg)" for="animal-weight" :error="animalErrors.weight" required>
          <input id="animal-weight" v-model.number="animalWeight" type="number" min="0" step="0.1" class="field-control tabular" :aria-invalid="Boolean(animalErrors.weight)" />
        </FormField>
      </div>
      <FormField label="Status" for="animal-status" required>
        <select id="animal-status" v-model="animalStatus" class="field-control">
          <option v-for="(status, value) in ANIMAL_STATUS" :key="value" :value="value">{{ status.label }}</option>
        </select>
      </FormField>
      <FormField label="Penjual" for="animal-vendor" :error="animalErrors.vendor_info">
        <input id="animal-vendor" v-model="animalVendor" type="text" class="field-control" placeholder="Nama dan kontak penjual" />
      </FormField>
    </FormDialog>

    <!-- Booking form -->
    <FormDialog
      v-model:open="showBookingForm"
      title="Daftarkan shohibul"
      description="Harga mengikuti paket yang dipilih, dan kuota paket berkurang satu."
      submit-label="Daftarkan"
      :submitting="createBookingMutation.isPending.value"
      @submit="submitBooking"
    >
      <FormField label="Shohibul" :error="bookingErrors.shohibul_id" required>
        <AsyncJamaahSelect v-model="shohibulId" placeholder="Cari nama jamaah…" />
      </FormField>
      <FormField label="Paket" for="booking-package" :error="bookingErrors.package_id" required>
        <select id="booking-package" v-model.number="bookingPackageId" class="field-control" :aria-invalid="Boolean(bookingErrors.package_id)">
          <option :value="0" disabled>Pilih paket</option>
          <option v-for="pkg in availablePackages" :key="pkg.id" :value="pkg.id">
            {{ pkg.name }}, {{ formatRupiah(pkg.price) }} (sisa {{ pkg.stock_remaining }})
          </option>
        </select>
      </FormField>
      <FormField label="Pembayaran" for="booking-payment" required>
        <select id="booking-payment" v-model="paymentStatus" class="field-control">
          <option v-for="(status, value) in PAYMENT_STATUS" :key="value" :value="value">{{ status.label }}</option>
        </select>
      </FormField>
      <div v-if="selectedPackagePrice > 0" class="flex items-center justify-between rounded-lg border border-batu bg-kapur/60 px-4 py-3 text-sm">
        <span class="text-jati-muda">Total</span>
        <span class="tabular font-semibold">{{ formatRupiah(selectedPackagePrice) }}</span>
      </div>
    </FormDialog>
  </div>
</template>
