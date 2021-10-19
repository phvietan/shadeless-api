import { ParsedPath } from "../database/parsed_path";
import { BotFuzzer } from "database/bot_fuzzer";
import ApiFuzzerSender from "./apiFuzzer_sender";

class ApiFuzzer {
  options: BotFuzzer;
  apiSender: ApiFuzzerSender;
  parsedPath: ParsedPath;

  constructor(parsedPath: ParsedPath, options: BotFuzzer) {
    this.apiSender = new ApiFuzzerSender(options);
    this.parsedPath = parsedPath;
    this.options = options;
  }

}

export default ApiFuzzer;
