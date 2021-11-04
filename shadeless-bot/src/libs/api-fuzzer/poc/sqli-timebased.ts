import ApiFuzzerPocGeneric, {
  ApiFuzzer,
  MyAxiosResponse,
} from '../api-fuzzer-poc-generic';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';
import { ParsedPacket } from 'libs/databases/parsedPacket.database';
import { ConfigService } from 'config/config.service';
import * as path from 'path';
import * as fs from 'fs';

const f = fs
  .readFileSync(
    path.join(
      new ConfigService().getConfig().wordlistDir,
      'sqli_timebased.txt',
    ),
    'utf-8',
  )
  .trim();
const wordlist = f.split('\n');

export default class SQLiTimeBased
  extends ApiFuzzerPocGeneric
  implements ApiFuzzer
{
  constructor(botFuzzer: BotFuzzer, packet: ParsedPacket) {
    super(botFuzzer, packet, SQLiTimeBased.name);
  }

  async poc() {
    const opt = await this.getAxiosOptionsFromPacket(this.packet);
    let res = await this.sendAllBodyInValueWordlist(opt, wordlist);
    res = [
      ...res,
      ...(await this.sendAllQuerystringInValueWordlist(opt, wordlist)),
    ];
    return res;
  }

  async detect(res: MyAxiosResponse[]) {
    this.logger.setPrefix('Detection SQLI-timebased:');
    let cntOver4000 = 0;
    let cntNull = 0;
    res.forEach((r) => {
      if (r === null) cntNull += 1;
      else {
        if (r.timetook >= 4000) {
          cntOver4000 += 1;
        }
      }
    });
    if (cntNull < 5 && cntOver4000 > 0) {
      this.logger.log(
        'Detected, server is quite stable, detected sqli, null < 5',
      );
      return 1;
    }
    if (cntNull < 40 && cntOver4000 > 0) {
      this.logger.log('Detected, Server is not so stable 1, 5 < null < 40');
      return 0.7;
    }
    if (cntNull < 100 && cntOver4000 > 0) {
      this.logger.log('Detected, Server is not so stable 2, 40 < null < 100');
    }
    this.logger.log('Server is not stable at all, null > 100');
    return 0;
  }
}
