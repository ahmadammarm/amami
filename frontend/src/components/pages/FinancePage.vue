<script setup lang="ts">
import { computed, ref } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { financeService } from '@/api/services/finance';
import { useAuthStore } from '@/stores/auth';
import type { CreateTransactionPayload, Transaction } from '@/types/finance';
import { toast } from 'vue-sonner';

const authStore = useAuthStore();
const queryClient = useQueryClient();

// Fetch Data
const { data: funds, isLoading: isLoadingFunds } = useQuery({
  queryKey: ['funds'],
  queryFn: financeService.getFunds,
});

// Pagination State
const currentPage = ref(1);
const pageSize = ref(10);

const { data: paginatedData, isLoading: isLoadingTransactions } = useQuery({
  queryKey: ['transactions', currentPage, pageSize],
  queryFn: () => financeService.getTransactions(currentPage.value, pageSize.value),
});

// Permissions
const canManage = computed(() => {
  return authStore.userRole === 'SUPER_ADMIN' || authStore.userRole === 'BENDAHARA';
});

// Form State
const showAddForm = ref(false);
const formState = ref<CreateTransactionPayload>({
  fund_id: 0,
  type: 'CREDIT',
  amount: 0,
  category: '',
  description: '',
});

// Detail Dialog State
const showDetailDialog = ref(false);
const selectedTransaction = ref<Transaction | null>(null);

const paginatedTransactions = computed(() => {
  return paginatedData.value?.data || [];
});

const totalPages = computed(() => {
  return paginatedData.value?.total_pages || 1;
});

const totalRecords = computed(() => {
  return paginatedData.value?.total || 0;
});

// Mutations
const createTxMutation = useMutation({
  mutationFn: financeService.createTransaction,
  onSuccess: () => {
    toast.success('Transaction added successfully!');
    queryClient.invalidateQueries({ queryKey: ['funds'] });
    queryClient.invalidateQueries({ queryKey: ['transactions'] });
    showAddForm.value = false;
    formState.value = { fund_id: 0, type: 'CREDIT', amount: 0, category: '', description: '' };
    showDetailDialog.value = false;
  },
  onError: (error: Error) => {
    import('axios').then(axios => {
      if (axios.default.isAxiosError(error)) {
        toast.error(error.response?.data?.error || 'Failed to process transaction');
      } else {
        toast.error('Failed to process transaction');
      }
    });
  },
});

const submitTransaction = () => {
  if (formState.value.fund_id === 0 || formState.value.amount <= 0 || !formState.value.category) {
    toast.error('Please fill in all required fields properly.');
    return;
  }
  createTxMutation.mutate(formState.value);
};

// Reversal Logic (Immutable Ledger rule)
const handleReverseTransaction = () => {
  if (!selectedTransaction.value) return;
  
  const tx = selectedTransaction.value;
  const confirmMessage = `WARNING: The amami system uses an Immutable Ledger. You cannot Edit or Delete transactions.\n\nAre you sure you want to create a Reversal Entry for ${formatCurrency(tx.amount)}?`;
  
  if (!window.confirm(confirmMessage)) return;

  const reversalPayload: CreateTransactionPayload = {
    fund_id: tx.fund_id,
    type: tx.type === 'CREDIT' ? 'DEBIT' : 'CREDIT',
    amount: tx.amount,
    category: 'CORRECTION / REVERSAL',
    description: `Reversal of transaction ID: ${tx.id}`,
    reference_id: tx.id
  };

  createTxMutation.mutate(reversalPayload);
};

// Actions
const openDetail = (tx: Transaction) => {
  selectedTransaction.value = tx;
  showDetailDialog.value = true;
};

const getFundName = (id: number) => {
  return funds.value?.find(f => f.id === id)?.name || `Fund #${id}`;
};

