<script setup lang="ts">
import { computed, type Component } from 'vue';
import { useRouter } from 'vue-router';
import { Wallet, HandCoins, Beef, PackageSearch, CalendarRange, Copy, MapPin, Phone, ArrowRight, LayoutDashboard, LogIn } from 'lucide-vue-next';
import heroBanner from '../../assets/hero-banner.jpg';
import { useAuthStore } from '../../stores/auth';
import { useMosqueProfileQuery } from '../../composables/useSettingsQuery';
import { feedback } from '../../composables/useFeedback';
import { ArchFrame, ARCH_PATH } from '../atoms/arch';
import BrandMark from '../molecules/BrandMark.vue';

interface Service {
  name: string;
  description: string;
  /** Who keeps this record, per the RBAC matrix. */
  keptBy: string;
  icon: Component;
}

const router = useRouter();
const authStore = useAuthStore();
const { data: profile } = useMosqueProfileQuery();

const currentYear = new Date().getFullYear();
const mosqueName = computed(() => profile.value?.name || 'Masjid kami');

const services: Service[] = [
  {
    name: 'Kas masjid',
    description: 'Setiap pemasukan dan pengeluaran dicatat permanen. Koreksi dilakukan dengan entri pembalik, bukan dihapus.',
    keptBy: 'Bendahara',
    icon: Wallet,
  },
  {
    name: 'Zakat & infaq',
    description: 'Penerimaan dari muzakki dan penyaluran ke mustahik, diprioritaskan dengan skor kebutuhan dari survei.',
    keptBy: 'Bendahara dan takmir',
    icon: HandCoins,
  },
  {
    name: 'Qurban',
    description: 'Paket, pendaftaran shohibul, dan data hewan dari kedatangan sampai penyembelihan.',
    keptBy: 'Takmir dan sekretaris',
    icon: Beef,
  },
  {
    name: 'Inventaris',
    description: 'Daftar aset masjid beserta kondisinya, dan peminjaman oleh jamaah sampai barang kembali.',
    keptBy: 'Takmir dan sekretaris',
    icon: PackageSearch,
  },
  {
    name: 'Agenda kegiatan',
    description: 'Jadwal kajian dan kegiatan, lengkap dengan dokumentasi foto untuk laporan pertanggungjawaban.',
    keptBy: 'Takmir dan sekretaris',
    icon: CalendarRange,
  },
];

const zakatFitrah = computed(() => {
  const amount = Number(profile.value?.zakat_fitrah_amount);
  if (!Number.isFinite(amount) || amount <= 0) return null;
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(amount);
});

const hasBankAccount = computed(() => Boolean(profile.value?.bank_account_number));

const goToPortal = () => {
  router.push(authStore.isAuthenticated ? { name: 'Dashboard' } : { name: 'Login' });
};

const copyAccountNumber = async () => {
  const number = profile.value?.bank_account_number;
  if (!number) return;
  try {
    await navigator.clipboard.writeText(number);
    feedback.success('Nomor rekening disalin', `${number} siap ditempel di aplikasi bank Anda.`);
  } catch {
    feedback.error('Nomor rekening belum disalin', 'Browser tidak mengizinkan akses clipboard. Salin nomornya secara manual.');
  }
};
</script>

