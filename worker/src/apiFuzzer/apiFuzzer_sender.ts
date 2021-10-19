import fs from 'fs';
import fsPromise from 'fs/promises';
import path from 'path';
import { sleep } from "../libs/helper";
import Bluebird from "bluebird";
import { BotFuzzer } from "../database/bot_fuzzer";
import axios, { AxiosError, AxiosRequestConfig, AxiosResponse, Method } from 'axios';
import { PacketRequest } from "../database/packet";

class BotFuzzerSender {
  options: BotFuzzer;
  wordlist: string[];

  constructor (options: BotFuzzer) {
    const wordlistFile = fs.readFileSync(path.join(__dirname, 'fuzzer/poc/wordlists/dir.txt'), 'utf8').trim();
    this.wordlist = wordlistFile.split('\n');
    this.options = options;
  }

  static createHeader(headers: string[]): Record<string,string> {
    const result: Record<string,string> = {};
    headers.forEach(header => {
      const delimiter = header.indexOf(': ');
      if (delimiter === -1) return;
      const key = header.slice(0, delimiter);
      const value = header.slice(delimiter + 2);
      if (key.toLowerCase() === 'content-length') return;
      result[key] = value;
    });
    return result;
  }

  static async createFromPacket(packet: PacketRequest): Promise<AxiosRequestConfig> {
    const body = await fsPromise.readFile(path.join(__dirname, '../../main/files/', packet.project, packet.requestBodyHash));
    return {
      method: packet.method as Method,
      withCredentials: true,
      responseType: 'text',
      baseURL: packet.origin,
      url: packet.path,
      timeout: 1500,
      data: body,
      headers: this.createHeader(packet.requestHeaders),
      transformResponse: [data => data],
    };
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

export default BotFuzzerSender;