// Utils
const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(amount);
};

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-gray-900">Financial Ledger</h1>
        <p class="text-sm text-gray-500">Manage mosque funds and track immutable transactions.</p>
      </div>
      <button
        v-if="canManage"
        @click="showAddForm = !showAddForm"
        class="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none ring-offset-background bg-blue-600 text-white hover:bg-blue-700 h-10 py-2 px-4"
      >
        {{ showAddForm ? 'Close Form' : 'New Transaction' }}
      </button>
    </div>

    <!-- Fund Summaries -->
    <div v-if="isLoadingFunds" class="text-gray-500">Loading funds...</div>
    <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      <div v-for="fund in funds" :key="fund.id" class="rounded-xl border bg-white text-gray-900 shadow-sm p-6">
        <h3 class="text-sm font-medium text-gray-500">{{ fund.name }}</h3>
        <div class="mt-2 text-3xl font-bold">{{ formatCurrency(fund.current_balance) }}</div>
        <p class="text-xs text-gray-400 mt-1">Code: {{ fund.code }}</p>
      </div>
    </div>

    <!-- Add Transaction Form -->
    <div v-if="showAddForm && canManage" class="rounded-xl border bg-gray-50 p-6 shadow-sm">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Record New Transaction</h3>
      <form @submit.prevent="submitTransaction" class="space-y-4">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="block text-sm font-medium text-gray-700">Fund</label>
            <select v-model="formState.fund_id" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2 border bg-white">
              <option :value="0" disabled>Select a fund...</option>
              <option v-for="fund in funds" :key="fund.id" :value="fund.id">{{ fund.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Type</label>
            <select v-model="formState.type" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2 border bg-white">
              <option value="CREDIT">Inflow (Credit)</option>
              <option value="DEBIT">Outflow (Debit)</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Amount (IDR)</label>
            <input type="number" v-model.number="formState.amount" min="1" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2 border" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Category</label>
            <input type="text" v-model="formState.category" placeholder="e.g., Sedekah Jumat" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2 border" />
          </div>
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700">Description (Optional)</label>
            <input type="text" v-model="formState.description" placeholder="Notes about this transaction..." class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2 border" />
          </div>
        </div>
        <div class="flex justify-end">
          <button type="submit" :disabled="createTxMutation.isPending.value" class="inline-flex items-center justify-center rounded-md text-sm font-medium bg-blue-600 text-white hover:bg-blue-700 h-10 py-2 px-4 disabled:opacity-50">
            {{ createTxMutation.isPending.value ? 'Saving...' : 'Save Transaction' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Ledger Table with Pagination -->
    <div class="rounded-xl border bg-white shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b flex justify-between items-center">
        <h3 class="text-lg font-medium text-gray-900">Transaction History</h3>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-sm text-left">
          <thead class="text-xs text-gray-700 uppercase bg-gray-50 border-b">
            <tr>
              <th class="px-6 py-3">Date</th>
              <th class="px-6 py-3">Fund</th>
              <th class="px-6 py-3">Category</th>
              <th class="px-6 py-3">Type</th>
              <th class="px-6 py-3">Amount</th>
              <th class="px-6 py-3 text-right">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="isLoadingTransactions">
              <td colspan="6" class="px-6 py-4 text-center text-gray-500">Loading transactions...</td>
            </tr>
            <tr v-else-if="paginatedTransactions.length === 0">
              <td colspan="6" class="px-6 py-4 text-center text-gray-500">No transactions found.</td>
            </tr>
            <tr v-else v-for="tx in paginatedTransactions" :key="tx.id" class="border-b hover:bg-gray-50 transition-colors">
              <td class="px-6 py-4 text-gray-500 whitespace-nowrap">{{ formatDate(tx.created_at) }}</td>
              <td class="px-6 py-4 text-gray-700 font-medium">{{ getFundName(tx.fund_id) }}</td>
              <td class="px-6 py-4 font-medium text-gray-900">{{ tx.category }}</td>
              <td class="px-6 py-4">
                <span :class="['px-2 py-1 rounded-full text-xs font-medium', tx.type === 'CREDIT' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700']">
                  {{ tx.type }}
                </span>
              </td>
              <td class="px-6 py-4 font-mono font-medium">
                {{ formatCurrency(tx.amount) }}
              </td>
              <td class="px-6 py-4 text-right">
                <button @click="openDetail(tx)" class="text-blue-600 hover:text-blue-900 font-medium text-sm">View Details</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <!-- Pagination Controls -->
      <div class="px-6 py-4 border-t flex items-center justify-between text-sm text-gray-500 bg-gray-50">
        <div>
          Showing <span class="font-medium">{{ totalRecords === 0 ? 0 : (currentPage - 1) * pageSize + 1 }}</span> to 
          <span class="font-medium">{{ Math.min(currentPage * pageSize, totalRecords) }}</span> of 
          <span class="font-medium">{{ totalRecords }}</span> entries
        </div>
        <div class="flex space-x-2">
          <button @click="currentPage--" :disabled="currentPage === 1" class="px-3 py-1 border rounded hover:bg-gray-100 disabled:opacity-50">Previous</button>
          <button @click="currentPage++" :disabled="currentPage === totalPages" class="px-3 py-1 border rounded hover:bg-gray-100 disabled:opacity-50">Next</button>
        </div>
      </div>
    </div>
  </div>

  <!-- Detail Dialog Modal -->
  <div v-if="showDetailDialog && selectedTransaction" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
    <div class="bg-white rounded-xl shadow-xl w-full max-w-md overflow-hidden">
      <div class="px-6 py-4 border-b flex justify-between items-center bg-gray-50">
        <h3 class="text-lg font-bold text-gray-900">Transaction Details</h3>
        <button @click="showDetailDialog = false" class="text-gray-400 hover:text-gray-600">&times;</button>
      </div>
      <div class="p-6 space-y-4">
        <div>
          <p class="text-sm text-gray-500">Transaction ID</p>
          <p class="font-mono text-xs mt-1">{{ selectedTransaction.id }}</p>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <p class="text-sm text-gray-500">Date</p>
            <p class="font-medium">{{ formatDate(selectedTransaction.created_at) }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500">Fund</p>
            <p class="font-medium">{{ getFundName(selectedTransaction.fund_id) }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500">Type</p>
            <span :class="['inline-block px-2 py-1 rounded-full text-xs font-medium mt-1', selectedTransaction.type === 'CREDIT' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700']">
              {{ selectedTransaction.type }}
            </span>
          </div>
          <div>
            <p class="text-sm text-gray-500">Amount</p>
            <p class="font-mono font-bold">{{ formatCurrency(selectedTransaction.amount) }}</p>
          </div>
          <div class="col-span-2">
            <p class="text-sm text-gray-500">Category</p>
            <p class="font-medium">{{ selectedTransaction.category }}</p>
          </div>
          <div class="col-span-2" v-if="selectedTransaction.reference_id">
            <p class="text-sm text-gray-500">Reference / Corrects ID</p>
            <p class="font-mono text-xs">{{ selectedTransaction.reference_id }}</p>
          </div>
        </div>
      </div>
      
      <!-- Immutable Ledger Actions -->
      <div class="px-6 py-4 border-t bg-gray-50 flex justify-end space-x-3">
        <button @click="showDetailDialog = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border rounded-md hover:bg-gray-50">
          Close
        </button>
        <!-- Replace Edit/Delete with Reverse -->
        <button v-if="canManage" @click="handleReverseTransaction" :disabled="createTxMutation.isPending.value" class="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 disabled:opacity-50 transition-colors shadow-sm">
          Reverse Transaction
        </button>
      </div>
    </div>
  </div>
</template>
