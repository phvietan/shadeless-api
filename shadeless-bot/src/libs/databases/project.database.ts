import { Collection, Db } from 'mongodb';

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
  value: string;
  type: BlacklistType;
};

export interface Project {
  _id?: string;
  created_at?: Date;
  updated_at?: Date;
  name: string;
  description: string;
  status: ProjectStatus;
  blacklist: Blacklist[];
  whitelist: '';
}

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
        throw new Error('WTF? DBO is undefined?');
      }
      this.instance = new ProjectDb(dbo);
    }
    return this.instance;
  }

  async getOneProjectByName(projectName: string): Promise<Project> {
    const document = await this.db.findOne({ name: projectName });
    if (!document) return undefined;
    return document as any as Project;
  }
}

// Get filter by project for blacklist/whitelist
export function getFilterByProjectForBW(project: Project): any {
  const blacklistExact = project.blacklist.filter(
    (b) => b.type === BlacklistType.BLACKLIST_VALUE,
  );
  const blacklistRegex = project.blacklist.find(
    (b) => b.type === BlacklistType.BLACKLIST_REGEX,
  );

  const filter = {
    project: project.name,
    origin: {
      $nin: blacklistExact,
      $regex: project.whitelist,
    },
  };
  if (blacklistRegex) {
    filter.origin['$not'] = {
      $regex: blacklistRegex.value,
    };
  }
  return filter;
}

export default ProjectDb;
