<script setup lang="ts">
import { ref, watchEffect } from 'vue';
import { 
  useMosqueProfileQuery, 
  useUpdateMosqueProfile
} from '../../composables/useSettingsQuery';
import { 
  Building2, Phone, MapPin, Fingerprint, Loader2, Edit3, Save, 
  Settings as SettingsIcon, ShieldCheck, X, Landmark, Wallet, Calendar
} from 'lucide-vue-next';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import * as zod from 'zod';
import Input from '../atoms/input/Input.vue';
import Button from '../atoms/button/Button.vue';
import SuccessDialog from '../atoms/SuccessDialog.vue';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '../atoms/alert-dialog';

// --- Data Fetching ---
const { data: profile, isLoading: isProfileLoading, error: profileError } = useMosqueProfileQuery();
const { mutate: updateProfile, isPending: isUpdatingProfile } = useUpdateMosqueProfile();

// --- State ---
const isEditing = ref(false);
const isDiscardDialogOpen = ref(false);
const successDialog = ref({
  open: false,
  title: '',
  description: '',
});

const handleReload = () => window.location.reload();

// --- Forms ---
const profileSchema = toTypedSchema(
  zod.object({
    name: zod.string().min(1, 'Mosque name is required'),
    address: zod.string().min(1, 'Address is required'),
    phone: zod.string().min(1, 'Phone number is required'),
    legal_yayasan_id: zod.string().min(1, 'Legal ID is required'),
    bank_name: zod.string().optional(),
    bank_account_name: zod.string().optional(),
    bank_account_number: zod.string().optional(),
    zakat_fitrah_amount: zod.string().optional(),
    active_hijri_year: zod.string().optional(),
  })
);

const profileForm = useForm({ 
  validationSchema: profileSchema,
  initialValues: {
    name: '',
    address: '',
    phone: '',
    legal_yayasan_id: '',
    bank_name: '',
    bank_account_name: '',
    bank_account_number: '',
    zakat_fitrah_amount: '',
    active_hijri_year: '',
  }
});
const [pName] = profileForm.defineField('name');
const [pAddress] = profileForm.defineField('address');
const [pPhone] = profileForm.defineField('phone');
const [pLegalId] = profileForm.defineField('legal_yayasan_id');
const [pBankName] = profileForm.defineField('bank_name');
const [pBankAccName] = profileForm.defineField('bank_account_name');
const [pBankAccNum] = profileForm.defineField('bank_account_number');
const [pZakatFitrah] = profileForm.defineField('zakat_fitrah_amount');
const [pHijriYear] = profileForm.defineField('active_hijri_year');

// --- Logic ---
const startEdit = () => {
  isEditing.value = true;
  if (profile.value) {
    profileForm.resetForm({ values: { ...profile.value } });
  }
};

const onProfileSubmit = profileForm.handleSubmit((values) => {
  updateProfile(values, { 
    onSuccess: () => {
      isEditing.value = false;
      successDialog.value = {
        open: true,
        title: 'Profile Updated',
        description: 'The mosque profile information has been successfully saved to the system.',
      };
    } 
  });
});

const handleCancel = () => {
  if (profileForm.meta.value.dirty) {
    isDiscardDialogOpen.value = true;
  } else {
    isEditing.value = false;
  }
};

const confirmDiscard = () => {
  isEditing.value = false;
  isDiscardDialogOpen.value = false;
  if (profile.value) {
    profileForm.resetForm({ values: { ...profile.value } });
  }
};

// Sync forms when data arrives
watchEffect(() => {
  if (profile.value && !isEditing.value) {
    profileForm.setValues({ ...profile.value });
  }
});
</script>

