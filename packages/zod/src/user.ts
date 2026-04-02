import { z } from "zod";
import { ZDateTime, ZHTTPError, ZUUID } from "./common.js";

export const ZUserRole = z.enum(["viewer", "analyst", "admin"]);
export const ZUserStatus = z.enum(["active", "inactive"]);

export const ZUser = z.object({
  id: ZUUID,
  authUserId: z.string(),
  email: z.string().email(),
  name: z.string(),
  role: ZUserRole,
  status: ZUserStatus,
  createdAt: ZDateTime,
  updatedAt: ZDateTime,
});

export const ZBootstrapUserBody = z.object({
  email: z.string().email(),
  name: z.string().min(2).max(100),
});

export const ZCreateUserBody = z.object({
  authUserId: z.string().min(1),
  email: z.string().email(),
  name: z.string().min(2).max(100),
  role: ZUserRole,
  status: ZUserStatus.default("active"),
});

export const ZUpdateUserBody = z
  .object({
    email: z.string().email().optional(),
    name: z.string().min(2).max(100).optional(),
    role: ZUserRole.optional(),
    status: ZUserStatus.optional(),
  });

export const ZUserList = z.array(ZUser);

export const ZUserErrorResponse = ZHTTPError;
