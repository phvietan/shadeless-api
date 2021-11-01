import { PacketRequest } from 'libs/databases/packet.database';
import ExpressFastifyOpenRedirect from './poc/express-fastify-open-redirect';
import FastifyDOSCVE202122964 from './poc/express-fastify-dos';
import { ApiFuzzer } from './api-fuzzer-poc-generic';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';
import { ParsedPacket } from 'libs/databases/parsedPacket.database';
import SQLiTimeBased from './poc/sqli-timebased';
import Bluebird from 'bluebird';

export default class ApiFuzzerSender {
  private pocs: ApiFuzzer[];
  private pocsName: string[];

  constructor(options: BotFuzzer, packet: ParsedPacket) {
    this.pocs = [
      new ExpressFastifyOpenRedirect(options, packet),
      new FastifyDOSCVE202122964(options, packet),
      new SQLiTimeBased(options, packet),
    ];
    this.pocsName = this.pocs.map((v) => v.constructor.name);
  }

  async sendPocs(): Promise<string[]> {
    const result = await Bluebird.map(
      this.pocs,
      async (poc) => {
        if (poc.condition && !poc.condition()) return 0;
        const res = await poc.poc();
        const ok = (await poc.detect(res)) >= 0.5;
        return ok;
      },
      { concurrency: 1 },
    );
    return result
      .map((res, idx) => {
        if (!res) return null;
        return this.pocsName[idx];
      })
      .filter((r) => r !== null);
  }
}
