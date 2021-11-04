import { Injectable } from '@nestjs/common';
import * as dotenv from 'dotenv';
import * as path from 'path';
dotenv.config();

type Config = {
  databaseUrl: string;
  bindAddress: string;
  frontendUrl: string;
  wordlistDir: string;
  wordlistFile: string;
  bodyDir: string;
};

@Injectable()
export class ConfigService {
  getConfig(): Config {
    const wordlistDir = path.join(__dirname, '../../wordlists/');
    return {
      databaseUrl: process.env.DATABASE_URL,
      bindAddress: process.env.BIND_ADDRESS,
      frontendUrl: process.env.FRONTEND_URL,
      wordlistDir,
      wordlistFile: path.join(wordlistDir, 'dir_test.txt'),
      bodyDir: path.join(__dirname, '../../../files/'),
    };
  }
}
