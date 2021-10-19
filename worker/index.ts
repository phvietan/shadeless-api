import Bluebird from 'bluebird';
import { MongoClient } from 'mongodb';
import PathFuzzer from './src/pathFuzzer/pathFuzzer';
import { initDatabase } from './src/database/index';
import ParsedPathDb from './src/database/parsed_path';
import ProjectDb from './src/database/project';
import BotPathDb from './src/database/bot_path';
import { sleep } from './src/libs/helper';

const url = "mongodb://localhost:27017/";

MongoClient.connect(url, function(err, db) {
  if (err || !db) throw err;
  const dbo = db.db("shadeless");
  initDatabase(dbo);
  main();
});

async function main() {
  const parsedPathDb = ParsedPathDb.getInstance();
  const botPathDb = BotPathDb.getInstance();
  const projectDb = ProjectDb.getInstance();

  await parsedPathDb.resetScanning();

  while (true) {
    const botPathRunning = await botPathDb.getRunningProject();
    await Bluebird.map(botPathRunning, async (bp) => {
      const project = await projectDb.getOneProjectByName(bp.project);
      const fuzzPaths = await parsedPathDb.getTodo(project);
      await Bluebird.map(fuzzPaths, async (path) => {
        console.log(`Fuzzing path: ${path.origin}${path.path}`);
        const pathFuzzer = new PathFuzzer(path, {
          ...bp,
        });
        await pathFuzzer.run();
        await sleep(3000);
      });
    });
    await sleep(5000);
  }
}

