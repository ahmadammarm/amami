<script setup lang="ts">
import { useRouter } from 'vue-router';
import { LogIn, Calendar, Info, ChevronRight, Activity, HandHeart, ShieldCheck, LayoutDashboard } from 'lucide-vue-next';
import heroBanner from '../../assets/hero-banner.jpg';
import { useAuthStore } from '../../stores/auth';

const router = useRouter();
const authStore = useAuthStore();
const currentYear = new Date().getFullYear();

// Static Data for initial version
const programs = [
  {
    title: 'Zakat Management',
    description: 'Transparent and precise collection and distribution of Zakat Fitrah and Zakat Maal.',
    icon: HandHeart,
    color: 'text-emerald-500'
  },
  {
    title: 'Qurban Operations',
    description: 'End-to-end handling of Eid al-Adha sacrifices, from donor booking to meat distribution.',
    icon: Activity,
    color: 'text-amber-500'
  },
  {
    title: 'Community Registry',
    description: 'Dynamic database of the Jamaah, categorizing volunteers, donors, and Mustahik.',
    icon: ShieldCheck,
    color: 'text-indigo-500'
  }
];

const newsArticles = [
  {
    id: 1,
    title: 'New Mosque Expansion Plan Approved',
    date: 'Jul 24, 2026',
    excerpt: 'The committee has finalized the blueprints for the new library and study hall extension.',
  },
  {
    id: 2,
    title: 'Record Zakat Collection this Ramadan',
    date: 'Jun 10, 2026',
    excerpt: 'Alhamdulillah, our community came together to support over 500 families this year.',
  },
  {
    id: 3,
    title: 'Youth Leadership Program Launch',
    date: 'May 15, 2026',
    excerpt: 'Register your children for the upcoming intensive leadership and religious studies camp.',
  }
];
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex flex-col font-sans">
    
    <!-- Navbar -->
    <header class="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-gray-100 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center gap-2">
            <span class="text-2xl font-bold text-primary">amami</span>
          </div>
          <nav class="hidden md:flex gap-8 text-sm font-medium text-gray-600">
            <a href="#history" class="hover:text-primary transition-colors">History</a>
            <a href="#programs" class="hover:text-primary transition-colors">Programs</a>
            <a href="#news" class="hover:text-primary transition-colors">News</a>
          </nav>
          <div class="flex items-center">
            <button @click="authStore.isAuthenticated ? router.push('/dashboard') : router.push('/login')" class="inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-semibold text-white bg-primary rounded-lg hover:bg-primary/90 transition-colors shadow-sm">
              <component :is="authStore.isAuthenticated ? LayoutDashboard : LogIn" class="w-4 h-4" />
              {{ authStore.isAuthenticated ? 'Dashboard' : 'Sign In' }}
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Hero Section -->
    <section class="relative bg-primary overflow-hidden">
      <!-- Decorative background elements -->
      <div class="absolute inset-0 opacity-10">
        <div class="absolute -top-24 -right-24 w-96 h-96 rounded-full bg-white blur-3xl"></div>
        <div class="absolute top-32 -left-20 w-72 h-72 rounded-full bg-white blur-3xl"></div>
      </div>
      
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 lg:py-32">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <!-- Left Column (Text & CTA) -->
          <div class="flex flex-col items-start text-left z-10">
            <h1 class="text-4xl md:text-5xl lg:text-6xl font-bold text-white tracking-tight mb-6">
              Digitizing Mosque Operations for a Better Community
            </h1>
            <p class="text-lg md:text-xl text-white/90 mb-10 leading-relaxed max-w-xl">
              Amami is a comprehensive Enterprise Resource Planning platform that transforms traditional mosque management into an auditable, data-driven framework.
            </p>
            <div class="flex flex-col sm:flex-row gap-4 w-full sm:w-auto">
              <button @click="authStore.isAuthenticated ? router.push('/dashboard') : router.push('/login')" class="px-8 py-3.5 bg-white text-primary font-bold rounded-xl hover:bg-gray-50 transition-colors shadow-lg shadow-primary/20 flex items-center justify-center gap-2 text-lg">
                {{ authStore.isAuthenticated ? 'Go to Dashboard' : 'Access Portal' }}
                <ChevronRight class="w-5 h-5" />
              </button>
              <a href="#history" class="px-8 py-3.5 bg-white/10 text-white font-semibold rounded-xl hover:bg-white/20 transition-colors backdrop-blur-sm border border-white/20 flex items-center justify-center text-lg">
                Learn More
              </a>
            </div>
          </div>
          
          <!-- Right Column (Image) -->
          <div class="relative z-10 w-full max-w-md lg:max-w-lg mx-auto rounded-2xl overflow-hidden shadow-2xl border border-white/10 transform hover:scale-[1.02] transition-transform duration-500">
            <div class="absolute inset-0 bg-gradient-to-tr from-primary/20 to-transparent mix-blend-overlay z-10 pointer-events-none"></div>
            <img :src="heroBanner" alt="Hero Banner" class="w-full aspect-square md:aspect-[4/3] object-cover rounded-2xl" />
          </div>
        </div>
      </div>
    </section>

    <!-- History Section -->
    <section id="history" class="py-20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-16 items-center">
          <div>
            <div class="inline-flex items-center justify-center p-3 bg-primary/10 rounded-xl mb-6 text-primary">
              <Info class="w-8 h-8" />
            </div>
            <h2 class="text-3xl font-bold text-gray-900 mb-6">Our History & Vision</h2>
            <p class="text-gray-600 mb-6 leading-relaxed text-lg">
              Established in the heart of the community, our mosque has served as a spiritual and social center for decades. As the community grew, so did the complexity of our operations.
            </p>
            <p class="text-gray-600 leading-relaxed text-lg mb-8">
              Amami was born out of the necessity to professionalize these operations. By enforcing strict Role-Based Access Control and maintaining an immutable financial ledger, we ensure absolute transparency and trust in managing community resources.
            </p>
          </div>
          <div class="relative">
            <!-- Placeholder for an image or illustration -->
            <div class="aspect-video bg-gray-100 rounded-2xl overflow-hidden shadow-xl border border-gray-200 flex items-center justify-center">
              <span class="text-gray-400 font-medium">Mosque Historic Image Placeholder</span>
            </div>
            <div class="absolute -bottom-6 -left-6 bg-white p-6 rounded-xl shadow-lg border border-gray-100 hidden md:block">
              <p class="text-4xl font-bold text-primary mb-1">20+</p>
              <p class="text-sm font-semibold text-gray-500">Years of Service</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Programs (Catalog) Section -->
    <section id="programs" class="py-20 bg-gray-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center max-w-3xl mx-auto mb-16">
          <h2 class="text-3xl font-bold text-gray-900 mb-4">Core Programs</h2>
          <p class="text-gray-600 text-lg">Explore the key modules that power our digital mosque management framework.</p>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div v-for="program in programs" :key="program.title" class="bg-white p-8 rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow group">
            <div class="w-14 h-14 bg-gray-50 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
              <component :is="program.icon" :class="['w-7 h-7', program.color]" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">{{ program.title }}</h3>
            <p class="text-gray-600 leading-relaxed">{{ program.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- News Section -->
    <section id="news" class="py-20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-end mb-12">
          <div>
            <h2 class="text-3xl font-bold text-gray-900 mb-4">Latest News</h2>
            <p class="text-gray-600 text-lg max-w-2xl">Stay updated with the latest announcements, articles, and community highlights.</p>
          </div>
          <button class="hidden sm:inline-flex items-center text-primary font-semibold hover:text-primary/80 transition-colors">
            View All News <ChevronRight class="w-4 h-4 ml-1" />
          </button>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <article v-for="news in newsArticles" :key="news.id" class="flex flex-col group cursor-pointer">
            <div class="aspect-video bg-gray-100 rounded-2xl mb-5 overflow-hidden">
              <div class="w-full h-full bg-gray-200 group-hover:bg-gray-300 transition-colors flex items-center justify-center text-gray-400">
                Image Placeholder
              </div>
            </div>
            <time class="text-sm font-semibold text-primary mb-2 flex items-center gap-2">
              <Calendar class="w-4 h-4" /> {{ news.date }}
            </time>
            <h3 class="text-xl font-bold text-gray-900 mb-3 group-hover:text-primary transition-colors">
              {{ news.title }}
            </h3>
            <p class="text-gray-600 line-clamp-3">
              {{ news.excerpt }}
            </p>
          </article>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="bg-gray-900 text-white py-12 mt-auto">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8 border-b border-gray-800 pb-8">
          <div class="col-span-1 md:col-span-2">
            <span class="text-2xl font-bold text-white mb-4 block">amami</span>
            <p class="text-gray-400 max-w-md">
              Aplikasi Manajemen Masjid. Empowering communities through transparent, efficient, and digital administration.
            </p>
          </div>
          <div>
            <h4 class="font-semibold mb-4 text-gray-100">Quick Links</h4>
            <ul class="space-y-2 text-gray-400 text-sm">
              <li><a href="#history" class="hover:text-white transition-colors">History</a></li>
              <li><a href="#programs" class="hover:text-white transition-colors">Programs</a></li>
              <li><a href="#news" class="hover:text-white transition-colors">News</a></li>
            </ul>
          </div>
          <div>
            <h4 class="font-semibold mb-4 text-gray-100">Portal</h4>
            <ul class="space-y-2 text-gray-400 text-sm">
              <li><button @click="router.push('/login')" class="hover:text-white transition-colors">Admin Login</button></li>
              <li><button @click="router.push('/login')" class="hover:text-white transition-colors">Jamaah Portal</button></li>
            </ul>
          </div>
        </div>
        <div class="flex flex-col md:flex-row justify-between items-center text-sm text-gray-500">
          <p>&copy; {{ currentYear }} amami. All rights reserved.</p>
          <div class="mt-4 md:mt-0 space-x-4">
            <a href="#" class="hover:text-gray-300">Privacy Policy</a>
            <a href="#" class="hover:text-gray-300">Terms of Service</a>
          </div>
        </div>
      </div>
    </footer>

  </div>
</template>
