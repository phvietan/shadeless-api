import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from "axios";
import ParsedPacketDb, { ParsedPacket } from "../database/parsed_packet";
import ParsedPathDb, { ParsedPath, PathStatus, PathType } from "../database/parsed_path";
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
    return HttpRequest.createOptsForDirPath(this.cacheParsedPacket, this.parsedPath.path);
  }

  private getParsedPathFromResponse(results: AxiosResponse<any>[]): ParsedPath[] {
    const cacheParsedPacket = this.cacheParsedPacket as ParsedPacket;
    return results.map(res => {
      return {
        type: PathType.NONE,
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
    const responses = await HttpRequest.sendAllOpts(requestsOptions);
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
