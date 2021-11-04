import ExpressFastifyOpenRedirect from './poc/express-fastify-open-redirect';
import FastifyDOSCVE202122964 from './poc/express-fastify-dos';
import { ApiFuzzer } from './api-fuzzer-poc-generic';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';
import { ParsedPacket } from 'libs/databases/parsedPacket.database';
import SQLiTimeBased from './poc/sqli-timebased';
import * as Bluebird from 'bluebird';
import { sleep } from 'libs/helper';
import * as fs from 'fs/promises';
export default class ApiFuzzerSender {
  options: BotFuzzer;
  private pocs: ApiFuzzer[];
  private pocsName: string[];

  constructor(options: BotFuzzer, packet: ParsedPacket) {
    this.options = options;
    this.pocs = [
      new ExpressFastifyOpenRedirect(options, packet),
      new FastifyDOSCVE202122964(options, packet),
      new SQLiTimeBased(options, packet),
    ];
    this.pocsName = this.pocs.map((v) => v.constructor.name);
  }

  async sendPocs(): Promise<[string, string[]]> {
    let logs = '';
    const result = await Bluebird.map(
      this.pocs,
      async (poc) => {
        if (poc.condition && !poc.condition()) return 0;
        const res = await poc.poc();
        const ok = (await poc.detect(res)) >= 0.5;
        const currentLog = await fs.readFile(poc.logDir, 'utf-8');
        logs += '\n' + currentLog;
        await sleep(this.options.sleepRequest);
        return ok;
      },
      { concurrency: 1 },
    );

    return [
      logs,
      result
        .map((res, idx) => {
          if (!res) return null;
          return this.pocsName[idx];
        })
        .filter((r) => r !== null),
    ];
  }
}
