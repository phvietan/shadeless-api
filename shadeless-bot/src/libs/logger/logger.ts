import { Logger } from '@nestjs/common';
import * as winston from 'winston';

type LoggerOptions = {
  prefix?: string;
  name: string;
  logDir: string;
};

export default class ShadelessLogger extends Logger {
  private w: winston.Logger;
  private prefix?: string;

  constructor(options?: LoggerOptions) {
    super();

    const prefix = options?.prefix;
    this.prefix = prefix;

    this.w = options
      ? winston.createLogger({
          level: 'info',
          transports: [
            new winston.transports.Console(),
            new winston.transports.File({ filename: options.logDir }),
          ],
        })
      : winston.createLogger({
          level: 'info',
          transports: [new winston.transports.Console()],
        });
  }

  spawn(options?: LoggerOptions): ShadelessLogger {
    return new ShadelessLogger(options);
  }

  setPrefix(prefix: string): ShadelessLogger {
    this.prefix = prefix;
    return this;
  }

  log(message: string) {
    if (this.prefix) {
      this.w.info(`${this.prefix} ${message}`);
    } else {
      this.w.info(message);
    }
  }

  static sanitizeLogDir(name: string): string {
    return name.replace(/\/|\\|\>|\x00|\||:|&/g, '_');
  }
}
