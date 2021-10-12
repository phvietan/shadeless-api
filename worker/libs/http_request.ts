import { ParsedPacketRequest } from "./database/parsed_packet";
import axios, { AxiosError, AxiosInstance, AxiosRequestConfig, AxiosResponse, Method } from 'axios';
import fsPromise from 'fs/promises';
import fs from 'fs';
import path from 'path';
import { randomHex } from "./helper";
import Bluebird from "bluebird";

class HttpRequest {
  wordlist: string[];

  constructor () {
    const wordlistFile = fs.readFileSync(path.join(__dirname, 'fuzzer/poc/wordlists/dir.txt'), 'utf8').trim();
    this.wordlist = wordlistFile.split('\n');
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

  static async createFromPacket(packet: ParsedPacketRequest): Promise<AxiosRequestConfig> {
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

  private getOptFromPacketAndNewPath(packet: ParsedPacketRequest, currentPath: string, newPath: string): AxiosRequestConfig {
    const path = (currentPath[currentPath.length - 1] === '/') ? currentPath : currentPath + '/';
    return {
      method: 'GET',
      withCredentials: true,
      responseType: 'text',
      transformResponse: [data => data],
      baseURL: packet.origin,
      url: path + newPath,
      timeout: 1500,
      headers: HttpRequest.createHeader(packet.requestHeaders),
    };
  }

  createOptsForDirPath(packet: ParsedPacketRequest, currentPath: string): AxiosRequestConfig[] {
    const opts = this.wordlist.map(w => this.getOptFromPacketAndNewPath(packet, currentPath, w));
    opts.push(this.getOptFromPacketAndNewPath(packet, currentPath, randomHex(32)));
    [opts[0], opts[opts.length - 1]] = [opts[opts.length - 1], opts[0]]; // Make the first option as the random path (404 request)
    return opts;
  }

  async sendAllOpts(requestsOptions: AxiosRequestConfig[]): Promise<AxiosResponse<any>[]> {
    let cnt = 0;
    const resps = await Bluebird.map(requestsOptions, async (opts) => {
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

const ins = new HttpRequest();
export default ins;
