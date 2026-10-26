<script setup lang="ts">
import { computed, ref } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { financeService } from '@/api/services/finance';
import { getErrorMessage } from '@/api/client';
import { queryKeys } from '@/api/queryKeys';
import { useAuthStore } from '@/stores/auth';
import { Permission } from '@/lib/permissions';
import { reversalSchema, transactionSchema } from '@/schemas/finance';
import type { Transaction } from '@/types/finance';
import { toast } from 'vue-sonner';

const authStore = useAuthStore();
const queryClient = useQueryClient();

// Fetch Data
const { data: funds, isLoading: isLoadingFunds } = useQuery({
  queryKey: queryKeys.finance.funds,
  queryFn: financeService.getFunds,
});

// Pagination State
const currentPage = ref(1);
const pageSize = ref(10);

const { data: paginatedData, isLoading: isLoadingTransactions } = useQuery({
  queryKey: computed(() => queryKeys.finance.transactions(currentPage.value, pageSize.value)),
  queryFn: () => financeService.getTransactions(currentPage.value, pageSize.value),
});

// Permissions
const canWrite = computed(() => authStore.can(Permission.LedgerWrite));

// Transaction Form
const showAddForm = ref(false);
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

// Detail & Reversal State
const showDetailDialog = ref(false);
const selectedTransaction = ref<Transaction | null>(null);
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

const paginatedTransactions = computed(() => paginatedData.value?.items ?? []);
const totalPages = computed(() => paginatedData.value?.total_pages || 1);
const totalRecords = computed(() => paginatedData.value?.total ?? 0);

/** A transaction can be reversed once, and reversal entries themselves cannot be reversed. */
const canReverseSelected = computed(() => {
  const tx = selectedTransaction.value;
  return canWrite.value && tx !== null && tx.reversal_of_id === null && tx.reversed_by_id === null;
});

const refreshLedger = () => {
  queryClient.invalidateQueries({ queryKey: queryKeys.finance.all });
  queryClient.invalidateQueries({ queryKey: queryKeys.dashboard.metrics });
};

// Mutations
const createTxMutation = useMutation({
  mutationFn: financeService.createTransaction,
  onSuccess: () => {
    toast.success('Transaction added successfully!');
    refreshLedger();
    showAddForm.value = false;
    resetTransactionForm();
  },
  onError: (error: unknown) => {
    toast.error(getErrorMessage(error, 'Failed to process transaction'));
  },
});

const reverseTxMutation = useMutation({
  mutationFn: ({ id, reason }: { id: string; reason: string }) => financeService.reverseTransaction(id, { reason }),
  onSuccess: () => {
    toast.success('Reversal entry recorded');
    refreshLedger();
    closeDetail();
  },
  onError: (error: unknown) => {
    toast.error(getErrorMessage(error, 'Failed to reverse transaction'));
  },
});

const submitTransaction = handleTransactionSubmit((values) => {
  createTxMutation.mutate(values);
});

// Reversal Logic (Immutable Ledger rule): corrections are new offsetting entries.
const submitReversal = handleReversalSubmit((values) => {
  if (!selectedTransaction.value) return;
  reverseTxMutation.mutate({ id: selectedTransaction.value.id, reason: values.reason });
});

// Actions
const openDetail = (tx: Transaction) => {
  selectedTransaction.value = tx;
  showReversalForm.value = false;
  resetReversalForm();
  showDetailDialog.value = true;
};

const closeDetail = () => {
  showDetailDialog.value = false;
  showReversalForm.value = false;
  selectedTransaction.value = null;
};

const getFundName = (id: number) => {
  return funds.value?.find(f => f.id === id)?.name || `Fund #${id}`;
};

