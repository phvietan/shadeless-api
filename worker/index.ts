import Bluebird from 'bluebird';
import PathFuzzer from './libs/fuzzer/path';
import { MongoClient, Db } from 'mongodb';
import { initDatabase } from './libs/database/index';
import ParsedPathDb from './libs/database/parsed_path';
import PacketDb from 'libs/database/packet';
import ProjectDb from 'libs/database/project';

const url = "mongodb://localhost:27017/";

MongoClient.connect(url, function(err, db) {
  if (err || !db) throw err;
  const dbo = db.db("shadeless");
  main(dbo);
});

async function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function main(dbo: Db) {
  initDatabase(dbo);
  const parsedPathDb = ParsedPathDb.getInstance();
  await parsedPathDb.resetScanning();

  const projectDb = ProjectDb.getInstance();
  const recentProject = await projectDb.getRecentHackingProject();

  while (true) {
    const fuzzPaths = await parsedPathDb.getTodo(recentProject);
    const options = await botPathDb.getOptions();
    // await Bluebird.map(fuzzPaths, async (path) => {
    //   console.log(`Fuzzing path: ${path.origin}${path.path}`);
    //   const pathFuzzer = new PathFuzzer(path);
    //   await pathFuzzer.run();
    //   await sleep(3000);
    // });
    await sleep(3000);
  }

}

