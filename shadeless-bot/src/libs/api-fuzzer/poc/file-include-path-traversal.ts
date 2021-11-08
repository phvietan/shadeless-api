import ApiFuzzerPocGeneric, { ApiFuzzer } from '../api-fuzzer-poc-generic';
import { AxiosResponse } from 'axios';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';
import { ParsedPacket } from 'libs/databases/parsedPacket.database';
import * as fs from 'fs';
import * as path from 'path';
import { ConfigService } from 'config/config.service';
import { randomHex } from 'libs/helper';

const lfiF = fs
  .readFileSync(
    path.join(new ConfigService().getConfig().wordlistDir, 'lfi.txt'),
    'utf-8',
  )
  .trim();
const payloads = lfiF.split('\n');

export default class FileInclusionPathTraverse
  extends ApiFuzzerPocGeneric
  implements ApiFuzzer
{
  cache: string[];

  constructor(botFuzzer: BotFuzzer, packet: ParsedPacket) {
    super(botFuzzer, packet, FileInclusionPathTraverse.name);
  }

  async poc() {
    const opt = await this.getAxiosOptionsFromPacket(this.packet);

    this.logger.log('Substituting LFI payload into body');
    const resBody = await this.sendAllBodyInValueWordlist(opt, payloads);

    this.logger.log('Substituting LFI into querystring');
    const resQs = await this.sendAllQuerystringInValueWordlist(opt, payloads);
    const responses = [...resBody, ...resQs];

    for (let i = 0; i < responses.length; i += 1) {
      const r = responses[i];
      if (r === null) return;
      if (this.isEtcPasswd(r.data as string)) {
        this.logger.log('Detected LFI');
        return 1;
      }
    }
    return 0;
  }
}
