import ApiFuzzerPocGeneric, {
  ApiFuzzer,
  MyAxiosResponse,
} from '../api-fuzzer-poc-generic';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';
import { ParsedPacket } from 'libs/databases/parsedPacket.database';
import { ConfigService } from 'config/config.service';
import * as path from 'path';
import * as fs from 'fs';

const fSleep = fs
  .readFileSync(
    path.join(new ConfigService().getConfig().wordlistDir, 'rce_timebased.txt'),
    'utf-8',
  )
  .trim();
const wordlistSleep = fSleep.split('\n');

const fRceEtcPasswd = fs
  .readFileSync(
    path.join(new ConfigService().getConfig().wordlistDir, 'rce_etcpasswd.txt'),
    'utf-8',
  )
  .trim();
const wordlistEtcPasswd = fRceEtcPasswd.split('\n');

export default class RCEInputNoSanitize
  extends ApiFuzzerPocGeneric
  implements ApiFuzzer
{
  constructor(botFuzzer: BotFuzzer, packet: ParsedPacket) {
    super(botFuzzer, packet, RCEInputNoSanitize.name);
  }

  async pocTimebased() {
    const opt = await this.getAxiosOptionsFromPacket(this.packet);
    const resBody = await this.sendAllBodyInValueWordlist(opt, wordlistSleep);
    const resQs = await this.sendAllQuerystringInValueWordlist(
      opt,
      wordlistSleep,
    );
    const responses = [...resBody, ...resQs];
    let cntOver4000 = 0;
    let cntNull = 0;
    responses.forEach((r) => {
      if (r === null) cntNull += 1;
      else {
        if (r.timetook >= 4000) {
          cntOver4000 += 1;
        }
      }
    });
    if (cntNull < 5 && cntOver4000 > 0) {
      this.logger.log('Detected, server is quite stable, null < 5');
      return 1;
    }
    if (cntNull < 40 && cntOver4000 > 0) {
      this.logger.log('Detected, Server is not so stable 1, 5 < null < 40');
      return 0.7;
    }
    if (cntNull < 100 && cntOver4000 > 0) {
      this.logger.log('Detected, Server is not so stable 2, 40 < null < 100');
    }
    if (cntNull >= 100)
      this.logger.log('Server is not stable at all, null > 100');
    else {
      this.logger.log('Not detected');
    }
    return 0;
  }

  async pocCatEtcPasswd() {
    const opt = await this.getAxiosOptionsFromPacket(this.packet);
    const resBody = await this.sendAllBodyInValueWordlist(
      opt,
      wordlistEtcPasswd,
    );
    const resQs = await this.sendAllQuerystringInValueWordlist(
      opt,
      wordlistEtcPasswd,
    );
    const responses = [...resBody, ...resQs];
    for (let i = 0; i < responses.length; i++) {
      if (this.isEtcPasswd(responses[i].data as string)) {
        return 1;
      }
    }
    return 0;
  }

  async poc() {
    const poc1 = await this.pocCatEtcPasswd();
    if (poc1 === 1) return 1;
    return this.pocTimebased();
  }
}
