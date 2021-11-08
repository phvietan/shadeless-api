import ApiFuzzerPocGeneric, { ApiFuzzer } from '../api-fuzzer-poc-generic';
import { BotFuzzer } from 'libs/databases/botFuzzer.database';
import { ParsedPacket } from 'libs/databases/parsedPacket.database';
import { ChildProcessWithoutNullStreams, spawn } from 'child_process';
import { getHeaders } from 'libs/helper';

export default class AutoArjun
  extends ApiFuzzerPocGeneric
  implements ApiFuzzer
{
  constructor(botFuzzer: BotFuzzer, packet: ParsedPacket) {
    super(botFuzzer, packet, AutoArjun.name);
  }

  async runCommand(cmd: ChildProcessWithoutNullStreams): Promise<number> {
    return new Promise((resolve, reject) => {
      cmd.stdout.on('data', (data) => {
        console.log(`stdout: ${data}`);
      });

      cmd.stderr.on('data', (data) => {
        console.log(`stderr: ${data}`);
      });

      cmd.on('error', (error) => {
        console.log(`error: ${error.message}`);
        reject(error);
      });

      cmd.on('close', (code) => {
        console.log(`child process exited with code ${code}`);
        resolve(code);
      });
    });
  }

  async poc() {
    const cmdArjun = spawn('arjun', [
      '-u',
      `${this.packet.origin}${this.packet.path}`,
      '-c',
      '250',
      '--headers',
      getHeaders(this.packet.requestHeaders),
    ]);

    return this.runCommand(cmdArjun);
  }
}
