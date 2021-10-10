import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from "axios";
import ParsedPacketDb, { ParsedPacket } from "../database/parsed_packet";
import ParsedPathDb, { ParsedPath, PathStatus } from "../database/parsed_path";
import HttpRequest from '../http_request';
import Bluebird from 'bluebird';
import PathFilter from "./path_filter";

class PathFuzzer {
  parsedPath: ParsedPath;
  cacheParsedPacket?: ParsedPacket;

  constructor(parsedPath: ParsedPath) {
    this.parsedPath = parsedPath;
  }

  private async getRequestsOptions(): Promise<AxiosRequestConfig[]> {
    this.cacheParsedPacket = await ParsedPacketDb.getInstance().getOneByRequestId(this.parsedPath.requestPacketId);
    if (!this.cacheParsedPacket) return [];
    return HttpRequest.createForDirPath(this.cacheParsedPacket, this.parsedPath.path);
  }

  private getParsedPathFromResponse(results: AxiosResponse<any>[]): ParsedPath[] {
    const cacheParsedPacket = this.cacheParsedPacket as ParsedPacket;
    return results.map(res => {
      return {
        type: '',
        requestPacketId: '',
        origin: cacheParsedPacket.origin,
        path: res.config.url as string,
        force: false,
        status: PathStatus.TODO,
        project: cacheParsedPacket.project,
        created_at: new Date(),
        updated_at: new Date(),
        error: '',
      };
    });
  }

  async run() {
    await ParsedPathDb.getInstance().update({
      _id: this.parsedPath._id,
    }, {
      status: PathStatus.SCANNING,
    });
    const requestsOptions = await this.getRequestsOptions();
    let cnt = 0;
    const responses = await Bluebird.map(requestsOptions, async (opts) => {
      cnt += 1;
      if (cnt % 30 === 0) {
        console.log(`Done ${cnt}/${requestsOptions.length}: ${opts.baseURL}${opts.url}`);
      }
      try {
        const resp = await axios.request(opts);
        return resp;
      }
      catch (err: any) {
        const error = err as AxiosError<any>;
        // The request was made and the server responded with a status code that falls out of the range of 2xx
        if (error.response) {
          const { response } = error;
          return response;
        } else if (error.request) { // The request was made but no response was received
          return null;
        } else {
          // Something happened in setting up the request that triggered an Error
          console.log('WTF error in config? ', error.message);
          return null;
        }
      }
    }, { concurrency: 3 });
    if (responses[0] === null) {
      const err = "Got error when GET random 404 page, please run again or recheck";
      console.log(err);
      await ParsedPathDb.getInstance().updateError({ "requestPacketId": this.parsedPath.requestPacketId }, err);
      await ParsedPathDb.getInstance().update({
        _id: this.parsedPath._id,
      }, {
        status: PathStatus.SCANNING,
      });
      return;
    }
    const pathFilter = new PathFilter(responses.filter(res => res !== null) as AxiosResponse<any>[]);
    const result = await pathFilter.filter(`Fuzzing path: ${this.parsedPath.origin}${this.parsedPath.path}`);
    const dbResult = this.getParsedPathFromResponse(result);
    await ParsedPathDb.getInstance().insertResult(dbResult);
    await ParsedPathDb.getInstance().update({
      _id: this.parsedPath._id,
    }, {
      status: PathStatus.DONE,
    });
    console.log("Done")
  }
}

export default PathFuzzer;
