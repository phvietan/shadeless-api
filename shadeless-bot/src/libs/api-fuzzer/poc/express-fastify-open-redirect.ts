import ApiFuzzerPocGeneric, { ApiFuzzer } from '../api-fuzzer-poc-generic';
import { AxiosResponse } from 'axios';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';
import { ParsedPacket } from 'libs/databases/parsedPacket.database';

export default class ExpressFastifyOpenRedirect
  extends ApiFuzzerPocGeneric
  implements ApiFuzzer
{
  constructor(botFuzzer: BotFuzzer, packet: ParsedPacket) {
    super(botFuzzer, packet, ExpressFastifyOpenRedirect.name);
  }

  async poc() {
    const opt = await this.getAxiosOptionsFromPacket(this.packet);
    let { url } = opt;
    if (url[url.length - 1] !== '/') url += '/';
    opt.method = 'GET';
    opt.url = url + '/google.com/%2e%2e'; // GET to -> <host>//google.com/%2e%2e
    return this.sendOneRequest(opt);
  }

  async detect(res: AxiosResponse) {
    this.logger.setPrefix('Detection ExpressFastifyOpenRedirect:');
    if (res && res.status >= 300 && res.status < 400) {
      if (res.headers['location'].includes('//google.com')) {
        this.logger.log('Detected');
        return 1;
      }
    }
    return 0;
  }
}
