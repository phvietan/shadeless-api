import { Collection, Db } from 'mongodb';

export interface BotPath {
  timeout: number;
  project: string;
  running: boolean;
  sleepRequest: number;
  asyncRequest: number;
}

class BotPathDb {
  db: Collection<Document>;
  private static instance?: BotPathDb = undefined;

  constructor(dbo: Db) {
    this.db = dbo.collection('bot_paths');
  }

  static getInstance(dbo?: Db): BotPathDb {
    if (!this.instance) {
      if (!dbo) {
        throw new Error('WTF? DBO is undefined?')
      }
      this.instance = new BotPathDb(dbo);
    }
    return this.instance;
  }

  async getOneByProject(project: string) {
    const document = await this.db.findOne({ project });
    if (!document) return undefined;
    return (document as any) as BotPath;
  }

  async getRunningProject() {
    const documents = await this.db.find({ running: true }).toArray();
    return (documents as any) as BotPath[];
  }
}

export default BotPathDb;
