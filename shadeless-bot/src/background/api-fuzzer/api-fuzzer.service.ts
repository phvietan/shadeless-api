import { Injectable } from '@nestjs/common';
import * as Bluebird from 'bluebird';
import AllDatabases from 'libs/databases/all.database';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';
import ParsedPacketDb, { ParsedPacket } from 'libs/databases/parsedPacket.database';
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
        const fuzzPaths = await parsedPacketDb.getTodo(project);
        await Bluebird.map(fuzzPaths, async (path) => {
          this.logger.log(`Fuzzing path: ${path.origin}${path.path}`);
          await this.runOne(botFuzzer, path);
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
      { _id: parsedPath._id },
      { status: FuzzStatus.SCANNING, logDir },
    );

    const taskLogger = this.logger.spawn({ name: 'ApiFuzzer', logDir });
    const pathFuzzerSender = new PathFuzzerSender(botPath, taskLogger);
    const responses = await pathFuzzerSender.prepareAndSendAll(
      await ParsedPacketDb.getInstance().getOneByRequestId(
        parsedPath.requestPacketId,
      ),
      parsedPath.path,
    );
    if (responses[0] === null) {
      const error = 'Got error when GET random 404 page';
      this.logger.log(error);
      await ParsedPathDb.getInstance().update(
        { _id: parsedPath._id },
        {
          status: FuzzStatus.DONE,
          error,
        },
      );
      return;
    } else {
      const result = await new PathFuzzerFilter(
        responses.filter((res) => res !== null) as AxiosResponse<any>[],
        taskLogger.setPrefix(
          `Fuzzing path: ${parsedPath.origin}${parsedPath.path}`,
        ),
      ).filter();
      await ParsedPathDb.getInstance().update(
        { _id: parsedPath._id },
        {
          result: result.map((res) => res.config.url),
          status: FuzzStatus.DONE,
        },
      );
      this.logger.log('Done');
    }
  }
}
