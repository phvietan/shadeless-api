import { Injectable } from '@nestjs/common';
import * as Bluebird from 'bluebird';
import AllDatabases from 'libs/databases/all.database';
import ParsedPacketDb from 'libs/databases/parsedPacket.database';
import ParsedPathDb, {
  ParsedPath,
  FuzzStatus,
} from 'libs/databases/parsedPath.database';
import ShadelessLogger from 'libs/logger/logger';
import { sleep } from 'libs/helper';
import PathFuzzerSender from 'libs/path-fuzzer/path-fuzzer-sender';
import { BotPath } from 'libs/databases/botPath.database';
import PathFuzzerFilter from 'libs/path-fuzzer/path-fuzzer-filter';
import { AxiosResponse } from 'axios';

@Injectable()
export class PathFuzzerService {
  logger = new ShadelessLogger();

  async runForever() {
    const { parsedPathDb, botPathDb, projectDb } = AllDatabases.getInstance();
    await parsedPathDb.resetScanning();
    while (true) {
      const botPathRunning = await botPathDb.getRunningProject();
      this.logger.log(`[BotPath]: Found ${botPathRunning.length} targets`);
      await Bluebird.map(botPathRunning, async (botPath) => {
        const project = await projectDb.getOneProjectByName(botPath.project);
        const fuzzPaths = await parsedPathDb.getTodo(project);
        await Bluebird.map(fuzzPaths, async (path) => {
          this.logger.log(`Fuzzing path: ${path.origin}${path.path}`);
          await this.runOne(botPath, path);
        });
      });
      await sleep(5000);
    }
  }

  private async runOne(botPath: BotPath, parsedPath: ParsedPath) {
    const logDir = `logs/path/${
      botPath.project
    }/${ShadelessLogger.sanitizeLogDir(
      parsedPath.origin + parsedPath.path,
    )}.txt`;
    await ParsedPathDb.getInstance().update(
      { _id: parsedPath._id },
      { status: FuzzStatus.SCANNING, logDir },
    );

    const taskLogger = this.logger.spawn({ name: 'PathFuzzer', logDir });
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
