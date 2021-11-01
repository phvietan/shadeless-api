import { Collection, Db } from 'mongodb';

export interface BotFuzzer {
  timeout: number;
  project: string;
  running: boolean;
  sleepRequest: number;
  asyncRequest: number;
}

class BotFuzzerDb {
  db: Collection<Document>;
  private static instance?: BotFuzzerDb = undefined;

  constructor(dbo: Db) {
    this.db = dbo.collection('bot_paths');
  }

  static getInstance(dbo?: Db): BotFuzzerDb {
    if (!this.instance) {
      if (!dbo) {
        throw new Error('WTF? DBO is undefined?');
      }
      this.instance = new BotFuzzerDb(dbo);
    }
    return this.instance;
  }

  async getOneByProject(project: string) {
    const document = await this.db.findOne({ project });
    if (!document) return undefined;
    return document as any as BotFuzzer;
  }

  async getRunningProject() {
    const documents = await this.db.find({ running: true }).toArray();
    return documents as any as BotFuzzer[];
  }
}

export default BotFuzzerDb;
