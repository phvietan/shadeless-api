import { Injectable } from '@nestjs/common';
import * as dotenv from 'dotenv';
import * as path from 'path';
dotenv.config();

type Config = {
  databaseUrl: string;
  bindAddress: string;
  frontendUrl: string;
  wordlistDir: string;
};

@Injectable()
export class ConfigService {
  private static instance?: ConfigService;

  constructor() {
    console.log('Initialized ConfigService');
  }

  static getInstance(): ConfigService {
    if (this.instance) return this.instance;
    this.instance = new ConfigService();
    return this.instance;
  }

  getConfig(): Config {
    return {
      databaseUrl: process.env.DATABASE_URL,
      bindAddress: process.env.BIND_ADDRESS,
      frontendUrl: process.env.FRONTEND_URL,
      wordlistDir: path.join(__dirname, '../../wordlists/dir.txt'),
    };
  }
}
