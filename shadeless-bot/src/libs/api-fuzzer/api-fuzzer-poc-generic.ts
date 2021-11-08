import * as fs from 'fs/promises';
import * as Bluebird from 'bluebird';
import axios, {
  AxiosError,
  AxiosRequestConfig,
  AxiosResponse,
  Method,
} from 'axios';
import { PacketRequest } from 'libs/databases/packet.database';
import {
  getHeaderMapFromHeaders,
  isArray,
  isNumber,
  isObject,
  isString,
  sleep,
} from 'libs/helper';
import ShadelessLogger from 'libs/logger/logger';
import { ConfigService } from 'config/config.service';
import * as path from 'path';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';
import { ParsedPacket } from 'libs/databases/parsedPacket.database';
import * as qs from 'qs';
export interface ApiFuzzer {
  logDir: string;
  logger: ShadelessLogger;

  condition?: () => Promise<boolean>;
  poc: () => Promise<any>;
}

export type MyAxiosResponse = AxiosResponse & { timetook: number };

// TODO: push wordlist dir into BotPath
export default class ApiFuzzerPocGeneric {
  logDir: string;
  options: BotFuzzer;
  logger: ShadelessLogger;
  packet: ParsedPacket;

  constructor(options: BotFuzzer, packet: ParsedPacket, name: string) {
    this.logDir = `logs/api/${options.project}/${ShadelessLogger.sanitizeLogDir(
      packet.origin + packet.path + packet.hash,
    )}/${name}.txt`;

    this.options = options;
    const logger = new ShadelessLogger({
      name,
      logDir: this.logDir,
      prefix: name,
    });
    this.packet = packet;
    this.logger = logger;
  }

  resultObjPayload: any[] = [];
  private recursiveSubstitutePayloadToObj(
    rootObj: any,
    obj: any,
    payload: string,
    reflectedParameters?: Record<string, string>,
  ) {
    for (const key of Object.keys(obj)) {
      if (isString(obj[key]) || isNumber(obj[key])) {
        if (reflectedParameters && !(key in reflectedParameters)) continue;
        const tmp = obj[key];
        obj[key] = payload;
        this.resultObjPayload.push(Object.assign({}, rootObj));
        obj[key] = tmp;
      }
      if (isArray(obj[key]) || isObject(obj[key])) {
        this.recursiveSubstitutePayloadToObj(rootObj, obj[key], payload);
      }
    }
  }

  substitutePayloadToObj(
    obj: any,
    payload: string,
    reflectedParameters?: Record<string, string>,
  ): any[] {
    this.resultObjPayload = [];
    this.recursiveSubstitutePayloadToObj(
      obj,
      obj,
      payload,
      reflectedParameters,
    );
    return this.resultObjPayload;
  }

  protected async getAxiosOptionsFromPacket(
    packet: PacketRequest,
  ): Promise<AxiosRequestConfig> {
    const conf = new ConfigService().getConfig();
    const body = await fs.readFile(
      path.join(conf.bodyDir, packet.project, packet.requestBodyHash),
      'utf-8',
    );
    const params = qs.parse(packet.querystring);
    return {
      method: packet.method as Method,
      params,
      withCredentials: true,
      responseType: 'text',
      transformResponse: [(data) => data],
      url: `${packet.origin}${packet.path}`,
      timeout: this.options.timeout,
      headers: getHeaderMapFromHeaders(packet.requestHeaders),
      data: body,
      maxRedirects: 0,
    };
  }

  protected async sendOneRequest(
    opt: AxiosRequestConfig<any>,
  ): Promise<MyAxiosResponse | null> {
    const before = Date.now();
    try {
      const resp: MyAxiosResponse = {
        ...(await axios.request(opt)),
        timetook: Date.now() - before,
      };
      return resp;
    } catch (err: any) {
      const error = err as AxiosError<any>;
      const after = Date.now();
      // The request was made and the server responded with a status code that falls out of the range of 2xx
      if (error.response) {
        const response: MyAxiosResponse = {
          ...error.response,
          timetook: after - before,
        };
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

  protected async sendAllQuerystringInValue(
    opt: AxiosRequestConfig<any>,
    payload: string,
    reflectedParameters?: Record<string, string>,
  ): Promise<MyAxiosResponse[]> {
    const { params } = opt;
    const listParamsWithPayload = this.substitutePayloadToObj(
      params,
      payload,
      reflectedParameters,
    );
    return Bluebird.map(listParamsWithPayload, async (params) => {
      const newOpt = Object.assign({}, opt);
      newOpt.params = params;
      return this.sendOneRequest(newOpt);
    });
  }

  protected async sendAllBodyInValue(
    opt: AxiosRequestConfig<any>,
    payload: string,
    reflectedParameters?: Record<string, string>,
  ): Promise<MyAxiosResponse[]> {
    const contentType = opt.headers['Content-Type'];
    if (!contentType) return [];
    let body: any = {};
    if (contentType.includes('json')) {
      body = JSON.parse(opt.data);
    }
    if (contentType.includes('x-www-form-urlencoded')) {
      body = Object.assign({}, qs.parse(opt.data));
    }
    const listParamsWithPayload = this.substitutePayloadToObj(
      body,
      payload,
      reflectedParameters,
    );
    return Bluebird.map(listParamsWithPayload, async (data) => {
      const newOpt = Object.assign({}, opt);
      if (contentType.includes('x-www-form-urlencoded')) {
        newOpt.data = qs.stringify(data);
      } else {
        newOpt.data = data;
      }
      return this.sendOneRequest(newOpt);
    });
  }

  protected async sendAllQuerystringInValueWordlist(
    opt: AxiosRequestConfig<any>,
    wordlist: string[],
    reflectedParameters?: Record<string, string>,
  ): Promise<MyAxiosResponse[]> {
    let cnt = 0;
    let result: MyAxiosResponse[] = [];
    await Bluebird.map(
      wordlist,
      async (word) => {
        sleep(this.options.sleepRequest);
        cnt += 1;
        if (cnt % 30 === 0) {
          this.logger.log(`Done ${cnt}/${wordlist.length}: ${word}`);
        }
        const resps = await this.sendAllQuerystringInValue(
          opt,
          word,
          reflectedParameters,
        );
        result = [...result, ...resps];
      },
      { concurrency: this.options.asyncRequest },
    );
    return result;
  }

  protected async sendAllBodyInValueWordlist(
    opt: AxiosRequestConfig<any>,
    wordlist: string[],
    reflectedParameters?: Record<string, string>,
  ): Promise<MyAxiosResponse[]> {
    let cnt = 0;
    let result: MyAxiosResponse[] = [];
    await Bluebird.map(
      wordlist,
      async (word) => {
        sleep(this.options.sleepRequest);
        cnt += 1;
        if (cnt % 30 === 0) {
          this.logger.log(`Done ${cnt}/${wordlist.length}: ${word}`);
        }
        const resps = await this.sendAllBodyInValue(
          opt,
          word,
          reflectedParameters,
        );
        result = [...result, ...resps];
      },
      { concurrency: this.options.asyncRequest },
    );
    return result;
  }

  isEtcPasswd(s: string): boolean {
    return s.includes('root:x:0:0:root:/root');
  }
}
