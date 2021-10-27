import { PacketRequest } from 'libs/databases/packet.database';
import ApiFuzzerPocGeneric, { ApiFuzzer } from '../api-fuzzer-poc-generic';
import { AxiosResponse } from 'axios';
import ShadelessLogger from 'libs/logger/logger';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';
import { ParsedPacket } from 'libs/databases/parsedPacket.database';

export default class ExpressFastifyOpenRedirect
  extends ApiFuzzerPocGeneric
  implements ApiFuzzer
{
  packet: ParsedPacket;

  constructor(botFuzzer: BotFuzzer, packet: ParsedPacket) {
    const logDir = `logs/api/${
      botFuzzer.project
    }/${ShadelessLogger.sanitizeLogDir(
      packet.origin + packet.path + packet.hash,
    )}/${ExpressFastifyOpenRedirect.name}.txt`;

    super(botFuzzer, ExpressFastifyOpenRedirect.name, logDir);
    this.packet = packet;
  }

  async poc() {
    const opt = await this.getAxiosOptionsFromPacket(this.packet);
    let { url } = opt;
    if (url[url.length - 1] !== '/') url += '/';
    opt.method = 'GET';
    opt.url = url + '/google.com/%2e%2e'; // GET to -> <host>//google.com/%2e%2e
    return this.sendOneRequest(opt);
  }

  async detect(resp: AxiosResponse[]) {
    const res = resp[0];
    if (res && res.status >= 300 && res.status < 400) {
      if (res.headers['location'].includes('//google.com')) return true;
    }
    return false;
  }
}
