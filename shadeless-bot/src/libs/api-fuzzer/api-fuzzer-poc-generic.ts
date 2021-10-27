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
import path from 'path';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';

export interface ApiFuzzer {
  poc: () => Promise<any>;
  detect: (resp: AxiosResponse[]) => Promise<boolean>;
}

// TODO: push wordlist dir into BotPath
export default class ApiFuzzerPocGeneric {
  options: BotFuzzer;
  logger: ShadelessLogger;

  constructor(options: BotFuzzer, name: string, logDir: string) {
    this.options = options;
    const logger = new ShadelessLogger({
      name,
      logDir,
      prefix: name,
    });
    this.logger = logger;
  }

  static queryStringToObject(params: string) {
    return JSON.parse(
      '{"' + decodeURI(params.replace(/&/g, '","').replace(/=/g, '":"')) + '"}',
    );
  }

  resultObjPayload: any[];
  private recursiveSubstitutePayloadToObj(
    rootObj: any,
    obj: any,
    payload: string,
  ) {
    for (const key of Object.keys(obj)) {
      if (isString(obj[key]) || isNumber(obj[key])) {
        const tmp = obj[key];
        obj[key] = payload;
        this.resultObjPayload.push(rootObj);
        obj[key] = tmp;
      }
      if (isArray(obj[key]) || isObject(obj[key])) {
        this.recursiveSubstitutePayloadToObj(rootObj, obj[key], payload);
      }
    }
  }

  substitutePayloadToObj(obj: any, payload: string): any[] {
    this.recursiveSubstitutePayloadToObj(obj, obj, payload);
    return this.resultObjPayload;
  }

  protected async getAxiosOptionsFromPacket(
    packet: PacketRequest,
  ): Promise<AxiosRequestConfig> {
    const conf = new ConfigService().getConfig();
    const body = await fs.readFile(
      path.join(conf.bodyDir, packet.project, packet.requestBodyHash),
    );
    const params = ApiFuzzerPocGeneric.queryStringToObject(packet.querystring);
    return {
      method: packet.method as Method,
      params,
      withCredentials: true,
      responseType: 'text',
      transformResponse: [(data) => data],
      baseURL: packet.origin,
      url: packet.path,
      timeout: this.options.timeout,
      headers: getHeaderMapFromHeaders(packet.requestHeaders),
      data: body,
    };
  }

  protected async sendOneRequest(
    opt: AxiosRequestConfig<any>,
  ): Promise<AxiosResponse<unknown, any> | null> {
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

  protected async sendAllQuerystringInValue(
    opt: AxiosRequestConfig<any>,
    payload: string,
  ): Promise<AxiosResponse<unknown, any>[]> {
    const { params } = opt;
    const listParamsWithPayload = this.substitutePayloadToObj(params, payload);
    return Bluebird.map(listParamsWithPayload, async (params) => {
      const newOpt = Object.assign({}, opt);
      newOpt.params = params;
      return this.sendOneRequest(newOpt);
    });
  }

  protected async sendAllQuerystringInValueWordlist(
    opt: AxiosRequestConfig<any>,
    wordlist: string[],
  ): Promise<AxiosResponse<unknown, any>[]> {
    let cnt = 0;
    let result: AxiosResponse<unknown, any>[] = [];
    await Bluebird.map(
      wordlist,
      async (word) => {
        sleep(this.options.sleepRequest);
        cnt += 1;
        if (cnt % 30 === 0) {
          this.logger.log(`Done ${cnt}/${wordlist.length}: ${word}`);
        }
        const resps = await this.sendAllQuerystringInValue(opt, word);
        result = [...result, ...resps];
      },
      { concurrency: 3 },
    );
    return result;
  }
}
