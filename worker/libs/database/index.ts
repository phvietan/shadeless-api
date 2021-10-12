import { Db } from 'mongodb';
import ParsedPacketDb from './parsed_packet';
import ParsedPathDb from './parsed_path';
import { BlacklistType, Project } from './project';

export function initDatabase(dbo: Db) {
  ParsedPathDb.getInstance(dbo);
  ParsedPacketDb.getInstance(dbo);
}

// Get filter by project for blacklist/whitelist
export function getFilterByProjectForBW(project: Project): any {
  const blacklistExact = project.blacklist.filter(b => b.type === BlacklistType.BLACKLIST_VALUE);
  const blacklistRegex = project.blacklist.find(b => b.type === BlacklistType.BLACKLIST_REGEX);

  const filter = {
    project: project.name,
    origin: {
      $nin: blacklistExact,
      $regex: project.whitelist,
    }
  }
  if (blacklistRegex) {
    filter.origin["$not"] = {
      $regex: blacklistRegex,
    }
  }
  return filter;
}
