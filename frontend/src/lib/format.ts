const rupiah = new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 });
const number = new Intl.NumberFormat('id-ID');
const date = new Intl.DateTimeFormat('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
const dateTime = new Intl.DateTimeFormat('id-ID', {
  day: 'numeric',
  month: 'short',
  year: 'numeric',
  hour: '2-digit',
  minute: '2-digit',
});

export const formatRupiah = (value: number): string => rupiah.format(value);

export const formatNumber = (value: number): string => number.format(value);

export const formatDate = (iso: string): string => date.format(new Date(iso));

export const formatDateTime = (iso: string): string => dateTime.format(new Date(iso));

/** Amount with its unit: rupiah for IDR, otherwise "2,5 KG". */
export const formatAmount = (value: number, unit: string): string =>
  unit === 'IDR' ? formatRupiah(value) : `${number.format(value)} ${unit}`;
