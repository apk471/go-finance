import { initContract } from "@ts-rest/core";
import { financialRecordsContract } from "./financial-records.js";
import { systemContract } from "./system.js";
import { usersContract } from "./users.js";

const c = initContract();

export const apiContract = c.router({
  System: systemContract,
  Users: usersContract,
  FinancialRecords: financialRecordsContract,
});
