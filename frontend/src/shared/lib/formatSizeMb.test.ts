import { describe, expect, it } from 'vitest';

import { formatSizeBytes, formatSizeMb } from './formatSizeMb';

const formatNumber = (value: number) => new Intl.NumberFormat('en').format(value);

describe('formatSizeMb', () => {
  it('shows sizes below one gigabyte in megabytes', () => {
    expect(formatSizeMb(10.5, formatNumber)).toBe('10.5 MB');
  });

  it('rounds megabytes to two decimals', () => {
    expect(formatSizeMb(0.12345, formatNumber)).toBe('0.12 MB');
  });

  it('switches to gigabytes at exactly 1024 megabytes', () => {
    expect(formatSizeMb(1024, formatNumber)).toBe('1 GB');
  });

  it('rounds gigabytes to two decimals and groups digits', () => {
    expect(formatSizeMb(1024000, formatNumber)).toBe('1,000 GB');
  });

  it('switches to terabytes at exactly 1024 gigabytes', () => {
    expect(formatSizeMb(1024 * 1024 * 1.5, formatNumber)).toBe('1.5 TB');
  });

  it('treats a missing size as zero', () => {
    expect(formatSizeMb(undefined, formatNumber)).toBe('0 MB');
  });
});

describe('formatSizeBytes', () => {
  it('converts bytes to the same units as megabytes', () => {
    expect(formatSizeBytes(10.5 * 1024 * 1024, formatNumber)).toBe('10.5 MB');
  });

  it('shows a two-terabyte disk in terabytes', () => {
    expect(formatSizeBytes(2 * 1024 ** 4, formatNumber)).toBe('2 TB');
  });

  it('treats a missing size as zero', () => {
    expect(formatSizeBytes(undefined, formatNumber)).toBe('0 MB');
  });
});
