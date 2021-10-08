export function randomHex(len: number): string {
  return [...Array(len)].map(() => Math.floor(Math.random() * 16).toString(16)).join('');
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
