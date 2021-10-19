import fs from 'fs';
import path from 'path';
import { getHeaderMapFromHeaders, randomHex, sleep } from "../libs/helper";
import Bluebird from "bluebird";
import { BotPath } from "../database/bot_path";
import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios';
import { PacketRequest } from 'database/packet';

class BotPathSender {
  wordlist: string[];
  options: BotPath;

  constructor (options: BotPath) {
    const wordlistFile = fs.readFileSync(path.join(__dirname, 'wordlists/dir.txt'), 'utf8').trim();
    this.wordlist = wordlistFile.split('\n');
    this.options = options;
  }

  private getOptFromPacketAndNewPath(packet: PacketRequest, currentPath: string,  newPath: string): AxiosRequestConfig {
    const path = (currentPath[currentPath.length - 1] === '/') ? currentPath : currentPath + '/';
    return {
      method: 'GET',
      withCredentials: true,
      responseType: 'text',
      transformResponse: [data => data],
      baseURL: packet.origin,
      url: path + newPath,
      timeout: this.options.timeout,
      headers: getHeaderMapFromHeaders(packet.requestHeaders),
    };
  }

  prepare(packet: PacketRequest, currentPath: string): AxiosRequestConfig[] {
    const opts = this.wordlist.map(w => this.getOptFromPacketAndNewPath(packet, currentPath, w));
    opts.push(this.getOptFromPacketAndNewPath(packet, currentPath, randomHex(32)));
    [opts[0], opts[opts.length - 1]] = [opts[opts.length - 1], opts[0]]; // Make the first option as the random path (404 request)
    return opts;
  }

  async sendAll(requestsOptions: AxiosRequestConfig[]): Promise<AxiosResponse<any>[]> {
    let cnt = 0;
    const resps = await Bluebird.map(requestsOptions, async (opts) => {
      sleep(this.options.sleepRequest);
      cnt += 1;
      if (cnt % 30 === 0) {
        console.log(`Done ${cnt}/${requestsOptions.length}: ${opts.baseURL}${opts.url}`);
      }
      try {
        const resp = await axios.request(opts);
        return resp;
      }
      catch (err: any) {
        const error = err as AxiosError<any>;
        // The request was made and the server responded with a status code that falls out of the range of 2xx
        if (error.response) {
          const { response } = error;
          return response;
        } else if (error.request) { // The request was made but no response was received
          return null;
        } else {
          // Something happened in setting up the request that triggered an Error
          console.log('WTF error in config? ', error.message);
          return null;
        }
      }
    }, { concurrency: 3 });
    return resps;
  }
}

export default BotPathSender;
