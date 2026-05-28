<script setup lang="ts">
import { ref, watchEffect, computed } from 'vue';
import { 
  useMosqueProfileQuery, 
  useUpdateMosqueProfile,
  useSMTPConfigQuery,
  useUpdateSMTPMutation,
  useTestSMTPMutation,
  useSystemHealthQuery,
  useAuditLogsQuery
} from '../../composables/useSettingsQuery';
import { 
  Building2, Phone, MapPin, Fingerprint, Loader2, Edit3, Save, 
  Settings as SettingsIcon, Mail, ShieldCheck, Activity, History, 
  CheckCircle2, AlertCircle, Server, Database, Send, RefreshCcw, X,
  ChevronLeft, ChevronRight
} from 'lucide-vue-next';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import * as zod from 'zod';
import { toast } from 'vue-sonner';
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
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '../atoms/tabs';

// --- Data Fetching ---
const { data: profile, isLoading: isProfileLoading, error: profileError } = useMosqueProfileQuery();
const { data: smtpConfig, isLoading: isSMTPLoading } = useSMTPConfigQuery();
const { data: health, isLoading: isHealthLoading, refetch: refetchHealth } = useSystemHealthQuery();

// Audit Logs Pagination State
const auditPage = ref(1);
const auditLimit = ref(10);
const { data: auditData, isLoading: isAuditLoading } = useAuditLogsQuery(auditPage, auditLimit);
const auditLogs = computed(() => auditData.value?.logs || []);
const auditTotal = computed(() => auditData.value?.total_count || 0);
const totalPages = computed(() => Math.ceil(auditTotal.value / auditLimit.value));

const { mutate: updateProfile, isPending: isUpdatingProfile } = useUpdateMosqueProfile();
const { mutate: updateSMTP, isPending: isUpdatingSMTP } = useUpdateSMTPMutation();
const { mutate: testSMTP, isPending: isTestingSMTP } = useTestSMTPMutation();

// --- State ---
const activeTab = ref('profile');
const isEditing = ref(false);
const isDiscardDialogOpen = ref(false);
const successDialog = ref({
  open: false,
  title: '',
  description: '',
});

const handleReload = () => window.location.reload();

// --- Forms ---

// 1. Mosque Profile Form
const profileSchema = toTypedSchema(
  zod.object({
    name: zod.string().min(1, 'Mosque name is required'),
    address: zod.string().min(1, 'Address is required'),
    phone: zod.string().min(1, 'Phone number is required'),
    legal_yayasan_id: zod.string().min(1, 'Legal ID is required'),
  })
);

const profileForm = useForm({ 
  validationSchema: profileSchema,
  initialValues: {
    name: '',
    address: '',
    phone: '',
    legal_yayasan_id: '',
  }
});
const [pName] = profileForm.defineField('name');
const [pAddress] = profileForm.defineField('address');
const [pPhone] = profileForm.defineField('phone');
const [pLegalId] = profileForm.defineField('legal_yayasan_id');

// 2. SMTP Config Form
const smtpSchema = toTypedSchema(
  zod.object({
    host: zod.string().min(1, 'Host is required'),
    port: zod.string().min(1, 'Port is required'),
    username: zod.string().min(1, 'Username is required'),
    password: zod.string().min(1, 'Password is required'),
    from: zod.string().email('Invalid email address'),
  })
);

const smtpForm = useForm({ 
  validationSchema: smtpSchema,
  initialValues: {
    host: '',
    port: '',
    username: '',
    password: '',
    from: '',
  }
});
const [sHost] = smtpForm.defineField('host');
const [sPort] = smtpForm.defineField('port');
const [sUser] = smtpForm.defineField('username');
const [sPass] = smtpForm.defineField('password');
const [sFrom] = smtpForm.defineField('from');

// --- Logic ---