<template>
  <div class="flex min-h-dvh flex-col bg-kapur text-jati">
    <!-- Navigation -->
    <header class="sticky top-0 z-40 border-b border-batu bg-kapur/85 backdrop-blur">
      <div class="mx-auto flex h-16 max-w-6xl items-center justify-between gap-3 px-4 sm:px-6">
        <RouterLink :to="{ name: 'Landing' }" class="rounded-md">
          <BrandMark />
        </RouterLink>
        <nav class="hidden items-center gap-8 text-sm font-medium text-jati-muda md:flex" aria-label="Navigasi halaman">
          <a href="#layanan" class="hover:text-sajadah">Layanan</a>
          <a href="#zakat" class="hover:text-sajadah">Zakat & infaq</a>
          <a href="#kontak" class="hover:text-sajadah">Kontak</a>
        </nav>
        <button
          type="button"
          class="inline-flex shrink-0 items-center gap-2 rounded-lg bg-sajadah px-3 py-2 text-sm font-semibold text-white transition-colors hover:bg-sajadah-tua sm:px-4"
          @click="goToPortal"
        >
          <component :is="authStore.isAuthenticated ? LayoutDashboard : LogIn" class="size-4" />
          <template v-if="authStore.isAuthenticated">Dashboard</template>
          <template v-else>Masuk<span class="hidden sm:inline">&nbsp;pengurus</span></template>
        </button>
      </div>
    </header>

    <main class="flex-1">
      <!-- Hero -->
      <section class="mx-auto grid max-w-6xl items-center gap-12 px-4 pt-12 pb-20 sm:px-6 lg:grid-cols-[1.1fr_0.9fr] lg:gap-16 lg:pt-20 lg:pb-28">
        <div>
          <p class="text-sm font-semibold uppercase tracking-[0.2em] text-kuningan">{{ mosqueName }}</p>
          <h1 class="mt-5 font-display text-4xl font-semibold leading-[1.05] tracking-wide text-balance text-jati min-[400px]:text-5xl sm:text-6xl lg:text-7xl">
            Amanah jamaah, tercatat rapi.
          </h1>
          <p class="mt-6 max-w-xl text-lg leading-relaxed text-jati-muda">
            Kas, zakat, qurban, inventaris, dan kegiatan {{ mosqueName }} dikelola pengurus dalam satu catatan yang
            terbuka, sehingga setiap rupiah dan setiap kegiatan dapat dipertanggungjawabkan.
          </p>
          <div class="mt-9 flex flex-col gap-3 sm:flex-row">
            <a
              href="#zakat"
              class="inline-flex items-center justify-center gap-2 rounded-lg bg-sajadah px-6 py-3 font-semibold text-white transition-colors hover:bg-sajadah-tua"
            >
              Salurkan zakat & infaq
              <ArrowRight class="size-4" />
            </a>
            <a
              href="#layanan"
              class="inline-flex items-center justify-center rounded-lg border border-batu-tua bg-white px-6 py-3 font-semibold text-jati transition-colors hover:border-sajadah hover:text-sajadah"
            >
              Lihat yang kami kelola
            </a>
          </div>
        </div>

        <!-- Signature: the mosque seen through a mihrab arch -->
        <figure class="mx-auto w-full max-w-sm lg:max-w-md">
          <div class="relative">
            <svg
              viewBox="0 0 40 50"
              class="absolute -top-4 -right-4 h-full w-full text-kuningan"
              aria-hidden="true"
            >
              <path :d="ARCH_PATH" fill="none" stroke="currentColor" stroke-width="1" vector-effect="non-scaling-stroke" />
            </svg>
            <ArchFrame :src="heroBanner" :alt="`Halaman ${mosqueName}`" class="relative w-full" />
          </div>
          <figcaption v-if="profile?.address" class="mt-4 flex items-start gap-2 text-sm text-jati-muda">
            <MapPin class="mt-0.5 size-4 shrink-0 text-kuningan" />
            {{ profile.address }}
          </figcaption>
        </figure>
      </section>

      <!-- Services -->
      <section id="layanan" class="border-t border-batu bg-white">
        <div class="mx-auto grid max-w-6xl gap-10 px-4 py-20 sm:px-6 lg:grid-cols-[0.8fr_1.2fr] lg:gap-16">
          <div class="lg:sticky lg:top-28 lg:self-start">
            <h2 class="font-display text-4xl font-semibold tracking-wide text-jati">Yang dikelola pengurus</h2>
            <p class="mt-4 max-w-sm leading-relaxed text-jati-muda">
              Setiap catatan dipegang oleh pengurus yang bertanggung jawab atasnya, dan tidak bisa diubah diam-diam.
            </p>
          </div>
          <ul class="shaf border-y border-batu">
            <li v-for="service in services" :key="service.name" class="grid gap-4 py-7 sm:grid-cols-[auto_1fr]">
              <span class="flex size-11 items-center justify-center rounded-full bg-sajadah-muda text-sajadah">
                <component :is="service.icon" class="size-5" />
              </span>
              <div>
                <div class="flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1">
                  <h3 class="text-lg font-semibold text-jati">{{ service.name }}</h3>
                  <p class="text-xs font-semibold uppercase tracking-[0.14em] text-kuningan">Dicatat oleh {{ service.keptBy }}</p>
                </div>
                <p class="mt-2 leading-relaxed text-jati-muda">{{ service.description }}</p>
              </div>
            </li>
          </ul>
        </div>
      </section>

      <!-- Zakat & infaq -->
      <section id="zakat" class="pola-girih bg-sajadah-tua text-white">
        <div class="mx-auto grid max-w-6xl items-center gap-10 px-4 py-20 sm:px-6 lg:grid-cols-2 lg:gap-16">
          <div>
            <h2 class="font-display text-4xl font-semibold tracking-wide">Zakat & infaq</h2>
            <p class="mt-4 max-w-md leading-relaxed text-white/75">
              Zakat dan infaq yang Anda titipkan tercatat sebagai penerimaan kas, lalu disalurkan kepada mustahik
              yang telah disurvei pengurus.
            </p>
            <p v-if="zakatFitrah" class="mt-8 text-sm text-white/60">
              Zakat fitrah<span v-if="profile?.active_hijri_year"> tahun {{ profile.active_hijri_year }} H</span>
            </p>
            <p v-if="zakatFitrah" class="tabular mt-1 font-display text-4xl font-semibold text-kuningan">
              {{ zakatFitrah }} <span class="font-sans text-base font-medium text-white/70">per jiwa</span>
            </p>
          </div>

          <div class="rounded-2xl bg-white p-7 text-jati shadow-lg">
            <template v-if="hasBankAccount && profile">
              <p class="text-xs font-semibold uppercase tracking-[0.16em] text-jati-muda">Rekening masjid</p>
              <p class="mt-3 text-lg font-semibold">{{ profile.bank_name }}</p>
              <p class="tabular mt-1 font-mono text-xl tracking-wider break-all text-sajadah-tua sm:text-2xl">{{ profile.bank_account_number }}</p>
              <p class="mt-1 text-sm text-jati-muda">a.n. {{ profile.bank_account_name }}</p>
              <button
                type="button"
                class="mt-6 inline-flex w-full items-center justify-center gap-2 rounded-lg bg-sajadah px-5 py-3 font-semibold text-white transition-colors hover:bg-sajadah-tua"
                @click="copyAccountNumber"
              >
                <Copy class="size-4" />
                Salin nomor rekening
              </button>
              <p class="mt-4 text-xs leading-relaxed text-jati-muda">
                Setelah transfer, kabari pengurus agar titipan Anda dicatat atas nama Anda.
              </p>
            </template>
            <template v-else>
              <p class="font-semibold">Rekening belum diumumkan</p>
              <p class="mt-2 text-sm leading-relaxed text-jati-muda">
                Untuk menitipkan zakat atau infaq, silakan hubungi pengurus melalui kontak di bawah, atau datang langsung
                ke sekretariat masjid.
              </p>
            </template>
          </div>
        </div>
      </section>
    </main>

    <!-- Contact -->
    <footer id="kontak" class="bg-jati text-white/70">
      <div class="mx-auto grid max-w-6xl gap-10 px-4 py-14 sm:px-6 md:grid-cols-[1.4fr_1fr_auto]">
        <div>
          <BrandMark variant="light" />
          <p class="mt-4 max-w-sm text-sm leading-relaxed">
            {{ mosqueName }} mengelola administrasinya dengan amami, aplikasi manajemen masjid.
          </p>
        </div>
        <div class="space-y-3 text-sm">
          <p class="text-xs font-semibold uppercase tracking-[0.16em] text-white/45">Kontak</p>
          <p v-if="profile?.address" class="flex items-start gap-2">
            <MapPin class="mt-0.5 size-4 shrink-0 text-kuningan" /> {{ profile.address }}
          </p>
          <p v-if="profile?.phone" class="flex items-center gap-2">
            <Phone class="size-4 shrink-0 text-kuningan" />
            <a :href="`tel:${profile.phone}`" class="hover:text-white">{{ profile.phone }}</a>
          </p>
        </div>
        <div class="text-sm">
          <p class="text-xs font-semibold uppercase tracking-[0.16em] text-white/45">Pengurus</p>
          <button type="button" class="mt-3 hover:text-white" @click="goToPortal">
            {{ authStore.isAuthenticated ? 'Buka dashboard' : 'Masuk ke dashboard' }}
          </button>
        </div>
      </div>
      <div class="border-t border-white/10">
        <p class="mx-auto max-w-6xl px-4 py-5 text-xs text-white/45 sm:px-6">
          &copy; {{ currentYear }} {{ mosqueName }}
        </p>
      </div>
    </footer>
  </div>
</template>
