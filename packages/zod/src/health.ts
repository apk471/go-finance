import { z } from "zod";
import { ZDateTime } from "./common.js";

const ZHealthCheck = z.object({
  status: z.string(),
  response_time: z.string(),
  error: z.string().optional(),
});

export const ZHealthResponse = z.object({
  status: z.enum(["healthy", "unhealthy"]),
  timestamp: ZDateTime,
  environment: z.string(),
  checks: z.object({
    database: ZHealthCheck,
    redis: ZHealthCheck.optional(),
  }),
});
