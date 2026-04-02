import { z } from "zod";

export const ZUUID = z.string().uuid();

export const ZDate = z.string().regex(/^\d{4}-\d{2}-\d{2}$/);

export const ZDateTime = z.string().datetime();

export const ZDecimalString = z.string();

export const ZFieldError = z.object({
  field: z.string(),
  error: z.string(),
});

export const ZAction = z.object({
  type: z.string(),
  message: z.string(),
  value: z.string(),
});

export const ZHTTPError = z.object({
  code: z.string(),
  message: z.string(),
  status: z.number().int(),
  override: z.boolean(),
  errors: z.array(ZFieldError).default([]),
  action: ZAction.nullable().optional(),
});
