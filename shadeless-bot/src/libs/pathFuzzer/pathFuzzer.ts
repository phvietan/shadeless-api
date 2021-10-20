import { AxiosRequestConfig, AxiosResponse } from 'axios';

import PathFilter from './pathFuzzer.filter';
import BotPathSender from './pathFuzzer.sender';

import ParsedPathDb, {
  ParsedPath,
  PathStatus,
  PathType,
} from 'libs/databases/parsedPath.database';
import { BotPath } from 'libs/databases/botPath.database';
import ParsedPacketDb from 'libs/databases/parsedPacket.database';

class PathFuzzer {
  options: BotPath;
  pathSender: BotPathSender;
  parsedPath: ParsedPath;

  constructor(parsedPath: ParsedPath, options: BotPath) {
    this.pathSender = new BotPathSender(options);
    this.parsedPath = parsedPath;
    this.options = options;
  }

  private async getRequestsOptions(): Promise<AxiosRequestConfig[]> {
    const parsedPacket = await ParsedPacketDb.getInstance().getOneByRequestId(
      this.parsedPath.requestPacketId,
    );
    if (!parsedPacket) return [];
    return this.pathSender.prepare(parsedPacket, this.parsedPath.path);
  }

  private getParsedPathFromResponse(
    results: AxiosResponse<any>[],
  ): ParsedPath[] {
    return results.map((res) => {
      return {
        type: PathType.NONE,
        requestPacketId: '',
        origin: this.parsedPath.origin,
        path: res.config.url as string,
        force: false,
        status: PathStatus.TODO,
        project: this.parsedPath.project,
        created_at: new Date(),
        updated_at: new Date(),
        error: '',
      };
    });
  }

  async pushResultDB(result: AxiosResponse<any>[]) {
    const dbResult = this.getParsedPathFromResponse(result);
    await ParsedPathDb.getInstance().insertResult(dbResult);
    await ParsedPathDb.getInstance().update(
      {
        _id: this.parsedPath._id,
      },
      { status: PathStatus.DONE },
    );
    console.log('Done');
  }

  async run() {
    await ParsedPathDb.getInstance().update(
      {
        _id: this.parsedPath._id,
        error: '',
      },
      { status: PathStatus.SCANNING },
    );
    const requestsOptions = await this.getRequestsOptions();
    const responses = await this.pathSender.sendAll(requestsOptions);
    if (responses[0] === null) {
      const err =
        'Got error when GET random 404 page, please run again or recheck';
      console.log(err);
      await ParsedPathDb.getInstance().updateError(
        { requestPacketId: this.parsedPath.requestPacketId },
        err,
      );
      await ParsedPathDb.getInstance().update(
        {
          _id: this.parsedPath._id,
        },
        { status: PathStatus.DONE },
      );
      return;
    }
    const pathFilter = new PathFilter(
      responses.filter((res) => res !== null) as AxiosResponse<any>[],
    );
    const result = await pathFilter.filter(
      `Fuzzing path: ${this.parsedPath.origin}${this.parsedPath.path}`,
    );
    this.pushResultDB(result);
  }
}

export default PathFuzzer;
