export function randomHex(len: number): string {
  return [...Array(len)]
    .map(() => Math.floor(Math.random() * 16).toString(16))
    .join('');
}

export function randomBetween(a: number, b: number): number {
  const mx = Math.max(a, b);
  const mn = Math.min(a, b);
  const range = mx - mn;
  return Math.floor(Math.random() * range + mn);
}

export function pickKinN<T>(k: number, arr: T[]): T[] {
  const copiedArr = arr.slice();
  const shuffled = copiedArr.sort(() => 0.5 - Math.random());
  return shuffled.slice(0, k);
}

export async function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export function getHeaderMapFromHeaders(
  headers: string[],
): Record<string, string> {
  const result: Record<string, string> = {};
  headers.forEach((header) => {
    const delimiter = header.indexOf(': ');
    if (delimiter === -1) return;
    const key = header.slice(0, delimiter);
    const value = header.slice(delimiter + 2);
    if (key.toLowerCase() === 'content-length') return;
    result[key] = value;
  });
  return result;
}

export function isArray(value: any): boolean {
  return Array.isArray(value);
}

export function isNumber(value: any): boolean {
  return typeof value === 'number';
}

export function isString(value: any): boolean {
  return typeof value === 'string';
}

export function isObject(value: any): boolean {
  if (typeof value === 'object' && !Array.isArray(value) && value !== null) {
    return true;
  }
  return false;
}

export function getHeaders(headers: string[]): string {
  return headers.slice(1).reduce((before, header) => {
    if (header.toLowerCase().includes('content-length')) return before;
    return before + header + '\n';
  });
}