<template>
  <div class="max-w-6xl mx-auto space-y-8 pb-12">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
      <div class="space-y-1">
        <div class="flex items-center gap-2 text-primary font-semibold text-sm tracking-wider uppercase">
          <SettingsIcon class="w-4 h-4" />
          Konfigurasi Masjid
        </div>
        <h1 class="text-3xl font-bold text-gray-900">Profil & Pengaturan Masjid</h1>
        <p class="text-gray-500 font-medium">Kelola identitas, rekening, dan preferensi operasional masjid Anda.</p>
      </div>
      
      <div class="flex items-center gap-3">
        <Button 
          v-if="!isEditing" 
          @click="startEdit"
          variant="outline"
          class="bg-white border-gray-200 hover:border-primary hover:text-primary transition-all shadow-sm rounded-xl h-11"
        >
          <Edit3 class="w-4 h-4 mr-2" />
          Edit Profil Masjid
        </Button>
      </div>
    </div>

    <!-- Main Content -->
    <div v-if="profileError" class="p-12 bg-red-50 rounded-2xl border border-red-100 text-center max-w-2xl mx-auto mt-8">
      <div class="w-16 h-16 bg-red-100 text-red-600 rounded-full flex items-center justify-center mx-auto mb-4">
        <X class="w-8 h-8" />
      </div>
      <h3 class="text-lg font-bold text-gray-900 mb-2">Connection Error</h3>
      <p class="text-red-600 mb-6">We couldn't reach the server to load your settings. Please check your connection and try again.</p>
      <Button variant="outline" @click="handleReload" class="border-red-200 text-red-600 hover:bg-red-100">
        Retry Connection
      </Button>
    </div>
    <div v-else-if="isProfileLoading" class="py-20 text-center bg-white rounded-2xl border border-gray-100 mt-8">
      <Loader2 class="w-10 h-10 animate-spin text-primary mx-auto mb-4" />
      <p class="text-gray-500">Memuat profil masjid...</p>
    </div>
    <div v-else class="space-y-6 mt-8">
      <form @submit.prevent="onProfileSubmit" class="space-y-6">
        <!-- 1. Identitas Masjid -->
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <div class="px-8 py-6 border-b border-gray-50 bg-gray-50/30">
            <h2 class="text-xl font-bold text-gray-900 flex items-center gap-3">
              <div class="p-2 bg-primary/10 rounded-lg text-primary">
                <Building2 class="w-5 h-5" />
              </div>
              Identitas Masjid
            </h2>
          </div>
          <div class="p-8">
            <div class="grid grid-cols-1 gap-8">
              <div class="space-y-2">
                <label class="text-sm font-bold text-gray-700">Nama Resmi Masjid</label>
                <div v-if="!isEditing" class="p-4 bg-gray-50/50 rounded-xl border border-gray-100 font-medium flex items-center gap-3">
                  <Building2 class="w-5 h-5 text-primary/40" /> {{ profile?.name || '-' }}
                </div>
                <div v-else>
                  <Input v-model="pName" class="h-12 rounded-xl" placeholder="Contoh: Masjid Agung Al-Hikmah" />
                  <p class="text-xs text-red-500 mt-1">{{ profileForm.errors.value.name }}</p>
                </div>
              </div>
              <div class="space-y-2">
                <label class="text-sm font-bold text-gray-700">Alamat Lengkap</label>
                <div v-if="!isEditing" class="p-4 bg-gray-50/50 rounded-xl border border-gray-100 font-medium flex items-start gap-3 leading-relaxed">
                  <MapPin class="w-5 h-5 text-primary/40 mt-1" /> {{ profile?.address || '-' }}
                </div>
                <div v-else>
                  <textarea v-model="pAddress" rows="3" class="w-full p-4 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary/20 outline-none text-sm transition-all" placeholder="Alamat lengkap masjid..."></textarea>
                  <p class="text-xs text-red-500 mt-1">{{ profileForm.errors.value.address }}</p>
                </div>
              </div>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                <div class="space-y-2">
                  <label class="text-sm font-bold text-gray-700">Nomor Kontak / Telepon</label>
                  <div v-if="!isEditing" class="p-4 bg-gray-50/50 rounded-xl border border-gray-100 font-medium flex items-center gap-3">
                    <Phone class="w-5 h-5 text-primary/40" /> {{ profile?.phone || '-' }}
                  </div>
                  <div v-else>
                    <Input v-model="pPhone" class="h-12 rounded-xl" placeholder="021-XXXXX" />
                    <p class="text-xs text-red-500 mt-1">{{ profileForm.errors.value.phone }}</p>
                  </div>
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-bold text-gray-700">ID Legal Yayasan / Kemenag</label>
                  <div v-if="!isEditing" class="p-4 bg-gray-50/50 rounded-xl border border-gray-100 font-medium flex items-center gap-3">
                    <Fingerprint class="w-5 h-5 text-primary/40" /> {{ profile?.legal_yayasan_id || '-' }}
                  </div>
                  <div v-else>
                    <Input v-model="pLegalId" class="h-12 rounded-xl" placeholder="Nomor Yayasan" />
                    <p class="text-xs text-red-500 mt-1">{{ profileForm.errors.value.legal_yayasan_id }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 2. Pengaturan Finansial & Bank -->
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <div class="px-8 py-6 border-b border-gray-50 bg-gray-50/30">
            <h2 class="text-xl font-bold text-gray-900 flex items-center gap-3">
              <div class="p-2 bg-emerald-100/50 rounded-lg text-emerald-600">
                <Landmark class="w-5 h-5" />
              </div>
              Informasi Rekening Bank Masjid
            </h2>
            <p class="text-sm text-gray-500 mt-1 ml-12">Digunakan untuk keperluan transfer zakat, qurban, dan infaq jamaah.</p>
          </div>
          <div class="p-8 grid grid-cols-1 md:grid-cols-3 gap-8">
            <div class="space-y-2">
              <label class="text-sm font-bold text-gray-700">Nama Bank</label>
              <div v-if="!isEditing" class="p-4 bg-gray-50/50 rounded-xl border border-gray-100 font-medium">
                {{ profile?.bank_name || 'Belum diatur' }}
              </div>
              <div v-else>
                <Input v-model="pBankName" class="h-12 rounded-xl" placeholder="Contoh: BSI (Bank Syariah Indonesia)" />
              </div>
            </div>
            <div class="space-y-2">
              <label class="text-sm font-bold text-gray-700">Nomor Rekening</label>
              <div v-if="!isEditing" class="p-4 bg-gray-50/50 rounded-xl border border-gray-100 font-medium font-mono text-lg tracking-wider">
                {{ profile?.bank_account_number || 'Belum diatur' }}
              </div>
              <div v-else>
                <Input v-model="pBankAccNum" class="h-12 rounded-xl" placeholder="1234567890" />
              </div>
            </div>
            <div class="space-y-2">
              <label class="text-sm font-bold text-gray-700">Atas Nama (A/N)</label>
              <div v-if="!isEditing" class="p-4 bg-gray-50/50 rounded-xl border border-gray-100 font-medium">
                {{ profile?.bank_account_name || 'Belum diatur' }}
              </div>
              <div v-else>
                <Input v-model="pBankAccName" class="h-12 rounded-xl" placeholder="Masjid Agung Al-Hikmah" />
              </div>
            </div>
          </div>
        </div>

        <!-- 3. Preferensi Operasional -->
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <div class="px-8 py-6 border-b border-gray-50 bg-gray-50/30">
            <h2 class="text-xl font-bold text-gray-900 flex items-center gap-3">
              <div class="p-2 bg-amber-100/50 rounded-lg text-amber-600">
                <SettingsIcon class="w-5 h-5" />
              </div>
              Pengaturan Operasional Modul
            </h2>
          </div>
          <div class="p-8 grid grid-cols-1 md:grid-cols-2 gap-8">
            <div class="space-y-2">
              <label class="text-sm font-bold text-gray-700">Besaran Zakat Fitrah (Per Jiwa)</label>
              <p class="text-xs text-gray-500 mb-2">Nilai ekuivalen 2.5kg atau 3.5 liter beras dalam Rupiah.</p>
              <div v-if="!isEditing" class="p-4 bg-gray-50/50 rounded-xl border border-gray-100 font-medium flex items-center gap-3 text-lg">
                <Wallet class="w-5 h-5 text-emerald-600" />
                Rp {{ profile?.zakat_fitrah_amount ? parseInt(profile.zakat_fitrah_amount).toLocaleString('id-ID') : '40.000' }}
              </div>
              <div v-else class="relative">
                <span class="absolute left-4 top-1/2 -translate-y-1/2 font-bold text-gray-500">Rp</span>
                <Input v-model="pZakatFitrah" type="number" class="h-12 rounded-xl pl-10" placeholder="45000" />
              </div>
            </div>
            <div class="space-y-2">
              <label class="text-sm font-bold text-gray-700">Tahun Hijriah Aktif</label>
              <p class="text-xs text-gray-500 mb-2">Digunakan sebagai referensi default untuk Qurban & Agenda.</p>
              <div v-if="!isEditing" class="p-4 bg-gray-50/50 rounded-xl border border-gray-100 font-medium flex items-center gap-3 text-lg">
                <Calendar class="w-5 h-5 text-amber-600" />
                {{ profile?.active_hijri_year || '1447' }} H
              </div>
              <div v-else class="relative">
                <Input v-model="pHijriYear" type="number" class="h-12 rounded-xl" placeholder="1447" />
                <span class="absolute right-4 top-1/2 -translate-y-1/2 font-bold text-gray-500">H</span>
              </div>
            </div>
          </div>
        </div>

        <div v-if="isEditing" class="pt-4 flex justify-end gap-3 sticky bottom-4">
          <div class="bg-white p-4 rounded-2xl shadow-xl border border-gray-100 flex gap-3 w-full justify-end max-w-md">
            <Button type="button" variant="ghost" @click="handleCancel" :disabled="isUpdatingProfile" class="h-12 px-8 rounded-xl font-bold bg-gray-100 hover:bg-gray-200">Batal</Button>
            <Button type="submit" :disabled="isUpdatingProfile" class="h-12 px-8 rounded-xl font-bold shadow-lg shadow-primary/20 w-full sm:w-auto">
                <Loader2 v-if="isUpdatingProfile" class="w-4 h-4 mr-2 animate-spin" />
                <Save v-else class="w-4 h-4 mr-2" /> Simpan Perubahan
            </Button>
          </div>
        </div>
      </form>
    </div>

    <!-- Success Dialog -->
    <SuccessDialog
      v-model:open="successDialog.open"
      :title="successDialog.title"
      :description="successDialog.description"
    />

    <!-- Safety Dialog -->
    <AlertDialog v-model:open="isDiscardDialogOpen">
      <AlertDialogContent class="rounded-2xl">
        <AlertDialogHeader>
          <AlertDialogTitle class="text-xl">Batalkan Perubahan?</AlertDialogTitle>
          <AlertDialogDescription class="text-base leading-relaxed">
            Anda memiliki perubahan yang belum disimpan. Yakin ingin membatalkannya?
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter class="gap-3">
          <AlertDialogCancel class="rounded-xl h-11 font-bold">Kembali Edit</AlertDialogCancel>
          <AlertDialogAction @click="confirmDiscard" class="bg-red-500 hover:bg-red-600 focus:ring-red-500 rounded-xl h-11 font-bold shadow-lg shadow-red-200">
            Ya, Batalkan
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>

    <!-- Privacy Note -->
    <div v-if="!isProfileLoading && !profileError" class="p-6 bg-blue-50/50 rounded-2xl border border-blue-100 flex gap-4 transition-all hover:shadow-sm">
      <div class="p-2 bg-blue-100 rounded-xl text-blue-600 h-fit"><ShieldCheck class="w-5 h-5" /></div>
      <div class="space-y-1">
        <h4 class="text-sm font-bold text-blue-900">Catatan Sistem</h4>
        <p class="text-xs text-blue-700 leading-relaxed font-medium">
          Informasi profil dan rekening yang diatur di sini akan tercetak secara otomatis pada kuitansi Zakat, Infaq, dan Pendaftaran Qurban jamaah. Pastikan nomor rekening dan nominal Zakat Fitrah sudah sesuai dengan ketetapan BAZNAS/Takmir setempat.
        </p>
      </div>
    </div>
  </div>
</template>
