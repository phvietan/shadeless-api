import { AxiosResponse } from 'axios';
import { randomBetween } from 'libs/helper';
import ShadelessLogger from 'libs/logger/logger';

function diffLength(a: string, b: string): number {
  const maxLength = Math.max(a.length, b.length);
  const minLength = Math.min(a.length, b.length);
  return maxLength / minLength;
}

// Caution: Alot of magic here
function contentSimilarityScore(a: string, b: string): number {
  if ((a.length > 50 || b.length > 50) && diffLength(a, b) >= 1.5) return 0; // Length too different, must be different
  let score = 0;
  for (let i = 0; i < 100; ++i) {
    const start = randomBetween(0, Math.max(a.length - 30, 0));
    const end = randomBetween(7, 30);
    const subs = a.slice(start, start + end);
    score += +b.includes(subs);
  }
  return score;
}

export default class PathFuzzerFilter {
  THRESHOLD_STATUS_CODE = 0.7;

  logger: ShadelessLogger;
  response404: AxiosResponse<any>;
  responses: AxiosResponse<any>[];
  qualified: AxiosResponse<any>[];

  constructor(responses: AxiosResponse<any>[], logger: ShadelessLogger) {
    this.logger = logger;
    this.qualified = [];
    this.responses = responses.slice(1);
    this.response404 = responses[0];
  }

  private filter404() {
    this.logger.log(`Found ${this.responses.length} responses, filtering ...`);
    this.responses = this.responses.filter(
      (res) => res.status !== 404 && res.status !== 429,
    );
    this.logger.log(
      `After filter 404 status code, ${this.qualified.length}/${this.responses.length}`,
    );
  }

  private filterDominantStatusCode() {
    const count: Record<number, any> = {};
    let maxFreq = 0;
    let rememberStatusCode = -1;
    this.responses.forEach((res) => {
      count[res.status] = count[res.status] || 0;
      count[res.status] += 1;
      if (maxFreq < count[res.status]) {
        maxFreq = count[res.status];
        rememberStatusCode = res.status;
      }
    });
    if (maxFreq > this.responses.length * this.THRESHOLD_STATUS_CODE) {
      this.responses = this.responses
        .map((res) => {
          if (res.status !== rememberStatusCode) {
            this.qualified.push(res);
            return null;
          }
          return res;
        })
        .filter((res) => res !== null) as AxiosResponse<any>[];
    }
    this.logger.log(
      `After filter dominant status code, ${this.qualified.length}/${this.responses.length}`,
    );
  }

  private filterSimilar404() {
    this.responses.forEach((res) => {
      const score = contentSimilarityScore(res.data, this.response404.data);
      if (score < 50) this.qualified.push(res);
    });
    this.logger.log(`After filter 404 similarity, ${this.qualified.length}`);
  }

  private filterCaptcha() {
    // If we reach captcha, our response will now be filled with lots of captcha page, so this is last filter
    const threshold = Math.floor(0.33 * this.qualified.length);
    const checkbox = new Array(this.qualified.length).fill(true);
    for (let i = this.qualified.length - 1; i >= 0; --i) {
      let matched = 0;
      for (let j = 0; j < this.qualified.length; ++j) {
        if (i === j) continue;
        const score = contentSimilarityScore(
          this.qualified[j].data,
          this.qualified[i].data,
        );
        matched += +(score > 50);
      }
      if (matched > threshold) {
        checkbox[i] = false;
      }
    }
    this.qualified = this.qualified.filter((_, index) => checkbox[index]);
    this.logger.log(`After captcha, ${this.qualified.length}`);
  }

  async filter() {
    this.filter404();
    this.filterDominantStatusCode();
    this.filterSimilar404();
    this.filterCaptcha();
    this.logger.log(
      `Result is: ${this.qualified.map((q) => q.config.url).toString()}`,
    );
    return this.qualified;
  }
}
