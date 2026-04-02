import { extendZodWithOpenApi } from "@anatine/zod-openapi";
import { z } from "zod";

extendZodWithOpenApi(z);

export * from "./utils.js";
export * from "./common.js";
export * from "./health.js";
export * from "./user.js";
export * from "./financial-record.js";
