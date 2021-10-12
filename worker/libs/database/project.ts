import { Collection, Db } from 'mongodb';
import PacketDb from './packet';

export enum ProjectStatus {
  TODO = 'todo',
  HACKING = 'hacking',
  DONE = 'done',
}
export enum BlacklistType {
  BLACKLIST_REGEX = 'regex',
  BLACKLIST_VALUE = 'value',
}
export type Blacklist = {
  value: string,
  type: BlacklistType,
}

export interface Project {
  _id?: string;
  created_at?: Date;
  updated_at?: Date;
  name: string;
  description: string;
  status: ProjectStatus;
  blacklist: Blacklist[],
  whitelist: '',
};

class ProjectDb {
  FUZZ_PATH_NUM = 1;
  db: Collection<Document>;
  private static instance?: ProjectDb = undefined;

  constructor(dbo: Db) {
    this.db = dbo.collection('projects');
  }

  static getInstance(dbo?: Db): ProjectDb {
    if (!this.instance) {
      if (!dbo) {
        throw new Error('WTF? DBO is undefined?')
      }
      this.instance = new ProjectDb(dbo);
    }
    return this.instance;
  }

  async getHackingProjects() {
    const documents = await this.db.find({
      status: ProjectStatus.HACKING,
    });
    return (documents as any) as Project[];
  }

  async getRecentHackingProject() {
    const packetDb = PacketDb.getInstance();
    const hackingProjects = await this.getHackingProjects();
    const timestampUpdateProject = new Array<number>(hackingProjects.length).fill(0);
    let mostRecentTime = -1, cacheIndex = -1;
    for (let i = 0; i < timestampUpdateProject.length; i++) {
      const project = hackingProjects[i];
      const recentLog = await packetDb.getRecentDocumentByProject(project.name);
      if (!recentLog) {
        timestampUpdateProject[i] = -1;
      } else {
        timestampUpdateProject[i] = new Date(recentLog.updated_at).getTime() / 1000;
      }
      if (mostRecentTime < timestampUpdateProject[i]) {
        mostRecentTime = timestampUpdateProject[i];
        cacheIndex = i;
      }
    }
    return hackingProjects[cacheIndex];
  }

}

export default ProjectDb;
