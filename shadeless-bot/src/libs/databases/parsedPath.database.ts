import { Collection, Db } from 'mongodb';
import { Project, getFilterByProjectForBW } from './project.database';

export enum FuzzStatus {
  TODO = 'todo',
  SCANNING = 'scanning',
  DONE = 'done',
}

export interface ParsedPath {
  _id?: string;
  requestPacketId: string;
  origin: string;
  path: string;
  status: FuzzStatus;
  project: string;
  force: boolean;
  created_at?: Date;
  updated_at?: Date;
  error: string;
}

class ParsedPathDb {
  FUZZ_PATH_NUM = 1;
  db: Collection<Document>;
  private static instance?: ParsedPathDb = undefined;

  constructor(dbo: Db) {
    this.db = dbo.collection('parsed_paths');
  }

  static getInstance(dbo?: Db): ParsedPathDb {
    if (!this.instance) {
      if (!dbo) {
        throw new Error('WTF? DBO is undefined?');
      }
      this.instance = new ParsedPathDb(dbo);
    }
    return this.instance;
  }

  async getTodo(project: Project) {
    const filter = getFilterByProjectForBW(project);
    const documents = await this.db
      .find({
        status: FuzzStatus.TODO,
        requestPacketId: { $ne: '' },
        ...filter,
      })
      .sort({ created_at: 1 })
      .limit(this.FUZZ_PATH_NUM)
      .toArray();
    return documents as any as ParsedPath[];
  }

  async update(query: any, value: any) {
    await this.db.updateOne(query, {
      $set: value,
    });
  }

  async resetScanning() {
    return this.db.updateMany(
      { status: FuzzStatus.SCANNING },
      {
        $set: { status: FuzzStatus.TODO },
      },
    );
  }
}

export default ParsedPathDb;
