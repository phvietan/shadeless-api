import { Injectable } from '@nestjs/common';
import * as Bluebird from 'bluebird';
import ApiFuzzerSender from 'libs/api-fuzzer/api-fuzzer-sender';
import AllDatabases from 'libs/databases/all.database';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';
import ParsedPacketDb, {
  ParsedPacket,
} from 'libs/databases/parsedPacket.database';
import { FuzzStatus } from 'libs/databases/parsedPath.database';
import { sleep } from 'libs/helper';
import ShadelessLogger from 'libs/logger/logger';

@Injectable()
export class ApiFuzzerService {
  logger = new ShadelessLogger();

  async runForever() {
    const { parsedPacketDb, botFuzzerDb, projectDb } =
      AllDatabases.getInstance();
    await parsedPacketDb.resetScanning();
    while (true) {
      const botFuzzerRunning = await botFuzzerDb.getRunningProject();
      this.logger.log(`Found ${botFuzzerRunning.length} targets`);
      await Bluebird.map(botFuzzerRunning, async (botFuzzer) => {
        const project = await projectDb.getOneProjectByName(botFuzzer.project);
        const fuzzApis = await parsedPacketDb.getTodo(project);
        await Bluebird.map(fuzzApis, async (parsedPacket) => {
          this.logger.log(
            `Fuzzing api: ${parsedPacket.origin}${parsedPacket.path}`,
          );
          await this.runOne(botFuzzer, parsedPacket);
        });
      });
      await sleep(5000);
    }
  }

  private async runOne(botFuzzer: BotFuzzer, parsedPacket: ParsedPacket) {
    const logDir = `logs/api/${
      botFuzzer.project
    }/${ShadelessLogger.sanitizeLogDir(
      parsedPacket.origin + parsedPacket.path + parsedPacket.hash,
    )}.txt`;
    await ParsedPacketDb.getInstance().update(
      { _id: parsedPacket._id },
      { status: FuzzStatus.SCANNING, logDir },
    );

    const pathFuzzerSender = new ApiFuzzerSender(botFuzzer, parsedPacket);
    const responses = await pathFuzzerSender.sendPocs();
    await ParsedPacketDb.getInstance().update(
      { _id: parsedPacket._id },
      { status: FuzzStatus.DONE, result: responses },
    );
  }
}