// Utils
const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(value);
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
        v-if="canWrite"
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
    <div v-if="showAddForm && canWrite" class="rounded-xl border bg-gray-50 p-6 shadow-sm">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Record New Transaction</h3>
      <form @submit="submitTransaction" class="space-y-4" novalidate>
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="block text-sm font-medium text-gray-700">Fund</label>
            <select v-model.number="fundId" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2 border bg-white">
              <option :value="0" disabled>Select a fund...</option>
              <option v-for="fund in funds" :key="fund.id" :value="fund.id">{{ fund.name }}</option>
            </select>
            <p v-if="transactionErrors.fund_id" class="mt-1 text-xs text-red-600">{{ transactionErrors.fund_id }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Type</label>
            <select v-model="txType" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2 border bg-white">
              <option value="CREDIT">Inflow (Credit)</option>
              <option value="DEBIT">Outflow (Debit)</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Amount (IDR)</label>
            <input type="number" v-model.number="amount" min="1" step="1" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2 border" />
            <p v-if="transactionErrors.amount" class="mt-1 text-xs text-red-600">{{ transactionErrors.amount }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Category</label>
            <input type="text" v-model="category" placeholder="e.g., Sedekah Jumat" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2 border" />
            <p v-if="transactionErrors.category" class="mt-1 text-xs text-red-600">{{ transactionErrors.category }}</p>
          </div>
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700">Description (Optional)</label>
            <input type="text" v-model="description" placeholder="Notes about this transaction..." class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2 border" />
            <p v-if="transactionErrors.description" class="mt-1 text-xs text-red-600">{{ transactionErrors.description }}</p>
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
              <td class="px-6 py-4 font-medium text-gray-900">
                {{ tx.category }}
                <span v-if="tx.reversed_by_id" class="ml-2 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-gray-200 text-gray-600">Reversed</span>
                <span v-if="tx.reversal_of_id" class="ml-2 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-amber-100 text-amber-700">Reversal</span>
              </td>
              <td class="px-6 py-4">
                <span :class="['px-2 py-1 rounded-full text-xs font-medium', tx.type === 'CREDIT' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700']">
                  {{ tx.type }}
                </span>
              </td>
              <td class="px-6 py-4 font-mono font-medium" :class="{ 'line-through text-gray-400': tx.reversed_by_id }">
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
          <button @click="currentPage++" :disabled="currentPage >= totalPages" class="px-3 py-1 border rounded hover:bg-gray-100 disabled:opacity-50">Next</button>
        </div>
      </div>
    </div>
  </div>

  <!-- Detail Dialog Modal -->
  <div v-if="showDetailDialog && selectedTransaction" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
    <div class="bg-white rounded-xl shadow-xl w-full max-w-md overflow-hidden">
      <div class="px-6 py-4 border-b flex justify-between items-center bg-gray-50">
        <h3 class="text-lg font-bold text-gray-900">Transaction Details</h3>
        <button @click="closeDetail" class="text-gray-400 hover:text-gray-600">&times;</button>
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
          <div class="col-span-2" v-if="selectedTransaction.description">
            <p class="text-sm text-gray-500">Description</p>
            <p class="font-medium">{{ selectedTransaction.description }}</p>
          </div>
          <div class="col-span-2" v-if="selectedTransaction.reversal_of_id">
            <p class="text-sm text-gray-500">Reverses Transaction</p>
            <p class="font-mono text-xs">{{ selectedTransaction.reversal_of_id }}</p>
          </div>
          <div class="col-span-2" v-if="selectedTransaction.reversed_by_id">
            <p class="text-sm text-gray-500">Reversed By</p>
            <p class="font-mono text-xs">{{ selectedTransaction.reversed_by_id }}</p>
          </div>
          <div class="col-span-2" v-if="selectedTransaction.reference_id">
            <p class="text-sm text-gray-500">Reference ID</p>
            <p class="font-mono text-xs">{{ selectedTransaction.reference_id }}</p>
          </div>
        </div>

        <!-- Reversal Entry Form -->
        <form v-if="showReversalForm" @submit="submitReversal" class="rounded-lg border border-red-200 bg-red-50 p-4 space-y-3" novalidate>
          <p class="text-sm text-red-800">
            The ledger is immutable. A reversal posts an opposite entry of
            <span class="font-semibold">{{ formatCurrency(selectedTransaction.amount) }}</span> to the same fund.
          </p>
          <div>
            <label class="block text-sm font-medium text-gray-700">Reason</label>
            <textarea v-model="reason" rows="2" placeholder="e.g., Duplicate entry" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm sm:text-sm p-2 border bg-white"></textarea>
            <p v-if="reversalErrors.reason" class="mt-1 text-xs text-red-600">{{ reversalErrors.reason }}</p>
          </div>
          <div class="flex justify-end gap-2">
            <button type="button" @click="showReversalForm = false" class="px-3 py-1.5 text-sm font-medium text-gray-700 bg-white border rounded-md hover:bg-gray-50">Cancel</button>
            <button type="submit" :disabled="reverseTxMutation.isPending.value" class="px-3 py-1.5 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 disabled:opacity-50">
              {{ reverseTxMutation.isPending.value ? 'Posting...' : 'Confirm Reversal' }}
            </button>
          </div>
        </form>
      </div>

      <!-- Immutable Ledger Actions -->
      <div class="px-6 py-4 border-t bg-gray-50 flex justify-end space-x-3">
        <button @click="closeDetail" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border rounded-md hover:bg-gray-50">
          Close
        </button>
        <!-- Replace Edit/Delete with Reverse -->
        <button v-if="canReverseSelected && !showReversalForm" @click="showReversalForm = true" class="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 disabled:opacity-50 transition-colors shadow-sm">
          Reverse Transaction
        </button>
      </div>
    </div>
  </div>
</template>
