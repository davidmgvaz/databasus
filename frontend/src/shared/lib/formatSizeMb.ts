const MB_IN_GB = 1024;
const MB_IN_TB = 1024 * 1024;
const BYTES_IN_MB = 1024 * 1024;

export const formatSizeMb = (
  sizeMb: number | undefined,
  formatNumber: (value: number) => string,
): string => {
  const size = sizeMb ?? 0;

  if (size >= MB_IN_TB) {
    return `${formatNumber(Number((size / MB_IN_TB).toFixed(2)))} TB`;
  }

  if (size >= MB_IN_GB) {
    return `${formatNumber(Number((size / MB_IN_GB).toFixed(2)))} GB`;
  }

  return `${formatNumber(Number(size.toFixed(2)))} MB`;
};

export const formatSizeBytes = (
  sizeBytes: number | undefined,
  formatNumber: (value: number) => string,
): string => formatSizeMb((sizeBytes ?? 0) / BYTES_IN_MB, formatNumber);