const startEdit = () => {
  isEditing.value = true;
  if (activeTab.value === 'profile' && profile.value) {
    profileForm.resetForm({ values: { ...profile.value } });
  } else if (activeTab.value === 'smtp' && smtpConfig.value) {
    smtpForm.resetForm({ values: { ...smtpConfig.value, password: '' } });
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

const onSMTPSubmit = smtpForm.handleSubmit((values) => {
  updateSMTP(values, { 
    onSuccess: () => {
      isEditing.value = false;
      successDialog.value = {
        open: true,
        title: 'SMTP Config Saved',
        description: 'Mail gateway configuration has been updated and is ready for use.',
      };
    } 
  });
});

const handleTestSMTP = smtpForm.handleSubmit((values) => {
  testSMTP(values, {
    onSuccess: () => {
      successDialog.value = {
        open: true,
        title: 'Connection Successful',
        description: 'Successfully established a handshake with the SMTP server.',
      };
    }
  });
});

const handleHealthCheck = async () => {
  const result = await refetchHealth();
  const currentHealth = result.data;
  
  if (currentHealth?.database && currentHealth?.smtp) {
    successDialog.value = {
      open: true,
      title: 'System All Green',
      description: 'All core infrastructure services are responding and healthy.',
    };
  } else if (!currentHealth?.database || !currentHealth?.smtp) {
    toast.error('Some services are currently unreachable. Check status badges for details.');
  }
};

const handleCancel = () => {
  const isDirty = activeTab.value === 'profile' ? profileForm.meta.value.dirty : smtpForm.meta.value.dirty;
  if (isDirty) {
    isDiscardDialogOpen.value = true;
  } else {
    isEditing.value = false;
  }
};

const confirmDiscard = () => {
  isEditing.value = false;
  isDiscardDialogOpen.value = false;
  // Revert values
  if (activeTab.value === 'profile' && profile.value) {
    profileForm.resetForm({ values: { ...profile.value } });
  } else if (activeTab.value === 'smtp' && smtpConfig.value) {
    smtpForm.resetForm({ values: { ...smtpConfig.value, password: '••••••••' } });
  }
};

// Sync forms when data arrives
watchEffect(() => {
  if (profile.value && !isEditing.value) {
    profileForm.setValues({ ...profile.value });
  }
  if (smtpConfig.value && !isEditing.value) {
    smtpForm.setValues({ ...smtpConfig.value, password: '••••••••' });
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
          System Settings
        </div>
        <h1 class="text-3xl font-bold text-gray-900">Workspace Settings</h1>
        <p class="text-gray-500 font-medium">Configure your mosque workspace and system parameters.</p>
      </div>
      
      <div class="flex items-center gap-3">
        <Button 
          v-if="!isEditing && (activeTab === 'profile' || activeTab === 'smtp')" 
          @click="startEdit"
          variant="outline"
          class="bg-white border-gray-200 hover:border-primary hover:text-primary transition-all shadow-sm rounded-xl h-11"
        >
          <Edit3 class="w-4 h-4 mr-2" />
          Edit {{ activeTab === 'profile' ? 'Profile' : 'Configuration' }}
        </Button>
      </div>
    </div>

    <!-- Main Tabs -->
    <Tabs v-model="activeTab" class="space-y-8" @update:model-value="isEditing = false">
      <TabsList class="bg-gray-100/50 p-1 rounded-xl h-12 inline-flex border border-gray-100 shadow-sm">
        <TabsTrigger value="profile" class="rounded-lg px-6 font-bold">
          <Building2 class="w-4 h-4 mr-2" />
          Mosque Profile
        </TabsTrigger>
        <TabsTrigger value="health" class="rounded-lg px-6 font-bold">
          <Activity class="w-4 h-4 mr-2" />
          System Health
        </TabsTrigger>
        <TabsTrigger value="smtp" class="rounded-lg px-6 font-bold">
          <Mail class="w-4 h-4 mr-2" />
          SMTP Config
        </TabsTrigger>
        <TabsTrigger value="audit" class="rounded-lg px-6 font-bold">
          <History class="w-4 h-4 mr-2" />
          Audit Logs
        </TabsTrigger>
      </TabsList>

      <!-- 1. Mosque Profile Content -->
      <TabsContent value="profile" class="mt-0 outline-none">
        <div v-if="profileError" class="p-12 bg-red-50 rounded-2xl border border-red-100 text-center max-w-2xl mx-auto">
          <div class="w-16 h-16 bg-red-100 text-red-600 rounded-full flex items-center justify-center mx-auto mb-4">
            <X class="w-8 h-8" />
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-2">Connection Error</h3>
          <p class="text-red-600 mb-6">We couldn't reach the server to load your settings. Please check your connection and try again.</p>
          <Button variant="outline" @click="handleReload" class="border-red-200 text-red-600 hover:bg-red-100">
            Retry Connection
          </Button>
        </div>
        <div v-else-if="isProfileLoading" class="py-20 text-center bg-white rounded-2xl border border-gray-100">
          <Loader2 class="w-10 h-10 animate-spin text-primary mx-auto mb-4" />
          <p class="text-gray-500">Loading profile...</p>
        </div>
        <div v-else class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <div class="px-8 py-6 border-b border-gray-50 bg-gray-50/30">
            <h2 class="text-xl font-bold text-gray-900 flex items-center gap-3">
              <div class="p-2 bg-primary/10 rounded-lg text-primary">
                <Building2 class="w-5 h-5" />
              </div>
              Mosque Information
            </h2>
          </div>
          <div class="p-8">
            <form @submit.prevent="onProfileSubmit" class="space-y-8">
              <div class="grid grid-cols-1 gap-8">
                <div class="space-y-2">
                  <label class="text-sm font-bold text-gray-700">Official Mosque Name</label>
                  <div v-if="!isEditing" class="p-4 bg-gray-50/50 rounded-xl border border-gray-100 font-medium flex items-center gap-3">
                    <Building2 class="w-5 h-5 text-primary/40" /> {{ profile?.name }}
                  </div>
                  <div v-else>
                    <Input v-model="pName" class="h-12 rounded-xl" placeholder="e.g. Masjid Agung Al-Hikmah" />
                    <p class="text-xs text-red-500 mt-1">{{ profileForm.errors.value.name }}</p>
                  </div>
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-bold text-gray-700">Physical Address</label>
                  <div v-if="!isEditing" class="p-4 bg-gray-50/50 rounded-xl border border-gray-100 font-medium flex items-start gap-3 leading-relaxed">
                    <MapPin class="w-5 h-5 text-primary/40 mt-1" /> {{ profile?.address }}
                  </div>
                  <div v-else>
                    <textarea v-model="pAddress" rows="3" class="w-full p-4 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary/20 outline-none text-sm transition-all" placeholder="Full address..."></textarea>
                    <p class="text-xs text-red-500 mt-1">{{ profileForm.errors.value.address }}</p>
                  </div>
                </div>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                  <div class="space-y-2">
                    <label class="text-sm font-bold text-gray-700">Contact Number</label>
                    <div v-if="!isEditing" class="p-4 bg-gray-50/50 rounded-xl border border-gray-100 font-medium flex items-center gap-3">
                      <Phone class="w-5 h-5 text-primary/40" /> {{ profile?.phone }}
                    </div>
                    <div v-else>
                      <Input v-model="pPhone" class="h-12 rounded-xl" placeholder="021-XXXXX" />
                      <p class="text-xs text-red-500 mt-1">{{ profileForm.errors.value.phone }}</p>
                    </div>
                  </div>
                  <div class="space-y-2">
                    <label class="text-sm font-bold text-gray-700">Legal Foundation ID</label>
                    <div v-if="!isEditing" class="p-4 bg-gray-50/50 rounded-xl border border-gray-100 font-medium flex items-center gap-3">
                      <Fingerprint class="w-5 h-5 text-primary/40" /> {{ profile?.legal_yayasan_id }}
                    </div>
                    <div v-else>
                      <Input v-model="pLegalId" class="h-12 rounded-xl" placeholder="Yayasan ID" />
                      <p class="text-xs text-red-500 mt-1">{{ profileForm.errors.value.legal_yayasan_id }}</p>
                    </div>
                  </div>
                </div>
              </div>
              <div v-if="isEditing" class="pt-8 flex justify-end gap-3 border-t border-gray-50">
                <Button type="button" variant="ghost" @click="handleCancel" :disabled="isUpdatingProfile" class="h-12 px-8 rounded-xl font-bold">Discard</Button>
                <Button type="submit" :disabled="isUpdatingProfile" class="h-12 px-8 rounded-xl font-bold shadow-lg shadow-primary/20">
                   <Loader2 v-if="isUpdatingProfile" class="w-4 h-4 mr-2 animate-spin" />
                   <Save v-else class="w-4 h-4 mr-2" /> Save Changes
                </Button>
              </div>
            </form>
          </div>
        </div>
      </TabsContent>

      <!-- 2. System Health Content -->
      <TabsContent value="health" class="mt-0 outline-none">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
          <div class="bg-white p-8 rounded-2xl border border-gray-100 shadow-sm space-y-6">
            <div class="flex items-center justify-between">
              <div class="p-3 bg-indigo-50 rounded-xl text-indigo-600"><Database class="w-6 h-6" /></div>
              <div :class="[health?.database ? 'bg-emerald-50 text-emerald-600' : 'bg-red-50 text-red-600', 'px-3 py-1 rounded-full text-xs font-bold flex items-center gap-1.5']">
                <CheckCircle2 v-if="health?.database" class="w-3.5 h-3.5" />
                <AlertCircle v-else class="w-3.5 h-3.5" />
                {{ health?.database ? 'CONNECTED' : 'DISCONNECTED' }}
              </div>
            </div>
            <div>
              <h3 class="text-lg font-bold text-gray-900">Database Engine</h3>
              <p class="text-sm text-gray-500 mt-1">PostgreSQL status and connection pool health.</p>
            </div>
            <Button variant="outline" class="w-full rounded-xl border-gray-200 h-11" @click="handleHealthCheck" :disabled="isHealthLoading">
              <RefreshCcw :class="['w-4 h-4 mr-2', isHealthLoading && 'animate-spin']" /> Check Connection
            </Button>
          </div>

          <div class="bg-white p-8 rounded-2xl border border-gray-100 shadow-sm space-y-6">
            <div class="flex items-center justify-between">
              <div class="p-3 bg-amber-50 rounded-xl text-amber-600"><Server class="w-6 h-6" /></div>
              <div :class="[health?.smtp ? 'bg-emerald-50 text-emerald-600' : 'bg-red-50 text-red-600', 'px-3 py-1 rounded-full text-xs font-bold flex items-center gap-1.5']">
                <CheckCircle2 v-if="health?.smtp" class="w-3.5 h-3.5" />
                <AlertCircle v-else class="w-3.5 h-3.5" />
                {{ health?.smtp ? 'READY' : 'UNREACHABLE' }}
              </div>
            </div>
            <div>
              <h3 class="text-lg font-bold text-gray-900">Mail Gateway</h3>
              <p class="text-sm text-gray-500 mt-1">Status of the configured SMTP server for outbound emails.</p>
            </div>
            <Button variant="outline" class="w-full rounded-xl border-gray-200 h-11" @click="activeTab = 'smtp'">
              <SettingsIcon class="w-4 h-4 mr-2" /> Configure Gateway
            </Button>
          </div>
        </div>
      </TabsContent>

      <!-- 3. SMTP Config Content -->
      <TabsContent value="smtp" class="mt-0 outline-none">
        <div v-if="isSMTPLoading" class="py-20 text-center bg-white rounded-2xl border border-gray-100">
          <Loader2 class="w-10 h-10 animate-spin text-primary mx-auto mb-4" />
          <p class="text-gray-500">Loading SMTP config...</p>
        </div>
        <div v-else class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <div class="px-8 py-6 border-b border-gray-50 bg-gray-50/30 flex items-center justify-between">
            <div>
              <h2 class="text-xl font-bold text-gray-900 flex items-center gap-3">
                <div class="p-2 bg-amber-100/50 rounded-lg text-amber-600"><Mail class="w-5 h-5" /></div>
                SMTP Configuration
              </h2>
              <p class="text-sm text-gray-500 mt-1 ml-12">Credentials for sending automated emails and receipts.</p>
            </div>
          </div>
          <div class="p-8">
            <form @submit.prevent="onSMTPSubmit" class="space-y-6">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-2">
                  <label class="text-sm font-bold text-gray-700">SMTP Host</label>
                  <Input v-model="sHost" :disabled="!isEditing" class="h-11 rounded-xl" placeholder="smtp.example.com" />
                  <p class="text-xs text-red-500 mt-1">{{ smtpForm.errors.value.host }}</p>
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-bold text-gray-700">Port</label>
                  <Input v-model="sPort" :disabled="!isEditing" class="h-11 rounded-xl" placeholder="587" />
                  <p class="text-xs text-red-500 mt-1">{{ smtpForm.errors.value.port }}</p>
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-bold text-gray-700">Username</label>
                  <Input v-model="sUser" :disabled="!isEditing" class="h-11 rounded-xl" placeholder="user@example.com" />
                  <p class="text-xs text-red-500 mt-1">{{ smtpForm.errors.value.username }}</p>
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-bold text-gray-700">Password</label>
                  <Input v-model="sPass" type="password" :disabled="!isEditing" class="h-11 rounded-xl" placeholder="••••••••" />
                  <p class="text-xs text-red-500 mt-1">{{ smtpForm.errors.value.password }}</p>
                </div>
                <div class="space-y-2 md:col-span-2">
                  <label class="text-sm font-bold text-gray-700">Sender (From Address)</label>
                  <Input v-model="sFrom" :disabled="!isEditing" class="h-11 rounded-xl" placeholder="no-reply@mosque.org" />
                  <p class="text-xs text-red-500 mt-1">{{ smtpForm.errors.value.from }}</p>
                </div>
              </div>
              <div v-if="isEditing" class="pt-6 flex flex-col md:flex-row justify-end gap-3 border-t border-gray-50">
                <Button type="button" variant="outline" @click="handleTestSMTP" :disabled="isTestingSMTP" class="h-11 px-6 rounded-xl border-amber-200 text-amber-600 hover:bg-amber-50">
                   <Loader2 v-if="isTestingSMTP" class="w-4 h-4 mr-2 animate-spin" />
                   <Send v-else class="w-4 h-4 mr-2" /> Test Handshake
                </Button>
                <div class="flex-1 md:hidden"></div>
                <Button type="button" variant="ghost" @click="handleCancel" :disabled="isUpdatingSMTP" class="h-11 px-8 rounded-xl font-bold">Discard</Button>
                <Button type="submit" :disabled="isUpdatingSMTP" class="h-11 px-8 rounded-xl font-bold shadow-lg shadow-primary/20">
                   <Loader2 v-if="isUpdatingSMTP" class="w-4 h-4 mr-2 animate-spin" />
                   <Save v-else class="w-4 h-4 mr-2" /> Save Config
                </Button>
              </div>
            </form>
          </div>
        </div>
      </TabsContent>

      <!-- 4. Audit Logs Content (Paginated) -->
      <TabsContent value="audit" class="mt-0 outline-none">
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden flex flex-col">
          <div class="px-8 py-6 border-b border-gray-50 bg-gray-50/30 flex items-center justify-between">
            <div>
              <h2 class="text-xl font-bold text-gray-900 flex items-center gap-3">
                <div class="p-2 bg-blue-100/50 rounded-lg text-blue-600"><History class="w-5 h-5" /></div>
                System Audit Logs
              </h2>
              <p class="text-sm text-gray-500 mt-1 ml-12">Total of {{ auditTotal }} activity records found.</p>
            </div>
            <div v-if="totalPages > 1" class="flex items-center gap-2">
              <Button 
                variant="outline" 
                size="icon" 
                class="h-8 w-8 rounded-lg" 
                :disabled="auditPage <= 1 || isAuditLoading"
                @click="auditPage--"
              >
                <ChevronLeft class="w-4 h-4" />
              </Button>
              <span class="text-xs font-bold text-gray-600 px-2">Page {{ auditPage }} of {{ totalPages }}</span>
              <Button 
                variant="outline" 
                size="icon" 
                class="h-8 w-8 rounded-lg" 
                :disabled="auditPage >= totalPages || isAuditLoading"
                @click="auditPage++"
              >
                <ChevronRight class="w-4 h-4" />
              </Button>
            </div>
          </div>
          <div class="p-0 overflow-x-auto">
             <table class="w-full text-left border-collapse min-w-[600px]">
               <thead>
                 <tr class="bg-gray-50/50 text-[10px] uppercase tracking-widest text-gray-400 font-bold border-b border-gray-50">
                   <th class="px-8 py-4">User</th>
                   <th class="px-8 py-4">Action</th>
                   <th class="px-8 py-4">Entity</th>
                   <th class="px-8 py-4 text-right">Timestamp</th>
                 </tr>
               </thead>
               <tbody class="divide-y divide-gray-50">
                 <template v-if="isAuditLoading">
                    <tr v-for="i in 5" :key="i">
                      <td colspan="4" class="px-8 py-4"><div class="h-4 bg-gray-100 animate-pulse rounded w-full"></div></td>
                    </tr>
                 </template>
                 <template v-else>
                   <tr v-for="log in auditLogs" :key="log.id" class="hover:bg-gray-50/30 transition-colors">
                     <td class="px-8 py-4 text-sm font-bold text-gray-700">{{ log.user }}</td>
                     <td class="px-8 py-4">
                       <span :class="[log.action === 'UPDATE' ? 'bg-indigo-50 text-indigo-600' : 'bg-gray-100 text-gray-600', 'px-2 py-0.5 rounded text-[10px] font-bold']">
                         {{ log.action }}
                       </span>
                     </td>
                     <td class="px-8 py-4 text-sm text-gray-500 font-medium">{{ log.entity }}</td>
                     <td class="px-8 py-4 text-xs text-gray-400 font-mono text-right">{{ log.timestamp }}</td>
                   </tr>
                   <tr v-if="!auditLogs?.length">
                     <td colspan="4" class="py-20 text-gray-400 italic text-center">No audit records found.</td>
                   </tr>
                 </template>
               </tbody>
             </table>
          </div>
          <!-- Bottom Pagination (Duplicate for UX) -->
          <div v-if="totalPages > 1" class="px-8 py-4 border-t border-gray-50 bg-gray-50/10 flex justify-between items-center">
             <p class="text-xs text-gray-400 font-medium">Showing {{ auditLogs.length }} of {{ auditTotal }} logs</p>
             <div class="flex items-center gap-2">
                <Button variant="ghost" size="sm" class="text-xs font-bold" :disabled="auditPage <= 1 || isAuditLoading" @click="auditPage = 1">First</Button>
                <Button variant="outline" size="sm" class="h-8 gap-1 rounded-lg px-3" :disabled="auditPage <= 1 || isAuditLoading" @click="auditPage--">
                   <ChevronLeft class="w-3.5 h-3.5" /> Previous
                </Button>
                <Button variant="outline" size="sm" class="h-8 gap-1 rounded-lg px-3" :disabled="auditPage >= totalPages || isAuditLoading" @click="auditPage++">
                   Next <ChevronRight class="w-3.5 h-3.5" />
                </Button>
             </div>
          </div>
        </div>
      </TabsContent>
    </Tabs>

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
          <AlertDialogTitle class="text-xl">Discard unsaved changes?</AlertDialogTitle>
          <AlertDialogDescription class="text-base leading-relaxed">
            You have modified settings. If you leave now, your progress will be lost.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter class="gap-3">
          <AlertDialogCancel class="rounded-xl h-11 font-bold">Back to Editing</AlertDialogCancel>
          <AlertDialogAction @click="confirmDiscard" class="bg-red-500 hover:bg-red-600 focus:ring-red-500 rounded-xl h-11 font-bold shadow-lg shadow-red-200">
            Discard Anyway
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>

    <!-- Privacy Note -->
    <div v-if="activeTab === 'profile'" class="p-6 bg-blue-50/50 rounded-2xl border border-blue-100 flex gap-4 transition-all hover:shadow-sm">
      <div class="p-2 bg-blue-100 rounded-xl text-blue-600 h-fit"><ShieldCheck class="w-5 h-5" /></div>
      <div class="space-y-1">
        <h4 class="text-sm font-bold text-blue-900">Security Note</h4>
        <p class="text-xs text-blue-700 leading-relaxed font-medium">
          Information updated here will be reflected across official mosque documents and digital receipts. 
        </p>
      </div>
    </div>
  </div>
</template>
