import ApiFuzzerPocGeneric, { ApiFuzzer } from '../api-fuzzer-poc-generic';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';
import { ParsedPacket } from 'libs/databases/parsedPacket.database';
import * as fs from 'fs/promises';
import { ConfigService } from 'config/config.service';
import * as path from 'path';
import { randomHex } from 'libs/helper';

export default class HtmlInjection
  extends ApiFuzzerPocGeneric
  implements ApiFuzzer
{
  constructor(botFuzzer: BotFuzzer, packet: ParsedPacket) {
    super(botFuzzer, packet, HtmlInjection.name);
  }

  async condition() {
    return this.packet.responseContentType.includes('html');
  }

  async poc() {
    const fHtmlInjection = await fs.readFile(
      path.join(new ConfigService().getConfig().wordlistDir, 'html_inject.txt'),
      'utf-8',
    );
    const wordlist = fHtmlInjection.trim().split('\n');
    const payload = wordlist.map((w) => w.replace(/{random}/g, randomHex(10)));

    const opt = await this.getAxiosOptionsFromPacket(this.packet);
    this.logger.log('Substituting in querystring for html injection');
    const resQs = await this.sendAllQuerystringInValueWordlist(opt, payload);
    this.logger.log('Substituting in body for html injection');
    const resBody = await this.sendAllBodyInValueWordlist(opt, payload);

    const responses = [...resQs, ...resBody];
    for (let i = 0; i < responses.length; i++) {
      if (responses[i] === null) continue;
      const contentType = responses[i].headers['content-type'] || '';
      if (!contentType.includes('html')) continue;
      const data = responses[i].data as string;
      for (let j = 0; j < payload.length; j++) {
        if (data.includes(payload[j])) {
          this.logger.log('Found html injection: ' + payload[j]);
          return 1;
        }
      }
    }
    return 0;
  }
}
