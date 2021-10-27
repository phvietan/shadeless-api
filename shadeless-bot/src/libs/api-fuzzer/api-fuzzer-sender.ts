import * as fs from 'fs';
import * as Bluebird from 'bluebird';
import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios';

import { BotPath } from 'libs/databases/botPath.database';
import { PacketRequest } from 'libs/databases/packet.database';
import ShadelessLogger from 'libs/logger/logger';

export default class ApiFuzzerSender {

  constructor(options: BotPath, logger: ShadelessLogger) {
    this.options = options;
  }

  async sendAll(packet: PacketRequest): Promise<boolean> {
    const requestsOptions = this.prepare(packet, currentPath);
    const resps = await this.sendAllRequests(requestsOptions);
    this.logger.log(
      `Done ${requestsOptions.length}/${requestsOptions.length}: ${requestsOptions[0].baseURL}`,
    );
    return resps;
  }
}
