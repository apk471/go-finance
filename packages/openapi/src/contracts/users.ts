import { initContract } from "@ts-rest/core";
import {
  ZBootstrapUserBody,
  ZHTTPError,
  ZUpdateUserBody,
  ZUser,
  ZUserList,
  ZUUID,
  ZCreateUserBody,
} from "@boilerplate/zod";
import { z } from "zod";
import { getSecurityMetadata } from "@/utils.js";

const c = initContract();

export const usersContract = c.router({
  getCurrentUser: {
    summary: "Get current user",
    path: "/api/v1/users/me",
    method: "GET",
    description: "Returns the authenticated application's current active user.",
    metadata: getSecurityMetadata(),
    responses: {
      200: ZUser,
      401: ZHTTPError,
      403: ZHTTPError,
    },
  },
  bootstrapUser: {
    summary: "Bootstrap initial admin user",
    path: "/api/v1/users/bootstrap",
    method: "POST",
    description: "Creates the very first application user as an admin.",
    metadata: getSecurityMetadata(),
    body: ZBootstrapUserBody,
    responses: {
      201: ZUser,
      400: ZHTTPError,
      403: ZHTTPError,
    },
  },
  listUsers: {
    summary: "List users",
    path: "/api/v1/users",
    method: "GET",
    description: "Lists all provisioned users. Admin only.",
    metadata: getSecurityMetadata(),
    responses: {
      200: ZUserList,
      401: ZHTTPError,
      403: ZHTTPError,
    },
  },
  getUserById: {
    summary: "Get user by id",
    path: "/api/v1/users/:id",
    method: "GET",
    description: "Fetches a single user by UUID. Admin only.",
    metadata: getSecurityMetadata(),
    pathParams: z.object({
      id: ZUUID,
    }),
    responses: {
      200: ZUser,
      401: ZHTTPError,
      403: ZHTTPError,
      404: ZHTTPError,
    },
  },
  createUser: {
    summary: "Create user",
    path: "/api/v1/users",
    method: "POST",
    description: "Creates a new user. Admin only.",
    metadata: getSecurityMetadata(),
    body: ZCreateUserBody,
    responses: {
      201: ZUser,
      400: ZHTTPError,
      401: ZHTTPError,
      403: ZHTTPError,
    },
  },
  updateUser: {
    summary: "Update user",
    path: "/api/v1/users/:id",
    method: "PATCH",
    description: "Updates an existing user. Admin only.",
    metadata: getSecurityMetadata(),
    pathParams: z.object({
      id: ZUUID,
    }),
    body: ZUpdateUserBody,
    responses: {
      200: ZUser,
      400: ZHTTPError,
      401: ZHTTPError,
      403: ZHTTPError,
      404: ZHTTPError,
    },
  },
});
