import * as fs from 'fs';
import * as Bluebird from 'bluebird';
import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios';

import { BotPath } from 'libs/databases/botPath.database';
import { PacketRequest } from 'libs/databases/packet.database';
import { getHeaderMapFromHeaders, randomHex, sleep } from 'libs/helper';
import { ConfigService } from 'config/config.service';
import ShadelessLogger from 'libs/logger/logger';

const { wordlistFile } = new ConfigService().getConfig();
const f = fs.readFileSync(wordlistFile, 'utf8').trim();
const wordlist = f.split('\n');

// TODO: push wordlist dir into BotPath
export default class PathFuzzerSender {
  wordlist: string[] = wordlist;
  options: BotPath;
  logger: ShadelessLogger;

  constructor(options: BotPath, logger: ShadelessLogger) {
    this.options = options;
    this.logger = logger;
  }

  private getAxiosOptions(
    packet: PacketRequest,
    currentPath: string,
    newPath: string,
  ): AxiosRequestConfig {
    const path =
      currentPath[currentPath.length - 1] === '/'
        ? currentPath
        : currentPath + '/';
    return {
      method: 'GET',
      withCredentials: true,
      responseType: 'text',
      transformResponse: [(data) => data],
      baseURL: packet.origin,
      url: path + newPath,
      timeout: this.options.timeout,
      headers: getHeaderMapFromHeaders(packet.requestHeaders),
    };
  }

  private prepare(
    packet: PacketRequest,
    currentPath: string,
  ): AxiosRequestConfig[] {
    const opts = this.wordlist.map((w) =>
      this.getAxiosOptions(packet, currentPath, w),
    );
    opts.push(this.getAxiosOptions(packet, currentPath, randomHex(32)));
    [opts[0], opts[opts.length - 1]] = [opts[opts.length - 1], opts[0]]; // Make the first option as the random path (404 request)
    return opts;
  }

  private async sendOneRequest(opt: AxiosRequestConfig<any>) {
    try {
      const resp = await axios.request(opt);
      return resp;
    } catch (err: any) {
      const error = err as AxiosError<any>;
      // The request was made and the server responded with a status code that falls out of the range of 2xx
      if (error.response) {
        const { response } = error;
        return response;
      } else if (error.request) {
        // The request was made but no response was received
        return null;
      } else {
        // Something happened in setting up the request that triggered an Error
        this.logger.log('WTF error in config? ' + error.message);
        return null;
      }
    }
  }

  private async sendAllRequests(requestsOptions: AxiosRequestConfig[]) {
    let cnt = 0;
    return Bluebird.map(
      requestsOptions,
      async (opt) => {
        sleep(this.options.sleepRequest);
        cnt += 1;
        if (cnt % 30 === 0) {
          this.logger.log(
            `Done ${cnt}/${requestsOptions.length}: ${opt.baseURL}${opt.url}`,
          );
        }
        return this.sendOneRequest(opt);
      },
      { concurrency: 3 },
    );
  }

  async prepareAndSendAll(
    packet: PacketRequest,
    currentPath: string,
  ): Promise<AxiosResponse<any>[]> {
    const requestsOptions = this.prepare(packet, currentPath);
    const resps = await this.sendAllRequests(requestsOptions);
    this.logger.log(
      `Done ${requestsOptions.length}/${requestsOptions.length}: ${requestsOptions[0].baseURL}`,
    );
    return resps;
  }
}
